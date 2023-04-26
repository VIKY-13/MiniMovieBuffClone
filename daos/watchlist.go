package daos

import (
	"golangmovietask/dtos"

	_ "github.com/lib/pq"
)


func GetMovieIdListOnUserWatchlist(user_id string) ([]dtos.RetrieveMovData,error){
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

func AddMovieToUserWatchlistDb(addToWatchList dtos.Favourite) error{
	statement,err := Db.Prepare("insert into watchlist(user_id,movie_id) values($1,$2)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(addToWatchList.User_id,addToWatchList.Movie_id)
	if err != nil{
		return err
	}
	return nil
}

func RemoveMovieFromWatchlistDb(removeFromWatchlist dtos.Favourite) error{
	statement,err := Db.Prepare("DELETE from watchlist where movie_id=$1 and user_id=$2")
	if err != nil{
		return err
	}
	_,err = statement.Exec(removeFromWatchlist.Movie_id,removeFromWatchlist.User_id)
	if err != nil{
		return err
	}
	return nil
}