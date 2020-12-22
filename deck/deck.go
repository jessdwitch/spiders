package deck

import (
	"math/rand"
	"strings"
)

type (
	// Cardlist : A collection of Cards
	Cardlist []*Card

	// Deck : A collection of Cards, arranged into draw, hand, and discard piles
	Deck struct {
		DrawPile    Cardlist
		DiscardPile Cardlist
		Hand        Cardlist
		ExhaustPile Cardlist
		// Count : Total cards in circulation (Draw + Discard + Hand)
		Count int
	}
	// IDEA: A version of Deck where card state (which pile it's in) is stored in an array of ints.
	//	It's sort of like a linked list. The indices correspond to the indices of the whole
	//	cardlist. The value stored at each index is the next card in the pile. We store the "head"
	//	of each pile.
	// 	Compare to an array of enums which correspond to pile. Want the draw pile? In the enum array
	//	that's always O(n of all cards), whereas with the linked one, it's O(n of pile). Drawing is
	//	3 calls: update the draw head, update the hand head, correct reference at the end of the
	//	drawn cards. Shuffles are much easier than an actual linked list: you pull the list of
	//	indices, shuffle, and put them back (the catch being, you have to check for self-linked
	//	nodes).

	// IndexOutOfBoundsError : Attempted to access an element outside the collection
	IndexOutOfBoundsError struct{}
)

// NewDeck : Create a new Deck from a list of Cards
func NewDeck(c Cardlist) *Deck {
	// Deep copy cards
	cards := make(Cardlist, len(c))
	for i := 0; i < len(c); i++ {
		card := *c[i]
		cards[i] = &card
	}
	cards.Shuffle()
	return &Deck{
		DrawPile:    cards,
		DiscardPile: Cardlist{},
		Hand:        Cardlist{},
		ExhaustPile: Cardlist{},
		Count:       len(cards),
	}
}

// NewDeckFromIDs : Create a new Decklist from their lookup IDs mapped to their quantity
func NewDeckFromIDs(ids map[CardID]int) (Deck, error) {
	panic("Not implemented")
}

// ResetDraw : Shuffle the Draw and Discard piles together
func (d *Deck) ResetDraw() {
	d.DrawPile = append(d.DrawPile, d.DiscardPile...)
	d.DrawPile.Shuffle()
	d.DiscardPile = d.DiscardPile[:0] // TODO: Revisit if a 0-slice is a good choice here
}

// DrawCards : DrawCards cards from the deck. Get as many as possible from the DrawCards. If it runs out shuffle
// and try again.
func (d *Deck) DrawCards(howMany int) error {
	// TODO: This function looks like shit lol
	if howMany > len(d.DrawPile)+len(d.DiscardPile) {
		howMany = len(d.DrawPile) + len(d.DiscardPile)
	}
	if howMany == 0 {
		return nil // TODO: Do we want an error here? Enough cards just don't exist
	}
	// Do we need as many or more than Draw can provide?
	if howMany >= len(d.DrawPile) {
		d.Hand = append(d.Hand, d.DrawPile...)
		howMany -= len(d.DrawPile)
		d.DrawPile = d.DrawPile[:0]
		d.ResetDraw()
		return d.DrawCards(howMany)
	}
	d.Hand = append(d.Hand, d.DrawPile[:howMany]...)
	d.DrawPile = d.DrawPile[howMany:]
	return nil
}

// Discard : Move a card from the Hand to the Discard
func (d *Deck) Discard(i int) error {
	if len(d.Hand) <= i {
		return &IndexOutOfBoundsError{}
	}
	d.DiscardPile.Insert(d.Hand[i], 0)
	d.Hand = append(d.Hand[:i], d.Hand[i+1:]...)
	return nil
}

// Peak : Get the top howMany cards without drawing them, less if unavailable.
func (c *Cardlist) Peak(howMany int) Cardlist {
	if howMany > len(*c) {
		return *c
	}
	return (*c)[:howMany]
}

// Insert : Add a card to the list
func (c *Cardlist) Insert(card *Card, i int) error {
	if i > len(*c) {
		return &IndexOutOfBoundsError{}
	}
	if i == 0 {
		result := append(Cardlist{card}, *c...)
		*c = result
		return nil
	}
	if i == len(*c) {
		result := append(*c, card)
		*c = result
		return nil
	}
	result := append((*c)[:i+1], (*c)[i:]...)
	result[i] = card
	*c = result
	return nil
}

// AddCard : Add a card to the deck. If toDiscard is true, it's added to the discard, else randomly
// in the draw.
func (d *Deck) AddCard(c *Card, toDiscard bool) error {
	if toDiscard {
		d.DiscardPile = append(d.DiscardPile, c)
	} else {
		if err := d.DrawPile.Insert(c, rand.Intn(len(d.DrawPile))); err != nil {
			return err
		}
	}
	d.Count++
	return nil
}

// RemoveCard : INCOMPLETE
// func (d *Deck) RemoveCard() error {

// }

// Shuffle : Shuffle these cards
func (c *Cardlist) Shuffle() {
	rand.Shuffle(len(*c), func(i, j int) { (*c)[i], (*c)[j] = (*c)[j], (*c)[i] })
}

func (c *Cardlist) String() string {
	result := []string{}
	for _, card := range *c {
		result = append(result, card.Name)
	}
	return strings.Join(result, "\t")
}

func (i *IndexOutOfBoundsError) Error() string {
	return "index out of bounds"
}
