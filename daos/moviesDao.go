package daos

import (
	"database/sql"
	"errors"
	"fmt"

	"golangmovietask/models"
	"golangmovietask/dtos"
	_ "github.com/lib/pq"
)

type DAO struct{
	Db *sql.DB
}

func Init(DbConn *sql.DB) *DAO{
	return &DAO{
		Db : DbConn,
	}
}

//To avoid Sql Injection problem we use separate function for each queries

func (m *DAO) GetMovieIdListOnCast(castMember string) []models.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := m.Db.Prepare("SELECT mc.movie_id FROM movie_cast mc INNER JOIN casts c ON mc.cast_member_id=c.cast_member_id AND c.name= $1 ;")
	rows,_ := statement.Query(castMember)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= m.RetriveDataOnMovie_id(movies_id)
	return movies
}

func (m *DAO) GetMovieIdListOnYear(year string) []models.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := m.Db.Prepare("SELECT DISTINCT movie_id FROM movie WHERE date_part('year', release_date) = $1;")
	rows,_ := statement.Query(year)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= m.RetriveDataOnMovie_id(movies_id)
	return movies
}

func (m *DAO) GetAllMovieIdList()([]models.RetrieveMovData,error){
	var movies_id []string
	var movie_id string
	statement,err := m.Db.Prepare("SELECT movie_id FROM movie ORDER BY release_date desc;")
	if err != nil{
		return nil,err
	}
	rows,_ := statement.Query()
	for rows.Next(){
		err = rows.Scan(&movie_id)
		if err != nil{
			return nil,err
		}
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= m.RetriveDataOnMovie_id(movies_id)
	return movies,nil
}


func (m *DAO) GetMovieIdListOnLanguage(language string) []models.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := m.Db.Prepare("SELECT DISTINCT movie_id FROM movie WHERE language_name = $1;")
	rows,_ := statement.Query(language)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= m.RetriveDataOnMovie_id(movies_id)
	return movies
}

func (m *DAO) GetMovieIdListOnCertification(certification string) []models.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := m.Db.Prepare("SELECT DISTINCT movie_id FROM movie WHERE certification = $1;")
	rows,_ := statement.Query(certification)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= m.RetriveDataOnMovie_id(movies_id)
	return movies
}

func (m *DAO) RetriveDataOnMovie_id(movies_id []string)([]models.RetrieveMovData ,error){
	var movies []models.RetrieveMovData 
	for i:=0;i<len(movies_id);i++{
		var reqMovie models.RetrieveMovData
		//new code
		//in the below statement we conntect 3 tables movie,running_time & genre of the movie 
		statement,err := m.Db.Prepare("SELECT m.*, mg.genre FROM movie m INNER JOIN movie_genre mg ON m.movie_id= mg.movie_id AND m.movie_id= $1 ;")
		if err!= nil{
			fmt.Println("error in sql query movie")
			return nil,err
		}
		rows,err := statement.Query(movies_id[i])
		if err!= nil{
			fmt.Println("error in sql query movie")
			return nil,err
		}
		var genre string
		for rows.Next(){
			err:= rows.Scan(&reqMovie.Movie_id,&reqMovie.Title,&reqMovie.Realease_date,&reqMovie.Languge_name,&reqMovie.Summary,&reqMovie.Certification,&reqMovie.Running_time.Hours,&reqMovie.Running_time.Minutes,&genre)
			if err != nil{
				fmt.Println(err)
				return nil,err
			}
			reqMovie.Genres = append(reqMovie.Genres, genre)
		}
		//here we'll e getting the photos url
		statement,err = m.Db.Prepare("SELECT DISTINCT photos_url FROM movie_photos WHERE movie_id =$1")
		var photo_url string
		if err!=nil{
			fmt.Println("error in photos query")
			return nil,err
		}
		rows,_ = statement.Query(movies_id[i])
		for rows.Next(){
			_ = rows.Scan(&photo_url)
			reqMovie.Photos = append(reqMovie.Photos, photo_url)
		}
		//here we'll be getting the trailer url
		statement,err = m.Db.Prepare("SELECT DISTINCT trailer_url FROM movie_trailer WHERE movie_id =$1")
		var trailer_url string
		if err!=nil{
			fmt.Println("error in trailer query")
			return nil,err
		}
		rows,_ = statement.Query(movies_id[i])
		for rows.Next(){
			_ = rows.Scan(&trailer_url)
			reqMovie.Trailers = append(reqMovie.Trailers, trailer_url)
		}
		//in here we get the cast data by connecting movie_cast and casts table 
		statement,err = m.Db.Prepare("SELECT c.name,c.role,c.cast_member_id,c.poster FROM movie_cast mc INNER JOIN casts c ON mc.movie_id= $1 AND mc.cast_member_id=c.cast_member_id;")
		if err!=nil{
			fmt.Println("error in cast query")
			return nil,err
		}
		rows,_ = statement.Query(movies_id[i])
		defer rows.Close()
		var castMember models.Cast //created to use this and append them to the movie.cast structure
		for rows.Next() {
			err = rows.Scan(&castMember.Name, &castMember.Role, &castMember.Cast_member_id, &castMember.Poster)
			if err!= nil{
				fmt.Println("error in scanning cast")
				return nil,err
			}
			reqMovie.Cast = append(reqMovie.Cast, castMember)
		}
		//here we get the crew data by connecting the movie_crew and crew table
		statement,err = m.Db.Prepare("SELECT cr.name,cr.role,cr.crew_member_id,cr.poster FROM movie_crew mcr INNER JOIN crew cr ON mcr.movie_id= $1 AND mcr.crew_member_id=cr.crew_member_id")
		if err!=nil{
			fmt.Println("error in crew query")
			return nil,err
		}
		rows,_ = statement.Query(movies_id[i])
		defer rows.Close()
		var crewMember models.Crew //created to use this and append them to the reqMovie.crew structure
		for rows.Next() {
			err = rows.Scan(&crewMember.Name, &crewMember.Role, &crewMember.Crew_member_id, &crewMember.Poster)
			if err!= nil{
				fmt.Println("error in scanning crew")
				return nil,err
			}
			reqMovie.Crew = append(reqMovie.Crew, crewMember)
		}
		//Getting user reviews
		statement,err = m.Db.Prepare("SELECT ROUND(AVG(rating),1) FROM reviews WHERE movie_id= $1 ;")
		if err != nil{
			fmt.Println("error i overall rating")
			return nil,err
		}
		_ = statement.QueryRow(movies_id[i]).Scan(&reqMovie.OverallUserRating)
		//getting all user reviews for this movie
		statement,err = m.Db.Prepare("SELECT u.firstname, r.review, r.rating FROM users u INNER JOIN reviews r ON r.movie_id=$1 AND  r.user_id=u.user_id;")
		if err != nil{
			fmt.Println("error in getting all user reviews")
			return nil,err
		}
		rows,_ = statement.Query(movies_id[i])
		for rows.Next(){
			var userReviews models.UserReviews
			_ = rows.Scan(&userReviews.Username,&userReviews.Review,&userReviews.Rating)
			reqMovie.UserReviews= append(reqMovie.UserReviews, userReviews)
		}
		//if there is more than 1 movie has to be retrieved, we get the data in a array for such cases
		movies = append(movies, reqMovie)
	}
	return movies,nil
}

