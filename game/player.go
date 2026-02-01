package game

// Wall pattern - standard Azul wall layout
// Each row has all 5 colors in a shifted pattern
var WallPattern = [5][5]TileColor{
	{Blue, Yellow, Red, Black, White},
	{White, Blue, Yellow, Red, Black},
	{Black, White, Blue, Yellow, Red},
	{Red, Black, White, Blue, Yellow},
	{Yellow, Red, Black, White, Blue},
}

// PatternLine represents one of the 5 pattern lines (staging area)
type PatternLine struct {
	Size   int       // 1-5 tiles capacity
	Color  TileColor // Color of tiles in this line (NoTile if empty)
	Filled int       // Number of tiles currently placed
}

// IsFull returns true if the pattern line is completely filled
func (p *PatternLine) IsFull() bool {
	return p.Filled == p.Size
}

// IsEmpty returns true if the pattern line has no tiles
func (p *PatternLine) IsEmpty() bool {
	return p.Filled == 0
}

// CanAccept returns true if this line can accept tiles of the given color
func (p *PatternLine) CanAccept(color TileColor) bool {
	if p.IsFull() {
		return false
	}
	if p.IsEmpty() {
		return true
	}
	return p.Color == color
}

// Add places tiles into the pattern line, returns overflow count
func (p *PatternLine) Add(color TileColor, count int) int {
	if p.IsEmpty() {
		p.Color = color
	}

	space := p.Size - p.Filled
	toPlace := min(count, space)
	p.Filled += toPlace

	return count - toPlace // overflow
}

// Clear empties the pattern line, returns the tiles that were there
func (p *PatternLine) Clear() (TileColor, int) {
	color := p.Color
	count := p.Filled
	p.Color = NoTile
	p.Filled = 0
	return color, count
}

// Clone creates a copy
func (p *PatternLine) Clone() *PatternLine {
	return &PatternLine{
		Size:   p.Size,
		Color:  p.Color,
		Filled: p.Filled,
	}
}

// PlayerBoard represents a player's personal board
type PlayerBoard struct {
	PatternLines [5]*PatternLine // 5 pattern lines (rows 1-5)
	Wall         [5][5]bool      // Which wall positions are filled
	FloorLine    []TileColor     // Negative point tiles
	Score        int
}

// Floor line penalties
var FloorPenalties = []int{-1, -1, -2, -2, -2, -3, -3}

// NewPlayerBoard creates an empty player board
func NewPlayerBoard() *PlayerBoard {
	pb := &PlayerBoard{
		FloorLine: make([]TileColor, 0, 7),
		Score:     0,
	}

	for i := 0; i < 5; i++ {
		pb.PatternLines[i] = &PatternLine{
			Size:   i + 1,
			Color:  NoTile,
			Filled: 0,
		}
	}

	return pb
}

// CanPlaceOnLine checks if a color can be placed on a specific pattern line
// Must also check that the wall doesn't already have that color in that row
func (pb *PlayerBoard) CanPlaceOnLine(lineIdx int, color TileColor) bool {
	if lineIdx < 0 || lineIdx >= 5 {
		return false
	}

	// Check if pattern line can accept this color
	if !pb.PatternLines[lineIdx].CanAccept(color) {
		return false
	}

	// Check if wall already has this color in this row
	wallCol := pb.GetWallColumn(lineIdx, color)
	if pb.Wall[lineIdx][wallCol] {
		return false
	}

	return true
}

// GetWallColumn returns the column index where a color goes in a given row
func (pb *PlayerBoard) GetWallColumn(row int, color TileColor) int {
	for col := 0; col < 5; col++ {
		if WallPattern[row][col] == color {
			return col
		}
	}
	return -1 // Should never happen
}

// PlaceTiles adds tiles to a pattern line, overflow goes to floor
func (pb *PlayerBoard) PlaceTiles(lineIdx int, color TileColor, count int) {
	if lineIdx == -1 {
		// Direct to floor line
		for i := 0; i < count; i++ {
			pb.FloorLine = append(pb.FloorLine, color)
		}
		return
	}

	overflow := pb.PatternLines[lineIdx].Add(color, count)

	// Overflow goes to floor
	for i := 0; i < overflow; i++ {
		pb.FloorLine = append(pb.FloorLine, color)
	}
}

// AddToFloor adds a tile to the floor line
func (pb *PlayerBoard) AddToFloor(color TileColor) {
	pb.FloorLine = append(pb.FloorLine, color)
}

