package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

// User represents the user wanting to generate
// a profile card.
type User struct {
	Firstname   string            `json:"firstname"`
	Lastname    string            `json:"lastname"`
	Role        string            `json:"role"`
	Team        string            `json:"team"`
	WorksAt     string            `json:"works_at"`
	City        string            `json:"city"`
	Hometown    string            `json:"hometown"`
	Quote       string            `json:"quote"`
	Photo       string            `json:"photo"`
	HasWorkedAt []string          `json:"has_worked_at"`
	Skills      []string          `json:"skills"`
	Interests   []string          `json:"interests"`
	Other       []string          `json:"other"`
	Errors      map[string]string `json:"errors"`
}

// main starts the web/file server and listens for requests.
func main() {
	port := "3000"
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r := mux.NewRouter()
	r.HandleFunc("/", ShowForm)
	r.HandleFunc("/gen", GenerateProfileCard).Methods("POST")
	r.HandleFunc("/preview", ShowProfileCardPage).Methods("POST")
	http.Handle("/", r)
	log.Println("Server is online under http://localhost:" + port + ".")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func ShowForm(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./public/templates/form.html")
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

// GenerateProfileCard generates a PDF of from the inputted
// multipart form data and returns it.
func GenerateProfileCard(w http.ResponseWriter, r *http.Request) {
	data, err := GenerateTemplateData(r)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	escapeNonArrayTemplateDate(&data)
	data.Preview = false
	pdfBytes, err := GeneratePDF(&data, "./public/templates/aoe-profile-card.html")
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s %s.pdf", data.User.Firstname, data.User.Lastname))
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
