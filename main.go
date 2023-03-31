package main

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	"html/template"
	// "io/ioutil"
	"net/http"
	"os"
	// "strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

/*
database name = MovieDetails
moviestructure from dto = movdata
post method = createdata
*/


var Db *sql.DB
// var err error
var templ *template.Template

func main() {
	//DB connection
	err := godotenv.Load(".env")
	checkErr(err)
	port := os.Getenv("PORT")
	DatabaseConnection()
	//server up process
	r:= mux.NewRouter()
	fmt.Println("starting server")
	r.HandleFunc("/movie/explore", ExploreMovies).Methods("GET")
	r.HandleFunc("/",welcome).Methods("GET")
	r.HandleFunc("/movie/rating/update", UpdateMovieRating).Methods("PUT")
	r.HandleFunc("/movie/rating/new", PostMovieRating).Methods("POST")
	r.HandleFunc("/user/login", UserLogin).Methods("POST")
	r.HandleFunc("/user/update", UpdateUserProfile).Methods("PUT")
	r.HandleFunc("/movie/create/{name}", HostAuthentication(PostNewMovieData)).Methods("POST")
	r.HandleFunc("/movie/getmoviebyname/{name}", GetMovieDataByName).Methods("GET")
	r.HandleFunc("/minimovibuff/endpoints", APIDocumentation).Methods("GET")
	r.HandleFunc("/movie/", GetMovieDataByQueryParams).Methods("GET")
	r.HandleFunc("/user/create",CreateNewUser).Methods("POST")
	r.HandleFunc("/user/watchlist/add",AddMovieToUserWatchlist).Methods("POST")
	r.HandleFunc("/user/watchlist/remove",RemoveMovieFromUserWatchlist).Methods("DELETE")
	r.HandleFunc("/user/watchlist",GetUserWatchlist).Methods("GET")
	r.HandleFunc("/user/favourites",GetUserFavourites).Methods("GET")
	r.HandleFunc("/user/favourite/add",AddUserFavourite).Methods("POST")
	r.HandleFunc("/user/favourite/remove",RemoveUserFavourite).Methods("DELETE")
	http.ListenAndServe("0.0.0.0:"+port, r)
}

func welcome(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"welcome")
}

func DatabaseConnection(){
	db_host,db_port,db_user,db_password,db_name := os.Getenv("DB_HOST"),os.Getenv("DB_PORT"),os.Getenv("DB_USER"),os.Getenv("DB_PASSWORD"),os.Getenv("DB_NAME")
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	db, err := sql.Open("postgres", postgresqlDbInfo) //make the connection
	checkErr(err)
	// defer db.Close()
	err = db.Ping() // checks
	checkErr(err)
	Db = db
	fmt.Println("Connected to the Database")
}

// error check function
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}