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
	var watchlist favourite
	json.NewDecoder(r.Body).Decode(&watchlist)
	_,err := Db.Exec("insert into watchlist(user_id,movie_id) values($1,$2)",watchlist.User_id,watchlist.Movie_id)
	checkErr(err)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(watchlist)
}

func RemoveMovieFromUserWatchlist(w http.ResponseWriter, r *http.Request){
	var watchlist favourite
	json.NewDecoder(r.Body).Decode(&watchlist)
	_,err := Db.Exec("DELETE from watchlist where movie_id=$1 and user_id=$2",watchlist.Movie_id,watchlist.User_id)
	checkErr(err)
	w.WriteHeader(http.StatusNoContent)
}

func GetUserWatchlist(w http.ResponseWriter, r *http.Request){
	var user favourite
	json.NewDecoder(r.Body).Decode(&user)
	movies:=GetMovieIdList("watchlist","user_id",user.User_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}