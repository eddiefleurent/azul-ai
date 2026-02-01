package game

import (
	"math/rand"
)

// Bag holds tiles to be drawn and discarded tiles
type Bag struct {
	tiles    []TileColor
	discards []TileColor
	rng      *rand.Rand
	seed     int64 // Original seed for cloning
}

// NewBag creates a bag with 20 tiles of each color (100 total)
func NewBag(seed int64) *Bag {
	b := &Bag{
		tiles:    make([]TileColor, 0, 100),
		discards: make([]TileColor, 0, 100),
		rng:      rand.New(rand.NewSource(seed)),
		seed:     seed,
	}

	// Add 20 of each color
	for _, color := range AllColors() {
		for i := 0; i < 20; i++ {
			b.tiles = append(b.tiles, color)
		}
	}

	b.Shuffle()
	return b
}

// Shuffle randomizes the tile order
func (b *Bag) Shuffle() {
	b.rng.Shuffle(len(b.tiles), func(i, j int) {
		b.tiles[i], b.tiles[j] = b.tiles[j], b.tiles[i]
	})
}

// Draw removes and returns n tiles from the bag
// If bag is empty, refills from discards first
func (b *Bag) Draw(n int) []TileColor {
	drawn := make([]TileColor, 0, n)

	for i := 0; i < n; i++ {
		if len(b.tiles) == 0 {
			b.RefillFromDiscards()
			if len(b.tiles) == 0 {
				break // No more tiles anywhere
			}
		}
		drawn = append(drawn, b.tiles[len(b.tiles)-1])
		b.tiles = b.tiles[:len(b.tiles)-1]
	}

	return drawn
}

// RefillFromDiscards moves all discards back to the bag and shuffles
func (b *Bag) RefillFromDiscards() {
	b.tiles = append(b.tiles, b.discards...)
	b.discards = b.discards[:0]
	b.Shuffle()
}

// Discard adds tiles to the discard pile
func (b *Bag) Discard(tiles []TileColor) {
	b.discards = append(b.discards, tiles...)
}

// TilesRemaining returns count of tiles in bag (not discards)
func (b *Bag) TilesRemaining() int {
	return len(b.tiles)
}

// TotalTilesInPlay returns tiles in bag + discards
func (b *Bag) TotalTilesInPlay() int {
	return len(b.tiles) + len(b.discards)
}

// Clone creates a deep copy of the bag (for AI simulation)
func (b *Bag) Clone() *Bag {
	newBag := &Bag{
		tiles:    make([]TileColor, len(b.tiles)),
		discards: make([]TileColor, len(b.discards)),
		rng:      rand.New(rand.NewSource(b.seed)),
		seed:     b.seed,
	}
	copy(newBag.tiles, b.tiles)
	copy(newBag.discards, b.discards)
	return newBag
}
