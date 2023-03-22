package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/souravlayek/storage-server/internal/api/router"
	"github.com/souravlayek/storage-server/internal/database"
	"github.com/souravlayek/storage-server/utils"
)

func main() {
	myArgs := os.Getenv("ENV_TYPE")
	fmt.Println(myArgs)
	if len(myArgs) < 1 || myArgs == "server" {
		utils.LoadENV()
	}
	hostname := os.Getenv("ENDPOINT")
	database.ConnectDB()
	r := router.CreateRouter()
	fmt.Println("Server has started successfully.")
	fmt.Println("Your API endpoint is: " + hostname)
	log.Fatal(http.ListenAndServe(":8000", r))
}
