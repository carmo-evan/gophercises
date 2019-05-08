package model

type PhoneNumber struct {
	Number string
	ID     uint `db:"id"`
}

type Driver interface {
	Wipe() error
	Seed() error
	GetAll() (*[]PhoneNumber, error)
	Update(*PhoneNumber) error
	Get(ID uint) (*PhoneNumber, error)
}

type Table struct {
	Driver Driver
}

func (t *Table) Wipe() error {
	return t.Driver.Wipe()
}
func (t *Table) Seed() error {
	return t.Driver.Seed()
}

func (t *Table) GetAll() (*[]PhoneNumber, error) {
	return t.Driver.GetAll()
}

func (t *Table) Update(p *PhoneNumber) error {
	return t.Driver.Update(p)
}

func (t *Table) Get(ID uint) (*PhoneNumber, error) {
	return t.Driver.Get(ID)
}
