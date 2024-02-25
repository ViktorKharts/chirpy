package database

import "os"

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	
	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
		Password: password,
	}

	for _, u := range dbStructure.Users {
		if u.Email == user.Email {
			return User{}, os.ErrExist
		}
	}
	
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	} 

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	users := []User{}

	dbStructure, err := db.loadDB()
	if err != nil {
		return users, err
	}

	for _, u := range dbStructure.Users {
		users = append(users, u)
	}

	return users, nil
} 

func (db *DB) GetUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{}
	for _, u := range dbStructure.Users {
		if u.Email == email {
			user = u
			break
		}
	}

	user, ok := dbStructure.Users[user.ID]	
	if !ok {
		return User{}, os.ErrNotExist 
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]	
	if !ok {
		return User{}, os.ErrNotExist 
	}

	user.Email = email
	user.Password = password
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	} 

	return user, nil
}

