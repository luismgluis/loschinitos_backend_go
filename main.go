package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/cors"
	//para los CORS
)

const (
	port = 9080        //9180 //8080
	host = "localhost" //127.0.0.1
)

var router *chi.Mux

func routers() *chi.Mux {

	router.Get("/", ping)

	router.Get("/clientes", AllClientes)
	router.Get("/alldata", AllClientes)
	router.Get("/cliente/{id}", ObtenerClienteByID)
	router.Post("/cliente/create", CrearCliente)
	router.Put("/cliente/update/{id}", ActualizarCliente)
	router.Delete("/cliente/{id}", EliminarCliente)

	return router
}

// server starting point
func ping(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "Pong"})
}

func main() {
	router = chi.NewRouter() //iniciamos el router
	//esto porque lo estamos invocando de otro dominio
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(middleware.Recoverer)
	routers() //pone las redirecciones
	http.ListenAndServe(":3000", Logger())
}
