package game

// TileColor represents the 5 tile colors in Azul
type TileColor int

const (
	Blue TileColor = iota
	Yellow
	Red
	Black
	White
	FirstPlayerMarker // Special marker, not a real tile
	NoTile            // Empty slot
)

const NumColors = 5

func (t TileColor) String() string {
	switch t {
	case Blue:
		return "B"
	case Yellow:
		return "Y"
	case Red:
		return "R"
	case Black:
		return "K"
	case White:
		return "W"
	case FirstPlayerMarker:
		return "1"
	case NoTile:
		return "."
	default:
		return "?"
	}
}

func (t TileColor) FullName() string {
	switch t {
	case Blue:
		return "Blue"
	case Yellow:
		return "Yellow"
	case Red:
		return "Red"
	case Black:
		return "Black"
	case White:
		return "White"
	case FirstPlayerMarker:
		return "First Player"
	case NoTile:
		return "Empty"
	default:
		return "Unknown"
	}
}

// ColorFromString parses a color from user input
func ColorFromString(s string) (TileColor, bool) {
	switch s {
	case "b", "B", "blue", "Blue":
		return Blue, true
	case "y", "Y", "yellow", "Yellow":
		return Yellow, true
	case "r", "R", "red", "Red":
		return Red, true
	case "k", "K", "black", "Black":
		return Black, true
	case "w", "W", "white", "White":
		return White, true
	default:
		return NoTile, false
	}
}

// AllColors returns all 5 tile colors
func AllColors() []TileColor {
	return []TileColor{Blue, Yellow, Red, Black, White}
}
