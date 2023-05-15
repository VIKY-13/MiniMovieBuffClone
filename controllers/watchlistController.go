package controllers

import(
	"encoding/json"
	"net/http"
	"golangmovietask/models"
)

// *controllers is refered from the movieController file where we have the struct and we use the same

func (wl *Controllers) AddMovieToUserWatchlist(w http.ResponseWriter, r *http.Request){
	var addToWatchList models.Favourite		//watchlist also uses same struct as favourite so we take that struct
	err := json.NewDecoder(r.Body).Decode(&addToWatchList)
	if err != nil{
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	wl.Service.AddMovieToUserWatchlistService(w,addToWatchList)
}

func (wl *Controllers) RemoveMovieFromUserWatchlist(w http.ResponseWriter, r *http.Request){
	var removeFromWatchlist models.Favourite	//we use favourite structure as it requires the same
	removeFromWatchlist.User_id = r.URL.Query().Get("user_id")
	removeFromWatchlist.Movie_id = r.URL.Query().Get("movie_id")
	wl.Service.RemoveMovieFromUserWatchlistService(w,removeFromWatchlist)
}

func (wl *Controllers) GetUserWatchlist(w http.ResponseWriter, r *http.Request){
	var userWatchlist models.Favourite
	userWatchlist.User_id = r.URL.Query().Get("user_id")	//we're getting only the user id and from userWatchlist we only use uer_id field
	wl.Service.GetUserWatchlistService(w,userWatchlist)
}
