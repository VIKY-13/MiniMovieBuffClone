package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"golang.org/x/text/cases"
)

func ExploreMovies(w http.ResponseWriter, r *http.Request){
	movies := GetMovieIdList("movie","","")
	w.Header().Set("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func GetMovieDataByQueryParams(w http.ResponseWriter, r *http.Request){
	tc := cases.Title(language.English)
	cast := tc.String(r.URL.Query().Get("cast"))
	language_name := tc.String(r.URL.Query().Get("language"))
	release_date := tc.String(r.URL.Query().Get("year"))
	certification := tc.String(r.URL.Query().Get("certification"))
	var movies []retrieveMovie
	if(language_name!=""){
		movies=GetMovieIdList("movie","language_name",language_name)
	}else if (len(release_date)!=0){
		movies=GetMovieIdList("movie","release_date",release_date)
	}else if (len(certification)!=0){
		movies=GetMovieIdList("movie","certification",certification)
	}else if (cast!= ""){
		movies =GetMovieIdList("casts","name",cast)
	}else{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Check the query")
		return
	}
	// fmt.Print(movies)
	w.Header().Set("content-type", "application/json")
	if movies!=nil{
		json.NewEncoder(w).Encode(movies)
	}else{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w,"we dont have any data related to this")
		return
	}
}

func GetMovieIdList(table string,column string,parameter string) []retrieveMovie{
	var movies_id []string
	var movie_id string
	var query string
	if column=="release_date"{
		query = fmt.Sprintf("SELECT distinct movie_id FROM "+table+" WHERE date_part('year', release_date) = '"+parameter+"';")
	}else if (parameter==""){
		query = fmt.Sprintf("SELECT movie_id FROM "+table+" ORDER BY release_date desc;")
	}else{
		query = fmt.Sprintf("select distinct movie_id from "+table+" where "+column+"='"+parameter+"';")
	}
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
		query := fmt.Sprintf("select m.movie_id,m.title,m.release_date,m.language_name, r.hours || ';' ||r.minutes as running_time from movie m inner join running_time r on m.movie_id='"+movies_id[i]+"' and m.movie_id=r.movie_id;")
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
	baseURL := fmt.Sprintf("https://www.moviebuff.com/%s.json", movieName)
	resp, err := http.Get(baseURL)
	checkErr(err)
	json.NewDecoder(resp.Body).Decode(&movie)
	//checking whether the data is already exist
	query := "SELECT COUNT(movie_id) FROM movie WHERE movie_id=$1;"
    var count int
    err = Db.QueryRow(query, movie.Movie_id).Scan(&count)
    checkErr(err)

    // If the count is 1, the data exists
    if count >= 1 {
        w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w,"data already exist")
		return
    }

	_, err = Db.Exec("INSERT INTO movie (movie_id,title,release_date,language_name,summary) VALUES ($1,$2,$3,$4,$5)", movie.Movie_id, movie.Title, movie.Realease_date, movie.Languge_name, movie.Summary)
	checkErr(err)
	_, err = Db.Exec("INSERT INTO running_time (movie_id,hours,minutes) VALUES ($1,$2,$3)", movie.Movie_id, movie.Running_time.Hours, movie.Running_time.Minutes)
	checkErr(err)
	for i := 0; i < len(movie.Cast); i++ {
		_, err = Db.Exec("INSERT INTO casts (movie_id,name,role,actor_id,poster) VALUES ($1,$2,$3,$4,$5)", movie.Movie_id, movie.Cast[i].Name, movie.Cast[i].Role, movie.Cast[i].Actor_id, movie.Cast[i].Poster)
		checkErr(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}


// getting the data of the requested movie
func GetMovieDataByName(w http.ResponseWriter, r *http.Request) {
	vars :=mux.Vars(r)
	tc := cases.Title(language.English)
	movieName := tc.String(vars["name"])
	// var reqmovie []retrieveMovie
	reqmovie:=GetMovieIdList("movie","title",movieName)
//here
	//to find the primary key
	// query := fmt.Sprint("select m.movie_id,m.title,m.release_date,m.language_name, r.hours || ';' ||r.minutes as running_time from movie m inner join running_time r on m.movie_id=(select movie_id from movie where title ='"+movieName+"') and m.movie_id = r.movie_id;")
	// rows, err := Db.Query(query)
	// checkErr(err)
	// defer rows.Close()
	// for rows.Next() {
	// 	// err = rows.Scan(&id.Uuid)
	// 	err = rows.Scan(&movie.Movie_id, &movie.Title, &movie.Realease_date, &movie.Languge_name,&movie.Running_time)
	// 	checkErr(err)
	// }
//here

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


	if reqmovie!=nil{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reqmovie)
	}else{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w,"check the movie name")
		return
	}
	
}

func UpdateMovieRating(w http.ResponseWriter, r *http.Request){
	var updaterating movierating
	json.NewDecoder(r.Body).Decode(&updaterating)
	_,err := Db.Query("update reviews set rating=$1,review=$2 where user_id=$3 and movie_id=$4",updaterating.Rating,updaterating.Review,updaterating.User_id,updaterating.Movie_id)
	checkErr(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updaterating)
}

func PostMovieRating(w http.ResponseWriter, r *http.Request){
	var ratingdata movierating
	json.NewDecoder(r.Body).Decode(&ratingdata)
	_,err := Db.Exec("insert into reviews(user_id,movie_id,rating,review) values($1,$2,$3,$4)",ratingdata.User_id,ratingdata.Movie_id,ratingdata.Rating,ratingdata.Review)
	checkErr(err)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ratingdata)
}