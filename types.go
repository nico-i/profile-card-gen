package main

// User represents the user wanting to generate
// a profile card.
type User struct {
	Firstname   string   `json:"firstname"`
	Lastname    string   `json:"lastname"`
	Role        string   `json:"role"`
	Team        string   `json:"team"`
	WorksAt     string   `json:"works_at"`
	City        string   `json:"city"`
	Hometown    string   `json:"hometown"`
	Quote       string   `json:"quote"`
	Photo       string   `json:"photo"`
	HasWorkedAt []string `json:"has_worked_at"`
	Skills      []string `json:"skills"`
	Interests   []string `json:"interests"`
	Other       []string `json:"other"`
}

// TemplateData contains all the necessary data
// to fill a profile card HTML template
type TemplateData struct {
	User     User
	BasePath string
}
