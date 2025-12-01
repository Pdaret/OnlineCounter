#!/bin/bash
set -e

echo "⚙️ Configuring systemd-resolved DNS..."

cat <<EOF > /etc/systemd/resolved.conf
[Resolve]
DNS=1.1.1.1 1.0.0.1
FallbackDNS=8.8.8.8 8.8.4.4
DNSStubListener=no
EOF

systemctl restart systemd-resolved

echo "✅ systemd-resolved DNS configuration applied."
