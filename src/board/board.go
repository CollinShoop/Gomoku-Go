package board

import "fmt"

var pieceMapping = [3]string{".", "1", "2"}

type board struct {
	Rows, Cols                               int
	pieces                                   []int
	lastMovePlayer, lastMoveRow, lastMoveCol int
	isOverCached                             bool
}

func NewBoard(rows, cols int) *board {
	b := new(board)
	b.Rows = rows
	b.Cols = cols
	b.pieces = make([]int, rows*cols)
	b.lastMovePlayer = -1
	b.lastMoveRow = -1
	b.lastMoveCol = -1
	b.isOverCached = false
	return b
}

func (b *board) Pieces() []int {
	piecesCopy := b.pieces
	return piecesCopy
}

// Get which player has a piece in a given position
// It returns a number, -1 for INVALID, 0 for NONE, 1 for PLAYER 1 and 2 for PLAYER 2
func (b *board) Get(col, row int) int {
	if col < 0 || row < 0 || col >= b.Cols || row >= b.Rows {
		return -1 // panic?
	}
	return b.pieces[row*b.Cols+col]
}

// IsOver returns a pair of values. The first indicates whether or not the game is over. The second value
// indicates the most recent player to have moved (ex player 1 or 2). If the game is over, this is also the winner.
func (b *board) IsOver() (bool, int) {
	if b.isOverCached {
		return true, b.lastMovePlayer
	}
	if b.lastMovePlayer == -1 {
		fmt.Printf("isOver: player=%d\n", b.lastMovePlayer)
		return false, b.lastMovePlayer
	}
	lineCount := b.maxLineCount(b.lastMovePlayer, b.lastMoveCol, b.lastMoveRow)
	// according to rules of Standard Gomoku, "overline" eg 6+ is illegal and does not count
	// for Freestyle, it's allowed.
	b.isOverCached = lineCount >= 5
	fmt.Printf("isOver: player=%d; lineCount=%d; isOver=%v\n", b.lastMovePlayer, lineCount, b.isOverCached)
	return b.isOverCached, b.lastMovePlayer
}

func (b *board) Move(player, col, row int) {
	if col < 0 || row < 0 || col >= b.Cols || row >= b.Rows {
		return // panic?
	}
	if player < 0 || player > 2 {
		return // panic?
	}
	isOver, _ := b.IsOver()
	if isOver {
		return // panic?
	}

	fmt.Printf("Move: player=%d; col=%d; row=%d;\n", player, col, row)
	b.pieces[row*b.Cols+col] = player
	// trace out in each direction to see if this piece will flip any others
	// this happens when placing this piece encloses the opponent
	b.flipCheck(player, col, row)

	b.lastMovePlayer = player
	b.lastMoveCol = col
	b.lastMoveRow = row
}

func (b *board) flipSingle(col, row int) {
	piece := b.Get(col, row)
	if piece == -1 || piece == 0 {
		return
	}
	b.pieces[row*b.Cols+col] = (piece % 2) + 1
}

func (b *board) flipCheck(player, col, row int) {
	b.directionalFlipCheck(player, col, row, -1, -1)
	b.directionalFlipCheck(player, col, row, -1, 0)
	b.directionalFlipCheck(player, col, row, -1, 1)
	b.directionalFlipCheck(player, col, row, 1, -1)
	b.directionalFlipCheck(player, col, row, 1, 0)
	b.directionalFlipCheck(player, col, row, 1, 1)
	b.directionalFlipCheck(player, col, row, 0, -1)
	b.directionalFlipCheck(player, col, row, 0, 1)
}

func (b *board) maxLineCount(player, col, row int) int {
	return maxInt(
		b.maxLineCountDirectional(player, col, row, -1, -1), // upper-left and bottom-right
		b.maxLineCountDirectional(player, col, row, -1, 1),  // upper-right and bottom-left
		b.maxLineCountDirectional(player, col, row, 0, 1),   // up-down
		b.maxLineCountDirectional(player, col, row, 1, 0))   // left-right
}

