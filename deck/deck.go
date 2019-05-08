package deck

import "strconv"

//go:generate stringer -type=Card

//Suits is a global string slice of suit names
var Suits = [4]string{"Spades", "Diamonds", "Clubs", "Hearts"}

//Card represents a single card in a deck with the label "A,J,Q,K..." and the underlying uint value (1,10,11,12...)
type Card struct {
	Suit  string
	Label string
	Value int
}

//Deck holds a slice of cards and creational methods according to the builder design pattern
type Deck struct {
	Cards []Card
}

//New creates and returns a standard deck (slice) of cards with all of the default settings
func New() Deck {
	var d Deck
	return d
}

func getLabelByValue(i int) string {
	switch i {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return strconv.Itoa(i)
	}
}

func (d Deck) Len() int           { return len(d.Cards) }
func (d Deck) Swap(i, j int)      { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] }
func (d Deck) Less(i, j int) bool { return d.Cards[i].Value < d.Cards[j].Value }

//Populate takes the number of decks we want to generate and creates the cards for each suite.
func (d Deck) Populate(numberOfDecks int) Deck {
	//initialize to empty slice
	d.Cards = []Card{}
	for _, suit := range Suits {
		cardsPerSuit := 13 * numberOfDecks
		for i := 1; i <= cardsPerSuit; i++ {
			value := i % 13
			if value == 0 {
				value = 13
			}
			var c Card
			c.Value = value
			c.Suit = suit
			c.Label = getLabelByValue(value)
			d.Cards = append(d.Cards, c)
		}
	}
	return d
}
