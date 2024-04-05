package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

type User struct {
	Id  string
	Name string
}

func userGetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		rows, err := db.Query("SELECT * FROM users WHERE id = ?", id)

		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		user := User{}

		if !rows.Next() {
			w.Write([]byte("User not found"))
			return
		}
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		res, err := json.Marshal(user)

		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(res)

	}
}

type Personality struct {
	Id string
	Extraversion float64
	Agreeableness float64
	Conscientiousness float64
	Neuroticism float64
	Openness float64
}

func personalityGetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		rows, err := db.Query("SELECT * FROM personalities WHERE id = ?", id)

		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		personality := Personality{}

		personality.Id = id

		if !rows.Next() {
			goto ret
		}

		if err := rows.Scan(&personality.Id, &personality.Extraversion, &personality.Agreeableness, &personality.Conscientiousness, &personality.Neuroticism, &personality.Openness); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		ret:

		res, err := json.Marshal(personality)

		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(res)

	}

}

type pingHandler struct{}


func (h *pingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	primaryUrl := os.Getenv("DATABASE_URL")
	dbName := "find-a-friend"
	authToken := os.Getenv("DATABASE_TOKEN")
	dir, err := os.MkdirTemp("", "libsql-*")
    if err != nil {
        fmt.Println("Error creating temporary directory:", err)
        os.Exit(1)
    }
    defer os.RemoveAll(dir)

    dbPath := filepath.Join(dir, dbName)

    connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
        libsql.WithAuthToken(authToken),
    )
    if err != nil {
        fmt.Println("Error creating connector:", err)
        os.Exit(1)
    }
    defer connector.Close()


	if err != nil {
		fmt.Println(err)
		return
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	mux := http.NewServeMux()

	mux.Handle("/ping", &pingHandler{})
	mux.Handle("/user/{id}", userGetHandler(db))
	mux.Handle("/personality/{id}", personalityGetHandler(db))


	fmt.Println("Starting server...")


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":" + port, mux)
}


