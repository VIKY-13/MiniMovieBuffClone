package daos

import (
	"database/sql"
	"errors"
	"fmt"

	// "golangmovietask/dtos"
	"golangmovietask/dtos"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func Init(DbConn *sql.DB) *sql.DB{
	Db = DbConn
	return Db
}

//To avoid Sql Injection problem we use separate function for each queries

func GetMovieIdListOnCast(castMember string) []dtos.RetrieveMovData{
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

func GetMovieIdListOnYear(year string) []dtos.RetrieveMovData{
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

func GetAllMovieIdList()[]dtos.RetrieveMovData{
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


func GetMovieIdListOnLanguage(language string) []dtos.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT distinct movie_id FROM movie WHERE language_name = $1;")
	rows,_ := statement.Query(language)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func GetMovieIdListOnCertification(certification string) []dtos.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT distinct movie_id FROM movie WHERE certification = $1;")
	rows,_ := statement.Query(certification)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func RetriveDataOnMovie_id(movies_id []string)([]dtos.RetrieveMovData ,error){
	var movies []dtos.RetrieveMovData 
	for i:=0;i<len(movies_id);i++{
		var reqMovie dtos.RetrieveMovData
		//new code
		//in the below statement we conntect 3 tables movie,running_time & genre of the movie 
		statement,err := Db.Prepare("select m.*, mg.genre from movie m inner join movie_genre mg on m.movie_id= mg.movie_id and m.movie_id= $1 ;")
		if err!= nil{
			fmt.Println("error in sql query movie")
			return movies,err
		}
		rows,err := statement.Query(movies_id[i])
		if err!= nil{
			fmt.Println("error in sql query movie")
			return movies,err
		}
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
		var castMember dtos.Cast //created to use this and append them to the movie.cast structure
		for rows.Next() {
			err = rows.Scan(&castMember.Name, &castMember.Role, &castMember.Cast_member_id, &castMember.Poster)
			if err!= nil{
				fmt.Println("error in scanning cast")
				return movies,err
			}
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
		var crewMember dtos.Crew //created to use this and append them to the reqMovie.crew structure
		for rows.Next() {
			err = rows.Scan(&crewMember.Name, &crewMember.Role, &crewMember.Crew_member_id, &crewMember.Poster)
			if err!= nil{
				fmt.Println("error in scanning crew")
				return movies,err
			}
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
			var userReviews dtos.UserReviews
			_ = rows.Scan(&userReviews.Username,&userReviews.Review,&userReviews.Rating)
			reqMovie.UserReviews= append(reqMovie.UserReviews, userReviews)
		}
		//if there is more than 1 movie has to be retrieved, we get the data in a array for such cases
		movies = append(movies, reqMovie)
	}
	return movies,nil
}

func GetMovieIdListOnName(movieName string) []dtos.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT distinct movie_id FROM movie WHERE title LIKE '%' || $1 || '%';")		// '||' will concate the string
	rows,_ := statement.Query(movieName)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}


func CheckMovieAlreadyExist(movie_id string)error{
	query := "SELECT COUNT(movie_id) FROM movie WHERE movie_id=$1;"
    var count int
    err := Db.QueryRow(query, movie_id).Scan(&count)
	if err != nil{
		return err
	}

    // If the count is 1, the data exists
    if count > 0 {
		return errors.New("movie aready exists")
    }
	return nil
}

func PostNewMovieDataToDb(movie dtos.MovData)error{
	//Preparing the query, Could've used inser all but for easy understanding purpose continued isering separately
	statement,err := Db.Prepare("INSERT INTO movie (movie_id,title,release_date,language_name,summary,certification,hours,minutes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)")
	if err!= nil{
		return err
	}
	//Executing the query
	_, err = statement.Exec(movie.Movie_id, movie.Title, movie.Realease_date, movie.Languge_name, movie.Summary, movie.Certification,movie.Running_time.Hours, movie.Running_time.Minutes)
	if err!= nil{
		return err
	}

	//insertig genre data of the movie
	for i := 0 ; i < len(movie.Genres) ; i++{
		statement,err = Db.Prepare("INSERT INTO movie_genre(movie_id,genre) VALUES ($1,$2)")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Movie_id,movie.Genres[i])
		if err!= nil{
			return err
		}
	}
	//inserting photos into the table
	for i := 0 ; i < len(movie.Photos) ; i++{
		statement,err = Db.Prepare("INSERT INTO movie_photos(movie_id,photos_url) VALUES ($1,$2)")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Movie_id,movie.Photos[i])
		if err!= nil{
			return err
		}
	}
	//insrting trailers url into the table
	for i := 0; i < len(movie.Trailers); i++{
		statement,err = Db.Prepare("INSERT INTO movie_trailer(movie_id,trailer_url) VALUES ($1,$2)")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Movie_id,movie.Trailers[i])
		if err!= nil{
			return err
		}
	}

	//for cast details

	for i := 0 ; i < len(movie.Cast) ; i++{
		//inserting only if there is no data of a cast member exist in the table 
		statement,err = Db.Prepare("INSERT INTO casts(cast_member_id, name, role, poster) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT * FROM casts WHERE cast_member_id = $5);")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Cast[i].Cast_member_id, movie.Cast[i].Name, movie.Cast[i].Role, movie.Cast[i].Poster,movie.Cast[i].Cast_member_id)
		if err!= nil{
			return err
		}
		//inseting cast members of a movie in the table
		statement,err = Db.Prepare("INSERT INTO movie_cast(movie_id,cast_member_id) VALUES ($1,$2)")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Movie_id,movie.Cast[i].Cast_member_id)
		if err!= nil{
			return err
		}
	}

	//for crew details

	for i := 0 ; i < len(movie.Crew) ; i++ {
		//inserting only if there is no data of a crew member exist in the table 
		statement,err = Db.Prepare("INSERT INTO crew(crew_member_id, name, role, poster) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT * FROM crew WHERE crew_member_id = $5);")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Crew[i].Crew_member_id, movie.Crew[i].Name, movie.Crew[i].Role, movie.Crew[i].Poster,movie.Crew[i].Crew_member_id)
		if err!= nil{
			return err
		}
		//inseting crew members of a movie in the table and this table contains the movieid and the crew memberid of that film
		statement,err = Db.Prepare("INSERT INTO movie_crew(movie_id,crew_member_id) VALUES ($1,$2)")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Movie_id,movie.Crew[i].Crew_member_id)
		if err!= nil{
			return err
		}
	}
	return nil
}

func UpdateMovieRatingInDb(updateRating dtos.MoieRating)error{
	statement,err := Db.Prepare("update reviews set rating=$1,review=$2 where user_id=$3 and movie_id=$4")
	if err != nil{
		return err
	}
	_,err = statement.Exec(updateRating.Rating,updateRating.Review,updateRating.User_id,updateRating.Movie_id)
	if err != nil{
		return err
	}
	return nil
}

func PostMovieRatingInDb(ratingData dtos.MoieRating)error{
	statement,err := Db.Prepare("INSERT INTO reviews(user_id,movie_id,rating,review) VALUES($1,$2,$3,$4)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(ratingData.User_id,ratingData.Movie_id,ratingData.Rating,ratingData.Review)
	if err != nil{
		return err
	}
	return nil
}