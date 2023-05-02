package daos

import (
	"golangmovietask/models"

	_ "github.com/lib/pq"
)


func GetMovieIdListOnUserWatchlist(user_id string) ([]models.RetrieveMovData,error){
	var movies_id []string
	var movie_id string
	statement,err := Db.Prepare("SELECT distinct movie_id FROM watchlist WHERE user_id = $1;")
	if err != nil{
		return nil,err
	}
	rows,err:= statement.Query(user_id)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		_ = rows.Scan(&movie_id)
		movies_id = append(movies_id, movie_id)
	}
	movies,_:= RetriveDataOnMovie_id(movies_id)
	return movies,nil
}

func AddMovieToUserWatchlistDb(addToWatchList models.Favourite) error{
	statement,err := Db.Prepare("INSERT INTO watchlist(user_id,movie_id) VALUES ($1,$2)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(addToWatchList.User_id,addToWatchList.Movie_id)
	if err != nil{
		return err
	}
	return nil
}

func RemoveMovieFromWatchlistDb(removeFromWatchlist models.Favourite) error{
	statement,err := Db.Prepare("DELETE FROM watchlist WHERE movie_id=$1 AND user_id=$2")
	if err != nil{
		return err
	}
	_,err = statement.Exec(removeFromWatchlist.Movie_id,removeFromWatchlist.User_id)
	if err != nil{
		return err
	}
	return nil
}