data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }

}

data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "tls_private_key" "deploy" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "deploy" {
  key_name   = "${var.project_name}-deploy-key"
  public_key = tls_private_key.deploy.public_key_openssh
}

resource "aws_security_group" "ssh" {
  name        = "${var.project_name}-ssh-sg"
  description = "SSH access for administration/deploy"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-ssh-sg"
  }
}

# SG compartilhado entre as 5 EC2 para tráfego interno (Kafka e Redis),
# usando regras self-referenced — não depende de IPs fixos.
resource "aws_security_group" "internal" {
  name        = "${var.project_name}-internal-sg"
  description = "Internal traffic between the project EC2 instances (Kafka, Redis)"
  vpc_id      = data.aws_vpc.default.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-internal-sg"
  }
}

resource "aws_security_group_rule" "internal_redis" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  security_group_id        = aws_security_group.internal.id
  source_security_group_id = aws_security_group.internal.id
}

resource "aws_security_group_rule" "internal_kafka_broker" {
  type                     = "ingress"
  from_port                = 9092
  to_port                  = 9092
  protocol                 = "tcp"
  security_group_id        = aws_security_group.internal.id
  source_security_group_id = aws_security_group.internal.id
}

resource "aws_security_group_rule" "internal_kafka_external" {
  type                     = "ingress"
  from_port                = 9094
  to_port                  = 9094
  protocol                 = "tcp"
  security_group_id        = aws_security_group.internal.id
  source_security_group_id = aws_security_group.internal.id
}

resource "aws_security_group" "kafka_ui" {
  name        = "${var.project_name}-kafka-ui-sg"
  description = "Access to the Kafka UI dashboard"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "Kafka UI"
    from_port   = 8082
    to_port     = 8082
    protocol    = "tcp"
    cidr_blocks = [var.admin_cidr]
  }

  tags = {
    Name = "${var.project_name}-kafka-ui-sg"
  }
}

resource "aws_security_group" "alb" {
  name        = "${var.project_name}-alb-sg"
  description = "Public HTTP access to the ALB"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "HTTP from anywhere"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-alb-sg"
  }
}

resource "aws_security_group" "app_server" {
  name_prefix = "${var.project_name}-app-server-sg-"
  description = "Server API - only accessible from the ALB"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "Traffic from ALB"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  tags = {
    Name = "${var.project_name}-app-server-sg"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "app_ai_server" {
  name        = "${var.project_name}-app-ai-server-sg"
  description = "Public access to the ai-server"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 8070
    to_port     = 8070
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-app-ai-server-sg"
  }
}

resource "aws_security_group" "app_socket_server" {
  name        = "${var.project_name}-app-socket-server-sg"
  description = "Public access to the socket-server"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 8765
    to_port     = 8765
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-app-socket-server-sg"
  }
}

resource "aws_security_group" "web" {
  name        = "${var.project_name}-web-sg"
  description = "Public HTTP access to the web app"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-web-sg"
  }
}

resource "aws_security_group" "rds" {
  name        = "${var.project_name}-rds-sg"
  description = "Allows access to the RDS Postgres instance"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "Postgres from the server EC2"
    from_port        = 5432
    to_port          = 5432
    protocol         = "tcp"
    security_groups  = [aws_security_group.app_server.id]
  }

  ingress {
    description = "Postgres from your local IP (DBeaver, TablePlus, etc)"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.admin_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-rds-sg"
  }
}
