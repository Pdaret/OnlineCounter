#!/bin/bash

set -e

# Variables (update if needed)
SERVICE_NAME="x-ui-monitor"
BINARY_NAME="x-ui-monitor"
BUILD_PATH="./cmd/app/" # change if your main.go is elsewhere
INSTALL_DIR="/usr/local/x-ui"
SERVICE_FILE="x-ui-monitor.service"

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

echo "✅ Deployment complete. Service status:"
sudo systemctl status $SERVICE_NAME.service
