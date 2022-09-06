package main

import (
	"github.com/nico-i/profile-card-gen/types"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc(
		"/", GenerateProfileCard,
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GenerateProfileCard(w http.ResponseWriter, r *http.Request) {
	/*
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Fatal(err)
		}

		userData := types.User{
			Firstname:   r.PostFormValue("firstname"),
			Lastname:    r.PostFormValue("lastname"),
			Role:        r.PostFormValue("role"),
			City:        r.PostFormValue("city"),
			Team:        r.PostFormValue("team"),
			WorksAt:     r.PostFormValue("works_at"),
			Hometown:    r.PostFormValue("hometown"),
			Quote:       r.PostFormValue("quote"),
			Photo:       r.MultipartForm.File["photo"],
			HasWorkedAt: r.Form["has_worked_at"],
			Skills:      r.Form["skills"],
			Interests:   r.Form["interests"],
			Other:       r.Form["other"],
		}
	*/

	userData := types.User{
		Firstname:   "Max",
		Lastname:    "Mustermann",
		Role:        "Chef",
		City:        "Berlin, DE",
		Team:        "The best team",
		WorksAt:     "Hamburg HQ",
		Hometown:    "Buxdehude",
		Quote:       "Good quote.exe",
		Photo:       nil,
		HasWorkedAt: []string{"for suxess", "AStA of the RheinMain\nUniversity of Applied Science"},
		Skills:      []string{"Web development", "Automation", "Graphic design"},
		Interests:   []string{"Photography", "Guitar", "Machine Learning"},
		Other:       []string{"Always open to conversation", "Excited to learn"},
	}

	data := types.Data{
		User:     userData,
		BasePath: CurrentURL(r),
	}

	pdfBytes, err := GeneratePDF(&data, "./public/templates/aoe-profile-card.html")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=personalCard.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(pdfBytes)
	if err != nil {
		return
	}
}

func CurrentURL(r *http.Request) string {
	return "http://" + r.Host + r.URL.Path
}

func ShowProfileCardPage(w http.ResponseWriter, r *http.Request) {
	userData := types.User{
		Firstname:   "Nico",
		Lastname:    "Ismaili",
		Role:        "Intern",
		City:        "Wiesbaden",
		Team:        "Active Wholebuy",
		WorksAt:     "HQ Wiesbaden",
		Hometown:    "Simmern, Germany",
		Quote:       "“In the beginning there was Nothing, but Nothing is unstable, so Something came about.” ― Exurb1a, The Bridge to Lucy Dunne",
		Photo:       nil,
		HasWorkedAt: []string{"for suxess", "AStA of the RheinMain\nUniversity of Applied Science"},
		Skills:      []string{"Web development", "Automation", "Graphic design"},
		Interests:   []string{"Photography", "Guitar", "Machine Learning"},
		Other:       []string{"Always open to conversation", "Excited to learn"},
	}

	data := types.Data{
		User:     userData,
		BasePath: CurrentURL(r),
	}

	tmpl, err := template.ParseFiles("./public/templates/aoe-profile-card.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		return
	}
}
