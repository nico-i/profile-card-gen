package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nico-i/profile-card-gen/types"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

// GenerateTemplateData extracts the relevant data from a multipart form and returns types.TemplateData
func GenerateTemplateData(r *http.Request) (types.TemplateData, error) {
	_, base64Str, err := HandleImageUpload(r, "photo")
	if err != nil {
		return types.TemplateData{}, err
	}
	return types.TemplateData{
		User: types.User{
			Firstname:   r.PostFormValue("firstname"),
			Lastname:    r.PostFormValue("lastname"),
			Role:        r.PostFormValue("role"),
			City:        r.PostFormValue("city"),
			Team:        r.PostFormValue("team"),
			WorksAt:     r.PostFormValue("works_at"),
			Hometown:    r.PostFormValue("hometown"),
			Quote:       r.PostFormValue("quote"),
			Photo:       base64Str,
			HasWorkedAt: r.Form["has_worked_at"],
			Skills:      r.Form["skills"],
			Interests:   r.Form["interests"],
			Other:       r.Form["other"],
		},
		BasePath: "http://" + r.Host + r.URL.Path,
	}, nil
}

// HandleImageUpload encodes an image uploaded through a
// multipart form into Base64 and returns the resulting string.
func HandleImageUpload(r *http.Request, formFileName string) (string, string, error) {
	// Maximum upload of 32 MB files
	file, handler, err := r.FormFile(formFileName)

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(file)

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	contentType := http.DetectContentType(data)

	switch contentType {
	case "image/png":
		fmt.Println("Image type is already PNG.")
	case "image/jpeg":
		fmt.Println("Image type is JPG, converting it to PNG.")
		img, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return "", "", fmt.Errorf("unable to decode jpeg: %w", err)
		}
		var buf bytes.Buffer

		if err := png.Encode(&buf, img); err != nil {
			return "", "", fmt.Errorf("unable to encode png")
		}
		data = buf.Bytes()
	default:
		return "", "", fmt.Errorf("unsupported content type: %s", contentType)
	}
	//convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(data)
	return handler.Filename, imgBase64Str, nil
}

// handleError returns a JSON with a short error message
// and the corresponding code if something goes wrong.
func handleError(w http.ResponseWriter, err error, errorCode int) {
	w.WriteHeader(errorCode)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["error"] = string(rune(errorCode))
	resp["message"] = err.Error()
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error happened while writing JSON. Err: %s", err)
	}
	return
}
