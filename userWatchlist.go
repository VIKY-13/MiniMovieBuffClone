package main

import (
	// "database/sql"
	"encoding/json"
	// "fmt"
	// "os"
	// "html/template"
	"io/ioutil"
	"net/http"
	// "strings"
	// "github.com/gorilla/mux"
	
	_ "github.com/lib/pq"
)

//watchlist
func AddMovieToUserWatchlist(w http.ResponseWriter, r *http.Request){
	resp,err := ioutil.ReadAll(r.Body)
	checkErr(err)
	var watchlist favourite
	json.Unmarshal(resp,&watchlist)
	_,err = Db.Exec("insert into watchlist(user_id,movie_id) values($1,$2)",watchlist.User_id,watchlist.Movie_id)
	checkErr(err)
	json.NewEncoder(w).Encode(watchlist)
}

func RemoveMovieFromUserWatchlist(w http.ResponseWriter, r *http.Request){
	resp,err := ioutil.ReadAll(r.Body)
	checkErr(err)
	var watchlist favourite
	json.Unmarshal(resp,&watchlist)
	_,err = Db.Exec("DELETE from watchlist where movie_id=$1 and user_id=$2",watchlist.Movie_id,watchlist.User_id)
	checkErr(err)
	json.NewEncoder(w).Encode("deleted")
}

func GetUserWatchlist(w http.ResponseWriter, r *http.Request){
	var movies_id []string
	var movie_id string
	var user favourite
	// Useremail := r.URL.Query().Get("useremail")
	resp,err := ioutil.ReadAll(r.Body)
	checkErr(err)
	json.Unmarshal(resp,&user)
	rows,err:=Db.Query("select distinct(movie_id) from watchlist where user_id='"+user.User_id+"';")
	checkErr(err)
	for rows.Next(){
		err = rows.Scan(&movie_id)
		checkErr(err)
		movies_id = append(movies_id, movie_id)
	}
	// movies:=RetriveDataOnMovie_id(movies_id)
	json.NewEncoder(w).Encode(movies_id)
}