package controllers

import (
	"encoding/json"
	"fmt"
	"golangmovietask/daos"
	"golangmovietask/dtos"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ExploreMovies(w http.ResponseWriter, r *http.Request){
	movies := daos.GetAllMovieIdList()
	w.Header().Set("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func GetMovieDataByQueryParams(w http.ResponseWriter, r *http.Request){
	tc := cases.Title(language.English)
	castMember := tc.String(r.URL.Query().Get("cast"))
	language := tc.String(r.URL.Query().Get("language"))
	year := tc.String(r.URL.Query().Get("year"))
	certification := tc.String(r.URL.Query().Get("certification"))
	var movies []dtos.RetrieveMovData
	if(language!=""){
		movies=daos.GetMovieIdListOnLanguage(language)
	}else if (year!=""){
		movies=daos.GetMovieIdListOnYear(year)
	}else if (len(certification)!=0){
		movies=daos.GetMovieIdListOnCertification(certification)
	}else if (castMember!= ""){
		movies = daos.GetMovieIdListOnCast(castMember)
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

// getting the data of the requested movie name
func GetMovieDataByName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"movie name")
	vars :=mux.Vars(r)
	tc := cases.Title(language.English)
	movieName := tc.String(vars["name"]) //will be getting the movie name from the url
	reqmovie:=daos.GetMovieIdListOnName(movieName) 
	if reqmovie!=nil{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reqmovie)
	}else{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w,"check the movie name")
		return
	}
	
}

func PostNewMovieData(w http.ResponseWriter, r *http.Request) {
	var movie dtos.MovData
	vars :=mux.Vars(r)
	movieName := vars["name"]
	baseURL := fmt.Sprintf("https://www.moviebuff.com/%s.json", movieName)
	resp, err := http.Get(baseURL)
	if err != nil{
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	json.NewDecoder(resp.Body).Decode(&movie)
	//checking whether the data is already exist
	err = daos.CheckMovieAlreadyExist(movie.Movie_id)
	if err != nil{
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w,err)
		return
	}
	err = daos.PostNewMovieDataToDb(movie)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return
	}	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovieRating(w http.ResponseWriter, r *http.Request){
	var updateRating dtos.MoieRating
	json.NewDecoder(r.Body).Decode(&updateRating)
	err := daos.UpdateMovieRatingInDb(updateRating)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateRating)
}

func PostMovieRating(w http.ResponseWriter, r *http.Request){
	var ratingData dtos.MoieRating
	json.NewDecoder(r.Body).Decode(&ratingData)
	err := daos.PostMovieRatingInDb(ratingData)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ratingData)
}