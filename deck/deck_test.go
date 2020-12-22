package deck_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/jessdwitch/spiders/deck"

	"github.com/stretchr/testify/assert"
)

func makeTestCard() *deck.Card {
	uniquifier := fmt.Sprint(rand.Int())
	return &deck.Card{
		Name:        "Dummy_" + uniquifier,
		Description: uniquifier,
	}
}

func makeCardlist(count int) deck.Cardlist {
	result := make(deck.Cardlist, count)
	for i := 0; i < count; i++ {
		result[i] = makeTestCard()
	}
	return result
}

func checkShuffle(t *testing.T, initial, shuffled deck.Cardlist) {
	assert.ElementsMatch(t, initial, shuffled)
	assert.NotSame(t, initial, shuffled)
}

func TestShuffle(t *testing.T) {
	c := makeCardlist(10)
	initial := c
	c.Shuffle()
	checkShuffle(t, initial, c)
}

func TestInsert(t *testing.T) {
	t.Run("Out of bounds", func(t *testing.T) {
		c := makeCardlist(10)
		target := makeTestCard()
		err := deck.IndexOutOfBoundsError{}
		assert.EqualError(t, c.Insert(target, 12), err.Error())
	})
	t.Run("Add to end", func(t *testing.T) {
		c := makeCardlist(10)
		target := makeTestCard()
		assert.NoError(t, c.Insert(target, 10))
		assert.Len(t, c, 11)
		assert.Equal(t, target, c[10])
	})
	t.Run("Add to beginning", func(t *testing.T) {
		c := makeCardlist(10)
		target := makeTestCard()
		assert.NoError(t, c.Insert(target, 0))
		assert.Len(t, c, 11)
		assert.Equal(t, target, c[0])
	})
	t.Run("Add to middle", func(t *testing.T) {
		c := makeCardlist(10)
		target := makeTestCard()
		assert.NoError(t, c.Insert(target, 5))
		assert.Len(t, c, 11)
		assert.Equal(t, target, c[5])
	})
}

func TestPeak(t *testing.T) {
	c := makeCardlist(4)
	assert.Equal(t, c[:3], c.Peak(3))
	assert.Equal(t, c, c.Peak(6))
}

func TestNewDeck(t *testing.T) {
	c := makeCardlist(10)
	initial := c
	d := deck.NewDeck(c)
	checkShuffle(t, initial, d.DrawPile)
	assert.Len(t, d.DrawPile, 10)
	assert.Len(t, d.DiscardPile, 0)
	assert.Len(t, d.Hand, 0)
	assert.Len(t, d.ExhaustPile, 0)
	assert.Equal(t, d.Count, 10)
}

func TestDraw(t *testing.T) {
	t.Run("Regular draw", func(t *testing.T) {
		c := makeCardlist(10)
		d := deck.NewDeck(c)
		initial := d.DrawPile
		assert.NoError(t, d.DrawCards(3))
		assert.Len(t, d.Hand, 3)
		assert.ElementsMatch(t, initial[:3], d.Hand)
		assert.Len(t, d.DrawPile, 7)
		assert.ElementsMatch(t, initial[3:], d.DrawPile)
		assert.Equal(t, d.Count, 10)
	})
	t.Run("Overdraw", func(t *testing.T) {
		c := makeCardlist(10)
		d := deck.NewDeck(c)
		initial := d.DrawPile
		assert.NoError(t, d.DrawCards(15))
		assert.Len(t, d.Hand, 10)
		assert.ElementsMatch(t, initial, d.Hand)
		assert.Len(t, d.DrawPile, 0)
		assert.Equal(t, d.Count, 10)
	})
}

func TestAddCard(t *testing.T) {
	c := makeCardlist(10)
	t.Run("Add to Draw", func(t *testing.T) {
		d := deck.NewDeck(c)
		target := makeTestCard()
		d.AddCard(target, false)
		assert.Len(t, d.DrawPile, 11)
		assert.Len(t, d.DiscardPile, 0)
		assert.Equal(t, d.Count, 11)
		assert.Contains(t, d.DrawPile, target)
	})
	t.Run("Add to Discard", func(t *testing.T) {
		d := deck.NewDeck(c)
		target := makeTestCard()
		d.AddCard(target, true)
		assert.Len(t, d.DrawPile, 10)
		assert.Len(t, d.DiscardPile, 1)
		assert.Equal(t, d.Count, 11)
		assert.Equal(t, target, d.DiscardPile[0])
	})
}

func TestDiscard(t *testing.T) {
	c := makeCardlist(10)
	d := deck.NewDeck(c)
	d.DrawCards(5)
	err := deck.IndexOutOfBoundsError{}
	assert.EqualError(t, d.Discard(7), err.Error())
	assert.Len(t, d.Hand, 5)
	target := d.Hand[2]
	assert.NoError(t, d.Discard(2))
	assert.Len(t, d.Hand, 4)
	assert.Len(t, d.DiscardPile, 1)
	assert.Equal(t, target, d.DiscardPile[0])
}

func TestReset(t *testing.T) {
	t.Run("Empty Discard", func(t *testing.T) {
		c := makeCardlist(10)
		d := deck.NewDeck(c)
		d.DrawCards(3)
		d.ResetDraw()
		assert.Len(t, d.DrawPile, 7)
		assert.Len(t, d.Hand, 3)
		assert.Len(t, d.DiscardPile, 0)
	})
	t.Run("With discard", func(t *testing.T) {
		c := makeCardlist(10)
		d := deck.NewDeck(c)
		d.DrawCards(3)
		d.Discard(2)
		d.ResetDraw()
		assert.Len(t, d.DrawPile, 8)
		assert.Len(t, d.Hand, 2)
		assert.Len(t, d.DiscardPile, 0)
	})
}
