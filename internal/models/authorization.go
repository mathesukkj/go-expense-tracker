package models

type LoginPayload struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	id uint
}
