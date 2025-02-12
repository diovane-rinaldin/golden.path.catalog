variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name prefix for resources"
  type        = string
  default     = "golden-path"
}

variable "bucket_name" {
  description = "Name for the S3 bucket"
  type        = string
  default     = "golden-path-images"
}

variable "kms_description" {
  description = "KMS key for API authentication"
  type        = string
  default     = "KMS key for API authentication"
}