func maxInt(value int, values ...int) int {
	max := value
	for i := 0; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}
	return max
}

// maxLineCountDirectional finds the number of player pieces continuously found in a given direction
// including pieces also in the opposite direction
func (b *board) maxLineCountDirectional(player, col, row, dcol, drow int) int {
	count := 1
	if dcol == 0 && drow == 0 {
		return count
	}
	// go in specified direction
	for icol, irow := col+dcol, row+drow; ; // start 1 iteration in to avoid flipping the first piece

	icol, irow = icol+dcol, irow+drow {
		piece := b.Get(icol, irow)
		if piece == player {
			count++
		} else {
			break
		}
	}
	dcol *= -1
	drow *= -1
	// go in the other direction
	for icol, irow := col+dcol, row+drow; ; // start 1 iteration in to avoid flipping the first piece

	icol, irow = icol+dcol, irow+drow {
		piece := b.Get(icol, irow)
		if piece == player {
			count++
		} else {
			break
		}
	}
	return count
}

// directionalFlipCheck will check for and flip any pieces as a result
//  of placing a given piece (player) at the col and row specified.
// This only checks in a single direction, indicated by dcol/drow (ex -1, -1) for upper-left direction.
func (b *board) directionalFlipCheck(player, col, row, dcol, drow int) {
	if dcol == 0 && drow == 0 {
		return
	}
	fmt.Printf("directionalFlipCheck: player=%d; col=%d; row=%d; dcol=%d; drow=%d\n", player, col, row, dcol, drow)
	var tcol, trow int
	for icol, irow := col, row; ; {
		icol += dcol
		irow += drow
		piece := b.Get(icol, irow)
		fmt.Printf("\tdirectionalFlipCheck: icol=%d; irow=%d; piece=%d\n", icol, irow, piece)
		if piece == -1 || piece == 0 {
			// wall or empty spot means no flip
			fmt.Printf("\tdirectionalFlipCheck: out of bounds or empty terminated.\n")
			return
		}
		if piece == player {
			// hit another piece the same as the player, this terminates the trace
			tcol = icol
			trow = irow
			break
		}
		// else piece is other player by default, continue
	}
	fmt.Printf("\tdirectionalFlipCheck: Starting flipping.\n")
	// at this point tcol and trow represents the terminal location after tracing
	// and running into the next matching player piece, so the next step is to flip
	// all pieces between (col, row) and (tcol, trow)
	for icol, irow := col+dcol, row+drow; // start 1 iteration in to avoid flipping the first piece
	icol != tcol || irow != trow; icol, irow = icol+dcol, irow+drow {
		fmt.Printf("\t\tdirectionalFlipCheck: Flipping col=%d;row=%d.\n", icol, irow)
		b.flipSingle(icol, irow)
	}
}

func (b *board) ToStr() string {
	toStr := ""
	toStr += fmt.Sprint("    ")
	for i := 0; i < b.Cols; i++ {
		toStr += fmt.Sprintf(" %X ", i)
	}
	toStr += fmt.Sprint("  \n    ")
	for i := 0; i < b.Cols; i++ {
		toStr += fmt.Sprint("___")
	}
	for i := 0; i < b.Cols*b.Rows; i++ {
		p := b.pieces[i]
		if i%b.Cols == 0 {
			if i > 0 {
				toStr += fmt.Sprint("│")
			}
			toStr += fmt.Sprintf("\n%X  │", i/b.Cols)
		}
		toStr += fmt.Sprintf(" %s ", pieceMapping[p])
	}
	toStr += fmt.Sprint("|\n    ")
	for i := 0; i < b.Cols; i++ {
		toStr += fmt.Sprint("‾‾‾")
	}
	return toStr
}

func (b *board) Print() {
	fmt.Println(b.ToStr())
}
