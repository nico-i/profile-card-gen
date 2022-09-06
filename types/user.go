package types

import (
	"mime/multipart"
)

type User struct {
	Firstname   string                  `json:"firstname"`
	Lastname    string                  `json:"lastname"`
	Role        string                  `json:"role"`
	Team        string                  `json:"team"`
	WorksAt     string                  `json:"works_at"`
	City        string                  `json:"city"`
	Hometown    string                  `json:"hometown"`
	Quote       string                  `json:"quote"`
	Photo       []*multipart.FileHeader `json:"photo"`
	HasWorkedAt []string                `json:"has_worked_at"`
	Skills      []string                `json:"skills"`
	Interests   []string                `json:"interests"`
	Other       []string                `json:"other"`
}
