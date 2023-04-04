output "queue_urls" {
  value = [for q in aws_sqs_queue.main : q.url]
}
