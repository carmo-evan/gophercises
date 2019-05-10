package deck

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

//go:generate stringer -type=Card

//Suits is a global map of suit names and their default sort order.
var Suits = map[string]int{
	"Spades":   1,
	"Diamonds": 2,
	"Clubs":    3,
	"Hearts":   4,
	"Joker":    5,
}

//Card represents a single card in a deck with the label "A,J,Q,K..." ,suit, and the underlying value (1,10,11,12...).
type Card struct {
	Suit  string
	Label string
	Value int
}

//Deck holds a slice of cards and creational methods according to the builder design pattern.
type Deck struct {
	Cards []Card
}

//New creates and returns a standard deck (slice) of cards with all of the default settings.
func New(numberOfDecks int) *Deck {
	d := &Deck{}
	return d.populate(numberOfDecks)
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

//Len returns the length of the underlying Cards slice. Satisfies the sort interface.
func (d *Deck) Len() int { return len(d.Cards) }

//Swap swtiches out card with index i for card with index j. Satisfies the sort interface.
func (d *Deck) Swap(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] }

//Less returns true if i is less than j and false otherwise. Satisfies the sort interface.
func (d *Deck) Less(i, j int) bool {

	if d.Cards[i].Suit == d.Cards[j].Suit {
		return d.Cards[i].Value < d.Cards[j].Value
	}
	return Suits[d.Cards[i].Suit] < Suits[d.Cards[j].Suit]
}

//Shuffle returns a pointer to the shuffled deck. The random factor has the current unix time as source.
func (d *Deck) Shuffle() *Deck {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := d.Len(); n > 0; n-- {
		randIndex := r.Intn(n)
		lastIndex := n - 1
		d.Swap(lastIndex, randIndex)
	}
	return d
}

//Sort is a wrapper around sort.Sort, which uses the default Less, Swap and Len methods on deck.
func (d *Deck) Sort() *Deck {
	sort.Sort(d)
	return d
}

//CustomSort takes a function satisfying the Less interface of the sort package and applies that to the deck.
func (d *Deck) CustomSort(f func(i, j int) bool) {
	sort.SliceStable(d.Cards, f)
}

/*populate takes the number of decks we want to generate and creates the cards for each suite.
The deck will be in default sorting (Ascending, Spades, Diamonds, Clubs, Hearts).*/
func (d *Deck) populate(numberOfDecks int) *Deck {
	d.Cards = []Card{}
	for suit := range Suits {
		if suit == "Joker" {
			continue
		}
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

//AddJokers adds an arbitrary number of Jokers to the deck.
func (d *Deck) AddJokers(numberOfJokers int) *Deck {
	for i := 0; i < numberOfJokers; i++ {
		c := Card{
			Label: "Joker",
			Suit:  "Joker",
			Value: 0,
		}
		d.Cards = append(d.Cards, c)
	}
	return d
}

/*Remove filters out cards like the Card r. If r has a suit, both suit and value have to match.
If it only has a value, all cards with that value will be removed. */
func (d *Deck) Remove(r Card) *Deck {
	for i, c := range d.Cards {

		if r.Suit == "" {
			r.Suit = c.Suit
		}

		if c.Value == r.Value && c.Suit == r.Suit {
			d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)
		}

	}
	return d
}

/*Contains searches for the existance of a type of card r. If r has a suit, both suit and value have to match.
If it only has a value, any card with that value will return true. */
func (d *Deck) Contains(r Card) bool {
	for _, c := range d.Cards {

		if r.Suit == "" {
			r.Suit = c.Suit
		}

		if c.Value == r.Value && c.Suit == r.Suit {
			return true
		}

	}
	return false
}
