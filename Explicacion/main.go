//https://www.gorillatoolkit.org/
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
	Para procesar peticiones http, el paquete http de Go cuenta con dos componentes principales que son el server mux y el handler.
	Un server mux es un enrutador de peticiones http. El segundo parametro de ListenAndServe es el sever mux.
*/
func main() {
	msg := mensaje{
		msg: "<h1>Hola mundo de nuevo</h1>",
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))
	//handleFunc utiliza el server mux interno que tiene el paquete http
	mux.Handle("/", fs)
	mux.HandleFunc("/prueba", prueba)
	mux.HandleFunc("/usuario", usuario)
	mux.Handle("/hola", msg)
	/*
		Crea un nuevo servidor (una estructura server) y luego la devuelve haciendo una llamada al metodo ListenAndServe de dicha estructura
		func ListenAndServe(addr string, handler Handler) error {
			server := &Server{Addr: addr, Handler: handler}
			return server.ListenAndServe()
		}
		http.ListenAndServe(":8080", mux)
		a continuacion creamos nuesta propia estructura server
	*/
	// con & indicamos que queremos crear un puntero de una estructura Server.
	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second, // Le digo que espere 10 segundos para leer la peticion, si pasa ese tiempo y no se leyó, se cancela la peticion
		WriteTimeout:   10 * time.Second, // Y que espere 10 segundos para la respuesta. Esto es muy importante ya que muchas peticiones de clientes lentos
		MaxHeaderBytes: 1 << 20,          // podrian afectar el funcionamiento de nuestro servidor
	} // MaxHeaderBytes: Aca le estamos pasando que el tamaño maximo de la cabecera es de 1mg. Con el operador binario << lo que hace es elevar 2 al cuadrado 20 veces y al dividir por 1024 tenemos 1 mega
	log.Println("Listening...")
	log.Fatal(server.ListenAndServe()) //Si hay un error, se para el servidor
}

func holaMundo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hola Mundo</h1>")
}

func prueba(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hola Mundo desde /prueba</h1>")
}

func usuario(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hola Mundo desde /usuario</h1>")
}

//Vamos a crear un handler por completo, es decir crear una estructura que implemente el metodo serveHTTP para que veamos como lo hacemos en el caso que ya el objeto que va a manejar la peticion sea un handler y no una funcion que va a trabajar como un handler
type mensaje struct {
	msg string
}

// utiliza un recibidor (m) la llamamos ServeHTTP (Con wr trae todo)
func (m mensaje) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, m.msg)
}
