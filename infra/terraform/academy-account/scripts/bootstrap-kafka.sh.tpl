#!/bin/bash
set -euxo pipefail

dnf install -y docker git
systemctl enable --now docker
usermod -aG docker ec2-user

DOCKER_COMPOSE_VERSION="v2.27.0"
mkdir -p /usr/local/lib/docker/cli-plugins
curl -sSL "https://github.com/docker/compose/releases/download/$${DOCKER_COMPOSE_VERSION}/docker-compose-linux-x86_64" \
  -o /usr/local/lib/docker/cli-plugins/docker-compose
chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

PRIVATE_IP=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)

git clone --branch ${repo_branch} --depth 1 ${repo_url} /opt/scd-server

cd /opt/scd-server/infra/kafka
cp .env.example .env
sed -i "s/^KAFKA_EXTERNAL_HOST=.*/KAFKA_EXTERNAL_HOST=$PRIVATE_IP/" .env

docker compose up -d
