# ============================================================
# INTENTIONALLY INSECURE TERRAFORM — for Trivy demo purposes
# ============================================================
# This Terraform configuration contains several misconfigurations
# that Trivy will detect. DO NOT use this in production.

provider "aws" {
  region = "eu-central-1"
}

# Flaw 1: S3 bucket without encryption
resource "aws_s3_bucket" "documents" {
  bucket = "trivy-demo-documents"
  acl    = "public-read"  # Flaw 2: Public read access

  tags = {
    Name        = "Document Storage"
    Environment = "demo"
  }
}

# Flaw 3: S3 bucket without versioning
# Flaw 4: S3 bucket without logging

# Flaw 5: Security group with unrestricted ingress
resource "aws_security_group" "app" {
  name        = "trivy-demo-app-sg"
  description = "Security group for demo application"

  # Flaw 6: Allow ALL inbound traffic from anywhere
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all inbound"
  }

  # Flaw 7: Allow ALL outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Flaw 8: RDS without encryption at rest
resource "aws_db_instance" "main" {
  identifier     = "trivy-demo-db"
  engine         = "postgres"
  engine_version = "15.3"
  instance_class = "db.t3.micro"

  allocated_storage = 20
  storage_encrypted = false  # Flaw: No encryption

  db_name  = "trivy_demo"
  username = "admin"
  password = "SuperSecret123!"  # Flaw 9: Hardcoded password

  publicly_accessible = true  # Flaw 10: DB publicly accessible

  skip_final_snapshot = true

  # Flaw 11: No backup retention
  backup_retention_period = 0
}

# Flaw 12: IAM policy with wildcard actions
resource "aws_iam_policy" "app_policy" {
  name        = "trivy-demo-app-policy"
  description = "Application policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "*"
        Resource = "*"
      }
    ]
  })
}

# Flaw 13: CloudWatch log group without encryption
resource "aws_cloudwatch_log_group" "app" {
  name              = "/app/trivy-demo"
  retention_in_days = 7
}
