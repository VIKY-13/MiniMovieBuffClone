package controllers

import (
	"encoding/json"
	"golangmovietask/models"
	"golangmovietask/services"
	"net/http"
	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)


//will be using the same structure for all the files in the cotrollers package
type Controllers struct{
	Service *services.Service
}

//This function returns all the movies in the Db
func (m *Controllers) ExploreMovies(w http.ResponseWriter, r *http.Request){
	movies,err := m.Service.ExploreMovieService(w)
	if err != nil{
		return
	}
	json.NewEncoder(w).Encode(movies)
}

func (m *Controllers) GetMovieDataByQueryParams(w http.ResponseWriter, r *http.Request){
	tc := cases.Title(language.English)
	castMember := tc.String(r.URL.Query().Get("cast"))
	language := tc.String(r.URL.Query().Get("language"))
	year := tc.String(r.URL.Query().Get("year"))
	certification := tc.String(r.URL.Query().Get("certification"))
	m.Service.GetMovieDataByQueryParamsService(w,castMember,language,year,certification)
}

// getting the data of the requested movie name
func (m *Controllers) GetMovieDataByName(w http.ResponseWriter, r *http.Request) {
	vars :=mux.Vars(r)
	tc := cases.Title(language.English)
	movieName := tc.String(vars["name"]) //will be getting the movie name from the url
	m.Service.GetMovieDataByNameService(w,movieName)	
}

func (m *Controllers) PostNewMovieData(w http.ResponseWriter, r *http.Request) {
	vars :=mux.Vars(r)
	movieName := vars["name"]
	m.Service.PostNewMovieDataService(w,movieName)
}

func (m *Controllers) UpdateMovieRating(w http.ResponseWriter, r *http.Request){
	var updateRating models.MovieRating
	err := json.NewDecoder(r.Body).Decode(&updateRating)
	if err != nil{
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	m.Service.UpdateMovieRatingService(w,updateRating)
}

func (m *Controllers) PostMovieRating(w http.ResponseWriter, r *http.Request){
	var ratingData models.MovieRating
	err := json.NewDecoder(r.Body).Decode(&ratingData)
	if err != nil{
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	m.Service.PostMovieRatingService(w,ratingData)
}