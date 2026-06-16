variable "aws_region" {
  description = "Região AWS onde o RDS será criado"
  type        = string
  default     = "sa-east-1"
}

variable "aws_profile" {
  description = "Profile do ~/.aws/credentials com as credenciais da sua conta pessoal AWS"
  type        = string
  default     = "personal"
}

variable "project_name" {
  description = "Prefixo usado no nome dos recursos"
  type        = string
  default     = "scd"
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

# IPs públicos das EC2 da conta AWS Academy que podem acessar o RDS.
# Como a sessão da Academy é temporária, esses IPs mudam a cada novo lab —
# atualize esta lista e rode `terraform apply` de novo quando isso acontecer.
variable "allowed_cidr_blocks" {
  description = "Lista de CIDRs (ex: \"3.91.10.20/32\") autorizados a conectar no RDS"
  type        = list(string)
}
