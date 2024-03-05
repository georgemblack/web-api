# Web API

The Web API serves two primary functions:

- Read/write data to Google Cloud Firestore
- Trigger builds (using Web Build service)

To start in Codespaces (until I figure out how to do this automatically):

```
export GOOGLE_APPLICATION_CREDENTIALS=/workspaces/web-api/google-application-credentials.json
echo "$GOOGLE_APPLICATION_CREDENTIALS_CONTENTS" > google-application-credentials.json
```

## Development

Development is done via GitHub Codespaces.

## Environment

The following env vars are also used (and set to dummy values by default):

```
USERNAME          # for auth
PASSWORD          # for auth
```

## Infrastructure

The Web API runs as a service on Google Cloud Run, and must be given an IAM role that can:

- Read/write to Firestore
- Invoke Web Build service
