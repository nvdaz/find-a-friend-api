package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

type WrappedDB struct {
	Db        *sql.DB
	connector *libsql.Connector
	tmpDir    string
}

func (wrappedDb *WrappedDB) Close() {
	wrappedDb.Db.Close()
	wrappedDb.connector.Close()
	os.RemoveAll(wrappedDb.tmpDir)
}

func NewDB() (*WrappedDB, error) {
	dbName := "local.db"
	url := os.Getenv("DATABASE_URI")
	authToken := os.Getenv("DATABASE_TOKEN")

	tmpDir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(tmpDir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, url,
		libsql.WithAuthToken(authToken),
	)

	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)

	return &WrappedDB{db, connector, tmpDir}, nil
}
