package services

import(
	"encoding/json"
	"fmt"
	"net/http"
	// "golangmovietask/daos"
	"golangmovietask/models"

)

// *Services is refered from the movieService file where we have the struct and we use the same
func (f *Service) AddUserFavouriteService(w http.ResponseWriter,addFavourite models.Favourite){
	err := f.DAO.AddMovieToUserFavouriteDb(addFavourite)
	if err!= nil{
		fmt.Fprintln(w,"error while adding favourites")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (f *Service) GetUserFavouritesService(w http.ResponseWriter,userFavourite models.Favourite){
	if userFavourite.User_id == ""{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"query is empty and it shouldn't")
		return
	}
	movies := f.DAO.GetMovieIdListOnUserFavourite(userFavourite.User_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (f *Service) RemoveUserFavouriteService(w http.ResponseWriter,removeFavourite models.Favourite){
	err := f.DAO.RemoveMovieFromFavouriteDb(removeFavourite)
	if err!= nil{
		fmt.Fprintln(w,"error while removing favourites")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}