#!/bin/bash
set -e

echo "🧩 Universe Cleaner"

systemctl disable --now systemd-resolved
rm /etc/resolv.conf
cat <<EOF > /etc/resolv.conf
nameserver 1.1.1.1
nameserver 1.0.0.1
EOF
chattr +i /etc/resolv.conf
