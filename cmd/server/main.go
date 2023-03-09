package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/souravlayek/storage-server/internal/api/router"
	"github.com/souravlayek/storage-server/internal/database"
)

func main() {
	database.ConnectDB()
	r := router.CreateRouter()
	fmt.Println("Server is listening on :8000 ....")
	log.Fatal(http.ListenAndServe(":8000", r))
}
