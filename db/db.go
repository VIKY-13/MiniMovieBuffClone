package db

import(
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
	"golangmovietask/config"
)

func DatabaseConnection()(*sql.DB,error){
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, config.DB_NAME)
	db, err := sql.Open("postgres", postgresqlDbInfo) //make the connection
	if err != nil{
		return nil,err
	}
	return db,nil
}