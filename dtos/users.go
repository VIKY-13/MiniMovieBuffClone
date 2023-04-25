package dtos

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

//user login struct
type userlogin struct{
	Useremail string `json:"useremail"`
	Password string `json:"password"`
}

//favourite,watchlist structure
type favourite struct{
	User_id string `json:"user_id"`
	Movie_id string `json:"movie_id"`
}