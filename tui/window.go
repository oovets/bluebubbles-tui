package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bluebubbles-tui/models"
)

// WindowID uniquely identifies a chat window
type WindowID int

// ChatWindow represents a single chat view with its own messages and input
type ChatWindow struct {
	ID       WindowID
	Chat     *models.Chat  // Which chat is displayed (nil = empty window)
	Messages MessagesModel // Own viewport for messages
	Input    InputModel    // Own input field
	Focused  bool          // Has keyboard focus?

	// Calculated dimensions from layout
	x, y, width, height int
}

// NewChatWindow creates a new empty chat window
func NewChatWindow(id WindowID) *ChatWindow {
	return &ChatWindow{
		ID:       id,
		Messages: NewMessagesModel(),
		Input:    NewInputModel(),
		Focused:  false,
	}
}

// SetBounds sets the window position and size
func (w *ChatWindow) SetBounds(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height

	// Reserve space for input (InputHeight)
	messagesHeight := height - InputHeight
	if messagesHeight < 1 {
		messagesHeight = 1
	}

	// Update sub-component sizes (subtract padding only)
	w.Messages.SetSize(width-2, messagesHeight)
	w.Input.SetSize(width - 2)
}

// SetChat sets the chat displayed in this window.
// It copies the chat to avoid stale pointer issues when the chat list is reordered.
func (w *ChatWindow) SetChat(chat *models.Chat) {
	if chat != nil {
		chatCopy := *chat
		w.Chat = &chatCopy
		w.Messages.SetChatName(chatCopy.GetDisplayName())
		w.Messages.SetMessages(nil) // Clear stale messages before fresh load
	} else {
		w.Chat = nil
		w.Messages.SetChatName("")
		w.Messages.SetMessages(nil)
	}
}

// Update handles messages for this window
func (w *ChatWindow) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	if w.Focused {
		var cmd tea.Cmd
		w.Input, cmd = w.Input.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		w.Messages, cmd = w.Messages.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// View renders the window
func (w *ChatWindow) View() string {
	// Pick style based on focus
	var style lipgloss.Style
	if w.Focused {
		style = FocusedWindowStyle
	} else {
		style = UnfocusedWindowStyle
	}

	// Calculate content dimensions (inside padding)
	contentWidth := w.width - 2  // subtract padding
	contentHeight := w.height

	if contentWidth < 1 {
		contentWidth = 1
	}
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Handle empty window
	if w.Chat == nil {
		placeholder := lipgloss.NewStyle().
			Foreground(ColorAccent).
			Align(lipgloss.Center).
			Width(contentWidth).
			Height(contentHeight).
			Render("Select a chat\n(Enter in chat list)")

		return style.
			Width(w.width).
			Height(w.height).
			Render(placeholder)
	}

	// Calculate heights for messages and input
	inputHeight := InputHeight
	messagesHeight := contentHeight - inputHeight
	if messagesHeight < 1 {
		messagesHeight = 1
	}

	// Render messages
	messagesView := w.Messages.View()

	// Render input
	inputView := w.Input.View()

	// Stack messages and input
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Width(contentWidth).
			Height(messagesHeight).
			MaxHeight(messagesHeight).
			Render(messagesView),
		lipgloss.NewStyle().
			Width(contentWidth).
			Height(inputHeight).
			Render(inputView),
	)

	return style.
		Width(w.width).
		Height(w.height).
		Render(content)
}
