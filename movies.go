package main

import (
	// "database/sql"
	"encoding/json"
	// "errors"
	"fmt"

	// "io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ExploreMovies(w http.ResponseWriter, r *http.Request){
	movies := GetAllMovieIdList()
	w.Header().Set("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func GetMovieDataByQueryParams(w http.ResponseWriter, r *http.Request){
	tc := cases.Title(language.English)
	castMember := tc.String(r.URL.Query().Get("cast"))
	language_name := tc.String(r.URL.Query().Get("language"))
	year := tc.String(r.URL.Query().Get("year"))
	certification := tc.String(r.URL.Query().Get("certification"))
	var movies []retrieveMovData
	if(language_name!=""){
		movies=GetMovieIdList("movie","language_name",language_name)
	}else if (year!=""){
		movies=GetMovieIdListOnYear(year)
	}else if (len(certification)!=0){
		movies=GetMovieIdList("movie","certification",certification)
	}else if (castMember!= ""){
		movies = GetMovieIdListOnCast(castMember)
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

func GetMovieIdListOnCast(castMember string) []retrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("select mc.movie_id from movie_cast mc inner join casts c on mc.cast_member_id=c.cast_member_id and c.name= $1 ;")
	rows,_ := statement.Query(castMember)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func GetMovieIdListOnYear(year string) []retrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT distinct movie_id FROM movie WHERE date_part('year', release_date) = $1;")
	rows,_ := statement.Query(year)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func GetAllMovieIdList()[]retrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT movie_id FROM movie ORDER BY release_date desc;")
	rows,_ := statement.Query()
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func GetMovieIdList(table string,column string,parameter string) []retrieveMovData{
	var movies_id []string
	var movie_id string
	query := fmt.Sprintf("select distinct movie_id from "+table+" where "+column+"='"+parameter+"';")//this query is taken for the data based on the certification,language,cast
	rows, err := Db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&movie_id)
		checkErr(err)
		movies_id = append(movies_id,movie_id )
	}//we'll be getting the list of movie id's and based on that well retrieve the data
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func RetriveDataOnMovie_id(movies_id []string)([]retrieveMovData ,error){
	var movies []retrieveMovData 
	for i:=0;i<len(movies_id);i++{
		var reqMovie retrieveMovData
		//new code
		//in the below statement we conntect 3 tables movie,running_time & genre of the movie 
		statement,err := Db.Prepare("select m.*,rt.hours,rt.minutes, mg.genre from movie m inner join running_time rt on m.movie_id=rt.movie_id  and m.movie_id= $1 inner join movie_genre mg on m.movie_id= mg.movie_id;")
		if err!= nil{
			fmt.Println("error in sql query movie")
			return movies,err
		}
		rows,err := statement.Query(movies_id[i])
		checkErr(err)
		var genre string
		for rows.Next(){
			err:= rows.Scan(&reqMovie.Movie_id,&reqMovie.Title,&reqMovie.Realease_date,&reqMovie.Languge_name,&reqMovie.Summary,&reqMovie.Certification,&reqMovie.Running_time.Hours,&reqMovie.Running_time.Minutes,&genre)
			if err != nil{
				fmt.Println(err)
				return movies,err
			}
			reqMovie.Genres = append(reqMovie.Genres, genre)
		}
		//here we'll e getting the photos url
		statement,err = Db.Prepare("select distinct photos_url from movie_photos where movie_id =$1")
		var photo_url string
		if err!=nil{
			fmt.Println("error in photos query")
			return movies,err
		}
		rows,_ = statement.Query(movies_id[i])
		for rows.Next(){
			_ = rows.Scan(&photo_url)
			reqMovie.Photos = append(reqMovie.Photos, photo_url)
		}
		//here we'll be getting the trailer url
		statement,err = Db.Prepare("select distinct trailer_url from movie_trailer where movie_id =$1")
		var trailer_url string
		if err!=nil{
			fmt.Println("error in trailer query")
			return movies,err
		}
		rows,_ = statement.Query(movies_id[i])
		for rows.Next(){
			_ = rows.Scan(&trailer_url)
			reqMovie.Trailers = append(reqMovie.Trailers, trailer_url)
		}
		//in here we get the cast data by connecting movie_cast and casts table 
		statement,err = Db.Prepare("select c.name,c.role,c.cast_member_id,c.poster from movie_cast mc inner join casts c on mc.movie_id= $1 and mc.cast_member_id=c.cast_member_id;")
		if err!=nil{
			fmt.Println("error in cast query")
			return movies,err
		}
		rows,_ = statement.Query(movies_id[i])
		defer rows.Close()
		var castMember cast //created to use this and append them to the movie.cast structure
		for rows.Next() {
			err = rows.Scan(&castMember.Name, &castMember.Role, &castMember.Cast_member_id, &castMember.Poster)
			checkErr(err)
			reqMovie.Cast = append(reqMovie.Cast, castMember)
		}
		//here we get the crew data by connecting the movie_crew and crew table
		statement,err = Db.Prepare("select cr.name,cr.role,cr.crew_member_id,cr.poster from movie_crew mcr inner join crew cr on mcr.movie_id= $1 and mcr.crew_member_id=cr.crew_member_id")
		if err!=nil{
			fmt.Println("error in crew query")
			return movies,err
		}
		rows,_ = statement.Query(movies_id[i])
		defer rows.Close()
		var crewMember crew //created to use this and append them to the reqMovie.crew structure
		for rows.Next() {
			err = rows.Scan(&crewMember.Name, &crewMember.Role, &crewMember.Crew_member_id, &crewMember.Poster)
			checkErr(err)
			reqMovie.Crew = append(reqMovie.Crew, crewMember)
		}
		//Getting user reviews
		statement,err = Db.Prepare("SELECT ROUND(AVG(rating),1) FROM reviews WHERE movie_id= $1 ;")
		if err != nil{
			fmt.Println("error i overall rating")
			return movies,err
		}
		_ = statement.QueryRow(movies_id[i]).Scan(&reqMovie.OverallUserRating)
		//getting all user reviews for this movie
		statement,err = Db.Prepare("select u.firstname, r.review, r.rating from users u inner join reviews r on r.movie_id=$1 and  r.user_id=u.user_id;")
		if err != nil{
			fmt.Println("error in getting all user reviews")
			return movies,err
		}
		rows,_ = statement.Query(movies_id[i])
		for rows.Next(){
			var userReviews userreviews
			_ = rows.Scan(&userReviews.Username,&userReviews.Review,&userReviews.Rating)
			reqMovie.UserReviews= append(reqMovie.UserReviews, userReviews)
		}
		//if there is more than 1 movie has to be retrieved, we get the data in a array for such cases
		movies = append(movies, reqMovie)
	}
	return movies,nil
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
//THIS IS NEWLY ADDED PROGRAM
	//Preparing the query, Could've used inser all but for easy understanding purpose continued isering separately
	statement,err := Db.Prepare("INSERT INTO movie (movie_id,title,release_date,language_name,summary,certification) VALUES ($1,$2,$3,$4,$5,$6)")
	if err!= nil{
		fmt.Fprintf(w,"error in sql query movie")
		return 
	}
	//Executing the query
	_, err = statement.Exec(movie.Movie_id, movie.Title, movie.Realease_date, movie.Languge_name, movie.Summary, movie.Certification)
	if err!= nil{
		fmt.Println(err)
		fmt.Fprintf(w,"error in sql queryexecution movie")
	}
	//running time data into the running time table
	statement,err = Db.Prepare("INSERT INTO running_time (movie_id,hours,minutes) VALUES ($1,$2,$3)")
	if err!= nil{
		fmt.Fprintf(w,"error in sql query running_time")
		return 
	}
	_,err = statement.Exec(movie.Movie_id, movie.Running_time.Hours, movie.Running_time.Minutes)
	if err!= nil{
		fmt.Fprintf(w,"error in sql query execution running_time")
		return 
	}
	//insertig genre data of the movie
	for _,genre := range movie.Genres{
		statement,err = Db.Prepare("insert into movie_genre(movie_id,genre) values ($1,$2)")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query movie_genre")
			return 
		}
		_,err = statement.Exec(movie.Movie_id,genre)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query execution genre")
			return 
		}
	}
	//inserting photos into the table
	for _,photos_url := range movie.Photos{
		statement,err = Db.Prepare("insert into movie_photos(movie_id,photos_url) values ($1,$2)")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query movie_photos")
			return 
		}
		_,err = statement.Exec(movie.Movie_id,photos_url)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query execution photos")
			return 
		}
	}
	//insrting trailers url into the table
	for _,trailer_url := range movie.Trailers{
		statement,err = Db.Prepare("insert into movie_trailer(movie_id,trailer_url) values ($1,$2)")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query movie_trailer")
			return 
		}
		_,err = statement.Exec(movie.Movie_id,trailer_url)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query exe trailer")
			return 
		}
	}

	//for cast details

	for i := 0; i < len(movie.Cast); i++{
		//inserting only if there is no data of a cast member exist in the table 
		statement,err = Db.Prepare("INSERT INTO casts(cast_member_id, name, role, poster) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT * FROM casts WHERE cast_member_id = $5);")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query casts")
			return 
		}
		_,err = statement.Exec(movie.Cast[i].Cast_member_id, movie.Cast[i].Name, movie.Cast[i].Role, movie.Cast[i].Poster,movie.Cast[i].Cast_member_id)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query exe casts 1")
			return 
		}
		//inseting cast members of a movie in the table
		statement,err = Db.Prepare("insert into movie_cast(movie_id,cast_member_id) values($1,$2)")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query movie_cast")
			return 
		}
		_,err = statement.Exec(movie.Movie_id,movie.Cast[i].Cast_member_id)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query exe casts 2")
			return 
		}
	}
	
	//for crew details

	for i := 0; i < len(movie.Cast); i++{
		//inserting only if there is no data of a cast member exist in the table 
		statement,err = Db.Prepare("INSERT INTO crew(crew_member_id, name, role, poster) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT * FROM crew WHERE crew_member_id = $5);")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query crew")
			return 
		}
		_,err = statement.Exec(movie.Crew[i].Crew_member_id, movie.Crew[i].Name, movie.Crew[i].Role, movie.Crew[i].Poster,movie.Crew[i].Crew_member_id)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query exe crew 1")
			return 
		}
		//inseting cast members of a movie in the table and this table contains the movieid and the crew memberid of that film
		statement,err = Db.Prepare("insert into movie_crew(movie_id,crew_member_id) values($1,$2)")
		if err!= nil{
			fmt.Fprintf(w,"error in sql query movie_crew")
			return 
		}
		_,err = statement.Exec(movie.Movie_id,movie.Crew[i].Crew_member_id)
		if err!= nil{
			fmt.Fprintf(w,"error in sql query exe crew 2")
			return 
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}


// getting the data of the requested movie
func GetMovieDataByName(w http.ResponseWriter, r *http.Request) {
	vars :=mux.Vars(r)
	tc := cases.Title(language.English)
	movieName := tc.String(vars["name"]) //will be getting the movie name from the url
	reqmovie:=GetMovieIdList("movie","title",movieName) 
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