package daos

import (
	"golangmovietask/models"

	_ "github.com/lib/pq"
)

func AddMovieToUserFavouriteDb(addFavourite models.Favourite) error{
	statement,err := Db.Prepare("INSERT INTO favourite(user_id,movie_id) VALUES($1,$2)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(addFavourite.User_id,addFavourite.Movie_id)
	if err != nil{
		return err
	}
	return nil
}

func GetMovieIdListOnUserFavourite(user_id string) []models.RetrieveMovData{
	var movies_id []string
	var movie_id string
	statement,_ := Db.Prepare("SELECT DISTINCT movie_id FROM favourite WHERE user_id = $1;")
	rows,_ := statement.Query(user_id)
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies
}

func RemoveMovieFromFavouriteDb(removeFavourite models.Favourite) error{
	statement,err := Db.Prepare("DELETE FROM favourite WHERE movie_id= $1 AND user_id= $2")
	if err != nil{
		return err
	}
	_,err = statement.Exec(removeFavourite.Movie_id,removeFavourite.User_id)
	if err != nil{
		return err
	}
	return nil
}