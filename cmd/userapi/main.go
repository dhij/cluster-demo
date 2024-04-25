package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ianschenck/envflag"
)

type User struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func main() {
	var (
		login    = envflag.String("MYSQL_USER", "root", "address where mysql db is listening")
		password = envflag.String("MYSQL_PASSWORD", "password", "address where mysql db is listening")
		host     = envflag.String("MYSQL_HOST", "mysql", "address where mysql db is listening")
		port     = envflag.String("MYSQL_PORT", "3306", "address where mysql db is listening")
		database = envflag.String("MYSQL_DATABASE", "cluster_demo", "address where mysql db is listening")
		dbDriver = "mysql"
	)
	envflag.Parse()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *login, *password, *host, *port, *database)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalf("connecting to MySQL: %s", err)
	}
	defer db.Close()
	log.Println("successfully connected to MySQL")

	uh := userHandler{
		ctx: context.Background(),
		db:  db,
	}
	http.Handle("/users", uh)
	http.ListenAndServe(":8080", nil)
}

type userHandler struct {
	ctx context.Context
	db  *sql.DB
}

func (uh userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.getUsers(w, r)
	case "POST":
		uh.createUser(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

var (
	insertUserQuery = "INSERT INTO users (uuid, name) VALUES (?, ?);"
	getUsersQuery   = "SELECT uuid, name from users;"
)

func (uh userHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	stmt, err := uh.db.PrepareContext(uh.ctx, insertUserQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(uh.ctx, &u.UUID, &u.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user added"))
}

func (uh userHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := uh.db.QueryContext(uh.ctx, getUsersQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UUID, &u.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
