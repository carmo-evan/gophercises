package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Sqlx struct{}

func getSqlxDb() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	return sqlx.Connect("postgres", psqlInfo)
}

func (s *Sqlx) Wipe() error {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		return err
	}

	tx := db.MustBegin()
	tx.MustExec(`TRUNCATE phone_numbers`)
	return tx.Commit()
}

func (s *Sqlx) Seed() error {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		return err
	}
	tx := db.MustBegin()

	sqlStatement := `
	INSERT INTO phone_numbers (number)
	VALUES ($1),($2),($3),($4),($5),($6),($7),($8)
	RETURNING number`

	tx.MustExec(sqlStatement, "1234567890", "123 456 7891", "(123) 456 7892", "(123) 456-7893", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892")
	return tx.Commit()
}

func (s *Sqlx) GetAll() (*[]PhoneNumber, error) {

	db, err := getSqlxDb()

	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(`SELECT id,number from phone_numbers`)

	if err != nil {
		return nil, err
	}

	var phoneNumbers []PhoneNumber
	for rows.Next() {
		var p PhoneNumber
		err := rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, p)
	}
	return &phoneNumbers, nil
}

func (s *Sqlx) Update(p *PhoneNumber) error {
	return nil
}
