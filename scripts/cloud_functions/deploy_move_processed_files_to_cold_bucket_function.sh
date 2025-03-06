#!/usr/bin/env bash

set -e
set -o pipefail
set -u

echo "############# Deploying the Cloud Function qatar-world-cup-move-processed-files-to-cold-bucket"

gcloud functions deploy qatar-world-cup-move-processed-files-to-cold-bucket \
  --quiet \
  --gen2 \
  --region=europe-west1 \
  --runtime=go121 \
  --source=functions/move_processed_files_to_cold_bucket_function \
  --entry-point=MoveProcessedFileToColdBucket \
  --run-service-account="$SERVICE_ACCOUNT" \
  --trigger-event-filters="type=google.cloud.audit.log.v1.written" \
  --trigger-event-filters="serviceName=bigquery.googleapis.com" \
  --trigger-event-filters="methodName=google.cloud.bigquery.v2.JobService.InsertJob" \
  --trigger-event-filters="resourceName=projects/${PROJECT_ID}/datasets/qatar_fifa_world_cup/tables/world_cup_team_players_stat" \
  --trigger-location=europe-west1 \
  --trigger-service-account="$SERVICE_ACCOUNT"
