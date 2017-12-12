package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type e map[string]string

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

func main() {
	var err error

	url, ok := os.LookupEnv("DATABASE_URL")

	if !ok {
		log.Fatalln("$DATABASE_URL is required")
	}

	db, err = connect(url)

	if err != nil {
		log.Fatalf("Connection error: %s", err.Error())
	}

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			switch r.Method {
			case "GET":
				users, err := users(db)

				if err != nil {
					errorResponse(r, w, err)
				}

				respond(r, w, http.StatusOK, users)
			case "POST":
				user, err := unmarshalUser(r)

				if err != nil {
					errorResponse(r, w, err)
					return
				}

				created, err := createUser(db, user.Username)

				if err != nil {
					errorResponse(r, w, err)
					return
				}

				respond(r, w, http.StatusOK, created)
			default:
				respond(r, w, http.StatusNotFound, nil)
			}
		} else {
			respond(r, w, http.StatusNotFound, nil)
		}
	})

	log.Printf("Starting server on port %s\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id       SERIAL,
      username VARCHAR(64) NOT NULL UNIQUE,
      CHECK (CHAR_LENGTH(TRIM(username)) > 0)
    );
  `)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func users(db *sql.DB) ([]User, error) {
	rows, err := db.Query(
		`SELECT id, username FROM users ORDER BY username`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0, 5)

	for rows.Next() {
		u := User{}

		err = rows.Scan(&u.Id, &u.Username)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func unmarshalUser(r *http.Request) (*User, error) {
	defer r.Body.Close()

	var user User

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func createUser(db *sql.DB, username string) (*User, error) {
	created := User{}

	row := db.QueryRow(
		`INSERT INTO users (username) VALUES ($1) RETURNING id, username`,
		username,
	)

	err := row.Scan(&created.Id, &created.Username)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func respond(r *http.Request, w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		bytes, err := json.Marshal(data)

		if err != nil {
			errorResponse(r, w, err)
			return
		}

    response(r, w, status, bytes)
	} else {
    response(r, w, status, nil)
	}
}

func errorResponse(r *http.Request, w http.ResponseWriter, err error) {
	bytes, _ := json.Marshal(e{"message": err.Error()})

	response(r, w, http.StatusInternalServerError, bytes)
}

func response(r *http.Request, w http.ResponseWriter, status int, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

  if bytes != nil {
    bytes = append(bytes, '\n')
    w.Write(bytes)
  }

	log.Printf(
		"\"%s %s %s\" %d %d\n",
		r.Method, r.URL.Path, r.Proto, status, len(bytes),
	)
}
