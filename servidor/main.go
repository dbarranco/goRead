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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if (len(r.URL.Query()["email"]) == 0) || (len(r.URL.Query()["passwd"]) == 0) {
		w.WriteHeader(http.StatusBadRequest)
		buffer := []byte("Parámetros incorrectos: " + r.URL.Query().Encode())
		log.Println("Parámetros incorrectos: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}
	email := r.URL.Query()["email"][0]
	passwd := r.URL.Query()["passwd"][0]

	res, err := db.Exec("INSERT INTO Usuario (email, passwd) VALUES(?,?)", email, passwd)
	log.Println("Resultado consulta: ", res)
	if err != nil {
		log.Println("Error al crear usuario: ", err)
		w.WriteHeader(http.StatusBadRequest)
		buffer := []byte(err.Error())
		w.Write(buffer)
	}
	buffer := []byte("OK")
	w.Write(buffer)
}

func handler_lista(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if len(r.URL.Query()["token"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		buffer := []byte("Parámetros incorrectos: " + r.URL.Query().Encode())
		log.Println("Parámetros incorrectos: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}
	token, err := strconv.Atoi(r.URL.Query()["token"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		buffer := []byte("Parámetros incorrectos: " + r.URL.Query().Encode())
		log.Println("Parámetros incorrectos: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}
	id := mapa[token]

	//id := mapa[id]
	if id == 0 {
		w.WriteHeader(http.StatusForbidden)
		buffer := []byte("Usuario incorrecto: " + r.URL.Query().Encode())
		log.Println("Usuario incorrecto: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}

	log.Println("Usuario ", strconv.Itoa(id), " pide lista de libros.")
	rows, err := db.Query("SELECT Libros.idLibro, Titulo,Descripcion,Creador,Idioma,Ano FROM userLibros,Libros WHERE userLibros.idUsuario=" + strconv.Itoa(id) + " AND userLibros.idLibro=Libros.idLibro")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		buffer := []byte("Error del servidor: " + err.Error())
		log.Println("Error del servidor: " + err.Error())
		w.Write(buffer)
		return
	}

	defer rows.Close()

	//Principio de la tabla:
	buffer := []byte("<table class=\"table\"><thead><tr><th>Título</th><th>Autor</th><th>Año</th><th>Idioma</th><th>Descripción</th></tr></thead><tbody>")
	w.Write(buffer)

	for rows.Next() {
		var id string
		var titulo string
		var descripcion string
		var creador string
		var idioma string
		var ano string
		if err := rows.Scan(&id, &titulo, &descripcion, &creador, &idioma, &ano); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			buffer := []byte("Error del servidor: " + err.Error())
			w.Write(buffer)
			log.Println("Error del servidor: " + err.Error())
			return
		}

		buffer := []byte(fmt.Sprintf("<tr><td><a href=\"/finalRead.html?libro=%s\">%s</a></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", id, titulo, creador, ano, idioma, descripcion))
		w.Write(buffer)
	}
	buffer = []byte("</tbody></table>")
	w.Write(buffer)
	err = rows.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		buffer := []byte("Error del servidor: " + err.Error())
		log.Println("Error del servidor: " + err.Error())
		w.Write(buffer)
		return
	}
}

func handler_descubrir(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if len(r.URL.Query()["token"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		buffer := []byte("Parámetros incorrectos: " + r.URL.Query().Encode())
		log.Println("Parámetros incorrectos: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}
	token, err := strconv.Atoi(r.URL.Query()["token"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		buffer := []byte("Parámetros incorrectos: " + r.URL.Query().Encode())
		log.Println("Parámetros incorrectos: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}
	id := mapa[token]

	//id := mapa[id]
	if id == 0 {
		w.WriteHeader(http.StatusForbidden)
		buffer := []byte("Usuario incorrecto: " + r.URL.Query().Encode())
		log.Println("Usuario incorrecto: " + r.URL.Query().Encode())
		w.Write(buffer)
		return
	}

	log.Println("Usuario ", strconv.Itoa(id), " pide lista de libros.")
	rows, err := db.Query("SELECT * FROM Libros")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		buffer := []byte("Error del servidor: " + err.Error())
		log.Println("Error del servidor: " + err.Error())
		w.Write(buffer)
		return
	}

	defer rows.Close()

	//Principio de la tabla:
	buffer := []byte("<table class=\"table\"><thead><tr><th>Título</th><th>Autor</th><th>Año</th><th>Idioma</th><th>Descripción</th></tr></thead><tbody>")
	w.Write(buffer)

	for rows.Next() {
		var id string
		var titulo string
		var descripcion string
		var creador string
		var idioma string
		var ano string
		if err := rows.Scan(&id, &titulo, &descripcion, &creador, &idioma, &ano); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			buffer := []byte("Error del servidor: " + err.Error())
			w.Write(buffer)
			log.Println("Error del servidor: " + err.Error())
			return
		}

		buffer := []byte(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", titulo, creador, ano, idioma, descripcion))
		w.Write(buffer)
	}
	buffer = []byte("</tbody></table>")
	w.Write(buffer)
	err = rows.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		buffer := []byte("Error del servidor: " + err.Error())
		log.Println("Error del servidor: " + err.Error())
		w.Write(buffer)
		return
	}
}

func main() {
	var err error
	// MySQL
	db, err = sql.Open("mysql", "app:app@tcp(150.214.182.97:3306)/goread")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos: ", err)
	}

	// HTTP
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.Handle("/libros/", http.FileServer(http.Dir(".")))

	//http.Handle("/libros", http.StripPrefix("/libros", http.FileServer(http.Dir("./libros"))))

	http.HandleFunc("/login", handler_login)
	http.HandleFunc("/signup", handler_signup)
	http.HandleFunc("/lista", handler_lista)
	http.HandleFunc("/descubrir", handler_descubrir)
	http.ListenAndServe(":80", nil)
}
