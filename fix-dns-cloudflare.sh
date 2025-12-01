#!/bin/bash
set -e

echo "🧩 Fixing /etc/resolv.conf"

# Remove symlink if it exists
if [ -L /etc/resolv.conf ]; then
    rm -f /etc/resolv.conf
fi

# Create a new resolv.conf with Cloudflare DNS
cat <<EOF > /etc/resolv.conf
nameserver 1.1.1.1
nameserver 1.0.0.1
EOF

# Make immutable
chattr +i /etc/resolv.conf

echo "✅ /etc/resolv.conf fixed and locked."
