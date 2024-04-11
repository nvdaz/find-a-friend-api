package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tursodatabase/go-libsql"
)

func upsertInterests(tx *sql.Tx, userId string, interests []Interest) error {
	if len(interests) == 0 {
		return nil
	}

	stmt, err := tx.Prepare("INSERT INTO interests (user, interest, intensity, skill) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, interest := range interests {
		if _, err := stmt.Exec(userId, interest.Interest, interest.Intensity, interest.Skill); err != nil {
			return err
		}
	}

	return nil
}

type Interest struct {
	Interest  string  `json:"interest"`
	Intensity float64 `json:"intensity"`
	Skill     float64 `json:"skill"`
}

type User struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Personality Personality `json:"personality"`
	Interests   []Interest  `json:"interests"`
}
type Personality struct {
	Extraversion      float64 `json:"extraversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Neuroticism       float64 `json:"neuroticism"`
	Openness          float64 `json:"openness"`
}

type PostUser struct {
	Id          string      `param:"id"`
	Name        string      `json:"name"`
	Personality Personality `json:"personality"`
	Interests   []Interest  `json:"interests"`
}

func postUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := PostUser{}

		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
		}

		tx, err := db.Begin()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error starting transaction")
		}
		defer tx.Rollback()

		_, err = tx.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.Id, user.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error inserting user")
		}

		_, err = tx.Exec(
			`INSERT INTO personalities (id, extraversion, agreeableness, conscientiousness, neuroticism, openness)
			VALUES (?, ?, ?, ?, ?, ?)`,
			user.Id, user.Personality.Extraversion, user.Personality.Agreeableness,
			user.Personality.Conscientiousness, user.Personality.Neuroticism, user.Personality.Openness,
		)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error inserting personality")
		}

		if err := upsertInterests(tx, user.Id, user.Interests); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error inserting interests")
		}

		if err := tx.Commit(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error committing transaction")
		}

		return c.NoContent(http.StatusCreated)

	}
}

type UpdateUser struct {
	Id          string       `param:"id"`
	Name        *string      `json:"name"`
	Personality *Personality `json:"personality"`
	Interests   *[]Interest  `json:"interests"`
}

func updateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := UpdateUser{}

		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
		}

		tx, err := db.Begin()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error starting transaction")
		}
		defer tx.Rollback()

		row := tx.QueryRow("SELECT id FROM users WHERE id = ?", user.Id)

		var id string
		if err := row.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "error querying user")
		}

		if user.Name != nil {
			_, err = tx.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, user.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error updating user")
			}
		}

		if user.Personality != nil {
			_, err = tx.Exec(
				`UPDATE personalities
			SET extraversion = ?, agreeableness = ?, conscientiousness = ?, neuroticism = ?, openness = ?
			WHERE id = ?`,
				user.Personality.Extraversion, user.Personality.Agreeableness,
				user.Personality.Conscientiousness, user.Personality.Neuroticism, user.Personality.Openness, user.Id,
			)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error updating personality")
			}
		}

		if user.Interests != nil {
			_, err = tx.Exec("DELETE FROM interests WHERE user = ?", user.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error deleting interests")
			}

			if err := upsertInterests(tx, user.Id, *user.Interests); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error inserting interests")
			}
		}

		if err := tx.Commit(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error committing transaction")
		}

		return c.NoContent(http.StatusOK)
	}
}

type GetUser struct {
	Id string `param:"id"`
}

func getUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		get_user := GetUser{}

		if err := c.Bind(&get_user); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
		}

		row := db.QueryRow(
			`SELECT u.id, u.name, p.extraversion, p.agreeableness, p.conscientiousness, p.neuroticism, p.openness
        	FROM users u
        	LEFT JOIN personalities p ON u.id = p.id
        	WHERE u.id = ?`,
			get_user.Id,
		)

		user := User{}

		if err := row.Scan(&user.Id, &user.Name, &user.Personality.Extraversion,
			&user.Personality.Agreeableness, &user.Personality.Conscientiousness,
			&user.Personality.Neuroticism, &user.Personality.Openness); err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "error querying user")
		}

		rows, err := db.Query("SELECT interest, intensity, skill FROM interests WHERE user = ?", user.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error querying interests")
		}
		defer rows.Close()

		var interests []Interest
		for rows.Next() {
			var interest Interest
			if err := rows.Scan(&interest.Interest, &interest.Intensity, &interest.Skill); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error scanning interests")
			}
			interests = append(interests, interest)
		}

		user.Interests = interests

		return c.JSON(http.StatusOK, user)
	}
}

func getAllUsers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := db.Query(
			`SELECT u.id, u.name, p.extraversion, p.agreeableness, p.conscientiousness, p.neuroticism, p.openness
			FROM users u
			LEFT JOIN personalities p ON u.id = p.id`,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error querying users")
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			user := User{}
			if err := rows.Scan(&user.Id, &user.Name, &user.Personality.Extraversion,
				&user.Personality.Agreeableness, &user.Personality.Conscientiousness,
				&user.Personality.Neuroticism, &user.Personality.Openness); err != nil {
				fmt.Println("Error scanning users in getAllUsers:", user.Id)
				fmt.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "error scanning users")
			}

			rows, err := db.Query("SELECT interest, intensity, skill FROM interests WHERE user = ?", user.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error querying interests")
			}
			defer rows.Close()

			var interests []Interest
			for rows.Next() {
				var interest Interest
				if err := rows.Scan(&interest.Interest, &interest.Intensity, &interest.Skill); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "error scanning interests")
				}
				interests = append(interests, interest)
			}

			user.Interests = interests
			users = append(users, user)
		}

		return c.JSON(http.StatusOK, users)
	}

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
		libsql.WithSyncInterval(time.Minute),
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

	e := echo.New()

	e.Use(middleware.CORS())
	e.POST("/user", postUser(db))
	e.GET("/user/:id", getUser(db))
	e.POST("/user/:id", updateUser(db))
	e.GET("/users", getAllUsers(db))

	fmt.Println("Starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
