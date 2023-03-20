package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"

	// "os"
	// "html/template"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	// "strings"
)
func GetMovieDataByQueryParams(w http.ResponseWriter, r *http.Request){
	cast := r.URL.Query().Get("cast")
	language_name := r.URL.Query().Get("language_name")
	release_date := r.URL.Query().Get("release_date")
	certification := r.URL.Query().Get("certification")
	var movies []retrieveMovie
	if(language_name!=""){
		movies=GetMovieIdList("movie","language_name",language_name)
	}else if (len(release_date)!=0){
		movies=GetMovieIdList("movie","release_date",release_date)
	}else if (len(certification)!=0){
		movies=GetMovieIdList("movie","certification",certification)
	}else if(cast!= ""){
		movies =GetMovieIdList("casts","name",cast)
	}else{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Check the query")
	}
	// fmt.Print(movies)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func GetMovieIdList(table string,column string,parameter string) []retrieveMovie{
	var movies_id []string
	var movie_id string
	query := fmt.Sprintf("select m.movie_id from "+table+" m inner join running_time r on "+column+"='"+parameter+"'and m.movie_id=r.movie_id;")
	rows, err := Db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&movie_id)
		checkErr(err)
		movies_id = append(movies_id,movie_id )
	}
	movies := RetriveDataOnMovie_id(movies_id)
	return movies
}

func RetriveDataOnMovie_id(movies_id []string)[]retrieveMovie{
	var movies []retrieveMovie
	var movie retrieveMovie
	for i:=0;i<len(movies_id);i++{
		query := fmt.Sprintf("select m.movie_id,m.title,m.release_date,m.language_name, r.hours || ';' ||r.minutes as running_time from movie m inner join running_time r on m.movie_id='"+movies_id[i]+"';")
		rows,err := Db.Query(query)
		checkErr(err)
		defer rows.Close()
		for rows.Next(){
			err = rows.Scan(&movie.Movie_id, &movie.Title, &movie.Realease_date, &movie.Languge_name,&movie.Running_time)
			checkErr(err)
		}
		movies = append(movies, movie)
	}
	return movies
}

func PostNewMovieData(w http.ResponseWriter, r *http.Request) {
	var movie movdata
	vars :=mux.Vars(r)
	movieName := vars["name"]
	// movieName := r.URL.Query().Get("name")
	baseURL := fmt.Sprintf("https://www.moviebuff.com/%s.json", movieName)
	resp, err := http.Get(baseURL)
	checkErr(err)
	reqBody, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	json.Unmarshal(reqBody, &movie)
	_, err = Db.Exec("INSERT INTO movie (movie_id,title,release_date,language_name,summary) VALUES ($1,$2,$3,$4,$5)", movie.Movie_id, movie.Title, movie.Realease_date, movie.Languge_name, movie.Summary)
	checkErr(err)
	_, err = Db.Exec("INSERT INTO running_time (movie_id,hours,minutes) VALUES ($1,$2,$3)", movie.Movie_id, movie.Running_time.Hours, movie.Running_time.Minutes)
	checkErr(err)
	for i := 0; i < len(movie.Cast); i++ {
		_, err = Db.Exec("INSERT INTO casts (movie_id,name,role,actor_id,poster) VALUES ($1,$2,$3,$4,$5)", movie.Movie_id, movie.Cast[i].Name, movie.Cast[i].Role, movie.Cast[i].Actor_id, movie.Cast[i].Poster)
		checkErr(err)
	}
	w.Header().Set("Content-Type", "application/json")
	// fmt.Println("request headers:", r.Header)
	// fmt.Println("response header:",w.Header())
    // userid := w.Header().Get("X-User-Id")
    // fmt.Println("user id:", userid)
	json.NewEncoder(w).Encode(movie)
}


// getting the data of the requested movie
func GetMovieDataByName(w http.ResponseWriter, r *http.Request) {
	vars :=mux.Vars(r)
	movieName := vars["name"]
	// movieName := r.URL.Query().Get("name")
	// var id retrieveId
	// var id string
	var movie retrieveMovie
	//to find the primary key
	query := fmt.Sprint("select m.movie_id,m.title,m.release_date,m.language_name, r.hours || ';' ||r.minutes as running_time from movie m inner join running_time r on m.movie_id=(select movie_id from movie where title ='"+movieName+"') and m.movie_id = r.movie_id;")
	rows, err := Db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		// err = rows.Scan(&id.Uuid)
		err = rows.Scan(&movie.Movie_id, &movie.Title, &movie.Realease_date, &movie.Languge_name,&movie.Running_time)
		checkErr(err)
	}
	//to get the data from the movie table
	// query = fmt.Sprint("select * from movie where movie_id='" + id + "'")//
	// rows, err = Db.Query(query)
	// checkErr(err)
	// defer rows.Close()
	// for rows.Next() {
	// 	err = rows.Scan(&movie.Movie_id, &movie.Title, &movie.Realease_date, &movie.Languge_name, &movie.Summary)
	// }
	// //to get the data from the running_time table
	// query = fmt.Sprint("select hours,minutes from running_time where movie_id='" + id + "'")//
	// rows, err = Db.Query(query)
	// checkErr(err)
	// defer rows.Close()
	// for rows.Next() {
	// 	err = rows.Scan(&movie.Running_time.Hours, &movie.Running_time.Minutes)
	// 	checkErr(err)
	// }
	//to get the data from the casts
	// query = fmt.Sprint("select name,role,actor_id,poster from casts where movie_id='" + id + "'")//
	// rows, err = Db.Query(query)
	// checkErr(err)
	// defer rows.Close()
	// var casts cast //created to use this and append them to the movie.cast structure
	// for rows.Next() {
	// 	err = rows.Scan(&casts.Name, &casts.Role, &casts.Actor_id, &casts.Poster)
	// 	checkErr(err)
	// 	movie.Cast = append(movie.Cast, casts)
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}