func (m *DAO) GetMovieIdListOnName(movieName string) ([]models.RetrieveMovData,error){
	var movies_id []string
	var movie_id string
	statement,_ := m.Db.Prepare("SELECT DISTINCT movie_id FROM movie WHERE title LIKE '%' || $1 || '%';")		// '||' will concate the string
	rows,_ := statement.Query(movieName)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,err:= m.RetriveDataOnMovie_id(movies_id)
	if err != nil{
		return nil,err
	}
	return movies,nil
}


func (m *DAO) CheckMovieAlreadyExist(movie_id string)error{
	query := "SELECT COUNT(movie_id) FROM movie WHERE movie_id=$1;"
    var count int
    err := m.Db.QueryRow(query, movie_id).Scan(&count)
	if err != nil{
		return err
	}

    // If the count is 1, the data exists
    if count > 0 {
		return errors.New("movie aready exists")
    }
	return nil
}

func (m *DAO) PostNewMovieDataToDb(movie dtos.MovData)error{
	//Preparing the query, Could've used inser all but for easy understanding purpose continued isering separately
	statement,err := m.Db.Prepare("INSERT INTO movie (movie_id,title,release_date,language_name,summary,certification,hours,minutes) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)")
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
		statement,err = m.Db.Prepare("INSERT INTO movie_genre(movie_id,genre) VALUES ($1,$2)")
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
		statement,err = m.Db.Prepare("INSERT INTO movie_photos(movie_id,photos_url) VALUES ($1,$2)")
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
		statement,err = m.Db.Prepare("INSERT INTO movie_trailer(movie_id,trailer_url) VALUES ($1,$2)")
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
		statement,err = m.Db.Prepare("INSERT INTO casts(cast_member_id, name, role, poster) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT * FROM casts WHERE cast_member_id = $5);")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Cast[i].Cast_member_id, movie.Cast[i].Name, movie.Cast[i].Role, movie.Cast[i].Poster,movie.Cast[i].Cast_member_id)
		if err!= nil{
			return err
		}
		//inseting cast members of a movie in the table
		statement,err = m.Db.Prepare("INSERT INTO movie_cast(movie_id,cast_member_id) VALUES ($1,$2)")
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
		statement,err = m.Db.Prepare("INSERT INTO crew(crew_member_id, name, role, poster) SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT * FROM crew WHERE crew_member_id = $5);")
		if err!= nil{
			return err
		}
		_,err = statement.Exec(movie.Crew[i].Crew_member_id, movie.Crew[i].Name, movie.Crew[i].Role, movie.Crew[i].Poster,movie.Crew[i].Crew_member_id)
		if err!= nil{
			return err
		}
		//inseting crew members of a movie in the table and this table contains the movieid and the crew memberid of that film
		statement,err = m.Db.Prepare("INSERT INTO movie_crew(movie_id,crew_member_id) VALUES ($1,$2)")
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

func (m *DAO) UpdateMovieRatingInDb(updateRating models.MovieRating)error{
	statement,err := m.Db.Prepare("update reviews set rating=$1,review=$2 where user_id=$3 and movie_id=$4")
	if err != nil{
		return err
	}
	_,err = statement.Exec(updateRating.Rating,updateRating.Review,updateRating.User_id,updateRating.Movie_id)
	if err != nil{
		return err
	}
	return nil
}

func (m *DAO) PostMovieRatingInDb(ratingData models.MovieRating)error{
	statement,err := m.Db.Prepare("INSERT INTO reviews(user_id,movie_id,rating,review) VALUES($1,$2,$3,$4)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(ratingData.User_id,ratingData.Movie_id,ratingData.Rating,ratingData.Review)
	if err != nil{
		return err
	}
	return nil
}