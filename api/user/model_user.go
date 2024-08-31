package user

import "github.com/golang-jwt/jwt/v4"

type (
	User struct {
		UserID    int32  `json:"user_id"`
		Nama      string `json:"nama"`
		Email     string `json: "email"`
		Password  string `json:"password"`
		CreatedAt string `json:"created_at"`
		Role      string `json:"role"`
		Photo     string `json:"photo"`
	}
	UserRequest struct {
		Nama     string `json:"Nama" validate:"required"`
		Email    string `json:"Email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
		Photo    string `json:"photo"`
	}

	LoginRequest struct {
		Nama     string `json:"nama" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	Claims struct {
		Username string `json:"username"`
		Role     string `json:"role"` // Tambahkan field role
		jwt.RegisteredClaims
	}
)
