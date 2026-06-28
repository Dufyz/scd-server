variable "aws_region" {
  description = "Região AWS onde tudo será criado"
  type        = string
  default     = "sa-east-1"
}

variable "aws_profile" {
  description = "Profile do ~/.aws/credentials com as credenciais da sua conta AWS"
  type        = string
  default     = "personal"
}

variable "project_name" {
  description = "Prefixo usado no nome dos recursos"
  type        = string
  default     = "scd"
}

variable "instance_type" {
  description = "Tipo de instância EC2 (free tier)"
  type        = string
  default     = "t3.micro"
}

variable "kafka_instance_type" {
  description = "Tipo de instância EC2 do Kafka (precisa de pelo menos 2 GB RAM)"
  type        = string
  default     = "t3.small"
}

variable "server_instance_count" {
  description = "Number of server EC2 instances behind the ALB (min 2 for HA)"
  type        = number
  default     = 2
}

variable "admin_cidr" {
  description = "CIDR (seu IP, ex: \"200.10.20.30/32\") autorizado a conectar via SSH, no painel do Kafka UI e direto no RDS (ex: DBeaver)"
  type        = string
}

variable "repo_url" {
  description = "URL do repositório git (usado pelo user_data para clonar infra/kafka e infra/redis)"
  type        = string
  default     = "https://github.com/Dufyz/scd-server.git"
}

variable "repo_branch" {
  description = "Branch que o user_data deve clonar"
  type        = string
  default     = "main"
}

variable "root_volume_size_gb" {
  description = "Tamanho do volume raiz (EBS) de cada EC2, em GB"
  type        = number
  default     = 8
}

variable "db_name" {
  description = "Nome do banco de dados Postgres"
  type        = string
  default     = "postgres"
}

variable "db_username" {
  description = "Usuário administrador do Postgres"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Senha do usuário administrador do Postgres"
  type        = string
  sensitive   = true
}

variable "db_instance_class" {
  description = "Classe da instância RDS do primary (db.t3.micro é elegível ao free tier)"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage_gb" {
  description = "Armazenamento alocado em GB"
  type        = number
  default     = 20
}

variable "db_replica_instance_class" {
  description = "Classe da instância RDS da read replica"
  type        = string
  default     = "db.t3.micro"
}
