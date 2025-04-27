# 🧹 AWS EBS Snapshot Cleanup Lambda Function

A serverless Go-based Lambda function that automatically cleans up unused EBS snapshots to optimize AWS costs.

## 📋 Overview

This Lambda function identifies and removes unnecessary EBS snapshots that match any of these criteria:
- 🔍 Snapshots not associated with any volume
- 🗑️ Snapshots whose associated volume no longer exists
- 💾 Snapshots whose associated volume is not attached to any running EC2 instance

By regularly cleaning up these unused snapshots, this function helps reduce unnecessary AWS storage costs.

## ✅ Prerequisites

- Go 1.16 or later
- AWS CLI configured with appropriate permissions
- AWS account with Lambda management permissions
- IAM role with necessary EC2 permissions

## 📁 Project Structure

```
.
├── main.go              # Lambda function code
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── README.md            # This documentation
```

## 🚀 Installation

### 1. Clone the repository

```bash
git clone https://github.com/D-singh121/EBS_Snapshot_Cleanup_lambda_withGO.git
cd Serverless_Cost_optimizer_lambda
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Build the Lambda function

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
```

### 4. Create a deployment package

```bash
zip lambda-function.zip bootstrap
```

## 📦 Deployment

### Using AWS Management Console

1. Open the AWS Lambda console
2. Click "Create function"
3. Select "Author from scratch"
4. Enter a function name (e.g., "ebs-snapshot-cleanup")
5. Select "Provided runtime" from the dropdown
6. Select architecture: "x86_64"
7. Create or select an execution role with required permissions
8. Click "Create function"
9. In the "Code" tab, upload the zip file
10. Configure timeout setting (recommended: at least 15 Sec)

### Using AWS CLI

```bash
aws lambda create-function \
  --function-name ebs-snapshot-cleanup \
  --runtime provided.al2 \
  --handler bootstrap \
  --zip-file fileb://lambda-function.zip \
  --role arn:aws:iam::<your-account-id>:role/<your-lambda-execution-role> \
  --timeout 15
```

## 🔒 IAM Permissions

The Lambda execution role needs these permissions:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeSnapshots",
                "ec2:DescribeInstances",
                "ec2:DescribeVolumes",
                "ec2:DeleteSnapshot"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "arn:aws:logs:*:*:*"
        }
    ]
}
```

## ⏰ Scheduling the Lambda Function

To run the function automatically:

1. Open AWS Lambda console and select your function
2. Click "Add trigger"
3. Select "EventBridge (CloudWatch Events)"
4. Create a new rule with a schedule expression:
   - Daily execution: `rate(1 day)`
   - Weekly execution: `rate(7 days)`
5. Click "Add"

## 📊 Monitoring and Logging

- Function logs all actions to CloudWatch Logs
- View logs in AWS Lambda console or CloudWatch console
- Consider CloudWatch Alarms for issue alerts

## 🛠️ Customization

Potential modifications:
- Add snapshot filtering criteria
- Implement age-based retention policies
- Add tagging support to exclude specific snapshots

## 💰 Cost Considerations

- Lambda charges are minimal for scheduled functions
- Cost savings from deleting unused snapshots outweigh Lambda costs
- Monitor AWS billing to verify savings

## 🔐 Security Best Practices

- Use principle of least privilege for Lambda execution role
- Implement encryption for sensitive data
- Regularly review and update the function

## ❓ Troubleshooting

Common issues:
- **⏱️ Timeout errors**: Increase Lambda timeout for many snapshots
- **🚫 Permission errors**: Verify execution role permissions
- **🌐 Network errors**: Check Lambda network access if using VPC

