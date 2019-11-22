#!/bin/sh

set -e

sanitize() {
  if [ -z "${1}" ]
  then
    >&2 echo "Unable to find ${2}. Did you configure your workflow correctly?"
    exit 1
  fi
}

sanitize "${INPUT_IMAGENAME}" "image-name"
sanitize "${INPUT_GCLOUDPROJECTID}" "gcloud-project-id"
sanitize "${GCLOUD_AUTH}" "gcloud-auth"

cd ${GITHUB_WORKSPACE}

# Set project
gcloud config set project ${INPUT_GCLOUDPROJECTID}

# Auth w/service account
echo ${GCLOUD_AUTH} | base64 --decode > ./key.json
gcloud auth activate-service-account --key-file=./key.json
rm ./key.json

# Submit build
gcloud builds submit --tag gcr.io/${INPUT_GCLOUDPROJECTID}/${INPUT_IMAGENAME}:latest
