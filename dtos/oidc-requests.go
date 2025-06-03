package dtos

type SignUpRequest struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required" validate:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required,min=8"`
}
