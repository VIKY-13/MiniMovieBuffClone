package models

//we get only thes
type MovData struct {
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Release_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  RunningTime  `json:"running_time"`
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
	Name     string       `json:"name"`
	Role     string       `json:"role"`
	Crew_member_id string `json:"uuid"`
	Poster   string       `json:"poster"`
}
type Cast struct {
	Name     string       `json:"name"`
	Role     string       `json:"role"`
	Cast_member_id string `json:"uuid"`
	Poster   string       `json:"poster"`
}

//rating struct used to upload data in db
type MovieRating struct{
	User_id string  `json:"user_id"`
	Movie_id string `json:"movie_id"`
	Rating float32  `json:"rating"`
	Review string   `json:"review"`
}

// this struct is used in retrieving the moviedata with reviews
type UserReviews struct{
	Username string  `json:"username"`
	Rating	 float64 `json:"rating"`
	Review	 string	 `json:"review"`
}

//retrieve movie data struct
type RetrieveMovData struct{
	Movie_id      string        `json:"uuid"`
	Title         string        `json:"title"`
	Realease_date string        `json:"release_date"`
	Languge_name  string        `json:"language_name"`
	Running_time  RunningTime   `json:"running_time"`
	Summary       string        `json:"summary"`
	Certification string	    `json:"certification"`
	Genres		  []string	    `json:"genres"`
	Photos        []string	    `json:"photos"`		
	Trailers      []string      `json:"trailers"`
	Cast          []Cast        `json:"cast"`
	Crew		  []Crew	    `json:"crew"`
	OverallUserRating float32   `json:"overalluserrating"`
	UserReviews   []UserReviews `json:"userreviews"`
}

