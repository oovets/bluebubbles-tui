package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bluebubbles-tui/models"
)

type chatItem struct {
	chat models.Chat
}

// minimalDelegate is a very compact list delegate with no extra spacing
type minimalDelegate struct{}

func (d minimalDelegate) Height() int                             { return 1 }
func (d minimalDelegate) Spacing() int                            { return 0 }
func (d minimalDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d minimalDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	c, ok := item.(chatItem)
	if !ok {
		return
	}

	title := c.Title()
	width := m.Width()
	if len(title) > width {
		title = title[:width-1] + "…"
	}

	// Highlight selected item
	if index == m.Index() {
		title = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("212")).
			Width(width).
			Render(title)
	}

	fmt.Fprint(w, title)
}


func (c chatItem) FilterValue() string {
	return c.chat.DisplayName
}

func (c chatItem) Title() string {
	name := c.chat.GetDisplayName()
	if c.chat.UnreadCount > 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("● ") + name
	}
	return name
}

func (c chatItem) Description() string {
	return ""
}

type ChatListModel struct {
	list   list.Model
	chats  []models.Chat
	width  int
	height int
}

func NewChatListModel() ChatListModel {
	l := list.New([]list.Item{}, minimalDelegate{}, ChatListWidth, 10)
	l.Title = "CHATS"
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowPagination(true)

	return ChatListModel{
		list: l,
	}
}

func (m *ChatListModel) SetChats(chats []models.Chat) {
	m.chats = chats
	items := make([]list.Item, len(chats))
	for i, chat := range chats {
		items[i] = chatItem{chat: chat}
	}
	m.list.SetItems(items)
	// Ensure we start at the top
	m.list.Select(0)
}

func (m *ChatListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

func (m *ChatListModel) SelectedChat() *models.Chat {
	if len(m.chats) == 0 {
		return nil
	}
	idx := m.list.Index()
	if idx < 0 || idx >= len(m.chats) {
		return nil
	}
	return &m.chats[idx]
}

func (m ChatListModel) Update(msg tea.Msg) (ChatListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ChatListModel) View() string {
	return m.list.View()
}
