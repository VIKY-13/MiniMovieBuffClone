package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/google/uuid"
	
	_ "github.com/lib/pq"
)


//user creation
func CreateNewUser(w http.ResponseWriter, r *http.Request){
	var newuser user
	uuid := uuid.New()
	// fmt.Println("ID Generated:",myUUID.String())
	reqBody,err:=ioutil.ReadAll(r.Body)
	checkErr(err)
	json.Unmarshal(reqBody,&newuser)
	_,err = Db.Exec("insert into users(user_id,firstname,lastname,email,password,age,phone_no) values($1,$2,$3,$4,$5,$6,$7)",uuid,newuser.Firstname,newuser.Lastname,strings.TrimSpace(newuser.Email),newuser.Password,newuser.Age,newuser.Phone_no)
	if err!=nil{
		fmt.Fprintf(w,"enter a valid username & password")
		fmt.Println(err)
	}
	for i:=0;i<len(newuser.Genre);i++{
		_,err = Db.Exec("insert into usergenre(user_id,genre) values($1,$2)",uuid,newuser.Genre[i])
		checkErr(err)
	}
	for i:=0;i<len(newuser.Language);i++{
		_,err = Db.Exec("insert into userlanguagepreferance(user_id,language) values($1,$2)",uuid,newuser.Language[i])
		checkErr(err)
	}
	json.NewEncoder(w).Encode(newuser)
}