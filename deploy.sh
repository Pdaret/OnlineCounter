#!/bin/bash

set -e

# Variables (update if needed)
SERVICE_NAME="x-ui-monitor"
BINARY_NAME="x-ui-monitor"
BUILD_PATH="./cmd/app/" # change if your main.go is elsewhere
INSTALL_DIR="/usr/local/x-ui"
SERVICE_FILE="x-ui-monitor.service"
XUI_DB_FILE="./x-ui-toyo.db# Change path if needed

echo "🌐 Installing Nginx..."
sudo apt update
sudo apt install -y nginx

echo "📦 Installing 3x-ui via script..."
bash <(curl -Ls https://raw.githubusercontent.com/mhsanaei/3x-ui/master/install.sh) <<EOF
# pressing Enter to accept defaults
EOF

echo "🗂️ Copying x-ui.db to /etc/x-ui..."
if [ -f "$XUI_DB_FILE" ]; then
    sudo rm -f /etc/x-ui/x-ui.db
    sudo cp "$XUI_DB_FILE" /etc/x-ui/x-ui.db
    sudo chmod 600 /etc/x-ui/x-ui.db
    sudo x-ui restart
    echo "✅ x-ui.db replaced successfully."
else
    echo "❌ x-ui.db file not found at $XUI_DB_FILE"
fi

echo "🔧 Building Go project..."
go build -o "$BINARY_NAME" "$BUILD_PATH"

echo "📂 Creating install directory if not exists..."
sudo mkdir -p "$INSTALL_DIR"

echo "🚚 Copying binary to $INSTALL_DIR..."
sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "📝 Installing systemd service..."
sudo cp "$SERVICE_FILE" /etc/systemd/system/$SERVICE_NAME.service

echo "🔄 Reloading systemd..."
sudo systemctl daemon-reexec
sudo systemctl daemon-reload

echo "📌 Enabling service to start at boot..."
sudo systemctl enable $SERVICE_NAME.service

echo "🚀 Starting service..."
sudo systemctl restart $SERVICE_NAME.service

echo "✅ Full deployment complete!"
