package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golangmovietask/models"
	// "golangmovietask/services"
	"net/http"

)

// *Services is refered from the movieService file where we have the struct and we use the same
func (u *Service) CreateNewUserService(w http.ResponseWriter,newUser models.User){
	//checking whether the user already exists
	err := u.DAO.CheckUserAlreadyExist(newUser.Email)
	if err != nil{
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w,"user already exist")
		return 
	}
	newUser.Password,err = HashPassword(newUser.Password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return 
	}
	err = u.DAO.AddNewUserToDb(newUser)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"error in updating Db")
		return 
	}
	w.WriteHeader(http.StatusCreated)
}

func (u *Service) UpdateUserProfileService(w http.ResponseWriter,existingUserUpdateData models.User){
	//As we don't have an OTP and E-Mail functionality for the password changeing purpose we dont allow the user to update the password nd email, if we needed we can add the Email change.
	password,err := u.DAO.GetUserPassword(existingUserUpdateData.User_id)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,err)
		return 
	}
	err = ComparePassword(password,existingUserUpdateData.Password)	//only if the password matches we can update the user details
	if err == nil{
		hashedPassword , err := HashPassword(existingUserUpdateData.Password)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w,err)
			return 
		}
		existingUserUpdateData.Password = hashedPassword
		err = u.DAO.UpdateExistingUser(existingUserUpdateData)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("error in updating Db")
			return 
		}
		existingUserUpdateData.Password = ""
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingUserUpdateData)
	}else{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w,"user password not matched")
		return 
	}
}

func (u *Service) UserLoginService(w http.ResponseWriter,userLoginCredentials models.UserLogin){
	var verifyUserLoginCredentials models.User		//this variable is for getting the data from the DB and doing the calculations
	//If user exist we get the user data
	verifyUserLoginCredentials,err := u.DAO.GetUserData(userLoginCredentials.Useremail,verifyUserLoginCredentials)
	if err == sql.ErrNoRows{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w,"No user found")
		return
	}
	err = ComparePassword(verifyUserLoginCredentials.Password,userLoginCredentials.Password)
	if err == nil {
		verifyUserLoginCredentials.Password=""
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(verifyUserLoginCredentials)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w,"password not matched")
		return
	}
}