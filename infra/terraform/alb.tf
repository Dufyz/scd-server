resource "aws_lb" "server" {
  name               = "${var.project_name}-server-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = data.aws_subnets.default.ids

  tags = {
    Name = "${var.project_name}-server-alb"
  }
}

resource "aws_lb_target_group" "server" {
  name     = "${var.project_name}-server-tg"
  port     = 3000
  protocol = "HTTP"
  vpc_id   = data.aws_vpc.default.id

  health_check {
    path                = "/health"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 15
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  tags = {
    Name = "${var.project_name}-server-tg"
  }
}

resource "aws_lb_target_group_attachment" "server" {
  count            = var.server_instance_count
  target_group_arn = aws_lb_target_group.server.arn
  target_id        = aws_instance.server[count.index].id
  port             = 3000
}

resource "aws_lb_listener" "server_http" {
  load_balancer_arn = aws_lb.server.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.server.arn
  }
}
