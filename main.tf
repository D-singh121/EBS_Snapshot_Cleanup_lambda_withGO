# main.tf

##############################
# 🔧 AWS Provider Configuration
##############################
provider "aws" {
  region = "ap-south-1"
}

###################################
# 📦 S3 Bucket to Store Lambda Code
###################################
resource "random_id" "bucket_suffix" {
  byte_length = 4
}
# create s3 bucket with random suffix for uniqueness of name globally on aws
resource "aws_s3_bucket" "lambda_bucket" {
  bucket        = "ebs-snapshot-cleanup-code-bucket-${random_id.bucket_suffix.hex}"
  force_destroy = true
}
# upload zip file to the bucket
resource "aws_s3_object" "lambda_zip" {
  bucket = aws_s3_bucket.lambda_bucket.id
  key    = "lambda-function.zip"
  source = "lambda-function.zip"
  etag   = filemd5("lambda-function.zip")
}

###############################################
# 🔐 IAM Role and Policy for Lambda Execution
###############################################
# creating an IAM role
resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda-ebs-cleanup-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# attaching policy to the role
resource "aws_iam_role_policy" "lambda_policy" {
  name = "lambda-ebs-policy"
  role = aws_iam_role.lambda_exec_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ec2:DescribeSnapshots",
          "ec2:DescribeInstances",
          "ec2:DescribeVolumes",
          "ec2:DeleteSnapshot"
        ],
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow",
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject"
        ],
        Resource = "*"
      },
    ]
  })
}

###################################
# 🧠 Lambda Function Configuration
###################################
resource "aws_lambda_function" "ebs_snapshot_cleanup" {
  function_name    = "ebs-snapshot-cleanup"
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  s3_bucket        = aws_s3_bucket.lambda_bucket.id
  s3_key           = aws_s3_object.lambda_zip.key
  source_code_hash = filebase64sha256("lambda-function.zip")
  timeout          = 15
  publish          = true
}

###################################
# ⏰ Schedule Lambda with EventBridge
###################################
# Schedule Rule to Run Lambda Daily
resource "aws_cloudwatch_event_rule" "daily_schedule" {
  name                = "daily-ebs-cleanup"
  schedule_expression = "rate(5 minutes)"
  #   schedule_expression = "rate(1 day)"
}

# Attach Lambda function as Target to the Event Rule
resource "aws_cloudwatch_event_target" "trigger_lambda" {
  rule      = aws_cloudwatch_event_rule.daily_schedule.name
  target_id = "ebsSnapshotCleanupLambda" # It is just a unique identifier for the target within the context of the CloudWatch rule.
  arn       = aws_lambda_function.ebs_snapshot_cleanup.arn
}

###################################################
# ✅ Allow EventBridge to Trigger the Lambda Function
###################################################
resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ebs_snapshot_cleanup.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.daily_schedule.arn
}
