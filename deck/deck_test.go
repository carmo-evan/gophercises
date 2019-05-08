package deck_test

import (
	"gophercises/deck"
	"log"
	"testing"
)

func TestPopulate(t *testing.T) {
	numberOfDecks := 2
	expectedNumberOfCards := numberOfDecks * 52
	d := deck.New().Populate(numberOfDecks)
	log.Printf("%v", d)
	if d.Len() != numberOfDecks*52 {
		t.Errorf("Expected %v cards, got %v cards.", expectedNumberOfCards, d.Len())
	}
}
