package controllers

import(
	"encoding/json"
	"fmt"
	"net/http"
	"golangmovietask/daos"
	"golangmovietask/dtos"

)

func AddMovieToUserWatchlist(w http.ResponseWriter, r *http.Request){
	var addToWatchList dtos.Favourite		//watchlist also uses same struct as favourite so we take that struct
	json.NewDecoder(r.Body).Decode(&addToWatchList)
	err := daos.AddMovieToUserWatchlistDb(addToWatchList)
	if err!= nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error while adding movie to watchlist")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func RemoveMovieFromUserWatchlist(w http.ResponseWriter, r *http.Request){
	var removeFromWatchlist dtos.Favourite	//we use favourite structure as it requires the same
	json.NewDecoder(r.Body).Decode(&removeFromWatchlist)
	err := daos.RemoveMovieFromWatchlistDb(removeFromWatchlist)
	if err!= nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error while removing movie from watchlist")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetUserWatchlist(w http.ResponseWriter, r *http.Request){
	var userWatchlist dtos.Favourite
	json.NewDecoder(r.Body).Decode(&userWatchlist)
	movies,err := daos.GetMovieIdListOnUserWatchlist(userWatchlist.User_id)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error while Getting movie from watchlist")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}
