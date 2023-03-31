package main


// type retrieveId struct {
// 	Uuid string
// }

type movdata struct {
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  running_time `json:"running_time"`
	Summary       string       `json:"summary"`
	Cast          []cast       `json:"cast"`
}

type running_time struct {
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
}

type cast struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Actor_id string `json:"uuid"`
	Poster   string `json:"poster"`
}

type retrieveMovie struct{
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  string	   `json:"running_time"`
}
//these structs are for APIDocumentations
type EndpointsHead struct{
	Movie []EndpointDescriptions `json:"movie"`
	User []EndpointDescriptions `json:"user"`
	Favourite []EndpointDescriptions `json:"favourite"`
	Watchlist []EndpointDescriptions `json:"watchlist"`
}

type EndpointDescriptions struct{
	Method string `json:"method"`
	Endpoints string	`json:"endpoints"`
	Description string	`json:"description"`
	Parameters []string	`json:"parameters"`
	Samplereqres SampleReqRes `json:"samplereqres"`
}

type SampleReqRes struct{
	ReqURL string `json:"requrl"`
	ReqBody string `json:"reqbody"`
	Response string `json:"response"`
}

// type parameters struct{
// 	Params string
// }

type documentationparsedata struct{
	Title string
	Endpointsdata EndpointsHead
}

//create user structures
type user struct{
	User_id string `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Email string `json:"email"`
	Password string `json:"password"`
	Age string `json:"age"`
	Phone_no string `json:"phone_no" maxLength:"10"`
}

//favourite,watchlist structure
type favourite struct{
	User_id string `json:"user_id"`
	Movie_id string `json:"movie_id"`
}

//user login struct
type userlogin struct{
	Useremail string `json:"useremail"`
	Password string `json:"password"`
}

//rating struct
type movierating struct{
	User_id string `json:"user_id"`
	Movie_id string `json:"movie_id"`
	Rating float32 `json:"rating"`
	Review string `json:"review"`
}