package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	"strings"
	"github.com/google/uuid"
	
	_ "github.com/lib/pq"
)


//user creation
func CreateNewUser(w http.ResponseWriter, r *http.Request){
	var newuser user
	newuser.User_id = uuid.New().String()
	json.NewDecoder(r.Body).Decode(&newuser)
	//checking whether the user already exists
	query := "SELECT COUNT(email) FROM users WHERE email=$1;"
    var count int
    err := Db.QueryRow(query, newuser.Email).Scan(&count)
    checkErr(err)

    // If the count is 1, the data exists
    if count >= 1 {
        w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w,"user already exist")
		return
    }
	Db.Exec("insert into users(user_id,firstname,lastname,email,password,age,phone_no) values($1,$2,$3,$4,$5,$6,$7)",newuser.User_id,newuser.Firstname,newuser.Lastname,strings.TrimSpace(newuser.Email),newuser.Password,newuser.Age,newuser.Phone_no)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newuser)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request){
	var existinguser user
	var existinguserupdate user
	json.NewDecoder(r.Body).Decode(&existinguserupdate)
	row,err := Db.Query("select * from users where email='"+existinguserupdate.Email+"';")
	checkErr(err)
	for row.Next(){
		err = row.Scan(&existinguser.User_id,&existinguser.Firstname,&existinguser.Lastname,&existinguser.Email,&existinguser.Password,&existinguser.Age,&existinguser.Phone_no)
		checkErr(err)
	}
	json.NewEncoder(w).Encode(existinguser)
	_,err = Db.Exec("update users set firstname=$1,lastname=$2,password=$3,age=$4,phone_no=$5 where user_id=$6",existinguserupdate.Firstname,existinguserupdate.Lastname,existinguserupdate.Password,existinguserupdate.Age,existinguserupdate.Phone_no,existinguser.User_id)
	checkErr(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existinguserupdate)
}