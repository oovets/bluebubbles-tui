# BlueBubbles TUI - iMessage Terminal Client

A sleek, real-time terminal user interface (TUI) for BlueBubbles, allowing you to send and receive iMessages directly from your terminal.

## Features

- Browse and read iMessage conversations with contact names
- Send messages to any chat (press Enter)
- Real-time message delivery via WebSocket (Socket.IO) with auto-reconnect
- New message indicators - chats with unread messages are highlighted in red and moved to the top
- Full keyboard navigation with Tab/Arrow keys
- Contact name lookup - shows real names instead of phone numbers
- Smart chat sorting by most recent activity
- Toggle chat list visibility and message timestamps

## Prerequisites

- Go 1.24+
- BlueBubbles server running on macOS with iMessage synced
- Network access to your BlueBubbles server (HTTP/HTTPS)
- macOS contacts synced to BlueBubbles (for contact name display)

## Dependencies

The following Go packages are used (automatically installed with `go mod tidy`):

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - UI components (list, textarea, viewport)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/gorilla/websocket` - WebSocket communication
- `github.com/google/uuid` - UUID generation for message IDs
- `github.com/tidwall/gjson` - JSON parsing

## Installation

```bash
cd ~/bluebubbles-tui
go build -o bluebubbles-tui .
```

## Configuration

Set environment variables or create a config file:

### Environment Variables (Easiest)

```bash
export BB_SERVER_URL="https://xxx.xxx.xxx.xxx:1234"
export BB_PASSWORD="your-api-password"
```

### Config File (Optional)

Create `~/.config/bluebubbles-tui/bluebubbles.yaml`:

```yaml
server_url: "https://xxx.xxx.xxx.xxx:1234"
password: "your-api-password"
message_limit: 50
chat_limit: 50
```

## Usage

```bash
./bluebubbles-tui
```

### Keyboard Shortcuts

#### Navigation

| Key | Action |
|-----|--------|
| `Tab` | Toggle focus between chat list and current window |
| `Escape` | Return to chat list from any window |
| `←` | Move to window on the left (or chat list if leftmost) |
| `→` | Move to window on the right |
| `↑` / `↓` or `k` / `j` | Navigate chat list / scroll messages |
| `g` (chat list) | Jump to top of chat list |
| `G` (chat list) | Jump to bottom of chat list |
| `Enter` (chat list) | Open selected chat in the focused window |
| `Enter` (input) | Send message |
| `Shift+Enter` (input) | New line in message |

#### Split Windows

| Key | Action |
|-----|--------|
| `Ctrl+F` | Split focused window horizontally (side by side) |
| `Ctrl+G` | Split focused window vertically (stacked) |
| `Ctrl+W` | Close focused window |

Up to 4 windows can be open simultaneously. Navigate to the chat list and press `Enter` to open a chat in whichever window is currently focused.

#### Toggles

| Key | Action |
|-----|--------|
| `Ctrl+S` | Toggle chat list visibility |
| `Ctrl+T` | Toggle message timestamps |
| `q` / `Ctrl+C` | Quit |

## Architecture

- **models/types.go** - Data structures (Chat, Message, Handle)
- **api/client.go** - REST API client for BlueBubbles server
- **ws/client.go** - WebSocket client for real-time updates (Socket.IO)
- **tui/app.go** - Main TUI model and orchestration
- **tui/chatlist.go** - Chat list component
- **tui/simplelist.go** - Custom scrollable list widget (no auto-centering)
- **tui/messages.go** - Message thread viewport
- **tui/input.go** - Message input box
- **config/config.go** - Configuration loading

## How It Works

1. **Contact Lookup**: Fetches all contacts from BlueBubbles and maps phone numbers to display names
2. **Chat Loading**: Loads all chats and sorts them by most recent message activity
3. **API Connection**: Connects to BlueBubbles REST API to fetch chats and messages with contact names enriched
4. **WebSocket**: Attempts to establish real-time WebSocket connection (Socket.IO) for live updates
5. **Message Sending**: Uses the `/api/v1/message/text` endpoint with Apple Script method and unique tempGuid
6. **Real-time Updates**: Receives new messages via WebSocket with auto-reconnect; incoming messages for inactive chats are highlighted in red and moved to the top of the list
7. **TUI Rendering**: Bubble Tea handles all terminal rendering and event loop

## Troubleshooting

### Connection fails with "certificate signed by unknown authority"
BlueBubbles uses self-signed HTTPS certificates. The client automatically skips TLS verification - this is expected.

### Seeing phone numbers instead of contact names
1. Ensure contacts are synced in your BlueBubbles server
2. Check the web interface contacts: https://xxx.xxx.xxx.xxx:1234/web
3. The client fetches contacts from BlueBubbles - if they're not there, names won't appear
4. Try restarting BlueBubbles server

### No chats showing
1. Ensure your BlueBubbles server has synced iMessages
2. Check the web interface: https://xxx.xxx.xxx.xxx:1234
3. Restart BlueBubbles if needed
4. Verify credentials are correct
5. Most active chats (with recent messages) will appear at the top

### Message sending fails
1. Verify you have an active chat selected (press Enter on a chat)
2. Make sure the input box is focused (press Tab to navigate)
3. Press Enter to send (not Ctrl+D)
4. Check the log file (~/.bluebubbles-tui.log) for API errors

### Messages not updating in real-time
1. WebSocket connection may have failed - check network/firewall rules
2. Real-time updates require the WebSocket connection to be active
3. Check firewall/network rules between this client and BlueBubbles server

## Building from Source

```bash
export PATH="/usr/local/go/bin:$PATH"
cd ~/bluebubbles-tui
go mod tidy
go build
```
