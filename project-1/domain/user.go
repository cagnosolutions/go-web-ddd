package domain

type User struct {
	Id        int    `json:"id",sql:"id"`
	FirstName string `json:"id",sql:"first_name"`
	LastName  string `json:"id",sql:"last_name"`
	Email     string `json:"id",sql:"email"`
}

var CreateUserTable = `CREATE TABLE user (
		id INTEGER NOT NULL,		
		first_name TEXT NOT NULL DEFAULT '',
		last_name TEXT NOT NULL DEFAULT '',
		email TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (id)
	  );`
