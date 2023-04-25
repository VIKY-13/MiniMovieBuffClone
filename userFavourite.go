package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
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
	var removeFavourite favourite
	json.NewDecoder(r.Body).Decode(&removeFavourite)
	_,err := Db.Exec("DELETE from favourite where movie_id= $1 and user_id= $2",removeFavourite.Movie_id,removeFavourite.User_id)
	if err!= nil{
		fmt.Fprintln(w,"error while removing favourites")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AddUserFavourite(w http.ResponseWriter, r *http.Request){
	var addFavourite favourite
	json.NewDecoder(r.Body).Decode(&addFavourite)
	_,err := Db.Exec("insert into favourite(user_id,movie_id) values($1,$2)",addFavourite.User_id,addFavourite.Movie_id)
	if err!= nil{
		fmt.Fprintln(w,"error while adding favourites")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetUserFavourites(w http.ResponseWriter, r *http.Request){
	var userFavourite favourite
	json.NewDecoder(r.Body).Decode(&userFavourite)
	movies:=GetMovieIdListOnUserFavourite(userFavourite.User_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func GetMovieIdListOnUserFavourite(user_id string) []retrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT distinct movie_id FROM favourite WHERE user_id = $1;")
	rows,_ := statement.Query(user_id)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}