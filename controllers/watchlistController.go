package controllers

import(
	"encoding/json"
	"net/http"
	"golangmovietask/models"
)

func (wl *Controllers) AddMovieToUserWatchlist(w http.ResponseWriter, r *http.Request){
	var addToWatchList models.Favourite		//watchlist also uses same struct as favourite so we take that struct
	json.NewDecoder(r.Body).Decode(&addToWatchList)
	wl.Service.AddMovieToUserWatchlistService(w,addToWatchList)
}

func (wl *Controllers) RemoveMovieFromUserWatchlist(w http.ResponseWriter, r *http.Request){
	var removeFromWatchlist models.Favourite	//we use favourite structure as it requires the same
	json.NewDecoder(r.Body).Decode(&removeFromWatchlist)
	wl.Service.RemoveMovieFromUserWatchlistService(w,removeFromWatchlist)
}

func (wl *Controllers) GetUserWatchlist(w http.ResponseWriter, r *http.Request){
	var userWatchlist models.Favourite
	json.NewDecoder(r.Body).Decode(&userWatchlist)
	wl.Service.GetUserWatchlistService(w,userWatchlist)
}
