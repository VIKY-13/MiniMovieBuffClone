package main

import (
	// "database/sql"
	"encoding/json"
	// "fmt"
	// "os"
	// "html/template"
	// "io/ioutil"
	"net/http"
	// "strings"
	// "github.com/gorilla/mux"
	
	_ "github.com/lib/pq"
)

//watchlist
func AddMovieToUserWatchlist(w http.ResponseWriter, r *http.Request){
	var addToWatchList favourite
	json.NewDecoder(r.Body).Decode(&addToWatchList)
	_,err := Db.Exec("insert into watchlist(user_id,movie_id) values($1,$2)",addToWatchList.User_id,addToWatchList.Movie_id)
	CheckInternalServerError(w,err,"error while adding movie to watchlist")
	w.WriteHeader(http.StatusCreated)
}

func RemoveMovieFromUserWatchlist(w http.ResponseWriter, r *http.Request){
	var removeFromWatchlist favourite	//we use favourite structure as it requires the same
	json.NewDecoder(r.Body).Decode(&removeFromWatchlist)
	_,err := Db.Exec("DELETE from watchlist where movie_id=$1 and user_id=$2",removeFromWatchlist.Movie_id,removeFromWatchlist.User_id)
	CheckInternalServerError(w,err,"error while removing movie from watchlist")
	w.WriteHeader(http.StatusNoContent)
}

func GetUserWatchlist(w http.ResponseWriter, r *http.Request){
	var userWatchlist favourite
	json.NewDecoder(r.Body).Decode(&userWatchlist)
	movies := GetMovieIdListOnUserWatchlist(userWatchlist.User_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func GetMovieIdListOnUserWatchlist(user_id string) []retrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT distinct movie_id FROM watchlist WHERE user_id = $1;")
	rows,_:= statement.Query(user_id)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}