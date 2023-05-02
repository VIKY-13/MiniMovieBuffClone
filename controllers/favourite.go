package controllers

import(
	"encoding/json"
	"fmt"
	"net/http"
	"golangmovietask/daos"
	"golangmovietask/models"

)

func AddUserFavourite(w http.ResponseWriter, r *http.Request){
	var addFavourite models.Favourite
	json.NewDecoder(r.Body).Decode(&addFavourite)
	err := daos.AddMovieToUserFavouriteDb(addFavourite)
	if err!= nil{
		fmt.Fprintln(w,"error while adding favourites")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetUserFavourites(w http.ResponseWriter, r *http.Request){
	var userFavourite models.Favourite
	json.NewDecoder(r.Body).Decode(&userFavourite)
	movies := daos.GetMovieIdListOnUserFavourite(userFavourite.User_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func RemoveUserFavourite(w http.ResponseWriter, r *http.Request){
	var removeFavourite models.Favourite
	json.NewDecoder(r.Body).Decode(&removeFavourite)
	err := daos.RemoveMovieFromFavouriteDb(removeFavourite)
	if err!= nil{
		fmt.Fprintln(w,"error while removing favourites")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}