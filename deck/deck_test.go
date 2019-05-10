package deck_test

import (
	"fmt"
	"gophercises/deck"
	"log"
	"sort"
	"testing"
)

var testCases = []struct {
	WithJoker      bool
	NumberOfJokers int
}{
	{true, 2},
	{false, 0},
}

func TestRemove(t *testing.T) {
	//TODO: table of test cases
	r := deck.Card{Value: 1, Suit: "Hearts"}
	d := deck.New(1).Remove(r)

	if d.Contains(r) {
		log.Printf("%v", d)
		t.Errorf(`Found card %v that should have been filtered out`, r)
	}
}

func TestNew(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf(`test New() With Joker %v`, tc.WithJoker), func(t *testing.T) {
			numberOfDecks := 2
			expectedNumberOfCards := numberOfDecks * (52 + tc.NumberOfJokers)
			d := deck.New(numberOfDecks)
			if tc.WithJoker {
				d.AddJokers(tc.NumberOfJokers * numberOfDecks)
			}
			if d.Len() != expectedNumberOfCards {
				log.Printf("%v", d)
				t.Errorf("Expected %v cards, got %v cards.", expectedNumberOfCards, d.Len())
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf(`test Shuffle() With Joker %v`, tc.WithJoker), func(t *testing.T) {
			d := deck.New(1)
			if tc.WithJoker {
				d.AddJokers(2)
			}
			d.Shuffle()
			if areCardsSequential(d, tc.WithJoker) {
				log.Printf("%v", d.Cards)
				t.Errorf("All cards were sequential.")
			}
		})
	}
}

func TestDefaultSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf(`test Sort() With Joker %v`, tc.WithJoker), func(t *testing.T) {

			d := deck.New(1).Shuffle()
			if tc.WithJoker {
				d.AddJokers(2)
			}
			d.Sort()
			if !areCardsSequential(d, tc.WithJoker) {
				log.Printf("%v", d.Cards)
				t.Errorf("All cards were not sequential.")
			}
		})
	}
}

func TestCustomSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf(`test CustomSort() - With Joker %v`, tc.WithJoker), func(t *testing.T) {
			d := deck.New(1)

			if tc.WithJoker {
				d.AddJokers(2)
			}

			d.Shuffle()

			d.CustomSort(func(i, j int) bool {
				return d.Cards[i].Value < d.Cards[j].Value
			})

			if !sort.SliceIsSorted(d.Cards, func(i, j int) bool {
				return d.Cards[i].Value < d.Cards[j].Value
			}) {
				log.Printf("%v", d.Cards)
				t.Errorf("All cards were not sorted by value.")
			}
		})
	}
}

func areCardsSequential(d *deck.Deck, withJoker bool) bool {
	var lastSeen deck.Card
	//let's loop all suits and all cards in the suit
	for suit := range deck.Suits {
		if !withJoker && suit == "Joker" {
			continue
		}
		for i := 1; i <= 13; i++ {
			c := d.Cards[i]
			// if current card's value is not one more than the last card's, it's shuffled
			if lastSeen.Value != 0 && c.Value != lastSeen.Value+1 && c.Suit == lastSeen.Suit {
				return false
			}
			lastSeen = c
		}
	}

	//if all cards were sequential, we never returned out of the loop. Fail.
	return true
}
