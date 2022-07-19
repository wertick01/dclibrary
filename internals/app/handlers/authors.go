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

type AuthorsHandler struct {
	processor *processors.AuthorsProcessor
}

func NewAuthorsHandler(processor *processors.AuthorsProcessor) *AuthorsHandler {
	handler := new(AuthorsHandler)
	handler.processor = processor
	return handler
}

func (handler *AuthorsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newAuthor models.Authors

	err := json.NewDecoder(r.Body).Decode(&newAuthor)
	if err != nil {
		WrapError(w, err)
		return
	}

	_, err = handler.processor.CreateAuthor(newAuthor)
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

func (handler *AuthorsHandler) List(w http.ResponseWriter, r *http.Request) {
	//vars := r.URL.Query() ЗАЧЕМ ОНО ТУТ НАДО
	list, err := handler.processor.ListAuthors()

	if err != nil {
		WrapError(w, err)
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   list,
	}

	WrapOK(w, m)
}

func (handler *AuthorsHandler) Find(w http.ResponseWriter, r *http.Request) {
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

	user, err := handler.processor.FindAuthor(id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   user,
	}

	WrapOK(w, m)
}

func (handler *AuthorsHandler) FindById(w http.ResponseWriter, r *http.Request) {
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

	books, err := handler.processor.AuthorsBooks(id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   books,
	}

	WrapOK(w, m)
}
