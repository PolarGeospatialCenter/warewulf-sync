package cmd

import (
	"log"
	"os/exec"
	"strings"

	"github.com/PolarGeospatialCenter/inventory-client/pkg/api/client"
	"github.com/PolarGeospatialCenter/warewulf-sync/pkg/warewulf"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

var nodeSyncCmd = &cobra.Command{
	Use:   "node-sync",
	Short: "Watch SQS queue for node events and sync to warewulf",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// load desired state from yaml
		db, err := warewulf.LoadYaml(cfg.GetString("config_path"))
		if err != nil {
			log.Fatalf("Unable to load files data: %v", err)
		}

		log.Print(db)
		// wait for sqs event
		// on event, load db and sync
		//
		// create aws session
		awsConfig := &aws.Config{}
		awsConfig = awsConfig.WithRegion(cfg.GetString("aws.region"))
		sess, err := session.NewSession(awsConfig)
		if err != nil {
			log.Fatalf("unable to create aws session: %v", err)
		}

		// Create a SQS service client.
		svc := sqs.New(sess)
		msgCh := getMsgCh(svc, cfg.GetString("sqs.queue_url"))
		for _ = range msgCh {
			inv, err := client.NewInventoryApiDefaultConfig("")
			if err != nil {
				log.Fatalf("Unable to connecto to inventory API: %v", err)
			}

			err = db.LoadNodesFromInventory(inv.NodeConfig(), cfg.GetString("system"))
			if err != nil {
				log.Fatalf("Unable to load nodes from inventory: %v", err)
			}

			// load warewulf
			wwdb, err := warewulf.LoadWwshDB()
			if err != nil {
				log.Fatalf("Unable to load warewulf database: %v", err)
			}
			log.Print(wwdb)

			syncCommands := make([][]string, 0)
			syncCommands = append(syncCommands, BuildSyncCommands(MakeSyncableMap(wwdb.Nodes), MakeSyncableMap(db.Nodes))...)
			if len(db.Nodes) > 0 {
				syncCommands = append(syncCommands, []string{"wwsh", "pxe", "-v", "--nodhcp"})
			}

			for _, cmd := range syncCommands {
				log.Print(cmd)
				c := exec.Command(cmd[0], cmd[1:]...)

				stdErrOut, err := c.CombinedOutput()
				if err != nil {
					log.Fatalf("Error executing '%s': %v", strings.Join(cmd, " "), err)
				}
				log.Printf("Result: %s", stdErrOut)
			}
		}
	},
}

func getMsgCh(svc *sqs.SQS, queueUrl string) chan sqs.Message {
	msgCh := make(chan sqs.Message)
	go func() {
		defer close(msgCh)
		for {
			result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl: aws.String(queueUrl),
				AttributeNames: aws.StringSlice([]string{
					"SentTimestamp",
				}),
				MaxNumberOfMessages: aws.Int64(1),
				MessageAttributeNames: aws.StringSlice([]string{
					"All",
				}),
				WaitTimeSeconds: aws.Int64(cfg.GetInt64("sqs.timeout")),
			})
			if err != nil {
				log.Fatalf("Unable to receive message from queue: %v", err)
			}

			for _, msg := range result.Messages {
				msgCh <- *msg
				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueUrl),
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					log.Printf("Unable to delete message from queue: %v", err)
				}

			}
		}
	}()
	return msgCh
}
