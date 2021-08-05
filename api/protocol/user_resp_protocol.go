package protocol

import "time"

type LoginResp struct {
	UserId         string    `json:"user_id"`
	Token          string    `json:"token"`
	InvitationCode string    `json:"invitation_code"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	Phone          string    `json:"phone"`
	LastLogin      time.Time `json:"last_login"`
	LoginCount     int32     `json:"login_count"`
}
