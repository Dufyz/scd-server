locals {
  app_bootstrap = file("${path.module}/scripts/bootstrap-docker.sh.tpl")

  kafka_bootstrap = templatefile("${path.module}/scripts/bootstrap-kafka.sh.tpl", {
    repo_url    = var.repo_url
    repo_branch = var.repo_branch
  })

  redis_bootstrap = templatefile("${path.module}/scripts/bootstrap-redis.sh.tpl", {
    repo_url    = var.repo_url
    repo_branch = var.repo_branch
  })
}

resource "aws_instance" "server" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.deploy.key_name
  subnet_id              = data.aws_subnets.default.ids[0]
  vpc_security_group_ids = [aws_security_group.ssh.id, aws_security_group.internal.id, aws_security_group.app_server.id]
  user_data              = local.app_bootstrap

  tags = {
    Name = "${var.project_name}-server"
  }
}

resource "aws_instance" "ai_server" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.deploy.key_name
  subnet_id              = data.aws_subnets.default.ids[0]
  vpc_security_group_ids = [aws_security_group.ssh.id, aws_security_group.internal.id, aws_security_group.app_ai_server.id]
  user_data              = local.app_bootstrap

  tags = {
    Name = "${var.project_name}-ai-server"
  }
}

resource "aws_instance" "socket_server" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.deploy.key_name
  subnet_id              = data.aws_subnets.default.ids[0]
  vpc_security_group_ids = [aws_security_group.ssh.id, aws_security_group.internal.id, aws_security_group.app_socket_server.id]
  user_data              = local.app_bootstrap

  tags = {
    Name = "${var.project_name}-socket-server"
  }
}

resource "aws_instance" "kafka" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.deploy.key_name
  subnet_id              = data.aws_subnets.default.ids[0]
  vpc_security_group_ids = [aws_security_group.ssh.id, aws_security_group.internal.id, aws_security_group.kafka_ui.id]
  user_data              = local.kafka_bootstrap

  tags = {
    Name = "${var.project_name}-kafka"
  }
}

resource "aws_instance" "redis" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.deploy.key_name
  subnet_id              = data.aws_subnets.default.ids[0]
  vpc_security_group_ids = [aws_security_group.ssh.id, aws_security_group.internal.id]
  user_data              = local.redis_bootstrap

  tags = {
    Name = "${var.project_name}-redis"
  }
}
