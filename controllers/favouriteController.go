package controllers

import(
	"encoding/json"
	// "fmt"
	"net/http"
	// "golangmovietask/daos"
	"golangmovietask/models"

)

// *controllers is refered from the movieController file where we have the struct and we use the same

func (f *Controllers) AddUserFavourite(w http.ResponseWriter, r *http.Request){
	var addFavourite models.Favourite
	err := json.NewDecoder(r.Body).Decode(&addFavourite)
	if err != nil{
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	f.Service.AddUserFavouriteService(w,addFavourite)
}

func (f *Controllers) GetUserFavourites(w http.ResponseWriter, r *http.Request){
	var userFavourite models.Favourite
	userFavourite.User_id = r.URL.Query().Get("user_id")	//we're getting only the user id and from userFavourite we only use uer_id field
	f.Service.GetUserFavouritesService(w,userFavourite)
}

func (f *Controllers) RemoveUserFavourite(w http.ResponseWriter, r *http.Request){
	var removeFavourite models.Favourite
	removeFavourite.User_id = r.URL.Query().Get("user_id")
	removeFavourite.Movie_id = r.URL.Query().Get("movie_id")
	f.Service.RemoveUserFavouriteService(w,removeFavourite)
}