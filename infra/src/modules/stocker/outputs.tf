output "command_queue_url" {
  value = aws_sqs_queue.command_queue.url
}

output "quote_available_topic_arn" {
  value = aws_sns_topic.quote_available.arn
}
