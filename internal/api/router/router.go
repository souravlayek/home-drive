package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/souravlayek/storage-server/internal/api/handler"
)

func CreateRouter() *mux.Router {
	fs := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "media/test.jpg")
	})
	router.HandleFunc("/api/upload", handler.UploadFileHandler).Methods("POST")
	router.HandleFunc("/s/{id}", handler.DownloadFile).Methods("GET")
	return router
}
