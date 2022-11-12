resource "aws_dynamodb_table" "stock_requests" {
  name         = "${var.environment}-stock-requests"
  hash_key     = "Id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "Id"
    type = "S"
  }
}

resource "aws_sns_topic" "request_quote_command" {
  name = "${var.environment}-request-quote-command"
}

module "bookdepository_stocker" {
  source = "./modules/stocker"

  name        = "bookdepository"
  environment = var.environment
  topic_arn   = aws_sns_topic.request_quote_command.arn
}

module "booktopia_stocker" {
  source = "./modules/stocker"

  name        = "booktopia"
  environment = var.environment
  topic_arn   = aws_sns_topic.request_quote_command.arn
}

module "quote_aggregator" {
  source = "./modules/quote-aggregator"

  environment = var.environment
  topics = {
    "bookdepository" = module.bookdepository_stocker.quote_available_topic_arn,
    "booktopia"      = module.booktopia_stocker.quote_available_topic_arn
  }
}
