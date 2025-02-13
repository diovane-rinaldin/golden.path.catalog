provider "aws" {
  region = var.aws_region
}

resource "aws_iam_user" "service_user" {
  name = "${var.project_name}-service-user"
}

resource "aws_iam_user_policy_attachment" "service_user_dynamodb" {
  user       = aws_iam_user.service_user.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}

resource "aws_iam_user_policy_attachment" "service_user_s3" {
  user       = aws_iam_user.service_user.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_access_key" "service_user_key" {
  user = aws_iam_user.service_user.name
}

# DynamoDB Tables
resource "aws_dynamodb_table" "technology" {
  name           = "${var.project_name}-technology"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
  
  attribute {
    name = "id"
    type = "S"
  }
  
  attribute {
    name = "name"
    type = "S"
  }
  
  global_secondary_index {
    name               = "name-index"
    hash_key           = "name"
    projection_type    = "ALL"
  }
}

resource "aws_dynamodb_table" "components" {
  name           = "${var.project_name}-components"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
  
  attribute {
    name = "id"
    type = "S"
  }
  
  attribute {
    name = "technology_id"
    type = "S"
  }
  
  attribute {
    name = "name"
    type = "S"
  }
  
  global_secondary_index {
    name               = "technology-index"
    hash_key           = "technology_id"
    projection_type    = "ALL"
  }
  
  global_secondary_index {
    name               = "name-index"
    hash_key           = "name"
    projection_type    = "ALL"
  }
}

resource "aws_dynamodb_table" "services" {
  name           = "${var.project_name}-services"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
  
  attribute {
    name = "id"
    type = "S"
  }
  
  attribute {
    name = "component_id"
    type = "S"
  }
  
  attribute {
    name = "name"
    type = "S"
  }
  
  global_secondary_index {
    name               = "component-index"
    hash_key           = "component_id"
    projection_type    = "ALL"
  }
  
  global_secondary_index {
    name               = "name-index"
    hash_key           = "name"
    projection_type    = "ALL"
  }
}

# S3 Bucket
resource "aws_s3_bucket" "images" {
  bucket = var.bucket_name
}

resource "aws_s3_bucket_public_access_block" "images" {
  bucket = aws_s3_bucket.images.id
  
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "public_read" {
  bucket = aws_s3_bucket.images.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.images.arn}/*"
      },
    ]
  })
}

# KMS Key for API Authentication
resource "aws_kms_key" "api_auth" {
  description             = var.kms_description
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_kms_alias" "api_auth" {
  name          = "alias/${var.project_name}-api-auth"
  target_key_id = aws_kms_key.api_auth.key_id
}

# IAM Role for Backend Service
resource "aws_iam_role" "backend_service" {
  name = "${var.project_name}-backend-service"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

# IAM Policy for Backend Service
resource "aws_iam_role_policy" "backend_service" {
  name = "${var.project_name}-backend-service-policy"
  role = aws_iam_role.backend_service.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          aws_dynamodb_table.technology.arn,
          aws_dynamodb_table.components.arn,
          aws_dynamodb_table.services.arn,
          "${aws_dynamodb_table.technology.arn}/index/*",
          "${aws_dynamodb_table.components.arn}/index/*",
          "${aws_dynamodb_table.services.arn}/index/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.images.arn}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:GenerateDataKey"
        ]
        Resource = aws_kms_key.api_auth.arn
      }
    ]
  })
}