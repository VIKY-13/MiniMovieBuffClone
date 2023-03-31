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

//favourites
func RemoveUserFavourite(w http.ResponseWriter, r *http.Request){
	var fav favourite
	json.NewDecoder(r.Body).Decode(&fav)
	_,err := Db.Exec("DELETE from favourite where movie_id=$1 and user_id=$2",fav.Movie_id,fav.User_id)
	checkErr(err)
	w.WriteHeader(http.StatusNoContent)
}

func AddUserFavourite(w http.ResponseWriter, r *http.Request){
	var fav favourite
	json.NewDecoder(r.Body).Decode(&fav)
	_,err := Db.Exec("insert into favourite(user_id,movie_id) values($1,$2)",fav.User_id,fav.Movie_id)
	checkErr(err)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fav)
}

func GetUserFavourites(w http.ResponseWriter, r *http.Request){
	var user favourite
	json.NewDecoder(r.Body).Decode(&user)
	movies:=GetMovieIdList("favourite","user_id",user.User_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}