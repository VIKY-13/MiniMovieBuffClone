package services

import(
	"encoding/json"
	"fmt"
	"net/http"
	"golangmovietask/models"

)

func (wl *Service) AddMovieToUserWatchlistService(w http.ResponseWriter,addToWatchList models.Favourite){
	err := wl.DAO.AddMovieToUserWatchlistDb(addToWatchList)
	if err!= nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error while adding movie to watchlist")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (wl *Service) GetUserWatchlistService(w http.ResponseWriter, userWatchlist models.Favourite){
	movies,err := wl.DAO.GetMovieIdListOnUserWatchlist(userWatchlist.User_id)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error while Getting movie from watchlist")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (wl *Service) RemoveMovieFromUserWatchlistService(w http.ResponseWriter,removeFromWatchlist models.Favourite){
	err := wl.DAO.RemoveMovieFromWatchlistDb(removeFromWatchlist)
	if err!= nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error while removing movie from watchlist")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}