variable "aws_region" {
  description = "Região AWS usada na conta AWS Academy"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "Profile do ~/.aws/credentials com as credenciais temporárias da sessão AWS Academy"
  type        = string
  default     = "academy"
}

variable "project_name" {
  description = "Prefixo usado no nome dos recursos"
  type        = string
  default     = "scd"
}

# AWS Academy normalmente só libera t2.micro/t3.micro nas labs com LabRole.
variable "instance_type" {
  description = "Tipo de instância EC2 (free tier)"
  type        = string
  default     = "t2.micro"
}

variable "admin_cidr" {
  description = "CIDR (seu IP, ex: \"200.10.20.30/32\") autorizado a conectar via SSH e no painel do Kafka UI"
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
