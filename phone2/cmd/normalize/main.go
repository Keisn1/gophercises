package main

import (
	"database/sql"
	"log"
	"phone"
	"phone/db"
)

func main() {
	if !db.CheckIfFileExists(db.DbPath) {
		db.CreateTable()
		numbers := phone.GetNumbers("numbers.txt")
		db.InsertNumbers(numbers)
	}

	dbConn, err := sql.Open("sqlite3", db.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("select numberId, phoneNumber from phoneNumbers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var nbrIds []int
	var nbrs []string
	for rows.Next() {
		var nbrId int
		var nbr string
		err = rows.Scan(&nbrId, &nbr)
		if err != nil {
			log.Fatal(err)
		}
		nbrIds = append(nbrIds, nbrId)
		nbrs = append(nbrs, nbr)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for idx, nbr := range nbrs {
		nbrId := nbrIds[idx]
		newNbr := phone.NormalizeNumber(nbr)
		phone.Update(newNbr, nbr, nbrId)
	}
}
