output "dynamodb_technology_table_arn" {
  value = aws_dynamodb_table.technology.arn
}

output "dynamodb_components_table_arn" {
  value = aws_dynamodb_table.components.arn
}

output "dynamodb_services_table_arn" {
  value = aws_dynamodb_table.services.arn
}

output "s3_bucket_url" {
  value = "https://${aws_s3_bucket.images.bucket_regional_domain_name}"
}

output "kms_key_arn" {
  value = aws_kms_key.api_auth.arn
}

output "dynamodb_endpoint" {
  value = "https://dynamodb.${var.aws_region}.amazonaws.com"
}