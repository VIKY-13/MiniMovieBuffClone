package daos

import (
	"errors"
	"golangmovietask/dtos"
	"strings"

	_ "github.com/lib/pq"
)

func AddNewUserToDb(newUser dtos.User)error{
	statement,err := Db.Prepare("INSERT INTO users(user_id,firstname,lastname,email,password,age,phone_no) VALUES ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(newUser.User_id,newUser.Firstname,newUser.Lastname,strings.TrimSpace(newUser.Email),string(newUser.Password),newUser.Age,newUser.Phone_no)
	if err != nil{
		return err
	}
	return nil
}

func CheckUserAlreadyExist(newUserEmail string)error{
	query := "SELECT COUNT(email) FROM users WHERE email=$1;"
    var count int
    err := Db.QueryRow(query, newUserEmail).Scan(&count)
	if err != nil{
		return err
	}
	//if the count is greater than 0 then the user already exists
    if count > 0 {
		return errors.New("user already exist")
    }
	return nil
}

func UpdateExistingUser(existingUserUpdateData dtos.User)error{
	statement,err := Db.Prepare("UPDATE users SET firstname=$1,lastname=$2,password=$3,age=$4,phone_no=$5 WHERE user_id=$6")
	if err != nil{
		return nil
	}
	_,err = statement.Exec(existingUserUpdateData.Firstname,existingUserUpdateData.Lastname,existingUserUpdateData.Password,existingUserUpdateData.Age,existingUserUpdateData.Phone_no,existingUserUpdateData.User_id)
	if err != nil{
		return err
	}
	return err
}

func GetUserData(userLoginEmail string,verifyUserLoginCredentials dtos.User)(dtos.User,error){
	statement,err := Db.Prepare("SELECT user_id,firstname,lastname,phone_no,age,email,password FROM users WHERE email= $1;")
	if err != nil{
		return verifyUserLoginCredentials,err
	}
	err = statement.QueryRow(userLoginEmail).Scan(&verifyUserLoginCredentials.User_id,&verifyUserLoginCredentials.Firstname,&verifyUserLoginCredentials.Lastname,&verifyUserLoginCredentials.Phone_no,&verifyUserLoginCredentials.Age,&verifyUserLoginCredentials.Email,&verifyUserLoginCredentials.Password)
	if err != nil{
		return verifyUserLoginCredentials,err
	}
	return verifyUserLoginCredentials,nil
}