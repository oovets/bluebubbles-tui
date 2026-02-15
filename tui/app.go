package tui

import (
	"encoding/json"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bluebubbles-tui/api"
	"github.com/bluebubbles-tui/models"
	"github.com/bluebubbles-tui/ws"
)

type focusRegion int

const (
	focusChatList focusRegion = iota
	focusMessages
	focusInput
)

// Message types for Bubble Tea
type (
	chatsLoadedMsg     []models.Chat
	messagesLoadedMsg  []models.Message
	sendSuccessMsg     struct{}
	sendErrMsg         error
	wsEventMsg         models.WSEvent
	wsConnectSuccessMsg struct{}
	wsConnectFailMsg   error
	errMsg             error
	msgWithChatGUID    struct {
		msg      models.Message
		chatGUID string
	}
)

type AppModel struct {
	// Sub-components
	chatList ChatListModel
	messages MessagesModel
	input    InputModel

	// State
	activeChat      *models.Chat
	loading         bool
	err             error
	wsConnected     bool
	lastRefreshTime time.Time

	// Clients
	apiClient *api.Client
	wsClient  *ws.Client

	// Terminal dimensions
	width  int
	height int

	// Focus tracking
	focused focusRegion
}

func NewAppModel(client *api.Client, wsClient *ws.Client) AppModel {
	return AppModel{
		chatList:  NewChatListModel(),
		messages:  NewMessagesModel(),
		input:     NewInputModel(),
		apiClient: client,
		wsClient:  wsClient,
		focused:   focusChatList,
		width:     80,
		height:    24,
	}
}

func (m AppModel) Init() tea.Cmd {
	cmds := []tea.Cmd{
		loadChatsCmd(m.apiClient),
		m.input.Focus(),
	}

	// Try to connect WebSocket for real-time updates
	if m.wsClient != nil {
		cmds = append(cmds, connectWSCmd(m.wsClient))
	}

	return tea.Batch(cmds...)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case chatsLoadedMsg:
		m.chatList.SetChats([]models.Chat(msg))
		if len(msg) > 0 {
			m.activeChat = &msg[0]
			m.messages.SetChatName(msg[0].GetDisplayName())
			return m, loadMessagesCmd(m.apiClient, msg[0].GUID)
		}
		return m, nil

	case messagesLoadedMsg:
		m.messages.SetMessages([]models.Message(msg))
		return m, nil

	case sendSuccessMsg:
		m.input.Clear()
		if m.activeChat != nil {
			return m, loadMessagesCmd(m.apiClient, m.activeChat.GUID)
		}
		return m, nil

	case sendErrMsg:
		m.err = msg
		return m, nil

	case wsConnectSuccessMsg:
		// WebSocket connected, start listening for events
		m.wsConnected = true
		return m, waitForWSEventCmd(m.wsClient)

	case wsConnectFailMsg:
		m.err = msg
		// Fall back to polling every 10 seconds
		return m, nil

	case wsEventMsg:
		// Handle WebSocket event
		return m.handleWSEvent(models.WSEvent(msg))

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		// Handle global keys first
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Cycle focus: chatList -> input -> messages -> chatList
			switch m.focused {
			case focusChatList:
				m.focused = focusInput
				m.input.textarea.Focus()
			case focusInput:
				m.focused = focusMessages
				m.input.textarea.Blur()
			case focusMessages:
				m.focused = focusChatList
				m.input.textarea.Blur()
			}
			return m, nil
		case "enter":
			if m.focused == focusChatList {
				// Select chat and load messages
				selected := m.chatList.SelectedChat()
				if selected != nil {
					m.activeChat = selected
					m.messages.SetChatName(selected.GetDisplayName())
					return m, loadMessagesCmd(m.apiClient, selected.GUID)
				}
				return m, nil
			} else if m.focused == focusInput {
				// Send message when pressing Enter in input
				text := m.input.GetText()
				if text != "" && m.activeChat != nil {
					return m, sendMessageCmd(m.apiClient, m.activeChat.GUID, text)
				}
				return m, nil
			}
		}
		// Fall through to delegate to sub-component
	}

	// Delegate to focused sub-component
	var cmd tea.Cmd
	switch m.focused {
	case focusChatList:
		m.chatList, cmd = m.chatList.Update(msg)
	case focusInput:
		m.input, cmd = m.input.Update(msg)
	case focusMessages:
		m.messages, cmd = m.messages.Update(msg)
	}

	return m, cmd
}

