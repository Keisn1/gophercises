package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DbPath string = "instance/db.sqlite"

func CheckIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func UpdateNbr(nbrId int, newNbr string) error {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `update phoneNumbers
set phoneNumber = (?)
where
numberId = (?);
`

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newNbr, nbrId)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func DeleteId(id int) error {
	dbConn, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `delete from phoneNumbers where (numberId) =  (?);`
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func CheckNbrAlreadyExists(nbrId int, nbr string) bool {
	dbConn, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("select numberId from phoneNumbers where phoneNumber = (?)", nbr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp int
		rows.Scan(&tmp)
		if tmp != nbrId {
			return true
		}
	}

	return false
}

func InsertNumbers(numbers []string) error {
	if !CheckIfFileExists(DbPath) {
		log.Fatal("Db doesn't exist")
	}

	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `insert into phoneNumbers (phoneNumber) values (?)`
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, nbr := range numbers {
		_, err := stmt.Exec(nbr)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func CreateTable() error {
	if CheckIfFileExists(DbPath) {
		return nil
	}

	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
create table phoneNumbers (
  'numberId' integer primary key autoincrement,
  'phoneNumber' text not null
);
`
	db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}
