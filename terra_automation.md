
# 🧹 EBS Snapshot Cleanup Automation with AWS Lambda

This Terraform module provisions an AWS Lambda function that automatically deletes old EBS snapshots. The Lambda function is packaged as a ZIP file and scheduled using Amazon EventBridge (CloudWatch Events) to run periodically. This setup includes secure storage, IAM roles, and necessary permissions.

---

## 📁 Project Structure

```
.
├── go.mod 
├── go.sum
├── main.go
├── main.tf                # Terraform configuration file
└── lambda-function.zip    # Pre-built Lambda deployment package
```

---

## 🚀 Features

- Deploys a Lambda function for EBS snapshot cleanup.
- Creates an S3 bucket to store Lambda code.
- Uses a randomized suffix for global uniqueness of the S3 bucket name.
- Attaches necessary IAM policies for Lambda to interact with EC2, logs, and S3.
- Schedules Lambda to run every 5 minutes using EventBridge.
- Grants EventBridge permission to invoke the Lambda function.

---

## 🛠️ Prerequisites

- [Terraform](https://www.terraform.io/downloads)
- AWS CLI configured with appropriate credentials
- A zipped Lambda deployment package (`lambda-function.zip`) present in the project root

---

## 📦 Resources Created

| Resource Type                 | Name                                |
|------------------------------|-------------------------------------|
| AWS Provider                 | ap-south-1                          |
| S3 Bucket                    | ebs-snapshot-cleanup-code-bucket-* |
| S3 Object                    | lambda-function.zip                 |
| IAM Role                     | lambda-ebs-cleanup-role             |
| IAM Role Policy              | lambda-ebs-policy                   |
| Lambda Function              | ebs-snapshot-cleanup                |
| CloudWatch Event Rule        | daily-ebs-cleanup                   |
| CloudWatch Event Target      | ebsSnapshotCleanupLambda            |
| Lambda Permission            | AllowExecutionFromEventBridge       |

---
---

## 🏗️ Build and Package the Lambda Function

Before applying the Terraform configuration, you need to build and zip the Lambda function:

### 1. Build for Amazon Linux (Lambda Environment)

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
```

### 2. Zip the Executable

```bash
zip lambda-function.zip bootstrap
```

> This zip file will be uploaded to S3 and used as the source for the Lambda function.

---

## 📋 Usage

### 1. Initialize Terraform

```bash
terraform init
```

### 2. Validate the Configuration

```bash
terraform validate
```

### 3. Apply the Configuration

```bash
terraform apply
```

When prompted, confirm with `yes`.

---

## 🔐 IAM Policy Details

The Lambda execution role has the following permissions:

- **EC2 Access**: 
  - Describe and delete snapshots
  - Describe volumes and instances

- **CloudWatch Logs**: 
  - Create log groups and streams
  - Publish log events

- **S3 Access**:
  - Upload, retrieve, and delete objects

---

## 📅 Scheduling

The Lambda function is scheduled to run every **5 minutes** by default. You can modify the schedule in the `aws_cloudwatch_event_rule.daily_schedule` resource:

```hcl
schedule_expression = "rate(1 day)"  # For daily runs
```

---

## 🧹 Cleanup

To destroy all resources provisioned by this project:

```bash
terraform destroy
```

---

## 📄 Notes

- Make sure `lambda-function.zip` exists in your project directory before applying.
- The Lambda handler is set to `bootstrap` and assumes a custom runtime (`provided.al2023`) latest for golang.

---