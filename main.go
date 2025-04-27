package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
)

func handleRequest(ctx context.Context) error {
	// Load AWS SDK configuration

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %v", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Get all EBS snapshots owned by current account
	snapshotsResp, err := ec2Client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		OwnerIds: []string{"self"},
	})
	if err != nil {
		return fmt.Errorf("failed to describe snapshots: %v", err)
	}

	// Get all active EC2 instance IDs
	instancesResp, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to describe instances: %v", err)
	}

	// Create a set of active instance IDs
	activeInstanceIDs := make(map[string]struct{})
	for _, reservation := range instancesResp.Reservations {
		for _, instance := range reservation.Instances {
			activeInstanceIDs[*instance.InstanceId] = struct{}{}
		}
	}

	// Iterate through each snapshot and check if it should be deleted
	for _, snapshot := range snapshotsResp.Snapshots {
		snapshotID := *snapshot.SnapshotId

		// Check if snapshot is not attached to any volume
		if snapshot.VolumeId == nil || *snapshot.VolumeId == "" {
			// Delete the snapshot
			_, err := ec2Client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
				SnapshotId: &snapshotID,
			})
			if err != nil {
				log.Printf("Failed to delete snapshot %s: %v", snapshotID, err)
				continue
			}
			fmt.Printf("Deleted EBS snapshot %s as it was not attached to any volume.\n", snapshotID)
		} else {
			// Check if the volume still exists
			volumeID := *snapshot.VolumeId
			volumeResp, err := ec2Client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
				VolumeIds: []string{volumeID},
			})

			if err != nil {
				// Check if the error is because the volume was not found
				var apiErr smithy.APIError
				if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidVolume.NotFound" {
					// The volume associated with the snapshot is not found (it might have been deleted)
					_, err := ec2Client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
						SnapshotId: &snapshotID,
					})
					if err != nil {
						log.Printf("Failed to delete snapshot %s: %v", snapshotID, err)
						continue
					}
					fmt.Printf("Deleted EBS snapshot %s as its associated volume was not found.\n", snapshotID)
				} else {
					log.Printf("Error describing volume %s: %v", volumeID, err)
				}
				continue
			}

			// Check if the volume is not attached to any running instance
			if len(volumeResp.Volumes) > 0 && len(volumeResp.Volumes[0].Attachments) == 0 {
				_, err := ec2Client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
					SnapshotId: &snapshotID,
				})
				if err != nil {
					log.Printf("Failed to delete snapshot %s: %v", snapshotID, err)
					continue
				}
				fmt.Printf("Deleted EBS snapshot %s as it was taken from a volume not attached to any running instance.\n", snapshotID)
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
