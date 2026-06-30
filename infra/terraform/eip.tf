# Elastic IPs fixos — sobrevivem a reinícios e recriações de instância.
# Limite padrão da AWS: 5 EIPs por região. Este stack usa os 5 para serviços
# expostos (server×2, web, socket-server, ai-server). Kafka e Redis ficam com
# IP efêmero até o aumento de quota (solicitado para 10).

resource "aws_eip" "server" {
  count  = var.server_instance_count
  domain = "vpc"

  tags = {
    Name = "${var.project_name}-server-${count.index}-eip"
  }
}

resource "aws_eip_association" "server" {
  count         = var.server_instance_count
  instance_id   = aws_instance.server[count.index].id
  allocation_id = aws_eip.server[count.index].id
}

resource "aws_eip" "ai_server" {
  domain = "vpc"

  tags = {
    Name = "${var.project_name}-ai-server-eip"
  }
}

resource "aws_eip_association" "ai_server" {
  instance_id   = aws_instance.ai_server.id
  allocation_id = aws_eip.ai_server.id
}

resource "aws_eip" "socket_server" {
  domain = "vpc"

  tags = {
    Name = "${var.project_name}-socket-server-eip"
  }
}

resource "aws_eip_association" "socket_server" {
  instance_id   = aws_instance.socket_server.id
  allocation_id = aws_eip.socket_server.id
}

resource "aws_eip" "web" {
  domain = "vpc"

  tags = {
    Name = "${var.project_name}-web-eip"
  }
}

resource "aws_eip_association" "web" {
  instance_id   = aws_instance.web.id
  allocation_id = aws_eip.web.id
}
