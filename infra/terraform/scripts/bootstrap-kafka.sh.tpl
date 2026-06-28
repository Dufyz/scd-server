#!/bin/bash
set -euxo pipefail

dnf install -y docker git
systemctl enable --now docker

# Aguarda o daemon estar pronto
timeout 30 bash -c 'until docker info > /dev/null 2>&1; do sleep 1; done'

usermod -aG docker ec2-user

IMDS_TOKEN=$(curl -sX PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
PRIVATE_IP=$(curl -s -H "X-aws-ec2-metadata-token: $IMDS_TOKEN" http://169.254.169.254/latest/meta-data/local-ipv4)

git clone --branch ${repo_branch} --depth 1 ${repo_url} /opt/scd-server

cd /opt/scd-server/infra/kafka
cp .env.example .env
sed -i "s/^KAFKA_EXTERNAL_HOST=.*/KAFKA_EXTERNAL_HOST=$PRIVATE_IP/" .env

docker compose up -d

# Garante que o Kafka sobe automaticamente após reboot
cat > /etc/systemd/system/kafka.service << 'EOF'
[Unit]
Description=Kafka
After=docker.service
Requires=docker.service

[Service]
WorkingDirectory=/opt/scd-server/infra/kafka
ExecStart=/usr/bin/docker compose up
ExecStop=/usr/bin/docker compose down
Restart=always

[Install]
WantedBy=multi-user.target
EOF

systemctl enable kafka
