package game

import (
	"fmt"
	"time"
)

// Game represents the full game state
type Game struct {
	Players       []*PlayerBoard
	Factories     []*Factory
	Center        *Center
	Bag           *Bag
	CurrentPlayer int
	FirstPlayer   int // Who goes first next round
	Round         int
	GameOver      bool
	NumPlayers    int
}

// NewGame creates a new game with the specified number of players
func NewGame(numPlayers int) *Game {
	if numPlayers < 2 {
		numPlayers = 2
	}
	if numPlayers > 4 {
		numPlayers = 4
	}

	// Number of factories based on player count
	numFactories := numPlayers*2 + 1

	g := &Game{
		Players:       make([]*PlayerBoard, numPlayers),
		Factories:     make([]*Factory, numFactories),
		Center:        NewCenter(),
		Bag:           NewBag(time.Now().UnixNano()),
		CurrentPlayer: 0,
		FirstPlayer:   0,
		Round:         1,
		GameOver:      false,
		NumPlayers:    numPlayers,
	}

	for i := 0; i < numPlayers; i++ {
		g.Players[i] = NewPlayerBoard()
	}

	for i := 0; i < numFactories; i++ {
		g.Factories[i] = NewFactory()
	}

	g.SetupRound()
	return g
}

// NewGameWithSeed creates a game with a specific random seed (for reproducibility)
func NewGameWithSeed(numPlayers int, seed int64) *Game {
	g := NewGame(numPlayers)
	g.Bag = NewBag(seed)
	g.SetupRound()
	return g
}

// SetupRound prepares factories for a new round
func (g *Game) SetupRound() {
	// Reset center
	g.Center = NewCenter()

	// Fill each factory with 4 tiles
	for _, f := range g.Factories {
		f.Tiles = f.Tiles[:0]
		tiles := g.Bag.Draw(4)
		f.Fill(tiles)
	}
}

// IsRoundOver returns true if all factories and center are empty
func (g *Game) IsRoundOver() bool {
	for _, f := range g.Factories {
		if !f.IsEmpty() {
			return false
		}
	}
	return g.Center.IsEmpty()
}

// Move represents a player action
type Move struct {
	FactoryIdx int       // -1 for center
	Color      TileColor // Which color to take
	LineIdx    int       // Which pattern line to place on (-1 for floor)
}

func (m Move) String() string {
	source := "center"
	if m.FactoryIdx >= 0 {
		source = fmt.Sprintf("factory %d", m.FactoryIdx)
	}

	dest := "floor"
	if m.LineIdx >= 0 {
		dest = fmt.Sprintf("line %d", m.LineIdx+1)
	}

	return fmt.Sprintf("Take %s from %s, place on %s", m.Color.FullName(), source, dest)
}

// GetValidMoves returns all legal moves for the current player
func (g *Game) GetValidMoves() []Move {
	moves := make([]Move, 0)
	player := g.Players[g.CurrentPlayer]

	// Moves from factories
	for fIdx, factory := range g.Factories {
		for _, color := range factory.GetColors() {
			for _, lineIdx := range player.GetValidPlacements(color) {
				moves = append(moves, Move{
					FactoryIdx: fIdx,
					Color:      color,
					LineIdx:    lineIdx,
				})
			}
		}
	}

	// Moves from center
	for _, color := range g.Center.GetColors() {
		for _, lineIdx := range player.GetValidPlacements(color) {
			moves = append(moves, Move{
				FactoryIdx: -1,
				Color:      color,
				LineIdx:    lineIdx,
			})
		}
	}

	return moves
}

// ApplyMove executes a move and updates game state
func (g *Game) ApplyMove(move Move) error {
	player := g.Players[g.CurrentPlayer]
	var tiles []TileColor
	var tookFirstPlayer bool

	if move.FactoryIdx == -1 {
		// Taking from center
		tiles, tookFirstPlayer = g.Center.TakeColor(move.Color)
		if tookFirstPlayer {
			player.AddToFloor(FirstPlayerMarker)
			g.FirstPlayer = g.CurrentPlayer
		}
	} else {
		// Taking from factory
		var remaining []TileColor
		tiles, remaining = g.Factories[move.FactoryIdx].TakeColor(move.Color)
		g.Center.AddTiles(remaining)
	}

	if len(tiles) == 0 {
		return fmt.Errorf("no tiles of color %s at source", move.Color.FullName())
	}

	// Place tiles
	player.PlaceTiles(move.LineIdx, move.Color, len(tiles))

	// Check if round is over
	if g.IsRoundOver() {
		g.EndRound()
	} else {
		g.NextPlayer()
	}

	return nil
}

// NextPlayer advances to the next player
func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % g.NumPlayers
}

// EndRound handles end-of-round scoring and setup
func (g *Game) EndRound() {
	// Score each player
	for _, player := range g.Players {
		// Move tiles to wall
		discards := player.TileWall()
		g.Bag.Discard(discards)

		// Score floor line
		floorDiscards := player.ScoreFloorLine()
		g.Bag.Discard(floorDiscards)
	}

	// Check for game end
	for _, player := range g.Players {
		if player.HasCompletedRow() {
			g.EndGame()
			return
		}
	}

	// Setup next round
	g.Round++
	g.CurrentPlayer = g.FirstPlayer
	g.SetupRound()
}

// EndGame handles final scoring
func (g *Game) EndGame() {
	g.GameOver = true

	for _, player := range g.Players {
		player.ScoreEndGame()
	}
}

// GetWinner returns the winning player index (or -1 for tie)
func (g *Game) GetWinner() int {
	if !g.GameOver {
		return -1
	}

	maxScore := -1
	winner := -1
	tie := false

	for i, player := range g.Players {
		if player.Score > maxScore {
			maxScore = player.Score
			winner = i
			tie = false
		} else if player.Score == maxScore {
			tie = true
		}
	}

	if tie {
		// Tiebreaker: most completed horizontal lines
		maxRows := -1
		for i, player := range g.Players {
			if player.Score != maxScore {
				continue
			}
			rows := 0
			for row := 0; row < 5; row++ {
				complete := true
				for col := 0; col < 5; col++ {
					if !player.Wall[row][col] {
						complete = false
						break
					}
				}
				if complete {
					rows++
				}
			}
			if rows > maxRows {
				maxRows = rows
				winner = i
			}
		}
	}

	return winner
}

// Clone creates a deep copy of the game state (for AI)
func (g *Game) Clone() *Game {
	newG := &Game{
		Players:       make([]*PlayerBoard, g.NumPlayers),
		Factories:     make([]*Factory, len(g.Factories)),
		Center:        g.Center.Clone(),
		Bag:           g.Bag.Clone(),
		CurrentPlayer: g.CurrentPlayer,
		FirstPlayer:   g.FirstPlayer,
		Round:         g.Round,
		GameOver:      g.GameOver,
		NumPlayers:    g.NumPlayers,
	}

	for i, p := range g.Players {
		newG.Players[i] = p.Clone()
	}

	for i, f := range g.Factories {
		newG.Factories[i] = f.Clone()
	}

	return newG
}
