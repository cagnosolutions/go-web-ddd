package db

type DataRepository interface {
	Init() error
	Insert(v interface{}) (int, error)
	Update(v interface{}) (int, error)
	Delete(v interface{}) (int, error)
	FindAll(v interface{}) (int, error)
	FindOne(v interface{}, ident interface{}) (int, error)
}

var dropTbl = `DROP TABLE IF EXISTS user;`
var createTbl = `CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL DEFAULT -1,		
		first_name TEXT NOT NULL DEFAULT '',
		last_name TEXT NOT NULL DEFAULT '',
		email TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (id)
	  );`

var updateStmt = `UPDATE user SET first_name=?, last_name=?, email=? WHERE search`
var insertStmt = `INSERT INTO user(first_name, last_name, email) VALUES (?,?,?);`
var selectStmt = `SELECT * FROM user ORDER BY id;`

var updateAlt = `UPDATE table
SET column_1 = new_value_1,
    column_2 = new_value_2
WHERE
    search_condition 
ORDER column_or_expression
LIMIT row_count OFFSET offset;`

var createTblAlt = `CREATE TABLE database_name.table_name(
   column1 datatype PRIMARY KEY (one or more columns),
   column2 datatype,
   column3 datatype,
   .....
   columnN datatype
);`

var delete = `DELETE FROM table_name WHERE [condition];`
