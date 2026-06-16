output "server_public_ip" {
  value = aws_instance.server.public_ip
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
