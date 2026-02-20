package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputModel struct {
	textarea textarea.Model
}

func NewInputModel() InputModel {
	ta := textarea.New()
	ta.Placeholder = ""
	ta.ShowLineNumbers = false
	ta.CharLimit = 10000
	ta.SetWidth(50)
	ta.SetHeight(3)

	// Strip all colors/borders from the textarea
	plain := ta.FocusedStyle
	plain.Base = lipgloss.NewStyle()
	plain.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle = plain

	blurred := ta.BlurredStyle
	blurred.Base = lipgloss.NewStyle()
	blurred.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle = blurred

	return InputModel{
		textarea: ta,
	}
}

func (m *InputModel) SetSize(width int) {
	m.textarea.SetWidth(width)
}

func (m *InputModel) GetText() string {
	return m.textarea.Value()
}

func (m *InputModel) Clear() {
	m.textarea.Reset()
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	return m.textarea.View()
}

func (m InputModel) Focused() bool {
	return m.textarea.Focused()
}

func (m *InputModel) Focus() tea.Cmd {
	return m.textarea.Focus()
}

func (m *InputModel) Blur() {
	m.textarea.Blur()
}
