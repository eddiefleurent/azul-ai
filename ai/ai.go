package ai

import (
	"math"
	"math/rand"
	"time"

	"github.com/eddiefleurent/azul-ai/game"
)

// Player interface for both human and AI players
type Player interface {
	ChooseMove(g *game.Game, moves []game.Move) game.Move
	Name() string
}

// Difficulty levels
type Difficulty int

const (
	Easy   Difficulty = iota // Random moves
	Medium                   // Basic heuristics
	Hard                     // Minimax with pruning
)

// AIPlayer implements an AI opponent
type AIPlayer struct {
	difficulty Difficulty
	playerIdx  int
	rng        *rand.Rand
}

// NewAIPlayer creates a new AI player
func NewAIPlayer(difficulty Difficulty, playerIdx int) *AIPlayer {
	return &AIPlayer{
		difficulty: difficulty,
		playerIdx:  playerIdx,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ai *AIPlayer) Name() string {
	switch ai.difficulty {
	case Easy:
		return "AI (Easy)"
	case Medium:
		return "AI (Medium)"
	case Hard:
		return "AI (Hard)"
	default:
		return "AI"
	}
}

// ChooseMove selects the best move based on difficulty
func (ai *AIPlayer) ChooseMove(g *game.Game, moves []game.Move) game.Move {
	if len(moves) == 0 {
		return game.Move{}
	}

	switch ai.difficulty {
	case Easy:
		return ai.randomMove(moves)
	case Medium:
		return ai.heuristicMove(g, moves)
	case Hard:
		return ai.minimaxMove(g, moves)
	default:
		return ai.randomMove(moves)
	}
}

// randomMove picks a random legal move
func (ai *AIPlayer) randomMove(moves []game.Move) game.Move {
	return moves[ai.rng.Intn(len(moves))]
}

// heuristicMove uses simple rules to pick a good move
func (ai *AIPlayer) heuristicMove(g *game.Game, moves []game.Move) game.Move {
	bestScore := math.MinInt32
	var bestMoves []game.Move

	for _, move := range moves {
		score := ai.evaluateMove(g, move)
		if score > bestScore {
			bestScore = score
			bestMoves = []game.Move{move}
		} else if score == bestScore {
			bestMoves = append(bestMoves, move)
		}
	}

	// Random among best moves
	return bestMoves[ai.rng.Intn(len(bestMoves))]
}

// evaluateMove scores a move using heuristics
func (ai *AIPlayer) evaluateMove(g *game.Game, move game.Move) int {
	score := 0
	player := g.Players[ai.playerIdx]

	// Prefer completing pattern lines
	if move.LineIdx >= 0 {
		pl := player.PatternLines[move.LineIdx]
		tilesNeeded := pl.Size - pl.Filled

		// Count how many tiles we're taking
		var tileCount int
		if move.FactoryIdx == -1 {
			for _, t := range g.Center.Tiles {
				if t == move.Color {
					tileCount++
				}
			}
		} else {
			for _, t := range g.Factories[move.FactoryIdx].Tiles {
				if t == move.Color {
					tileCount++
				}
			}
		}

		// Bonus for completing a line
		if tileCount >= tilesNeeded {
			score += 50 + (move.LineIdx+1)*10 // Bigger lines = more points
		}

		// Penalty for overflow
		overflow := tileCount - tilesNeeded
		if overflow > 0 {
			score -= overflow * 15
		}

		// Prefer filling larger lines
		score += move.LineIdx * 5
	} else {
		// Floor placement is bad
		score -= 30
	}

	// Avoid taking first player marker early in round
	if move.FactoryIdx == -1 && g.Center.HasFirstPlayerTile {
		if len(g.Center.Tiles) < 5 {
			score -= 20 // Penalty for taking 1st player with few tiles
		}
	}

	// Bonus for colors that help complete rows/columns/color sets
	score += ai.evaluateBoardProgress(player, move)

	return score
}

// evaluateBoardProgress checks if a move helps with bonus scoring
func (ai *AIPlayer) evaluateBoardProgress(player *game.PlayerBoard, move game.Move) int {
	if move.LineIdx < 0 {
		return 0
	}

	score := 0
	row := move.LineIdx
	col := player.GetWallColumn(row, move.Color)

	// Check row completion progress
	rowFilled := 0
	for c := 0; c < 5; c++ {
		if player.Wall[row][c] {
			rowFilled++
		}
	}
	if rowFilled >= 3 {
		score += 10 // Close to completing row
	}

	// Check column completion progress
	colFilled := 0
	for r := 0; r < 5; r++ {
		if player.Wall[r][col] {
			colFilled++
		}
	}
	if colFilled >= 3 {
		score += 15 // Close to completing column (worth more)
	}

	// Check color completion progress
	colorCount := 0
	for r := 0; r < 5; r++ {
		c := player.GetWallColumn(r, move.Color)
		if player.Wall[r][c] {
			colorCount++
		}
	}
	if colorCount >= 3 {
		score += 20 // Close to completing color set (worth most)
	}

	return score
}

// minimaxMove uses minimax with alpha-beta pruning
func (ai *AIPlayer) minimaxMove(g *game.Game, moves []game.Move) game.Move {
	bestScore := math.MinInt32
	var bestMove game.Move

	// Limit search depth based on game state
	depth := 4
	if len(moves) > 20 {
		depth = 3
	}

	for _, move := range moves {
		// Clone and apply move
		clone := g.Clone()
		if err := clone.ApplyMove(move); err != nil {
			// Skip invalid moves that fail to apply
			continue
		}

		// Determine if the AI is the next player to move (maximizing)
		// After applying the move, check who the current player is in the cloned game
		isAINext := clone.CurrentPlayer == ai.playerIdx
		score := ai.minimax(clone, depth-1, math.MinInt32, math.MaxInt32, isAINext)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove
}

// minimax with alpha-beta pruning
func (ai *AIPlayer) minimax(g *game.Game, depth int, alpha, beta int, maximizing bool) int {
	// Terminal conditions
	if depth == 0 || g.GameOver || g.IsRoundOver() {
		return ai.evaluateState(g)
	}

	moves := g.GetValidMoves()
	if len(moves) == 0 {
		return ai.evaluateState(g)
	}

	if maximizing {
		maxEval := math.MinInt32
		anyValid := false
		for _, move := range moves {
			clone := g.Clone()
			if err := clone.ApplyMove(move); err != nil {
				// Skip invalid moves that fail to apply
				continue
			}
			anyValid = true

			eval := ai.minimax(clone, depth-1, alpha, beta, clone.CurrentPlayer == ai.playerIdx)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)

			if beta <= alpha {
				break
			}
		}
		if !anyValid {
			return ai.evaluateState(g)
		}
		return maxEval
	} else {
		minEval := math.MaxInt32
		anyValid := false
		for _, move := range moves {
			clone := g.Clone()
			if err := clone.ApplyMove(move); err != nil {
				// Skip invalid moves that fail to apply
				continue
			}
			anyValid = true

			eval := ai.minimax(clone, depth-1, alpha, beta, clone.CurrentPlayer == ai.playerIdx)
			minEval = min(minEval, eval)
			beta = min(beta, eval)

			if beta <= alpha {
				break
			}
		}
		if !anyValid {
			return ai.evaluateState(g)
		}
		return minEval
	}
}

// evaluateState scores the current game state for the AI player
func (ai *AIPlayer) evaluateState(g *game.Game) int {
	myPlayer := g.Players[ai.playerIdx]
	score := myPlayer.Score * 10

	// Evaluate pattern line progress
	for i, pl := range myPlayer.PatternLines {
		if pl.Filled > 0 {
			// Partial credit for progress
			progress := float64(pl.Filled) / float64(pl.Size)
			score += int(progress * float64(i+1) * 5)
		}
	}

	// Evaluate wall bonuses potential
	score += ai.evaluateWallPotential(myPlayer) * 2

	// Subtract opponent scores
	for i, player := range g.Players {
		if i != ai.playerIdx {
			score -= player.Score * 8
		}
	}

	// Penalty for floor tiles
	score -= len(myPlayer.FloorLine) * 5

	return score
}

// evaluateWallPotential estimates bonus scoring potential
func (ai *AIPlayer) evaluateWallPotential(player *game.PlayerBoard) int {
	score := 0

	// Row completion potential
	for row := 0; row < 5; row++ {
		filled := 0
		for col := 0; col < 5; col++ {
			if player.Wall[row][col] {
				filled++
			}
		}
		if filled >= 4 {
			score += 10
		} else if filled >= 3 {
			score += 5
		}
	}

	// Column completion potential
	for col := 0; col < 5; col++ {
		filled := 0
		for row := 0; row < 5; row++ {
			if player.Wall[row][col] {
				filled++
			}
		}
		if filled >= 4 {
			score += 20
		} else if filled >= 3 {
			score += 10
		}
	}

	// Color completion potential
	for _, color := range game.AllColors() {
		count := 0
		for row := 0; row < 5; row++ {
			col := player.GetWallColumn(row, color)
			if player.Wall[row][col] {
				count++
			}
		}
		if count >= 4 {
			score += 25
		} else if count >= 3 {
			score += 12
		}
	}

	return score
}
