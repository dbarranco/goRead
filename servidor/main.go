// servidor project main.go
/*
	Primitivas:
		- Login
		- Sign up
		- Lista usuario
		- Lista completa
		- Búsqueda
		- Leer (libro, página)
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var mapa map[int]int = make(map[int]int)
var online map[string]int = make(map[string]int)

func handler_login(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query()["email"][0]
	contra := r.URL.Query()["passwd"][0]

	if online[email] == 0 {
		id := -1
		err := db.QueryRow("SELECT idUsuario FROM Usuario WHERE email=? AND passwd=?", email, contra).Scan(&id)
		switch {
		case err == sql.ErrNoRows:
			buffer := make([]byte, 1)
			buffer[0] = '0'
			log.Println("Error de autenticación del usuario: ", email)
			w.WriteHeader(http.StatusForbidden)
			w.Write(buffer)
		case err != nil:
			log.Println(err)
		default:
			fmt.Println("Sesión iniciada para el usuario ", email)
			var token int = int(rand.Int31())
			for (token == 0) || (mapa[token] == 0) {
				token = int(rand.Int31())
				mapa[token] = id
				online[email] = token
			}
			buffer := []byte(strconv.Itoa(token))
			w.Write(buffer)
		}
	} else {
		buffer := []byte(strconv.Itoa(online[email]))
		w.Write(buffer)
	}

}

func handler_signup(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query()["email"][0]
	passwd := r.URL.Query()["passwd"][0]

	res, err := db.Exec("INSERT INTO Usuario (email, passwd) VALUES(?,?)", email, passwd)
	log.Println("Resultado consulta: ", res)
	if err != nil {
		log.Println("Error al crear usuario: ", err)
	}
	buffer := []byte("OK")
	w.Write(buffer)
}

func handler_lista(w http.ResponseWriter, r *http.Request) {

}

func handler_buscar(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var err error
	// MySQL
	db, err = sql.Open("mysql", "app:app@tcp(150.214.182.97:3306)/goread")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos: ", err)
	}

	// HTTP
	http.Handle("/", http.FileServer(http.Dir("../web")))
	http.Handle("/libros/", http.FileServer(http.Dir(".")))

	//http.Handle("/libros", http.StripPrefix("/libros", http.FileServer(http.Dir("./libros"))))

	http.HandleFunc("/login", handler_login)
	http.HandleFunc("/signup", handler_signup)
	http.HandleFunc("/lista", handler_lista)
	http.HandleFunc("/buscar", handler_buscar)
	http.ListenAndServe(":80", nil)
}
