package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bluebubbles-tui/models"
)

type ChatListModel struct {
	list   SimpleListModel
	chats  []models.Chat
	width  int
	height int
}

func NewChatListModel() ChatListModel {
	return ChatListModel{
		list: NewSimpleListModel(),
	}
}

func (m *ChatListModel) SetChats(chats []models.Chat) {
	m.chats = chats
	m.list.SetItems(chats)
}

func (m *ChatListModel) SetSize(width, height int) {
	if m.width == width && m.height == height {
		return
	}
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

func (m *ChatListModel) SelectedChat() *models.Chat {
	return m.list.SelectedItem()
}

func (m ChatListModel) Update(msg tea.Msg) (ChatListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ChatListModel) View() string {
	return m.list.View()
}
