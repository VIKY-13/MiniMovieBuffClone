package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golangmovietask/daos"
	"golangmovietask/dtos"
	"golangmovietask/services"
	"net/http"

	"github.com/google/uuid"
)

//user creation
func CreateNewUser(w http.ResponseWriter, r *http.Request){
	var newUser dtos.User
	newUser.User_id = uuid.New().String()
	json.NewDecoder(r.Body).Decode(&newUser)
	//checking whether the user already exists
	err := daos.CheckUserAlreadyExist(newUser.Email)
	if err != nil{
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w,"user already exist")
		return 
	}
	newUser.Password,err = services.HashPassword(newUser.Password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return 
	}
	err = daos.AddNewUserToDb(newUser)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error in updating Db")
		return 
	}
	w.WriteHeader(http.StatusCreated)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request){
	var existingUserUpdateData dtos.User

	//after login only we'll be able to update, consideering that perspective we alredy have the user data which we return when the user logs in which contains the user_id too.
	json.NewDecoder(r.Body).Decode(&existingUserUpdateData)
	hashedPassword , err := services.HashPassword(existingUserUpdateData.Password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return 
	}
	existingUserUpdateData.Password = hashedPassword
	err = daos.UpdateExistingUser(existingUserUpdateData)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error in updating Db")
		return 
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingUserUpdateData)
}



//USER LOGIN
func UserLogin(w http.ResponseWriter, r *http.Request){
	var userLoginCredentials dtos.UserLogin		//this variable is for getting the data from the user
	var verifyUserLoginCredentials dtos.User		//this variable is for getting the data from the DB and doing the calculations
	err := json.NewDecoder(r.Body).Decode(&userLoginCredentials)
	if err != nil{
		fmt.Println("error in mapping data")
		return
	}
    
	//If user exist we get the user data
	verifyUserLoginCredentials,err = daos.GetUserData(userLoginCredentials.Useremail,verifyUserLoginCredentials)
	if err == sql.ErrNoRows{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w,"No user found")
		return
	}
	err = services.ComparePassword(verifyUserLoginCredentials.Password,userLoginCredentials.Password)
	if err == nil {
		verifyUserLoginCredentials.Password=""
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(verifyUserLoginCredentials)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
}






