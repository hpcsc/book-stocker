output "stock_requests_arn" {
  value = aws_dynamodb_table.stock_requests.arn
}

output "request_quote_command_topic_arn" {
  value = aws_sns_topic.request_quote_command.arn
}

output "bookdepository_stocker_queue_url" {
  value = module.bookdepository_stocker.command_queue_url
}

output "bookdepository_quote_available_topic_arn" {
  value = module.bookdepository_stocker.quote_available_topic_arn
}

output "booktopia_stocker_queue_url" {
  value = module.booktopia_stocker.command_queue_url
}

output "booktopia_quote_available_topic_arn" {
  value = module.booktopia_stocker.quote_available_topic_arn
}

output "quote_aggregator_queue_urls" {
  value = module.quote_aggregator.queue_urls
}
