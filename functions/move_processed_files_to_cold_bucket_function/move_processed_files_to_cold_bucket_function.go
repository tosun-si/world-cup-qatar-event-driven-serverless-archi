package move_processed_files_to_cold_bucket_function

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

type AuditLogEntry struct {
	ProtoPayload *AuditLogProtoPayload `json:"protoPayload"`
	LogResource  *AuditLogResource     `json:"resource"`
}

// AuditLogProtoPayload represents AuditLog within the LogEntry.protoPayload
// See https://cloud.google.com/logging/docs/reference/audit/auditlog/rest/Shared.Types/AuditLog
type AuditLogProtoPayload struct {
	MethodName         string                   `json:"methodName"`
	ResourceName       string                   `json:"resourceName"`
	AuthenticationInfo map[string]interface{}   `json:"authenticationInfo"`
	Metadata           *AuditLogPayloadMetadata `json:"metadata"`
}

type AuditLogResource struct {
	Type   string                 `json:"type"`
	Labels map[string]interface{} `json:"labels"`
}

type AuditLogPayloadMetadata struct {
	Type            string                 `json:"@type"`
	TableDataChange map[string]interface{} `json:"tableDataChange"`
}

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("MoveProcessedFileToColdBucket", moveProcessedFilesToColdBucket)
}

// Function moveProcessedFilesToColdBucket accepts and handles a CloudEvent object
func moveProcessedFilesToColdBucket(ctx context.Context, e event.Event) error {
	ExpectedDataset := "qatar_fifa_world_cup"
	ExpectedTable := "tables/world_cup_team_players_stat"

	//rawSourceBucket := "event-driven-services-qatar-fifa-world-cup-stats-raw"
	//rawSourceObject := "input/stats/world_cup_team_players_stats_raw_ndjson.json"
	//domainSourceBucket := "event-driven-services-qatar-fifa-world-cup-stats"
	//domainSourceObject := "input/stats/world_cup_team_players_stats_domain.json"
	//
	//DestBucket := "event-driven-services-qatar-fifa-world-cup-stats-cold"
	//RawDestObject := "input/raw/world_cup_team_players_stats_raw_ndjson.json"
	//DomainDestObject := "input/domain/world_cup_team_players_stats_domain.json"

	log.Printf("Event Type: %s", e.Type())
	log.Printf("Subject: %s", e.Subject())

	// Decode the Cloud Audit Logging message embedded in the CloudEvent
	logentry := &AuditLogEntry{}
	if err := e.DataAs(logentry); err != nil {
		ferr := fmt.Errorf("event.DataAs: %w", err)
		log.Print(ferr)
		return ferr
	}

	log.Printf("API Method: %s", logentry.ProtoPayload.MethodName)
	log.Printf("Resource Name: %s", logentry.ProtoPayload.ResourceName)
	if v, ok := logentry.ProtoPayload.AuthenticationInfo["principalEmail"]; ok {
		log.Printf("Principal: %s", v)
	}

	//gcpProjectId := logentry.LogResource.Labels["project_id"]
	bqDataset := logentry.LogResource.Labels["dataset_id"]
	bqTable := logentry.ProtoPayload.ResourceName

	insertedRowBqAsString := logentry.ProtoPayload.Metadata.TableDataChange["insertedRowsCount"].(string)

	insertedRowBq, _ := strconv.Atoi(insertedRowBqAsString)

	if bqDataset == ExpectedDataset &&
		strings.HasSuffix(bqTable, ExpectedTable) &&
		insertedRowBq > 0 {

		// Apply the logic here.
		log.Printf("########## The logic will be invoked #############")
	}

	return nil
}
