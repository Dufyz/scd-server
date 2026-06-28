resource "terraform_data" "sync_kafka_clients" {
  triggers_replace = {
    kafka_brokers = "${aws_instance.kafka.private_ip}:9094"
  }

  provisioner "local-exec" {
    command = join(" ", [
      "bash",
      "${path.module}/scripts/sync-kafka-clients.sh",
      local_sensitive_file.ssh_key.filename,
      "${aws_instance.kafka.private_ip}:9094",
      aws_instance.ai_server.public_ip,
      aws_instance.socket_server.public_ip,
      join(",", aws_instance.server[*].public_ip),
    ])
  }

  depends_on = [
    aws_instance.kafka,
    aws_instance.ai_server,
    aws_instance.socket_server,
    aws_instance.server,
  ]
}

resource "local_sensitive_file" "ssh_key" {
  content         = tls_private_key.deploy.private_key_pem
  filename        = "${path.module}/.terraform/sync-ssh-key.pem"
  file_permission = "0600"
}
