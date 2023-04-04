resource "aws_dynamodb_table" "stock_requests" {
  name         = "${var.environment}-stock-requests"
  hash_key     = "Id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "Id"
    type = "S"
  }
}
