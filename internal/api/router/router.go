package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/souravlayek/storage-server/internal/api/handler"
)

func CreateRouter() *mux.Router {
	fs := http.FileServer(http.Dir("media"))
	fs2 := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.Handle("/temp/", http.StripPrefix("/temp/", fs2))
	router := mux.NewRouter()
	router.HandleFunc("/api/upload", handler.UploadFileHandler).Methods("POST")
	router.HandleFunc("/s/{id}", handler.DownloadFile).Methods("GET")
	return router
}
