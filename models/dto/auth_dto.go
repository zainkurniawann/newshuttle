package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserDataOnLoginDTO struct {
	UserID    int64  `json:"user_id"`
	UserUUID  string `json:"user_uuid"`
	Username  string `json:"user_username"`
	RoleCode  string `json:"user_role_code"`
	Password  string `json:"user_password"`
}

type DeviceTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
