package main

import (
	// "database/sql"
	// "encoding/json"
	// "fmt"
	// "html/template"
	// "io/ioutil"
	// "fmt"
	"net/http"
	"os"

	// "strings"

	// "github.com/gorilla/mux"
	// "github.com/joho/godotenv"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

//middleware
func HostAuthentication(next http.HandlerFunc)http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		requsername,reqpassword,ok := r.BasicAuth()
		if ok{
			err := godotenv.Load(".env")
			checkErr(err)
			username := os.Getenv("AUTH_USERNAME")
			password := os.Getenv("AUTH_PASSWORD")
			if username==requsername && password==reqpassword{
				next.ServeHTTP(w,r)
				return
			}
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}