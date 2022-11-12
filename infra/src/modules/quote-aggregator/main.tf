resource "aws_sqs_queue" "main" {
  for_each = var.topics
  name     = "${var.environment}-quote-aggregator-${each.key}"
}

locals {
  queue_arn_by_name = { for q in aws_sqs_queue.main : q.name => q.arn }
}

resource "aws_sns_topic_subscription" "main" {
  for_each  = var.topics
  topic_arn = each.value
  protocol  = "sqs"
  endpoint  = local.queue_arn_by_name["${var.environment}-quote-aggregator-${each.key}"]
}
