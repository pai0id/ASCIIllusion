package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pai0id/CgCourseProject/internal/asciiser"
)

type TUIModel struct {
	mtx       asciiser.ASCIImtx
	commands  map[string]Command
	resizeCmd ResizeWindowCommand
}

func NewTUIModel(mtx asciiser.ASCIImtx, commands map[string]Command, resizeCmd ResizeWindowCommand) TUIModel {
	return TUIModel{mtx: mtx, commands: commands, resizeCmd: resizeCmd}
}

func (m *TUIModel) AppendCommand(command string, f Command) {
	m.commands[command] = f
}

func (m *TUIModel) SetMatrix(mtx asciiser.ASCIImtx) {
	m.mtx = mtx
}

func (m TUIModel) Init() tea.Cmd {
	return nil
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if f, ok := m.commands[msg.String()]; ok {
			return f(m)
		}
	case tea.WindowSizeMsg:
		if m.resizeCmd != nil {
			return m.resizeCmd(m, msg.Width, msg.Height)
		}
	default:
		return m, nil
	}
	return m, nil
}

func (m TUIModel) View() string {
	var result string
	for i := range m.mtx {
		for j := range m.mtx[i] {
			result += fmt.Sprintf("%c", m.mtx[i][j])
		}
		result += "\n"
	}
	return result
}

func DrawASCIImtr(mtx asciiser.ASCIImtx) {
	for i := range mtx {
		for j := range mtx[i] {
			fmt.Printf("%c", mtx[i][j])
		}
		fmt.Println()
	}
}
