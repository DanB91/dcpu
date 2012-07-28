// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// A single card comes as an unsigned, 8-bit integer with
// the following bit layout:
//
//     bits 0-3: The card value (0-12).
//     bits 4-5: The suit (0-3).
//     bit    6: Card is in use or not.
//     bit    7: <reserved for future use>

// Known card suits
const Hearts   = 0;
const Diamonds = 1;
const Clubs    = 2;
const Spades   = 3;

// suitName returns a string representation of the given suit.
function suitName (s)
{
	var name = "";

	switch (s) {
	case Hearts:
		name = "&#x2665;"; // ♥
		break;
	case Diamonds:
		name = "&#x2666;"; // ♦
		break;
	case Clubs:
		name = "&#x2663;"; // ♣
		break;
	case Spades:
		name = "&#x2660;"; // ♠
		break;
	}

	return name;
}

// cardValue returns the value of the given card.
function cardValue (c)
{
	 return c & 15;
}

// cardSuit returns the suit of the given card.
function cardSuit (c)
{
	return (c >> 4) & 3;
}

// cardName returns a string representation of the given card.
function cardName (c)
{
	var num = (c&15) + 1;
	var name = suitName((c >> 4) & 3);

	if (num > 1 && num < 11) {
		return name + num;
	}

	switch (num) {
	case 1:
		name += "A";
		break;
	case 11:
		name += "J";
		break;
	case 12:
		name += "Q";
		break;
	case 13:
		name += "K";
		break;
	}

	return name;
}
