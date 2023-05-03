package services

import (
	"encoding/json"
	"fmt"
	"golangmovietask/daos"
	"golangmovietask/dtos"
	"golangmovietask/models"
	"net/http"

	// "github.com/gorilla/mux"
	// "golang.org/x/text/cases"
	// "golang.org/x/text/language"
)

type Service struct{
	DAO *daos.DAO
}

func (m *Service) ExploreMovieService(w http.ResponseWriter)([]models.RetrieveMovData,error){
	movies,err := m.DAO.GetAllMovieIdList()
	if err != nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil,err
	}
	w.Header().Set("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	return movies,nil
}

func (m *Service) GetMovieDataByQueryParamsService(w http.ResponseWriter,castMember,language,year,certification string){
	var movies []models.RetrieveMovData
	if(language!=""){
		movies= m.DAO.GetMovieIdListOnLanguage(language)
	}else if (year!=""){
		movies= m.DAO.GetMovieIdListOnYear(year)
	}else if (len(certification)!=0){
		movies= m.DAO.GetMovieIdListOnCertification(certification)
	}else if (castMember!= ""){
		movies = m.DAO.GetMovieIdListOnCast(castMember)
	}else{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Check the query")
		return
	}
	w.Header().Set("content-type", "application/json")
	if movies!=nil{
		json.NewEncoder(w).Encode(movies)
	}else{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w,"we dont have any data related to this")
		return
	}
}

func (m *Service) GetMovieDataByNameService(w http.ResponseWriter, movieName string){
	reqmovie,_:= m.DAO.GetMovieIdListOnName(movieName) 
	if reqmovie!=nil{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reqmovie)
	}else{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w,"check the movie name")
		return
	}
}

func (m *Service) PostNewMovieDataService(w http.ResponseWriter,movieName string){
	baseURL := fmt.Sprintf("https://www.moviebuff.com/%s.json", movieName)
	resp, err := http.Get(baseURL)
	if err != nil{
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var movie dtos.MovData
	json.NewDecoder(resp.Body).Decode(&movie)
	//checking whether the data is already exist
	err = m.DAO.CheckMovieAlreadyExist(movie.Movie_id)
	if err != nil{
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w,err)
		return
	}
	err = m.DAO.PostNewMovieDataToDb(movie)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return
	}	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (m *Service) UpdateMovieRatingService(w http.ResponseWriter,updateRating models.MovieRating){
	err := m.DAO.UpdateMovieRatingInDb(updateRating)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateRating)
}

func (m *Service) PostMovieRatingService(w http.ResponseWriter,ratingData models.MovieRating){
	err := m.DAO.PostMovieRatingInDb(ratingData)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ratingData)
}