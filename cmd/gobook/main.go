package main

import (
	"database/sql"
	"gobooks/internal/service"
	"gobooks/internal/web"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	bookService := service.NewBookService(db)
	bookHandlers := web.NewBookHandlers(bookService)

	router := http.NewServeMux()
	router.HandleFunc("GET /book", bookHandlers.GetBooks)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("GET /book/{id}", bookHandlers.GetBookByID)
	router.HandleFunc("PUT /book/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /book/{id}", bookHandlers.DeleteBook)

	http.ListenAndServe(":8080", router)
}
