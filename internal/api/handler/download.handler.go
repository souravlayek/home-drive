package handler

import (
	"compress/flate"
	"context"
	"fmt"
	"image"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/souravlayek/storage-server/internal/database"
	"github.com/souravlayek/storage-server/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func parseInt(str string) int {
	if str == "" {
		return 0
	}
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return int(res)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["id"]
	width := parseInt(r.URL.Query().Get("w"))
	height := parseInt(r.URL.Query().Get("h"))
	size := parseInt(r.URL.Query().Get("s"))
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
	fileOptions := strings.Split(filePath, ".")
	fileExtension := fileOptions[len(fileOptions)-2]
	inputFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer inputFile.Close()
	gzipReader := flate.NewReader(inputFile)
	defer gzipReader.Close()
	err = os.Mkdir("temp", fs.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpfile, err := os.CreateTemp("temp", "decompressed-*."+fileExtension)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if (fileExtension == "jpg" || fileExtension == "png") && (size != 0 || width != 0) {
		file, _, err := image.Decode(gzipReader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if size != 0 {
			width = size
			height = 0
		}
		resizedImage := imaging.Resize(file, int(width), int(height), imaging.Lanczos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageFormat := imaging.JPEG
		if fileExtension == "png" {
			imageFormat = imaging.PNG
		}
		err = imaging.Encode(tmpfile, resizedImage, imageFormat)
		if err != nil {
			http.Error(w, "failed to encode image", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = io.Copy(tmpfile, gzipReader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tmpfile.Close()
	err = os.Rename(tmpfile.Name(), "decompressed."+fileExtension)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "decompressed."+fileExtension)

	err = os.Remove("decompressed." + fileExtension)
	if err != nil {
		fmt.Println("Failed to delete decompressed file:", err)
	}
}
