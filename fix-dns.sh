#!/bin/bash
set -e

echo "🧩 Fixing /etc/resolv.conf"

# Remove symlink if it exists
if [ -L /etc/resolv.conf ]; then
    rm -f /etc/resolv.conf
fi

# Create a new resolv.conf with Google DNS
cat <<EOF > /etc/resolv.conf
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

# Make immutable
chattr +i /etc/resolv.conf

echo "✅ /etc/resolv.conf fixed and locked."
