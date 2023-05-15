package controllers

import (
	"encoding/json"
	"fmt"
	"golangmovietask/models"
	"net/http"

	"github.com/google/uuid"
)

//user creation
// *controllers is refered from the movieController file where we have the struct and we use the same
func (u *Controllers) CreateNewUser(w http.ResponseWriter, r *http.Request){
	var newUser models.User
	newUser.User_id = uuid.New().String()
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil{
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	u.Service.CreateNewUserService(w,newUser)
}

func (u *Controllers) UpdateUserProfile(w http.ResponseWriter, r *http.Request){
	var existingUserUpdateData models.User
	//after login only we'll be able to update, considering that perspective we alredy have the user data which we return when the user logs in which contains the user_id too.
	err := json.NewDecoder(r.Body).Decode(&existingUserUpdateData)
	if err != nil{
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
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