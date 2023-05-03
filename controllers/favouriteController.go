package controllers

import(
	"encoding/json"
	// "fmt"
	"net/http"
	// "golangmovietask/daos"
	"golangmovietask/models"

)

func (f *Controllers) AddUserFavourite(w http.ResponseWriter, r *http.Request){
	var addFavourite models.Favourite
	json.NewDecoder(r.Body).Decode(&addFavourite)
	f.Service.AddUserFavouriteService(w,addFavourite)
}

func (f *Controllers) GetUserFavourites(w http.ResponseWriter, r *http.Request){
	var userFavourite models.Favourite
	json.NewDecoder(r.Body).Decode(&userFavourite)
	f.Service.GetUserFavouritesService(w,userFavourite)
}

func (f *Controllers) RemoveUserFavourite(w http.ResponseWriter, r *http.Request){
	var removeFavourite models.Favourite
	json.NewDecoder(r.Body).Decode(&removeFavourite)
	f.Service.RemoveUserFavouriteService(w,removeFavourite)
}