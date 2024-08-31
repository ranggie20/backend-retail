package consumer

import (
	"database/sql"

	repo "github.com/online-bnsp/backend/repo/generated"
)

type Handler struct {
	db    *sql.DB
	model *repo.Queries
}

func New(db *sql.DB) *Handler {
	dbGenerated := repo.New(db)

	return &Handler{db, dbGenerated}
}
