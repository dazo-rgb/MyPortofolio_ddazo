package domain

// Project is the portfolio's read-only project entity.
type Project struct {
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Category         string   `json:"category"`
	Date             string   `json:"date"`
	Description      string   `json:"description"`
	Overview         string   `json:"overview"`
	Role             string   `json:"role"`
	Stack            []string `json:"stack"`
	Features         []string `json:"features"`
	TechnicalWork    []string `json:"technicalWork"`
	EngineeringNotes []string `json:"engineeringNotes"`
	Repository       string   `json:"repository"`
	Featured         bool     `json:"featured"`
}
