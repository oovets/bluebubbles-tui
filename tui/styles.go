package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	ChatListWidth = 25  // fixed width for left panel
	InputHeight   = 3   // input box + borders
)

// Color scheme
const (
	ColorPrimary   = lipgloss.Color("212")  // pink
	ColorSecondary = lipgloss.Color("86")   // green
	ColorAccent    = lipgloss.Color("242")  // gray
	ColorBorder    = lipgloss.Color("240")  // dark gray
)

var (
	// Panel styles
	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	ActivePanelStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 1)

	// Chat list styles
	ChatListItemStyle = lipgloss.NewStyle().
		Padding(0).
		Margin(0)

	ChatListItemSelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(ColorPrimary).
		Padding(0).
		Margin(0)

	// Message styles
	MyMessageStyle = lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Align(lipgloss.Right)

	TheirMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Align(lipgloss.Left)

	TimestampStyle = lipgloss.NewStyle().
		Foreground(ColorAccent).
		PaddingRight(1)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	// Input styles
	InputStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ColorBorder)
)

// CalculateLayout returns the optimal dimensions for each panel
func CalculateLayout(screenWidth, screenHeight int) (chatListWidth, messagesWidth, messagesHeight, inputHeight int) {
	// Subtract 4 for borders and padding
	chatListWidth = ChatListWidth
	messagesWidth = screenWidth - chatListWidth - 4 // -2 for left panel border, -2 for right panel border
	messagesHeight = screenHeight - InputHeight - 3 // -1 status bar, -2 borders
	inputHeight = InputHeight

	return
}
