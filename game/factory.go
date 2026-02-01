package game

// Factory represents a factory display (holds 4 tiles)
type Factory struct {
	Tiles []TileColor
}

// NewFactory creates an empty factory
func NewFactory() *Factory {
	return &Factory{
		Tiles: make([]TileColor, 0, 4),
	}
}

// Fill adds tiles to the factory
func (f *Factory) Fill(tiles []TileColor) {
	f.Tiles = append(f.Tiles, tiles...)
}

// IsEmpty returns true if factory has no tiles
func (f *Factory) IsEmpty() bool {
	return len(f.Tiles) == 0
}

// TakeColor removes all tiles of the given color and returns them
// Remaining tiles go to center
func (f *Factory) TakeColor(color TileColor) (taken []TileColor, remaining []TileColor) {
	taken = make([]TileColor, 0)
	remaining = make([]TileColor, 0)

	for _, t := range f.Tiles {
		if t == color {
			taken = append(taken, t)
		} else {
			remaining = append(remaining, t)
		}
	}

	f.Tiles = f.Tiles[:0] // Clear factory
	return taken, remaining
}

// HasColor returns true if factory contains the given color
func (f *Factory) HasColor(color TileColor) bool {
	for _, t := range f.Tiles {
		if t == color {
			return true
		}
	}
	return false
}

// GetColors returns unique colors present in this factory
func (f *Factory) GetColors() []TileColor {
	seen := make(map[TileColor]bool)
	colors := make([]TileColor, 0)

	for _, t := range f.Tiles {
		if !seen[t] {
			seen[t] = true
			colors = append(colors, t)
		}
	}
	return colors
}

// Clone creates a deep copy
func (f *Factory) Clone() *Factory {
	newF := &Factory{
		Tiles: make([]TileColor, len(f.Tiles)),
	}
	copy(newF.Tiles, f.Tiles)
	return newF
}

// Center represents the center of the table
type Center struct {
	Tiles              []TileColor
	HasFirstPlayerTile bool
}

// NewCenter creates an empty center with the first player marker
func NewCenter() *Center {
	return &Center{
		Tiles:              make([]TileColor, 0),
		HasFirstPlayerTile: true,
	}
}

// AddTiles adds tiles to the center (from factory overflow)
func (c *Center) AddTiles(tiles []TileColor) {
	c.Tiles = append(c.Tiles, tiles...)
}

// IsEmpty returns true if center has no tiles (ignoring first player marker)
func (c *Center) IsEmpty() bool {
	return len(c.Tiles) == 0
}

// TakeColor removes all tiles of the given color
// Also takes first player marker if present
func (c *Center) TakeColor(color TileColor) (taken []TileColor, tookFirstPlayer bool) {
	taken = make([]TileColor, 0)
	remaining := make([]TileColor, 0)

	for _, t := range c.Tiles {
		if t == color {
			taken = append(taken, t)
		} else {
			remaining = append(remaining, t)
		}
	}

	tookFirstPlayer = c.HasFirstPlayerTile
	c.HasFirstPlayerTile = false
	c.Tiles = remaining

	return taken, tookFirstPlayer
}

// HasColor returns true if center contains the given color
func (c *Center) HasColor(color TileColor) bool {
	for _, t := range c.Tiles {
		if t == color {
			return true
		}
	}
	return false
}

// GetColors returns unique colors present in center
func (c *Center) GetColors() []TileColor {
	seen := make(map[TileColor]bool)
	colors := make([]TileColor, 0)

	for _, t := range c.Tiles {
		if !seen[t] {
			seen[t] = true
			colors = append(colors, t)
		}
	}
	return colors
}

// Clone creates a deep copy
func (c *Center) Clone() *Center {
	newC := &Center{
		Tiles:              make([]TileColor, len(c.Tiles)),
		HasFirstPlayerTile: c.HasFirstPlayerTile,
	}
	copy(newC.Tiles, c.Tiles)
	return newC
}
