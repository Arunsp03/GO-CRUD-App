package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	UserName string
	Password string
}

type UserDB struct {
	db *sql.DB
}

func (userDB *UserDB) RegisterUser(userName string, password string) error {
	if userName == "" {
		return fmt.Errorf("Username is empty")
	}
	if password == "" {
		return fmt.Errorf("Password is empty")
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	query, err := userDB.db.Exec("insert into Users (userName,password) values(?,?)", userName, hashedPassword)
	if err != nil {
		return err
	}
	fmt.Println(query)
	return nil

}

func (userDB *UserDB) ValidateLogin(userName string, password string) (bool, error) {
	if userName == "" {
		return false, fmt.Errorf("Username is empty")
	}
	if password == "" {
		return false, fmt.Errorf("Password is empty")
	}
	query, err := userDB.db.Query("select password as Password from users where userName=?", userName)
	if err != nil {
		return false, fmt.Errorf(err.Error())
	}
	defer query.Close()
	var users []User
	for query.Next() {
		var user User
		if err := query.Scan(&user.Password); err != nil {
			return false, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return false, fmt.Errorf("No user with this username exists")
	}
	hashedPassword := users[0].Password
	isValidUser := CheckPasswordHash(password, hashedPassword)
	if isValidUser {
		return true, nil
	}
	return false, nil

}
