package api

import (
	"dclibrary.com/internals/app/handlers"
	"github.com/gorilla/mux"
)

func CreateRoutes(userHandler *handlers.UsersHandler, carsHandler *handlers.BooksHandler) *mux.Router {

	r := mux.NewRouter()                       //создадим роутер для обработки путей, он же будет основным роутером для нашего сервера
	r.HandleFunc("/api/login").Methods("POST") //каждая функция реализует один и тот же интерфейс
	r.HandleFunc("/api/registration").Methods("POST")
	r.HandleFunc("/api/send-code").Methods("POST")
	r.HandleFunc("/api/refresh-token").Methods("POST")
	r.HandleFunc("/api/upload").Methods("POST")
	r.HandleFunc("/api/img/upload").Methods("POST")
	r.HandleFunc("/api/img/destroy").Methods("POST")
	r.HandleFunc("/api/whoami").Methods("GET")

	r.HandleFunc("/api/users", userHandler.Create).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler.Find).Methods("GET")
	r.HandleFunc("/api/users/self").Methods("GET")
	r.HandleFunc("/api/users").Methods("POST")
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler.Change).Methods("PUT")
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler.Delete).Methods("DELETE")
	r.HandleFunc("/api/users/self").Methods("PUT")

	r.HandleFunc("/api/book").Methods("POST")
	r.HandleFunc("/api/book").Methods("GET")
	r.HandleFunc("/api/book/{id:[0-9]+}").Methods("GET")
	r.HandleFunc("/api/book/{id:[0-9]+}").Methods("PUT")
	r.HandleFunc("/api/book/{id:[0-9]+}").Methods("DELETE")
	r.HandleFunc("/api/book/reverse").Methods("POST")

	r.HandleFunc("/api/favorietes/list").Methods("GET")
	r.HandleFunc("/api/favorietes/add").Methods("POST")
	r.HandleFunc("/api/favorietes/delete").Methods("DELETE")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler() //оборачиваем 404, для обработки NotFound
	return r
}
