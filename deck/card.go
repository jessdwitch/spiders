package deck

var (
	// DummyCard : For when you need a placeholder card
	DummyCard = Card{
		// TODO: When we have card lookups, replace with const of Dummy card ID
		ID: 0,
		Name: "Dummy",
		Description: "This card does nothing!",
	}
)



type (
	// CardID : An identifier for pulling static card data. The preferred way to communicate card
	//	data between systems.
	CardID int

	// Card : Data for a card.
	Card struct {
		ID CardID
		Name string
		Description string
	}
)