package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type WeChatService struct {
	AppID    string
	Secret   string
	MockMode bool
}

type code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func NewWeChatService(appID, secret string, mockMode bool) *WeChatService {
	return &WeChatService{
		AppID:    appID,
		Secret:   secret,
		MockMode: mockMode,
	}
}

func (s *WeChatService) Code2Session(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("微信登录 code 不能为空")
	}
	if s.MockMode {
		return "mock_" + code, nil
	}
	if s.AppID == "" || s.Secret == "" {
		return "", fmt.Errorf("微信登录配置缺失，请设置 WECHAT_APP_ID 和 WECHAT_APP_SECRET")
	}

	endpoint := "https://api.weixin.qq.com/sns/jscode2session?appid=" +
		url.QueryEscape(s.AppID) +
		"&secret=" + url.QueryEscape(s.Secret) +
		"&js_code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"

	fmt.Printf("微信登录请求: %s\n", endpoint)

	client, err := createHTTPClient()
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 客户端失败: %v", err)
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("微信接口请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取微信响应失败: %v", err)
	}

	fmt.Printf("微信登录响应: %s\n", string(respBody))

	var result code2SessionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析微信响应失败: %v, 响应内容: %s", err, string(respBody))
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信登录失败: %s, errCode: %d", result.ErrMsg, result.ErrCode)
	}
	if result.OpenID == "" {
		return "", fmt.Errorf("微信登录失败: 未获取到 openid")
	}
	return result.OpenID, nil
}

func createHTTPClient() (*http.Client, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		fmt.Printf("无法获取系统证书池: %v，尝试创建新的证书池\n", err)
		rootCAs = x509.NewCertPool()
	}

	// 尝试加载系统证书文件
	certFiles := []string{
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/pki/tls/certs/ca-bundle.crt",
		"/etc/ssl/ca-bundle.pem",
		"/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
	}

	for _, certFile := range certFiles {
		if _, err := os.Stat(certFile); err == nil {
			certData, err := os.ReadFile(certFile)
			if err == nil {
				if rootCAs.AppendCertsFromPEM(certData) {
					fmt.Printf("成功加载证书文件: %s\n", certFile)
				}
			}
		}
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{Transport: transport}, nil
}
