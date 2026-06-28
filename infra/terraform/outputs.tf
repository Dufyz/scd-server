output "server_public_ips" {
  description = "Public IPs of all server instances (for SSH/deploy)"
  value       = aws_instance.server[*].public_ip
}

output "alb_dns_name" {
  description = "ALB DNS — use this as the public endpoint for the server API"
  value       = aws_lb.server.dns_name
}

output "ai_server_public_ip" {
  value = aws_instance.ai_server.public_ip
}

output "socket_server_public_ip" {
  value = aws_instance.socket_server.public_ip
}

output "kafka_public_ip" {
  value = aws_instance.kafka.public_ip
}

output "kafka_private_ip" {
  description = "Use este IP no KAFKA_BROKERS dos outros serviços (mesma VPC)"
  value       = aws_instance.kafka.private_ip
}

output "redis_public_ip" {
  value = aws_instance.redis.public_ip
}

output "redis_private_ip" {
  description = "Use este IP no REDIS_URL dos outros serviços (mesma VPC)"
  value       = aws_instance.redis.private_ip
}

output "ssh_private_key_pem" {
  description = "Chave privada SSH gerada para acessar as EC2 (salve em um arquivo .pem e/ou no secret SSH_PRIVATE_KEY do GitHub)"
  value       = tls_private_key.deploy.private_key_pem
  sensitive   = true
}

output "rds_endpoint" {
  description = "Endpoint (host:port) do RDS Postgres"
  value       = aws_db_instance.postgres.endpoint
}

output "rds_address" {
  description = "Host do RDS Postgres (sem a porta)"
  value       = aws_db_instance.postgres.address
}

output "database_url" {
  description = "Connection string completa para usar em DATABASE_URL (primary, leitura e escrita)"
  value       = "postgresql://${var.db_username}:${var.db_password}@${aws_db_instance.postgres.address}:5432/${var.db_name}?sslmode=require"
  sensitive   = true
}

output "rds_replica_endpoint" {
  description = "Endpoint (host:port) da read replica"
  value       = aws_db_instance.postgres_read_replica.endpoint
}

output "rds_replica_address" {
  description = "Host da read replica (sem a porta)"
  value       = aws_db_instance.postgres_read_replica.address
}

output "database_url_replica" {
  description = "Connection string completa para usar em DATABASE_URL_REPLICA (somente leitura)"
  value       = "postgresql://${var.db_username}:${var.db_password}@${aws_db_instance.postgres_read_replica.address}:5432/${var.db_name}?sslmode=require"
  sensitive   = true
}
