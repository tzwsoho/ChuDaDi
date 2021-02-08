package model

// Flushes ♠Spades/♥Hearts/♣Clubs/♦Diamonds
type Flushes int

const (
	// FlushesDiamonds ♦ Diamonds
	FlushesDiamonds Flushes = iota
	// FlushesClubs ♣ Clubs
	FlushesClubs
	// FlushesHearts ♥ Hearts
	FlushesHearts
	// FlushesSpades ♠ Spades
	FlushesSpades
)

func (f Flushes) String() string {
	switch f {
	case FlushesDiamonds:
		return "♦"

	case FlushesClubs:
		return "♣"

	case FlushesHearts:
		return "♥"

	case FlushesSpades:
		return "♠"
	}

	return ""
}
