package web

import (
	"encoding/json"
	"gobooks/internal/service"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	service *service.BookService
}

func NewBookHandlers(service *service.BookService) *BookHandlers {
	return &BookHandlers{service: service}
}

func (handler BookHandlers) GetBooks(res http.ResponseWriter, req *http.Request) {
	books, err := handler.service.GetBooks()

	if err != nil {
		http.Error(res, "Failed to get books", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Context-Type", "application/json")
	json.NewEncoder(res).Encode(books)
}

func (handler *BookHandlers) CreateBook(res http.ResponseWriter, req *http.Request) {
	var book service.Book

	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		http.Error(res, "invalid request payload", http.StatusBadRequest)
	}

	err = handler.service.CreateBook(&book)
	if err != nil {
		http.Error(res, "failed to create book", http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(book)
}

func (h *BookHandlers) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookById(id)
	if err != nil {
		http.Error(w, "failed to get book", http.StatusInternalServerError)
		return
	}
	if book == nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// UpdateBook lida com a requisição PUT /books/{id}.
func (h *BookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	var book service.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	book.ID = id

	if err := h.service.UpdateBook(&book); err != nil {
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// DeleteBook lida com a requisição DELETE /books/{id}.
func (h *BookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBook(id); err != nil {
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
