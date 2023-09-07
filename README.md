# teams-league-cloudrun-service-fastapi

Project showing a complete use case with a Cloud Run Service written with a Python module and multiple files. The
deployment of service is done with FastApi and Uvicorn.

## Use case with Cloud Functions

### Deploy the Cloud Function that map raw data to domain data and upload the result to GCS

```bash
gcloud functions deploy qatar-world-cup-stats-raw-to-domain-data-gcs \
  --gen2 \
  --region=europe-west1 \
  --runtime=python310 \
  --source=functions/world_cup_stats_raw_to_domain_function \
  --entry-point=raw_to_domain_data_and_upload_to_gcs \
  --run-service-account=sa-cloud-functions-dev@gb-poc-373711.iam.gserviceaccount.com \
  --trigger-event-filters="type=google.cloud.storage.object.v1.finalized" \
  --trigger-event-filters="bucket=event-driven-functions-qatar-fifa-world-cup-stats-raw" \
  --trigger-location=europe-west1 \
  --trigger-service-account=sa-cloud-functions-dev@gb-poc-373711.iam.gserviceaccount.com
```

### Deploy the Cloud Function that add Fifa ranking to each team stats in the domain data

```bash
gcloud functions deploy qatar-world-cup-add-fifa-ranking-to-stats-domain-bq \
  --gen2 \
  --region=europe-west1 \
  --runtime=python310 \
  --source=functions/world_cup_team_add_fifa_ranking_function \
  --entry-point=add_fifa_ranking_to_stats_domain_and_save_to_bq \
  --run-service-account=sa-cloud-functions-dev@gb-poc-373711.iam.gserviceaccount.com \
  --trigger-event-filters="type=google.cloud.storage.object.v1.finalized" \
  --trigger-event-filters="bucket=event-driven-functions-qatar-fifa-world-cup-stats" \
  --trigger-location=europe-west1 \
  --trigger-service-account=sa-cloud-functions-dev@gb-poc-373711.iam.gserviceaccount.com
```

## Use case with Cloud Run Services

### Map raw to domain Data

#### Build Docker image

```bash
docker build -f services/world_cup_stats_raw_to_domain_service/Dockerfile -t $SERVICE_NAME_RAW_TO_DOMAIN .
docker tag $SERVICE_NAME_RAW_TO_DOMAIN $LOCATION-docker.pkg.dev/$PROJECT_ID/$REPO_NAME/$SERVICE_NAME_RAW_TO_DOMAIN:$IMAGE_TAG
docker push $LOCATION-docker.pkg.dev/$PROJECT_ID/$REPO_NAME/$SERVICE_NAME_RAW_TO_DOMAIN:$IMAGE_TAG
```

#### Deploy service

```bash
gcloud run deploy "$SERVICE_NAME_RAW_TO_DOMAIN" \
  --image "$LOCATION-docker.pkg.dev/$PROJECT_ID/$REPO_NAME/$SERVICE_NAME_RAW_TO_DOMAIN:$IMAGE_TAG" \
  --region="$LOCATION" \
  --no-allow-unauthenticated \
  --set-env-vars PROJECT_ID="$PROJECT_ID"
```

#### Create Event Arc Trigger on Service

```bash
gcloud eventarc triggers create "$SERVICE_NAME_RAW_TO_DOMAIN" \
  --destination-run-service="$SERVICE_NAME_RAW_TO_DOMAIN" \
  --destination-run-region="$LOCATION" \
  --location="$LOCATION" \
  --event-filters="type=google.cloud.storage.object.v1.finalized" \
  --event-filters="bucket=event-driven-services-qatar-fifa-world-cup-stats-raw" \
  --service-account=sa-cloud-run-dev@gb-poc-373711.iam.gserviceaccount.com
```

#### Deploy the service and create the Event Arc Trigger with Cloud Build

```bash
gcloud builds submit \
  --project=$PROJECT_ID \
  --region=$LOCATION \
  --config deploy-cloud-run-service.yaml \
  --substitutions _REPO_NAME="$REPO_NAME",_SERVICE_NAME_RAW_TO_DOMAIN="$SERVICE_NAME_RAW_TO_DOMAIN",_IMAGE_TAG="$IMAGE_TAG" \
  --verbosity="debug" .
```