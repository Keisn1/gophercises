package phone

import (
	"bufio"
	"log"
	"os"
	"phone/db"
	"regexp"
)

func GetNumbers(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var numbers []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		numbers = append(numbers, scanner.Text())
	}

	return numbers
}

func NormalizeNumber(number string) string {
	r := regexp.MustCompile("[^0-9]")
	return r.ReplaceAllString(number, "")
}

func Update(newNbr, oldNbr string, nbrId int) error {
	if db.CheckNbrAlreadyExists(nbrId, newNbr) {
		db.DeleteId(nbrId)
	}

	db.UpdateNbr(nbrId, newNbr)
	return nil
}
