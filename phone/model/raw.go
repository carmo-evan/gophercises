package model

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Raw struct{}

const (
	host   = "localhost"
	port   = 3000
	user   = "postgres"
	dbname = "gophercises"
)

func getRawDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func (r *Raw) Wipe() error {
	db := getRawDb()

	defer db.Close()

	_, err := db.Exec(`TRUNCATE phone_numbers`)

	return err

}

func (r *Raw) Seed() error {

	db := getRawDb()

	defer db.Close()

	sqlStatement := `
	INSERT INTO phone_numbers (number)
	VALUES ($1),($2),($3),($4),($5),($6),($7),($8)
	RETURNING number`

	_, err := db.Exec(sqlStatement, "1234567890", "123 456 7891", "(123) 456 7892", "(123) 456-7893", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892")

	return err
}

func (r *Raw) GetAll() (*[]PhoneNumber, error) {

	db := getRawDb()

	defer db.Close()

	rows, err := db.Query(`SELECT ID,number from phone_numbers`)

	if err != nil {
		return nil, err
	}

	var phoneNumbers []PhoneNumber

	for rows.Next() {
		var number string
		var id uint
		err := rows.Scan(&id, &number)
		if err != nil {
			return nil, err
		}

		phoneNumbers = append(phoneNumbers, PhoneNumber{Number: number, ID: id})
	}

	return &phoneNumbers, nil
}

func (r *Raw) Update(p *PhoneNumber) error {
	return nil
}
