package main

import (
	// "database/sql"
	// "encoding/json"
	// "fmt"
	// "html/template"
	// "io/ioutil"
	// "fmt"
	"encoding/json"
	// "fmt"
	// "internal/syscall/windows"
	"io/ioutil"
	"net/http"
	"os"

	// "strings"

	// "github.com/gorilla/mux"
	// "github.com/joho/godotenv"
	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv"

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

//USER LOGIN
func UserLogin(w http.ResponseWriter, r *http.Request){
	var user userlogin
	var verifyuser userlogin
	resp,err := ioutil.ReadAll(r.Body)
	checkErr(err)
	json.Unmarshal(resp,&user)
	row,err := Db.Query("select password from users where email='"+user.Useremail+"';")
	checkErr(err)
	for row.Next(){
		err=row.Scan(&verifyuser.Password)
		checkErr(err)
	 }
	if user.Password==verifyuser.Password{
		w.WriteHeader(http.StatusOK)
	}else{
		w.WriteHeader(http.StatusUnauthorized)
	}
}