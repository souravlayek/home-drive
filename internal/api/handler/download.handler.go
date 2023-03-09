package handler

import (
	"compress/flate"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/souravlayek/storage-server/internal/database"
	"github.com/souravlayek/storage-server/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var myMetadata model.MetaData
	err = database.MetaDataCollection.FindOne(context.TODO(), bson.M{
		"_id": objectId,
	}).Decode(&myMetadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filePath := myMetadata.Path
	fmt.Println(filePath)
	fileOptions := strings.Split(filePath, ".")
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()
	gzipReader := flate.NewReader(inputFile)
	defer gzipReader.Close()

	tmpfile, err := os.CreateTemp("media", "decompressed-*."+fileOptions[len(fileOptions)-2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tmpfile, gzipReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpfile.Close()
	err = os.Rename(tmpfile.Name(), "decompressed."+fileOptions[len(fileOptions)-2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "decompressed."+fileOptions[len(fileOptions)-2])

	err = os.Remove("decompressed." + fileOptions[len(fileOptions)-2])
	if err != nil {
		fmt.Println("Failed to delete decompressed file:", err)
	}
}
