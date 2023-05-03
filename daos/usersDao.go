package daos

import (
	"errors"
	"golangmovietask/models"
	"strings"

	_ "github.com/lib/pq"
)

func (u *DAO) AddNewUserToDb(newUser models.User)error{
	statement,err := u.Db.Prepare("INSERT INTO users(user_id,firstname,lastname,email,password,age,phone_no) VALUES ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil{
		return err
	}
	_,err = statement.Exec(newUser.User_id,newUser.Firstname,newUser.Lastname,strings.TrimSpace(newUser.Email),string(newUser.Password),newUser.Age,newUser.Phone_no)
	if err != nil{
		return err
	}
	return nil
}

func (u *DAO) CheckUserAlreadyExist(newUserEmail string)error{
	query := "SELECT COUNT(email) FROM users WHERE email=$1;"
    var count int
    err := u.Db.QueryRow(query, newUserEmail).Scan(&count)
	if err != nil{
		return err
	}
	//if the count is greater than 0 then the user already exists
    if count > 0 {
		return errors.New("user already exist")
    }
	return nil
}

func (u *DAO) UpdateExistingUser(existingUserUpdateData models.User)error{
	statement,err := u.Db.Prepare("UPDATE users SET firstname=$1,lastname=$2,age=$3,phone_no=$4 WHERE user_id=$5")
	if err != nil{
		return nil
	}
	_,err = statement.Exec(existingUserUpdateData.Firstname,existingUserUpdateData.Lastname,existingUserUpdateData.Age,existingUserUpdateData.Phone_no,existingUserUpdateData.User_id)
	if err != nil{
		return err
	}
	return err
}

func (u *DAO) GetUserPassword(user_id string)(string,error){
	var password string
	statement,err := u.Db.Prepare("SELECT password FROM users WHERE user_id = $1")
	if err != nil{
		return "",err
	}
	err = statement.QueryRow(user_id).Scan(&password)
	if err != nil{
		return "",err
	}
	return password,nil
}


func (u *DAO) GetUserData(userLoginEmail string,verifyUserLoginCredentials models.User)(models.User,error){
	statement,err := u.Db.Prepare("SELECT user_id,firstname,lastname,phone_no,age,email,password FROM users WHERE email= $1;")
	if err != nil{
		return verifyUserLoginCredentials,err
	}
	err = statement.QueryRow(userLoginEmail).Scan(&verifyUserLoginCredentials.User_id,&verifyUserLoginCredentials.Firstname,&verifyUserLoginCredentials.Lastname,&verifyUserLoginCredentials.Phone_no,&verifyUserLoginCredentials.Age,&verifyUserLoginCredentials.Email,&verifyUserLoginCredentials.Password)
	if err != nil{
		return verifyUserLoginCredentials,err
	}
	return verifyUserLoginCredentials,nil
}