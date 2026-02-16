package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bluebubbles-tui/models"
)

type MessagesModel struct {
	viewport viewport.Model
	messages []models.Message
	chatName string
	width    int
	height   int
}

func NewMessagesModel() MessagesModel {
	vp := viewport.New(60, 15)
	vp.MouseWheelEnabled = true

	return MessagesModel{
		viewport: vp,
	}
}

func (m *MessagesModel) SetMessages(messages []models.Message) {
	m.messages = messages
	m.renderContent()
}

// AppendMessage adds a single message to the list
func (m *MessagesModel) AppendMessage(msg models.Message) {
	m.messages = append(m.messages, msg)
	m.renderContent()
}

func (m *MessagesModel) SetChatName(name string) {
	m.chatName = name
}

func (m *MessagesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	// Reserve 1 line for the chat name header
	m.viewport.Height = height - 1
	m.renderContent()
}

func (m *MessagesModel) renderContent() {
	if len(m.messages) == 0 {
		m.viewport.SetContent("(No messages yet)")
		return
	}

	var sb strings.Builder

	for _, msg := range m.messages {
		timeStr := msg.ParsedTime().Format("15:04")

		var sender string
		if msg.IsFromMe {
			sender = "You"
		} else if msg.Handle != nil && msg.Handle.DisplayName != "" {
			sender = msg.Handle.DisplayName
		} else if msg.Handle != nil {
			sender = msg.Handle.Address
		} else {
			sender = "Unknown"
		}

		text := msg.Text

		// Wrap text if needed
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			if i == 0 {
				if msg.IsFromMe {
					// Apply style to entire message including timestamp, sender, and text
					styledLine := fmt.Sprintf("[%s] **%s**: %s", timeStr, sender, line)
					sb.WriteString(MyMessageStyle.Render(styledLine))
				} else {
					styledLine := fmt.Sprintf("[%s] **%s**: %s", timeStr, sender, line)
					sb.WriteString(TheirMessageStyle.Render(styledLine))
				}
			} else {
				sb.WriteString("\n")
				if msg.IsFromMe {
					sb.WriteString(MyMessageStyle.Render(line))
				} else {
					sb.WriteString(TheirMessageStyle.Render(line))
				}
			}
		}
		sb.WriteString("\n")
	}

	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
}

func (m *MessagesModel) ScrollUp() {
	m.viewport.LineUp(3)
}

func (m *MessagesModel) ScrollDown() {
	m.viewport.LineDown(3)
}

func (m MessagesModel) Update(msg tea.Msg) (MessagesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m MessagesModel) View() string {
	header := ""
	if m.chatName != "" {
		header = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Render(m.chatName) + "\n"
	}

	return header + m.viewport.View()
}
