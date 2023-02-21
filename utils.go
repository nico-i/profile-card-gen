package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
)

// GenerateTemplateData extracts the relevant data from a multipart form and returns types.TemplateData
func GenerateTemplateData(r *http.Request) (TemplateData, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return TemplateData{}, err
	}

	data := TemplateData{
		User: User{
			Firstname: r.PostFormValue("firstname"),
			Lastname:  r.PostFormValue("lastname"),
			Role:      r.PostFormValue("role"),
			City:      r.PostFormValue("city"),
			Team:      r.PostFormValue("team"),
			WorksAt:   r.PostFormValue("works_at"),
			Hometown:  r.PostFormValue("hometown"),
			Quote:     r.PostFormValue("quote"),
			Adjective: r.PostFormValue("adjective"),
		},
		BasePath: "http://" + r.Host + r.URL.Path,
	}
	data.User.HasWorkedAt = deleteAndEscapeStringArr(r.Form["has_worked_at"])
	data.User.Skills = deleteAndEscapeStringArr(r.Form["skills"])
	data.User.Happy = deleteAndEscapeStringArr(r.Form["happy"])
	data.User.Burdens = deleteAndEscapeStringArr(r.Form["burden"])
	data.User.Heros = deleteAndEscapeStringArr(r.Form["heros"])

	if len(data.User.HasWorkedAt) == 0 {
		data.User.HasWorkedAt = nil
	}
	_, base64Str, err := HandleImageUpload(r, "photo")
	if err != nil {
		data.User.Errors["Photo"] = ""
		return TemplateData{}, err
	}
	data.User.Photo = base64Str
	return data, nil
}

// HandleImageUpload encodes an image uploaded through a
// multipart form into Base64 and returns the resulting string.
func HandleImageUpload(r *http.Request, formFileName string) (string, string, error) {
	// Maximum upload of 32 MB files
	file, handler, err := r.FormFile(formFileName)
	if err != nil {
		return "", "", err
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(file)
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

func deleteAndEscapeStringArr(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, html.EscapeString(str))
		}
	}
	return r
}

func escapeNonArrayTemplateDate(data *TemplateData) {
	value := reflect.ValueOf(data).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() == reflect.TypeOf("") {
			str := field.Interface().(string)
			field.SetString(html.EscapeString(str))
		}
	}
}
