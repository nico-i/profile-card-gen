package main

import (
	"testing"
)

// TestGeneratePDF is a functional / component test to test if the system works cohesively
func TestGeneratePDF(t *testing.T) {
	data := TemplateData{
		User: User{
			Firstname:   "Max",
			Lastname:    "Mustermann",
			Role:        "Chef",
			City:        "Berlin, DE",
			Team:        "The best team",
			WorksAt:     "Hamburg HQ",
			Hometown:    "Buxdehude",
			Quote:       "Good quote.exe",
			Photo:       "https://picsum.photos/500/500",
			HasWorkedAt: []string{"for suxess", "AStA of the RheinMain\nUniversity of Applied Science"},
			Skills:      []string{"Web development", "Automation", "Graphic design"},
			Interests:   []string{"Photography", "Guitar", "Machine Learning"},
			Other:       []string{"Always open to conversation", "Excited to learn"},
		},
		BasePath: "https://www.example.com",
	}
	_, err := GeneratePDF(&data, "./public/templates/aoe-profile-card.html")
	if err != nil {
		t.Fatal(err)
	}
}
