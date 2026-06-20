// internal/model/user.go
package model

import "time"

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	AvatarURL    *string    `json:"avatar_url"`
	SpaceCredits *int       `json:"space_credits"`
	CreatedAt    *time.Time `json:"created_at"`
}