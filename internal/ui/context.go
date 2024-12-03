package ui

import (
	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

type Context struct {
	mtx asciiser.ASCIImtx
	v   *visualiser.Visualiser
}

func NewContext(mtx asciiser.ASCIImtx, v *visualiser.Visualiser) *Context {
	return &Context{mtx: mtx, v: v}
}

// func (m *Context) setMatrix(mtx asciiser.ASCIImtx) {
// 	m.mtx = mtx
// }

// func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if f, ok := m.commands[msg.String()]; ok {
// 			m, cmd := f(m)
// 			mtx, err := m.v.Reconvert(m.mtx)
// 			if err != nil {
// 				fmt.Printf("Error reconverting matrix: %v\n", err)
// 				return m, tea.Quit
// 			}
// 			m.SetMatrix(mtx)
// 			return m, cmd
// 		}
// 	case tea.WindowSizeMsg:
// 		mtx, err := m.v.Resize(m.mtx, msg.Width, msg.Height)
// 		if err != nil {
// 			fmt.Printf("Error resizing matrix: %v\n", err)
// 			return m, tea.Quit
// 		}
// 		m.SetMatrix(mtx)
// 		return m, nil
// 	default:
// 		return m, nil
// 	}
// 	return m, nil
// }

// func (m TUIModel) View() string {
// 	var result string
// 	for i := range m.mtx {
// 		for j := range m.mtx[i] {
// 			result += fmt.Sprintf("%c", m.mtx[i][j])
// 		}
// 		result += "\n"
// 	}
// 	return result
// }

// func DrawASCIImtr(mtx asciiser.ASCIImtx) {
// 	for i := range mtx {
// 		for j := range mtx[i] {
// 			fmt.Printf("%c", mtx[i][j])
// 		}
// 		fmt.Println()
// 	}
// }
