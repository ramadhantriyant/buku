package models

import (
	"database/sql"

	"github.com/ramadhantriyant/buku/internal/database"
)

type Config struct {
	DB        *sql.DB
	Queries   *database.Queries
	JWTSecret string
}
