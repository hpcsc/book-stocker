resource "aws_sqs_queue" "command_queue" {
  name = "${var.environment}-${var.name}-request-quote-command"
}

resource "aws_sns_topic_subscription" "command_queue" {
  topic_arn = var.topic_arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.command_queue.arn
}

resource "aws_sns_topic" "quote_available" {
  name = "${var.environment}-${var.name}-quote-available"
}
