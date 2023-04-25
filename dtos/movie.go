package dtos

//we get only thes
type movdata struct {
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  running_time `json:"running_time"`
	Summary       string       `json:"summary"`
	Certification string	   `json:"certification"`
	Genres		  []string	   `json:"genres"`
	Photos        []string	   `json:"photos"`		
	Trailers      []string     `json:"trailers"`
	Cast          []cast       `json:"cast"`
	Crew		  []crew	   `json:"crew"`
}

type running_time struct {
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
}

type crew struct{
	Name     string `json:"name"`
	Role     string `json:"role"`
	Crew_member_id string `json:"uuid"`
	Poster   string `json:"poster"`
}
type cast struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Cast_member_id string `json:"uuid"`
	Poster   string `json:"poster"`
}

//rating struct used to upload data in db
type movierating struct{
	User_id string `json:"user_id"`
	Movie_id string `json:"movie_id"`
	Rating float32 `json:"rating"`
	Review string `json:"review"`
}

// this struct is used in retrieving the moviedata with reviews
type userreviews struct{
	Username string `json:"username"`
	Rating	 float64 `json:"rating"`
	Review	 string	`json:"review"`
}

//retrieve movie data struct
type retrieveMovData struct{
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  running_time `json:"running_time"`
	Summary       string       `json:"summary"`
	Certification string	   `json:"certification"`
	Genres		  []string	   `json:"genres"`
	Photos        []string	   `json:"photos"`		
	Trailers      []string     `json:"trailers"`
	Cast          []cast       `json:"cast"`
	Crew		  []crew	   `json:"crew"`
	OverallUserRating float32  `json:"overalluserrating"`
	UserReviews   []userreviews `json:"userreviews"`
}

