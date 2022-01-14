package main

import (
	"database/sql"
	"github.com/cagnosolutions/go-web-ddd/project-1/domain"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // I know it's weird, but this is what you do with db drivers
)

var sqliteFile = "project-1/resources/sqlite3/db/project-1.sqlite"

func main() {

	// open the database
	db, err := sql.Open("sqlite3", sqliteFile)
	if err != nil {
		log.Fatalf("sql.open=%s", err)
	}
	// make sure we don't forget to close
	defer db.Close()

	// initialize a user sql repo
	userRepo := domain.NewUserSQLRepository(db)

	// create if not exists
	userRepo.CreateTable()

}

func createTable(db *sql.DB, createTableSyntax string) {
	log.Printf("preparing statement:\n%s\n", createTableSyntax)
	// prepare the statement
	statement, err := db.Prepare(createTableSyntax)
	if err != nil {
		log.Fatalf("create.prepare: %s", err)
	}
	// execute the statement
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("create.execute: %s", err)
	}
	log.Println("successfully executed statement!")
}

func createFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
}

func removeFile(path string) {
	os.Remove(path)
}
