package payment

import (
	"github.com/go-playground/validator/v10"
	repo "github.com/online-bnsp/backend/repo/generated"
)

type Handler struct {
	validate *validator.Validate
	db       *repo.Queries
}

func NewHandler(validate *validator.Validate, db *repo.Queries) *Handler {
	return &Handler{validate, db}
}
