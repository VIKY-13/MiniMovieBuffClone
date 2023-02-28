package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

/*
database name = MovieDetails
moviestructure from dto = movdata
post method = createdata
*/

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "MovieDetails"
)

// type retrieveId struct {
// 	Uuid string
// }

type movdata struct {
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  running_time `json:"running_time`
	Summary       string       `json:"summary"`
	Cast          []cast       `json:"cast"`
}

type running_time struct {
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
}

type cast struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Actor_id string `json:"uuid"`
	Poster   string `json:"poster"`
}

type retrieveMovie struct{
	Movie_id      string       `json:"uuid"`
	Title         string       `json:"title"`
	Realease_date string       `json:"release_date"`
	Languge_name  string       `json:"language_name"`
	Running_time  string	   `json:"running_time`
}


var Db *sql.DB
var err error


func main() {
	//DB connection
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", postgresqlDbInfo) //make the connection
	checkErr(err)
	defer db.Close()
	err = db.Ping() // checks
	checkErr(err)
	Db = db
	fmt.Println("Connected to the Database")

	//server up process
	fmt.Println("starting server")
	http.HandleFunc("/movie/create", PostNewMovieData)
	http.HandleFunc("/movie/read", GetMovieDataByName)
	http.ListenAndServe(":8000", nil)
}


// post function, creating new data in the DB
func PostNewMovieData(w http.ResponseWriter, r *http.Request) {
	var movie movdata
	movieName := r.URL.Query().Get("name")
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
	json.NewEncoder(w).Encode(movie)
}


// getting the data of the requested movie
func GetMovieDataByName(w http.ResponseWriter, r *http.Request) {
	movieName := r.URL.Query().Get("name")
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
	json.NewEncoder(w).Encode(movie)
}


// error check function
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
