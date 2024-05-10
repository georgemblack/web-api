# Web API

The Web API serves two primary functions:

- Read/write data to Google Cloud Firestore
- Trigger builds (using Web Build service)

To start in Codespaces (until I figure out how to do this automatically):

```
export GOOGLE_APPLICATION_CREDENTIALS=/workspaces/web-api/google-application-credentials.json
echo "$GOOGLE_APPLICATION_CREDENTIALS_CONTENTS" > google-application-credentials.json
```
