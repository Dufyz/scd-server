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
