# Battle Minimum Working Example

A basic battle with 2 player characters and 2 generic enemies.

## Features

- [ ] There are 2 pawns representing player characters
- [ ] There are 2 pawns representing enemy AI characters
- [x] There is a data model which tracks all pawns health totals
- [x] Pawns can have status effects
- [ ] Status effects resolve at the beginning and end of turn
- [ ] Pawn health totals are visible
- [x] There is a data model representing the draw deck
- [x] There is a data model representing the discard pile
- [x] There is a data model representing exhausted cards
- [x] There is a data model representing the player's hand
- [ ] The player can see which cards are in their hand
- [ ] The player can see which cards are in their discard pile
- [ ] The player can see which cards are in their exhausted pile
- [ ] At the beginning of the battle, the player draws their initial hand of 5
- [ ] There are defined turns, alternating between the player and AI
- [ ] The player can assign exactly 1 card to each of their pawns
- [ ] The player can see what card, if any, is assigned to a pawn
- [ ] If a card is assigned pawn, the player can unassign the card
- [ ] Certain cards have restrictions on which pawns they can be assigned to
- [ ] When the player ends their turn, assigned cards execute in the order they were played
- [ ] When the player ends their turn, they draw back up to their current hand size (initially 5)
- [ ] On the AI turn, the AI assigns actions to their pawns. They need not come from a deck
- [ ] When the AI has finished assigning actions, it ends its turn
- [ ] When the AI ends its turn, enemy assigned actions execute in the order they were played
- [ ] When a pawn reaches 0 health it is removed from battle
- [ ] When all pawns of either side are removed, the battle ends
