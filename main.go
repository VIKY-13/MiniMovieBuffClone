package main

import (
	"fmt"
	"golangmovietask/config"
	"golangmovietask/controllers"
	"golangmovietask/daos"
	"golangmovietask/db"
	"golangmovietask/middlewares"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/*
database name = MovieDetails
moviestructure from dto = movdata
post method = createdata
*/

// var err error

func main() {
	//DB connection
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	Db,err :=db.DatabaseConnection()
	if err != nil{
		log.Fatal("failed to connect Database")
		return
	}
	fmt.Println("connected to Db")
	defer Db.Close()
	daos.Init(Db)
	//server up process
	r:= mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            next.ServeHTTP(w, r)
        })
    })
	fmt.Println("starting server")
	r.HandleFunc("/movie/explore", controllers.ExploreMovies).Methods("GET")
	r.HandleFunc("/",welcome).Methods("GET")
	r.HandleFunc("/movie/rating/update", controllers.UpdateMovieRating).Methods("PUT")
	r.HandleFunc("/movie/rating/new", controllers.PostMovieRating).Methods("POST")
	r.HandleFunc("/user/login", controllers.UserLogin).Methods("POST")
	r.HandleFunc("/user/update", controllers.UpdateUserProfile).Methods("PUT")
	r.HandleFunc("/movie/create/{name}", middlewares.HostAuthentication(controllers.PostNewMovieData)).Methods("POST")
	r.HandleFunc("/movie/getmoviebyname/{name}", controllers.GetMovieDataByName).Methods("GET")
	r.HandleFunc("/minimoviebuff/endpoints", controllers.APIDocumentation).Methods("GET")
	r.HandleFunc("/movie/", controllers.GetMovieDataByQueryParams).Methods("GET")
	r.HandleFunc("/user/create",controllers.CreateNewUser).Methods("POST")
	r.HandleFunc("/user/watchlist/add",controllers.AddMovieToUserWatchlist).Methods("POST")
	r.HandleFunc("/user/watchlist/remove",controllers.RemoveMovieFromUserWatchlist).Methods("DELETE")
	r.HandleFunc("/user/watchlist",controllers.GetUserWatchlist).Methods("GET")
	r.HandleFunc("/user/favourites",controllers.GetUserFavourites).Methods("GET")
	r.HandleFunc("/user/favourite/add",controllers.AddUserFavourite).Methods("POST")
	r.HandleFunc("/user/favourite/remove",controllers.RemoveUserFavourite).Methods("DELETE")
	err = http.ListenAndServe(config.HOST_PORT, r)
	checkErr(err)
}

func welcome(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"welcome")
}


// error check function
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}