// ScoreWallTile calculates points for placing a tile on the wall
// Points = 1 + adjacent tiles in row + adjacent tiles in column
func (pb *PlayerBoard) ScoreWallTile(row, col int) int {
	points := 0

	// Count horizontal adjacent
	hCount := 1
	for c := col - 1; c >= 0 && pb.Wall[row][c]; c-- {
		hCount++
	}
	for c := col + 1; c < 5 && pb.Wall[row][c]; c++ {
		hCount++
	}

	// Count vertical adjacent
	vCount := 1
	for r := row - 1; r >= 0 && pb.Wall[r][col]; r-- {
		vCount++
	}
	for r := row + 1; r < 5 && pb.Wall[r][col]; r++ {
		vCount++
	}

	// If only the tile itself (no adjacents), score 1
	// If horizontal line, add hCount
	// If vertical line, add vCount
	if hCount == 1 && vCount == 1 {
		points = 1
	} else {
		if hCount > 1 {
			points += hCount
		}
		if vCount > 1 {
			points += vCount
		}
		// If isolated but we got here, something's wrong - default to 1
		if points == 0 {
			points = 1
		}
	}

	return points
}

// TileWall moves completed pattern lines to wall and scores
// Returns tiles to be discarded
func (pb *PlayerBoard) TileWall() []TileColor {
	discards := make([]TileColor, 0)

	for row := 0; row < 5; row++ {
		pl := pb.PatternLines[row]
		if !pl.IsFull() {
			continue
		}

		color, count := pl.Clear()
		col := pb.GetWallColumn(row, color)

		// Place tile on wall
		pb.Wall[row][col] = true

		// Score it
		pb.Score += pb.ScoreWallTile(row, col)

		// Remaining tiles (count - 1) go to discard
		for i := 0; i < count-1; i++ {
			discards = append(discards, color)
		}
	}

	return discards
}

// ScoreFloorLine applies floor penalties and clears floor
// Returns tiles to be discarded
func (pb *PlayerBoard) ScoreFloorLine() []TileColor {
	discards := make([]TileColor, 0)

	for i, tile := range pb.FloorLine {
		if i < len(FloorPenalties) {
			pb.Score += FloorPenalties[i]
		}
		// First player marker doesn't go to discard
		if tile != FirstPlayerMarker {
			discards = append(discards, tile)
		}
	}

	// Score can't go below 0
	if pb.Score < 0 {
		pb.Score = 0
	}

	pb.FloorLine = pb.FloorLine[:0]
	return discards
}

// ScoreEndGame adds bonus points at end of game
func (pb *PlayerBoard) ScoreEndGame() {
	// Complete horizontal lines: +2 each
	for row := 0; row < 5; row++ {
		complete := true
		for col := 0; col < 5; col++ {
			if !pb.Wall[row][col] {
				complete = false
				break
			}
		}
		if complete {
			pb.Score += 2
		}
	}

	// Complete vertical lines: +7 each
	for col := 0; col < 5; col++ {
		complete := true
		for row := 0; row < 5; row++ {
			if !pb.Wall[row][col] {
				complete = false
				break
			}
		}
		if complete {
			pb.Score += 7
		}
	}

	// All 5 of one color: +10 each
	for _, color := range AllColors() {
		count := 0
		for row := 0; row < 5; row++ {
			col := pb.GetWallColumn(row, color)
			if pb.Wall[row][col] {
				count++
			}
		}
		if count == 5 {
			pb.Score += 10
		}
	}
}

// HasCompletedRow returns true if any wall row is complete (game end trigger)
func (pb *PlayerBoard) HasCompletedRow() bool {
	for row := 0; row < 5; row++ {
		complete := true
		for col := 0; col < 5; col++ {
			if !pb.Wall[row][col] {
				complete = false
				break
			}
		}
		if complete {
			return true
		}
	}
	return false
}

// Clone creates a deep copy
func (pb *PlayerBoard) Clone() *PlayerBoard {
	newPB := &PlayerBoard{
		Wall:      pb.Wall, // Arrays are copied by value
		FloorLine: make([]TileColor, len(pb.FloorLine)),
		Score:     pb.Score,
	}

	for i := 0; i < 5; i++ {
		newPB.PatternLines[i] = pb.PatternLines[i].Clone()
	}

	copy(newPB.FloorLine, pb.FloorLine)
	return newPB
}

// GetValidPlacements returns all valid (lineIdx, color) pairs for placing tiles
// lineIdx of -1 means floor line
func (pb *PlayerBoard) GetValidPlacements(color TileColor) []int {
	valid := make([]int, 0)

	for i := 0; i < 5; i++ {
		if pb.CanPlaceOnLine(i, color) {
			valid = append(valid, i)
		}
	}

	// Can always place on floor
	valid = append(valid, -1)

	return valid
}
