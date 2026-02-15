#!/bin/bash

# BlueBubbles TUI runner script

export PATH="/usr/local/go/bin:$PATH"
export BB_SERVER_URL="${BB_SERVER_URL:-https://192.168.0.159:8443}"
export BB_PASSWORD="${BB_PASSWORD:-}"

if [ -z "$BB_PASSWORD" ]; then
    echo "Error: BB_PASSWORD environment variable not set"
    echo "Usage: BB_PASSWORD='your-password' ./run.sh"
    exit 1
fi

echo "Starting BlueBubbles TUI..."
echo "Server: $BB_SERVER_URL"
echo "Logs: ~/.bluebubbles-tui.log"
echo ""

./bluebubbles-tui

echo ""
echo "=== Recent Logs ==="
tail -20 ~/.bluebubbles-tui.log
