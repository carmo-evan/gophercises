package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	//postgres driver
	_ "github.com/lib/pq"
)

type Gorm struct{}

func getGormDb() (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	return gorm.Open("postgres", psqlInfo)
}

func (g *Gorm) Update(p *PhoneNumber) error {

	db, err := getGormDb()

	if err != nil {
		return err
	}

	defer db.Close()

	db.Save(p)

	return nil

}

func (g *Gorm) Get(ID uint) (*PhoneNumber, error) {

	db, err := getGormDb()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	p := &PhoneNumber{}

	db.Where("ID = ?", ID).First(&p)

	return p, nil

}

func (g *Gorm) GetAll() (*[]PhoneNumber, error) {

	db, err := getGormDb()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	phoneNumbers := &[]PhoneNumber{}

	db.Find(phoneNumbers)

	return phoneNumbers, nil

}

func (g *Gorm) Seed() error {

	db, err := getGormDb()

	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&PhoneNumber{})

	inputs := []string{"1234567890", "123 456 7891", "(123) 456 7892", "(123) 456-7893", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892"}

	for _, number := range inputs {
		db.Create(&PhoneNumber{Number: number})
	}
	return nil
}
func (g *Gorm) Wipe() error {

	db, err := getGormDb()
	if err != nil {
		return err
	}
	defer db.Close()

	db.Unscoped().Delete(PhoneNumber{})

	return nil
}
