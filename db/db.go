package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewDB() (*sql.DB, error) {
	url := fmt.Sprintf("%s?authToken=%s", os.Getenv("DATABASE_URI"), os.Getenv("DATABASE_TOKEN"))

	return sql.Open("libsql", url)
}
