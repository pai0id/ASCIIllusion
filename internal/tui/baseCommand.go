package tui

import tea "github.com/charmbracelet/bubbletea"

type Command func(TUIModel) (TUIModel, tea.Cmd)

type ResizeWindowCommand func(TUIModel, int, int) (TUIModel, tea.Cmd)
