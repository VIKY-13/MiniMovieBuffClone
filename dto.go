package dto



type movdata struct{
	Uuid string `json:"uuid"`
	Title string `json:"title"`
	Realease_date string `json:"release_date"`
	Languge_name string `json:"language_name"`
	Running_time running_time `json:"running_time`
}

type running_time struct{
	Hours string `json:"hours"`	
	Minutes string `json:"minutes"`
}
