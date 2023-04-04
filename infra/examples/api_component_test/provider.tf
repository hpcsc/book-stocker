provider "aws" {
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    iam      = "http://fake-aws:5000"
    sqs      = "http://fake-aws:5000"
    sns      = "http://fake-aws:5000"
    dynamodb = "http://fake-aws:5000"
  }
}
