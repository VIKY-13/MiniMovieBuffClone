package dtos

type MovData struct {
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  RunningTime `json:"running_time"`
	Summary       string       `json:"summary"`
	Certification string	   `json:"certification"`
	Genres		  []string	   `json:"genres"`
	Photos        []string	   `json:"photos"`		
	Trailers      []string     `json:"trailers"`
	Cast          []Cast       `json:"cast"`
	Crew		  []Crew	   `json:"crew"`
}

type RunningTime struct {
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
}

type Crew struct{
	Name     string `json:"name"`
	Role     string `json:"role"`
	Crew_member_id string `json:"uuid"`
	Poster   string `json:"poster"`
}
type Cast struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Cast_member_id string `json:"uuid"`
	Poster   string `json:"poster"`
}