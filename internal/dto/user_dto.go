package dto

// Request Format
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
	Bio      string `json:"bio"`
}

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	FullName       string `json:"fullname" binding:"required"`
	Password       string `json:"password" binding:"required,min=6"`
	PasswordVerify string `json:"password_verify" binding:"required,eqfield=Password"`
}

// Data format
type SignInSuccess struct {
	Token string `json:"token"`
}

type RegisterSuccess struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}

type ProfileResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"FullName"`
	Username       string `json:"username"`
	ProfilePicture string `json:"ProfilePicture"`
}

type FormatUsers struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"ProfilePicture"`
	FullName       string `json:"fullname"`
}