func (m *AppModel) updateLayout() {
	chatListW, messagesW, messagesH, _ := CalculateLayout(m.width, m.height)
	m.chatList.SetSize(chatListW, m.height-3)
	m.messages.SetSize(messagesW, messagesH)
	m.input.SetSize(messagesW)
}

func (m AppModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Render chat list panel
	chatListStyle := PanelStyle
	if m.focused == focusChatList {
		chatListStyle = ActivePanelStyle
	}
	chatPanel := chatListStyle.
		Width(ChatListWidth).
		Height(m.height - 3).
		Render(m.chatList.View())

	// Render messages panel
	messagesStyle := PanelStyle
	if m.focused == focusMessages {
		messagesStyle = ActivePanelStyle
	}

	messagesWidth := m.width - ChatListWidth - 4
	messagesHeight := m.height - InputHeight - 3

	messagesView := m.messages.View()
	messagesPanel := messagesStyle.
		Width(messagesWidth).
		Height(messagesHeight).
		Render(messagesView)

	// Render input panel
	inputStyle := PanelStyle
	if m.focused == focusInput {
		inputStyle = ActivePanelStyle
	}
	inputPanel := inputStyle.
		Width(messagesWidth).
		Height(InputHeight).
		Render(m.input.View())

	// Stack messages and input
	rightPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		messagesPanel,
		inputPanel,
	)

	// Join panels horizontally
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		chatPanel,
		rightPanel,
	)

	// Render status bar
	statusLeft := "  [Tab] switch  [↑↓] scroll  [Ctrl+D] send  [q] quit"
	statusRight := ""

	if m.loading {
		statusRight = " ⟳ Loading... "
	} else if m.err != nil {
		statusRight = fmt.Sprintf(" ⚠ Error: %v ", m.err)
	} else if m.wsConnected {
		statusRight = " 🔗 Live "
	} else {
		statusRight = " 📡 Polling "
	}

	statusBar := StatusBarStyle.
		Width(m.width - 2).
		Render(statusLeft + lipgloss.NewStyle().
			Align(lipgloss.Right).
			Width(m.width - len(statusLeft) - 4).
			Render(statusRight))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		statusBar,
	)
}

// Command constructors

func loadChatsCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		chats, err := client.GetChats(50)
		if err != nil {
			return errMsg(fmt.Errorf("failed to load chats: %v", err))
		}
		return chatsLoadedMsg(chats)
	}
}

func loadMessagesCmd(client *api.Client, chatGUID string) tea.Cmd {
	return func() tea.Msg {
		messages, err := client.GetMessages(chatGUID, 50)
		if err != nil {
			return errMsg(fmt.Errorf("failed to load messages: %v", err))
		}
		// Already reversed by API client
		return messagesLoadedMsg(messages)
	}
}

func sendMessageCmd(client *api.Client, chatGUID, text string) tea.Cmd {
	return func() tea.Msg {
		if err := client.SendMessage(chatGUID, text); err != nil {
			return sendErrMsg(err)
		}
		return sendSuccessMsg{}
	}
}

func connectWSCmd(wsClient *ws.Client) tea.Cmd {
	return func() tea.Msg {
		if err := wsClient.Connect(); err != nil {
			return wsConnectFailMsg(fmt.Errorf("websocket connection failed: %v", err))
		}
		return wsConnectSuccessMsg{}
	}
}

func waitForWSEventCmd(wsClient *ws.Client) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-wsClient.Events
		if !ok {
			return errMsg(fmt.Errorf("websocket connection closed"))
		}
		return wsEventMsg(event)
	}
}

// handleWSEvent processes incoming WebSocket events
func (m *AppModel) handleWSEvent(event models.WSEvent) (tea.Model, tea.Cmd) {
	switch event.Type {
	case "new-message":
		// Parse incoming message
		var msg models.Message
		if err := json.Unmarshal(event.Data, &msg); err != nil {
			return m, nil
		}

		// Only add if it's for the currently active chat
		if m.activeChat != nil && msg.ChatGUID == m.activeChat.GUID {
			// Get current messages and append
			currentMsgs := m.messages.messages
			currentMsgs = append(currentMsgs, msg)
			m.messages.SetMessages(currentMsgs)
		}

		// Re-arm WebSocket listener
		return m, waitForWSEventCmd(m.wsClient)

	case "updated-message":
		// Message updated (read status, etc.)
		// Re-arm WebSocket listener
		return m, waitForWSEventCmd(m.wsClient)

	case "chat-read-status-changed":
		// Chat read status changed
		// Re-arm WebSocket listener
		return m, waitForWSEventCmd(m.wsClient)

	default:
		// Unknown event type, just re-arm
		return m, waitForWSEventCmd(m.wsClient)
	}
}
