package dto

type SignInRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceType  string `json:"device_type" binding:"required"`
	DeviceToken string `json:"device_token" binding:"required"`
}