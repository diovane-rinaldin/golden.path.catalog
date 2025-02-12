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
#Exemplo de saída
#dynamodb_components_table_arn = "arn:aws:dynamodb:us-east-1:<account-id>:table/golden-path_components"
#dynamodb_endpoint = "https://dynamodb.us-east-1.amazonaws.com"
#dynamodb_services_table_arn = "arn:aws:dynamodb:us-east-1:<account-id>:table/golden-path_services"
#dynamodb_technology_table_arn = "arn:aws:dynamodb:us-east-1:<account-id>:table/golden-path_technology"
#kms_key_arn = "arn:aws:kms:us-east-1:<account-id>:key/<key>"
#s3_bucket_url = "https://golden-path-images.s3.us-east-1.amazonaws.com"