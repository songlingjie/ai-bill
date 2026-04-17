package data

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	OpenID    string    `gorm:"size:128;uniqueIndex" json:"openId"`
	NickName  string    `gorm:"size:128" json:"nickName"`
	AvatarURL string    `gorm:"size:512" json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
