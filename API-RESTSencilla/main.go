package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Note es una estructura para la api
type Note struct {
	Title       string    `json:"title"` //notacion json
	Description string    `json:"description"`
	CreatedAt   time.Time `json: created-at`
}

var noteStore = make(map[string]Note)

var id int

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GetNote).Methods("GET")
	r.HandleFunc("/api/notes", PostNote).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNote).Methods("PUT") //contiene el id de la nota que vamos a actualizar
	r.HandleFunc("/api/notes/{id}", DeleteNote).Methods("DELETE")
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening http://localhost:8080 ...")
	server.ListenAndServe()
}

//GetNote lee las notas almacenadas
func GetNote(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	for _, v := range noteStore {
		notes = append(notes, v)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(notes) //convierte la estructura de notes a json
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// PostNote postea
func PostNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	//aca descodificamos el json que viene en el cuerpo de la peticion
	err := json.NewDecoder(r.Body).Decode(&note) //rellena nuestra estructura note
	if err != nil {
		panic(err)
	}
	note.CreatedAt = time.Now() //agregamos la fecha
	id++                        //incrementamos el id
	k := strconv.Itoa(id)       //Convertimos el id a string
	noteStore[k] = note         // le pasamos la nota al map

	//Devolvemos nuestra nota ya creada
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(note) //convierte la estructura de notes a json
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// PutNote actualiza
func PutNote(w http.ResponseWriter, r *http.Request) {
	// extrae las notas que vengan de la peticion y devuelve un slice con esas notas
	vars := mux.Vars(r)
	// contiene el id de la nota
	k := vars["id"] //el paquete mux convirtio todo lo que venia en el request en un slice de string que utiliza como indice el nombre de la variable
	// obtenemos los datos que el usuario nos manda, es decir los que se van a actualizar
	var noteUpdate Note
	// codificamos el json a una estructura de Note
	err := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil {
		panic(err)
	}
	// vemos si el id es correcto, si existe entonces lo actualiza
	if note, ok := noteStore[k]; ok {
		noteUpdate.CreatedAt = note.CreatedAt //a los datos que cargo el usuario le agregamos la fecha
		delete(noteStore, k)                  //el id que nos mando el usuario lo borramos para ahora pasarle la nota a actualizar
		noteStore[k] = note
	} else {
		log.Printf("No encontramos el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteNote borra
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	// extrae las notas que vengan de la peticion y devuelve un slice con esas notas
	vars := mux.Vars(r)
	// contiene el id de la nota
	k := vars["id"]
	// si el id existe lo borramos
	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
	} else {
		log.Printf("No encontramos el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)
}
