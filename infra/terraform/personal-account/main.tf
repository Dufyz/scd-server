data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_security_group" "rds" {
  name        = "${var.project_name}-rds-sg"
  description = "Allows access to the RDS Postgres instance from the AWS Academy account EC2 instances"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "Postgres from the Academy account EC2 instances"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidr_blocks
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

resource "aws_db_subnet_group" "default" {
  name       = "${var.project_name}-rds-subnet-group"
  subnet_ids = data.aws_subnets.default.ids

  tags = {
    Name = "${var.project_name}-rds-subnet-group"
  }
}

resource "aws_db_instance" "postgres" {
  identifier     = "${var.project_name}-postgres"
  engine         = "postgres"
  engine_version = "16"

  instance_class    = var.db_instance_class
  allocated_storage = var.db_allocated_storage_gb
  storage_type      = "gp3"

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  multi_az            = false
  publicly_accessible = true

  db_subnet_group_name   = aws_db_subnet_group.default.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = 1
  skip_final_snapshot     = true
  deletion_protection     = false

  tags = {
    Name = "${var.project_name}-postgres"
  }
}

# Read replica: réplica somente leitura, replicada de forma assíncrona a
# partir do primary acima. Endpoint separado, usado pelo server para SELECTs.
resource "aws_db_instance" "postgres_read_replica" {
  identifier          = "${var.project_name}-postgres-read-replica"
  replicate_source_db = aws_db_instance.postgres.identifier

  instance_class = var.db_replica_instance_class
  storage_type   = "gp3"

  publicly_accessible = true

  vpc_security_group_ids = [aws_security_group.rds.id]

  skip_final_snapshot = true
  deletion_protection = false

  tags = {
    Name = "${var.project_name}-postgres-read-replica"
  }
}
