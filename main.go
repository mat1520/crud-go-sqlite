package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

type PageData struct {
	Users []User
	User  User // para editar
}

var tmpl *template.Template
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear tabla si no existe
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL
    );`)
	if err != nil {
		log.Fatal(err)
	}

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	// Servir archivos estáticos desde /static/
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/users", createUserHandler)
	http.HandleFunc("/users/delete/", deleteUserHandler)
	http.HandleFunc("/users/edit/", editOrUpdateUserHandler)

	addr := ":8080"
	log.Printf("Servidor escuchando en http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, password_hash FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash); err != nil {
			http.Error(w, "DB scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	data := PageData{Users: users}
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Error al ejecutar la plantilla: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error: "+err.Error(), http.StatusBadRequest)
		return
	}
	username := strings.TrimSpace(r.FormValue("name"))
	password := r.FormValue("password") // temporal: sin hash
	if username == "" || password == "" {
		http.Error(w, "username y password requeridos", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users(username, password_hash) VALUES(?, ?)", username, password)
	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	idStr := path.Base(r.URL.Path) // extrae {id} de /users/delete/{id}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "DB delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func editOrUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := path.Base(r.URL.Path) // /users/edit/{id}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		row := db.QueryRow("SELECT id, username, password_hash FROM users WHERE id = ?", id)
		var u User
		if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash); err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data := PageData{User: u}
		if err := tmpl.ExecuteTemplate(w, "edit.html", data); err != nil {
			http.Error(w, "Error al ejecutar la plantilla: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// POST -> actualizar
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Parse error: "+err.Error(), http.StatusBadRequest)
			return
		}
		username := strings.TrimSpace(r.FormValue("name"))
		password := r.FormValue("password") // temporal: sin hash
		if username == "" || password == "" {
			http.Error(w, "username y password requeridos", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("UPDATE users SET username = ?, password_hash = ? WHERE id = ?", username, password, id)
		if err != nil {
			http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
