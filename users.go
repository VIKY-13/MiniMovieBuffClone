package main

import (
	"encoding/json"
	// "errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

//user creation
func CreateNewUser(w http.ResponseWriter, r *http.Request){
	var newUser user
	newUser.User_id = uuid.New().String()
	json.NewDecoder(r.Body).Decode(&newUser)
	//checking whether the user already exists
	query := "SELECT COUNT(email) FROM users WHERE email=$1;"
    var count int
    err := Db.QueryRow(query, newUser.Email).Scan(&count)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
    // If the count is 1, the user already exists
    if count > 0 {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w,"user already exist")
		return
    }
	hashedPassword, err := HashPassword(newUser.Password)
	CheckInternalServerError(w,err,"error while hashing")
	_,err = Db.Exec("insert into users(user_id,firstname,lastname,email,password,age,phone_no) values($1,$2,$3,$4,$5,$6,$7)",newUser.User_id,newUser.Firstname,newUser.Lastname,strings.TrimSpace(newUser.Email),string(hashedPassword),newUser.Age,newUser.Phone_no)
	CheckInternalServerError(w,err,"error while storing user into the db")
	w.WriteHeader(http.StatusCreated)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request){
	var existingUser_UserId string
	var existingUserUpdateData user
	json.NewDecoder(r.Body).Decode(&existingUserUpdateData)
	err := Db.QueryRow("select user_id from users where email=$1;",existingUserUpdateData.Email).Scan(&existingUser_UserId)
	CheckInternalServerError(w,err,"error in query")
	_,err = Db.Exec("update users set firstname=$1,lastname=$2,password=$3,age=$4,phone_no=$5 where user_id=$6",existingUserUpdateData.Firstname,existingUserUpdateData.Lastname,existingUserUpdateData.Password,existingUserUpdateData.Age,existingUserUpdateData.Phone_no,existingUser_UserId)
	CheckInternalServerError(w,err,"in updating")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingUserUpdateData)
}

func HashPassword(password string) (string, error) {
    // GenerateFromPassword returns a hash of the password using bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
    if err != nil {
        return "", err
    }
    // convert the hash to a string and return
    return string(hash), nil
}

//USER LOGIN
func UserLogin(w http.ResponseWriter, r *http.Request){
	var userLoginCredentials userlogin		//this variable is for getting the data from the user
	var verifyUserLoginCredentials user		//this variable is for getting the data from the DB and doing the calculations
	resp,err := ioutil.ReadAll(r.Body)
	checkErr(err)
	json.Unmarshal(resp,&userLoginCredentials)

	//Checking whether the user already exists or not

	query := "SELECT COUNT(email) FROM users WHERE email=$1;"
    var count int
    err = Db.QueryRow(query, userLoginCredentials.Useremail).Scan(&count)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("user check query")
		return
	}
    // If the count is not equal to 1 there is no such user 
    if count != 1 {
        w.WriteHeader(http.StatusNotFound)
		return
    }

	//If user exist we get the user data
	statement,err := Db.Prepare("select user_id,firstname,lastname,phone_no,age,email,password from users where email= $1;")
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("in retrieving")
		return
	}
	_ = statement.QueryRow(userLoginCredentials.Useremail).Scan(&verifyUserLoginCredentials.User_id,&verifyUserLoginCredentials.Firstname,&verifyUserLoginCredentials.Lastname,&verifyUserLoginCredentials.Phone_no,&verifyUserLoginCredentials.Age,&verifyUserLoginCredentials.Email,&verifyUserLoginCredentials.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(verifyUserLoginCredentials.Password), []byte(userLoginCredentials.Password)); err == nil {
		verifyUserLoginCredentials.Password=""
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	json.NewEncoder(w).Encode(verifyUserLoginCredentials)
}

func CheckInternalServerError(w http.ResponseWriter, err error, message string) {
    if err != nil {
        fmt.Printf("%s: %v", message, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

