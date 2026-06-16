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
