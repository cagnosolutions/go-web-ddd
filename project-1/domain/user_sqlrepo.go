package domain

import (
	"database/sql"
	"errors"
	"log"
)

type UserSQLRepository struct {
	db *sql.DB
}

func NewUserSQLRepository(db *sql.DB) *UserSQLRepository {
	return &UserSQLRepository{
		db: db,
	}
}

func (u *UserSQLRepository) Init() error {
	u.SQLCreateTable()
	return nil
}

func (u *UserSQLRepository) Insert(v interface{}) (int, error) {
	u.SQLInsert(v.(*User))
	return 0, nil
}

func (u *UserSQLRepository) Update(v interface{}) (int, error) {
	u.SQLUpdate(v.(*User))
	return 0, nil
}

func (u *UserSQLRepository) Delete(v interface{}) (int, error) {
	u.SQLDelete(v.(*User))
	return 0, nil
}

func (u *UserSQLRepository) FindAll(v interface{}) (int, error) {
	u.SQLFindAll(v.([]*User))
	return 0, nil
}

func (u *UserSQLRepository) FindOne(v interface{}, ident interface{}) (int, error) {
	x1, ok := ident.(int)
	if ok {
		u.SQLFindOneByID(v.(*User), x1)
		return 0, nil
	}
	x2, ok := ident.(string)
	if ok {
		u.SQLFindOneByEmail(v.(*User), x2)
		return 0, nil
	}
	return 0, errors.New("bad ident type")
}

func (u *UserSQLRepository) SQLShowTables() {
	// sql statement
	showTables := `SELECT tbl_name FROM sqlite_master WHERE type='table' AND tbl_name NOT LIKE 'sqlite_%';`
	// prepare to protect against sql-injection
	stmt, err := u.db.Prepare(showTables)
	if err != nil {
		log.Fatalf("preparing(%q): %s", showTables, err)
	}
	// execute
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("executing prepared statment: %s", err)
	}
}

func (u *UserSQLRepository) SQLDropTable() {
	// sql statement
	dropTable := `DROP TABLE IF EXISTS user;`
	// prepare to protect against sql-injection
	stmt, err := u.db.Prepare(dropTable)
	if err != nil {
		log.Fatalf("preparing(%q): %s", dropTable, err)
	}
	// execute
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("executing prepared statment: %s", err)
	}
}

func (u *UserSQLRepository) SQLCreateTable() {
	// sql statement
	createTable := `CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL DEFAULT -1,		
		first_name TEXT NOT NULL DEFAULT '',
		last_name TEXT NOT NULL DEFAULT '',
		email TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (id)
	  );`
	// prepare to protect against sql-injection
	stmt, err := u.db.Prepare(createTable)
	if err != nil {
		log.Fatalf("preparing(%q): %s", createTable, err)
	}
	// execute
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("executing prepared statment: %s", err)
	}
}

// Insert usage:
// user := User{FirstName:"John", LastName:"Doe", Email:"jdoe@example.com"}
// domain.UserSQLRepository.Insert(&user)
func (u *UserSQLRepository) SQLInsert(user *User) {
	// sql statement
	insertUser := `INSERT INTO user(first_name, last_name, email) VALUES (?,?,?);`
	// prepare to protect against sql-injection
	stmt, err := u.db.Prepare(insertUser)
	if err != nil {
		log.Fatalf("preparing(%q): %s", insertUser, err)
	}
	// execute with params
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email)
	if err != nil {
		log.Fatalf("executing prepared statment: %s", err)
	}
}

var update = `UPDATE table
SET column_1 = new_value_1,
    column_2 = new_value_2
WHERE
    search_condition 
ORDER column_or_expression
LIMIT row_count OFFSET offset;`

// Update usage:
// user := User{Id: 3, LastName:"Smith"}
// domain.UserSQLRepository.Update(&user)
func (u *UserSQLRepository) SQLUpdate(user *User) {
	// sql statement
	updateUser := `UPDATE user SET first_name=?, last_name=?, email=? WHERE id=?;`
	// prepare to protect against sql-injection
	stmt, err := u.db.Prepare(updateUser)
	if err != nil {
		log.Fatalf("preparing(%q): %s", updateUser, err)
	}
	// execute with params
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		log.Fatalf("executing prepared statment: %s", err)
	}
}

// Delete usage:
// user := User{Id: 3}
// domain.UserSQLRepository.Delete(&user)
func (u *UserSQLRepository) SQLDelete(user *User) {
	// sql statement
	deleteUser := `DELETE FROM user WHERE id=?;`
	// prepare to protect against sql-injection
	stmt, err := u.db.Prepare(deleteUser)
	if err != nil {
		log.Fatalf("preparing(%q): %s", deleteUser, err)
	}
	// execute with params
	_, err = stmt.Exec(user.Id)
	if err != nil {
		log.Fatalf("executing prepared statment: %s", err)
	}
}

// FindAll usage:
// var users []*User
// domain.UserSQLRepository.FindAll(&users)
func (u *UserSQLRepository) SQLFindAll(users []*User) {
	// sql statement
	findAllUsers := `SELECT * FROM user ORDER BY id;`
	// execute select query
	rows, err := u.db.Query(findAllUsers)
	if err != nil {
		log.Fatalf("query(%q): %s", findAllUsers, err)
	}
	defer rows.Close()
	// load data into users ptr
	for rows.Next() {
		// instantiate new user
		var user *User
		// scan row result into new user instance
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			log.Fatalf("results: %s", err)
		}
		// append user instance to user supplied set
		users = append(users, user)
	}
}

// FindOneByID usage:
// var user User
// domain.UserSQLRepository.FindOneByID(&user, 4)
func (u *UserSQLRepository) SQLFindOneByID(user *User, id int) {
	// sql statement
	findUser := `SELECT * FROM user WHERE id=? LIMIT 1;`
	// execute select query
	row := u.db.QueryRow(findUser, id)
	// scan row result into user ptr instance provided
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		log.Fatalf("results: %s", err)
	}
}

// FindOneByEmail usage:
// var user User
// domain.UserSQLRepository.FindOneByEmail(&user, "jdoe@example.com")
func (u *UserSQLRepository) SQLFindOneByEmail(user *User, email string) {
	// sql statement
	findUser := `SELECT * FROM user WHERE email_address=? LIMIT 1;`
	// execute select query
	row := u.db.QueryRow(findUser, email)
	// scan row result into user ptr instance provided
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		log.Fatalf("results: %s", err)
	}
}
