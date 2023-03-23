package handler

import (
	"compress/flate"
	"context"
	"encoding/json"
	"image"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/buckket/go-blurhash"

	"github.com/souravlayek/storage-server/internal/database"
	"github.com/souravlayek/storage-server/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadResponse struct {
	Url      string `json:"url"`
	Blurhash string `json:"blurhash"`
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileOptions := strings.Split(handler.Filename, ".")
	filePath := "media/" + fileOptions[0] + time.Now().Format("_2006_01_02_15_04_05") + "." + fileOptions[1] + ".gz"
	err = os.MkdirAll("media", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	// Check whether uploaded file is image from filename
	isFileImage := strings.Contains(handler.Filename, ".jpg") || strings.Contains(handler.Filename, ".jpeg") || strings.Contains(handler.Filename, ".png")
	var blurHashStr = ""
	if isFileImage {
		// Decode the image
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}

		// Generate blurhash
		blurHashStr, err = blurhash.Encode(4, 3, img)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	compressor, err := flate.NewWriter(f, flate.BestCompression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer compressor.Close()
	_, err = io.Copy(compressor, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// File written successfully
	myMetaData := model.MetaData{
		Id:       primitive.NewObjectID(),
		Name:     handler.Filename,
		Path:     filePath,
		BlurHash: blurHashStr,
	}
	res, err := database.MetaDataCollection.InsertOne(context.TODO(), myMetaData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fileIdHex := res.InsertedID.(primitive.ObjectID).Hex()
	hostname := os.Getenv("ENDPOINT")
	url := hostname + "/s/" + fileIdHex
	myResp := UploadResponse{
		Url:      url,
		Blurhash: blurHashStr,
	}
	json.NewEncoder(w).Encode(myResp)
}
