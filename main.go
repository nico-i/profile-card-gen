package main

import (
	"html/template"
	"log"
	"net/http"
)

// main starts the web/file server and listens for requests.
func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc(
		"/", GenerateProfileCard,
	)
	// TODO Add routing for "/show" to enable preview of template
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// GenerateProfileCard generates a PDF of from the inputted
// multipart form data and returns it.
func GenerateProfileCard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	data, err := GenerateTemplateData(r)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	pdfBytes, err := GeneratePDF(&data, "./public/templates/aoe-profile-card.html")
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=profile.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(pdfBytes)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

// ShowProfileCardPage shows the generated profile card
// as an HTML page.
func ShowProfileCardPage(w http.ResponseWriter, r *http.Request) {
	data, err := GenerateTemplateData(r)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("./public/templates/aoe-profile-card.html")
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}
