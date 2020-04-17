package handler

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 100 << 20 // 100 MB
const uploadPath = "./uploads"
const fileNameSize = 16

func (h *handler) getChats(w http.ResponseWriter, r *http.Request) {
	chats, err := h.GetChats()
	if err != nil {

		respondHTTPError(w, err, http.StatusBadRequest)
	}

	respond(w, chats, http.StatusOK)
}

func (h *handler) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		//renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		respondHTTPError(w, fmt.Errorf("FILE_TOO_BIG: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	file, _, err := r.FormFile("file")
	if err != nil {
		respondHTTPError(w, fmt.Errorf("INVALID_FILE:A  %v", err), http.StatusBadRequest)

		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		respondHTTPError(w, fmt.Errorf("INVALID_FILE:B %v", err), http.StatusBadRequest)
		return
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	detectedFileType := http.DetectContentType(fileBytes)

	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png", "video/mp4":
		break
	default:
		respondHTTPError(w, fmt.Errorf("INVALID_FILE_TYPE: %v", err), http.StatusBadRequest)
		return
	}

	fileName := uniqueToken(fileNameSize)
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		respondHTTPError(w, fmt.Errorf("CANT_READ_FILE_TYPE: %v", err), http.StatusBadRequest)
		return
	}
	createUploadFolder(uploadPath)
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		respondHTTPError(w, fmt.Errorf("CANT_WRITE_FILE: %v", err), http.StatusBadRequest)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {

		respondHTTPError(w, fmt.Errorf("CANT_WRITE_FILE: %v", err), http.StatusBadRequest)
		return
	}

	uploadSuccess := uploadSuccess{
		true,
		"files/" + fileName + fileEndings[0],
	}

	respond(w, uploadSuccess, http.StatusOK)
}

type uploadSuccess struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}
