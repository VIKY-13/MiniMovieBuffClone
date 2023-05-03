package controllers

import (
	"encoding/json"
	"fmt"
	"golangmovietask/models"
	"net/http"

	"github.com/google/uuid"
)

//user creation
func (u *Controllers) CreateNewUser(w http.ResponseWriter, r *http.Request){
	var newUser models.User
	newUser.User_id = uuid.New().String()
	json.NewDecoder(r.Body).Decode(&newUser)
	u.Service.CreateNewUserService(w,newUser)
}

func (u *Controllers) UpdateUserProfile(w http.ResponseWriter, r *http.Request){
	var existingUserUpdateData models.User
	//after login only we'll be able to update, considering that perspective we alredy have the user data which we return when the user logs in which contains the user_id too.
	json.NewDecoder(r.Body).Decode(&existingUserUpdateData)
	u.Service.UpdateUserProfileService(w,existingUserUpdateData)
}



//USER LOGIN
func (u *Controllers) UserLogin(w http.ResponseWriter, r *http.Request){
	var userLoginCredentials models.UserLogin		//this variable is for getting the data from the user
	err := json.NewDecoder(r.Body).Decode(&userLoginCredentials)
	if err != nil{
		fmt.Println("error in mapping data")
		return
	}
    u.Service.UserLoginService(w,userLoginCredentials)
}