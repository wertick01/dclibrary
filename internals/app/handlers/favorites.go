package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"dclibrary.com/internals/app/models"
	"dclibrary.com/internals/app/processors"
	"github.com/gorilla/mux"
)

type FavorietesHandler struct {
	processor *processors.FavorietesProcessor
}

func NewFavorietesHandler(processor *processors.BooksProcessor) *BooksHandler {
	handler := new(BooksHandler)
	handler.processor = processor
	return handler
}

func (handler *FavorietesHandler) AddFavorieteBook(w http.ResponseWriter, r *http.Request) {
	var favorieteBook *models.FavorieteBooks

	err := json.NewDecoder(r.Body).Decode(&favorieteBook)
	if err != nil {
		WrapError(w, err)
		return
	}

	_, err = handler.processor.AddFavorieteBook(favorieteBook)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)
}

func (handler *FavorietesHandler) AddFavorieteAuthor(w http.ResponseWriter, r *http.Request) {
	var favorieteAuthor *models.FavorieteAuthors

	err := json.NewDecoder(r.Body).Decode(&favorieteAuthor)
	if err != nil {
		WrapError(w, err)
		return
	}

	_, err = handler.processor.AddFavorieteAuthor(favorieteAuthor)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)
}

func (handler *FavorietesHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	//vars := r.URL.Query() ЗАЧЕМ ОНО ТУТ НАДО
	var user_id int64

	err := json.NewDecoder(r.Body).Decode(&user_id)
	if err != nil {
		WrapError(w, err)
		return
	}
	list, err := handler.processor.ListFavorieteBooks(user_id)

	if err != nil {
		WrapError(w, err)
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   list,
	}

	WrapOK(w, m)
}

func (handler *FavorietesHandler) ListAuthors(w http.ResponseWriter, r *http.Request) {
	//vars := r.URL.Query() ЗАЧЕМ ОНО ТУТ НАДО
	var user_id int64

	err := json.NewDecoder(r.Body).Decode(&user_id)
	if err != nil {
		WrapError(w, err)
		return
	}
	list, err := handler.processor.DeleteFromFavorieteAuthors(user_id)

	if err != nil {
		WrapError(w, err)
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   list,
	}

	WrapOK(w, m)
}

func (handler *FavorietesHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //переменные, обьявленные в ресурсах попадают в Vars и могут быть адресованы
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}

	deletedbook, err := handler.processor.DeleteFromFavorieteBooks(id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   deletedbook,
	}

	WrapOK(w, m)
}

func (handler *FavorietesHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //переменные, обьявленные в ресурсах попадают в Vars и могут быть адресованы
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}

	deletedbook, err := handler.processor.DeleteFromFavorieteAuthors(id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   deletedbook,
	}

	WrapOK(w, m)
}
