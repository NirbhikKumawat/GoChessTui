package ui

import (
	"fmt"
	"strings"

	"github.com/NirbhikKumawat/GoChess/chess"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type EnterMode int

const (
	SelectMode EnterMode = iota
	MoveMode
)

var pieceMap = map[uint8]map[uint8]string{
	chess.White: {chess.Pawn: "♟", chess.Knight: "♞", chess.Bishop: "♝", chess.Rook: "♜", chess.Queen: "♛", chess.King: "♚", chess.Empty: "  "},
	chess.Black: {chess.Pawn: "♙", chess.Knight: "♘", chess.Bishop: "♗", chess.Rook: "♖", chess.Queen: "♕", chess.King: "♔", chess.Empty: "  "},
}

type Model struct {
	board         *chess.Board
	cursorSq      uint8
	selectedSq    uint8
	validMoves    chess.MoveList
	statusMessage string
	enterMode     EnterMode
}

func InitialModel() Model {
	b, _ := chess.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return Model{
		board:         b,
		cursorSq:      12,
		selectedSq:    255,
		statusMessage: "White to move",
		validMoves:    b.GenerateLegalMoves(),
		enterMode:     SelectMode,
	}
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursorSq < 56 {
				m.cursorSq += 8
			}
		case "down":
			if m.cursorSq > 7 {
				m.cursorSq -= 8
			}
		case "left":
			if m.cursorSq%8 != 0 {
				m.cursorSq--
			}
		case "right":
			if m.cursorSq%8 != 7 {
				m.cursorSq++
			}
		case "enter":
			color := m.board.GetColorType(m.cursorSq)
			//piece:= m.board.GetPieceType(m.cursorSq)
			if m.enterMode == SelectMode {
				if m.board.SideToMove == color {
					m.selectedSq = m.cursorSq
					m.statusMessage = "Piece Selected"
					m.enterMode = MoveMode
				} else if color == 2 {
					m.statusMessage = "Square is empty"
				} else {
					m.statusMessage = "invalid selection"
				}
			} else if m.enterMode == MoveMode {
				if m.cursorSq == m.selectedSq {
					m.selectedSq = 255
					m.statusMessage = "Piece DeSelected"
					m.enterMode = SelectMode
					break
				}
				moveFound := false
				for i := 0; i < m.validMoves.Count; i++ {
					move := m.validMoves.Moves[i]
					if move.From() == m.selectedSq && move.To() == m.cursorSq {
						flags := move.Flags()
						if (flags >= 8 && flags <= 10) || (flags >= 12 && flags <= 14) {
							continue
						}
						m.board.MakeMove(move)
						m.validMoves = m.board.GenerateLegalMoves()
						m.selectedSq = 255
						m.statusMessage = "Moved square"
						m.enterMode = SelectMode
						moveFound = true
						break
					}
				}
				if !moveFound {
					m.selectedSq = 255
					m.statusMessage = "invalid move"
					m.enterMode = SelectMode
				}
			}

		}
	}
	return m, nil
}
func (m Model) ViewChessBoard() string {
	chessBoard := [][]string{
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
	}
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			sq := uint8(rank*8 + file)
			color := m.board.GetColorType(sq)
			ptype := m.board.GetPieceType(sq)
			if sq == m.cursorSq {
				chessBoard[7-rank][file] = fmt.Sprintf("[%s]", pieceMap[color][ptype])
			} else if sq == m.selectedSq {
				chessBoard[7-rank][file] = fmt.Sprintf("{%s}", pieceMap[color][ptype])
			} else {
				chessBoard[7-rank][file] = fmt.Sprintf(" %s ", pieceMap[color][ptype])
			}
		}
	}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(chessBoard...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})
	ranks := strings.Join([]string{" A", "B", "C", "D", "E", "F", "G", "H  "}, "     ")
	files := strings.Join([]string{" 8", "7", "6", "5", "4", "3", "2", "1 "}, "\n\n ")

	return lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(lipgloss.Center, files, t.Render()),
		ranks,
	)
}
func (m Model) View() string {
	var b strings.Builder
	b.WriteString("GoChess\n")
	b.WriteString(m.ViewChessBoard())
	b.WriteString("\n" + m.statusMessage + "\n")
	b.WriteString("Use arrow keys to move, Space/Enter to select. 'q' to quit.\n")
	return b.String()
}
