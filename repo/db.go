package repo

import (
	"context"
	"database/sql"

	repo "github.com/online-bnsp/backend/repo/generated"
	// repo "github.com/online-bnsp/backend/repo/generated"
)

type Database struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	*repo.Queries
	db repo.DBTX
}

func (d *Database) CreateUser(ctx context.Context, username string, hashedPassword string, rolesJSON string) (int, error) {
	var userID int
	err := d.db.QueryRowContext(ctx, `
        INSERT INTO users (username, password, roles) 
        VALUES (?, ?, ?)
    `, username, hashedPassword, rolesJSON).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
