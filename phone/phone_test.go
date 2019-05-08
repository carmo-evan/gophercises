package phone

import (
	"fmt"
	"gophercises/phone/model"
	"gophercises/phone/util"
	"log"
	"testing"
)

func TestNormalize(t *testing.T) {
	normalized, err := util.Normalize("(123) 456-7890")
	expected := "1234567890"

	if err != nil {
		t.Error(err)
	}
	if normalized != expected {
		t.Fatalf("Expected %v, returned %v", expected, normalized)
	}

	fmt.Println(normalized)

}

func TestUpdate(t *testing.T) {
	drivers := &[]model.Driver{&model.Sqlx{}, &model.Gorm{}, &model.Raw{}}

	for _, driver := range *drivers {
		table := &model.Table{Driver: driver}
		p := &model.PhoneNumber{}
		p.ID = 185
		p.Number = "01010101001"
		table.Update(p)
		newP := table.Get(p.ID)
		if newP.Number != p.Number {
			t.Errorf(`Driver: %v - Expected number to be %v be got %v`, driver, p.Number, newP.Number)
		}
	}

}

func TestSqlxGetAll(t *testing.T) {

	table := &model.Table{Driver: &model.Sqlx{}}

	phoneNumbers, err := table.GetAll()

	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", phoneNumbers)
}

func TestRawGetAll(t *testing.T) {

	table := &model.Table{Driver: &model.Raw{}}

	phoneNumbers, err := table.GetAll()

	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", phoneNumbers)
}

func TestGormUpdate(t *testing.T) {

	table := &model.Table{Driver: &model.Gorm{}}

	p := &model.PhoneNumber{}
	p.ID = 185
	p.Number = "01010101001"

	if err := table.Update(p); err != nil {
		t.Error(err)
	}

}

func TestGormGetAll(t *testing.T) {

	table := &model.Table{Driver: &model.Gorm{}}

	phoneNumbers, err := table.GetAll()

	if err != nil {
		t.Error(err)
	}
	log.Print(phoneNumbers)
}

func TestGormSeed(t *testing.T) {
	table := &model.Table{Driver: &model.Gorm{}}
	if err := table.Seed(); err != nil {
		t.Error(err)
	}
}
func TestGormWipe(t *testing.T) {
	table := &model.Table{Driver: &model.Gorm{}}
	if err := table.Wipe(); err != nil {
		t.Error(err)
	}
}
func TestRawSeed(t *testing.T) {
	table := &model.Table{Driver: &model.Raw{}}
	if err := table.Seed(); err != nil {
		t.Error(err)
	}
}
func TestRawWipe(t *testing.T) {
	table := &model.Table{Driver: &model.Raw{}}
	if err := table.Wipe(); err != nil {
		t.Error(err)
	}
}

func TestSqlxSeed(t *testing.T) {
	table := &model.Table{Driver: &model.Sqlx{}}
	if err := table.Seed(); err != nil {
		t.Error(err)
	}
}
func TestSqlxWipe(t *testing.T) {
	table := &model.Table{Driver: &model.Sqlx{}}
	if err := table.Wipe(); err != nil {
		t.Error(err)
	}
}
