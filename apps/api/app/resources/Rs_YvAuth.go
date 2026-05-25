package resources

import "github.com/revianto/yava/api/app/models"

type UserResponse struct {
	Id        int64   `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
}

func UserResource(u models.YvUser) UserResponse {
	return UserResponse{Id: u.Id, Email: u.Email, Name: u.Name, AvatarUrl: u.AvatarUrl}
}
