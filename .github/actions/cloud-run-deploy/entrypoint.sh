#!/bin/sh

set -e

sanitize() {
  if [ -z "${1}" ]
  then
    >&2 echo "Unable to find ${2}. Did you configure your workflow correctly?"
    exit 1
  fi
}

sanitize "${INPUT_SERVICENAME}" "image-name"
sanitize "${INPUT_GCLOUDPROJECTID}" "gcloudProjectId"
sanitize "${INPUT_GCLOUDSERVICEACCOUNT}" "gcloudServiceAccount"
sanitize "${INPUT_GCLOUDRUNTIMESERVICEACCOUNT}" "gcloudRuntimeServiceAccount"
sanitize "${GCLOUD_AUTH}" "GCLOUD_AUTH"

# Get version from package.json
PACKAGE_VERSION=$(cat package.json \
  | grep version \
  | head -1 \
  | awk -F: '{ print $2 }' \
  | sed 's/[",]//g' \
  | tr -d '[[:space:]]')

# Set project
gcloud config set project ${INPUT_GCLOUDPROJECTID}

# Auth w/service account
echo ${GCLOUD_AUTH} | base64 --decode > ./key.json
gcloud auth activate-service-account --key-file=./key.json
rm ./key.json

# Submit build
gcloud builds submit --tag gcr.io/${INPUT_GCLOUDPROJECTID}/${INPUT_SERVICENAME}:${PACKAGE_VERSION}

# Deploy
gcloud run deploy ${INPUT_SERVICENAME} \
  --concurrency 20 \
  --max-instances 1000 \
  --memory 128Mi \
  --platform managed \
  --allow-unauthenticated \
  --service-account ${INPUT_GCLOUDRUNTIMESERVICEACCOUNT} \
  --region us-east1 \
  --image gcr.io/${INPUT_GCLOUDPROJECTID}/${INPUT_SERVICENAME}:${PACKAGE_VERSION}