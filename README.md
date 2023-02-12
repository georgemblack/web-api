# Web API

The Web API serves two primary functions:

* Read/write data to Google Cloud Firestore
* Trigger builds (using Web Build service)

Environment-based can be found in the `config` dir.

## Development

Development is done via GitHub Codespaces.

## Environment

The following env vars are also used (and set to dummy values by default):

```
USERNAME          # for auth
PASSWORD          # for auth
TOKEN_SECRET
```

## Infrastructure

The Web API runs as a service on Google Cloud Run, and must be given an IAM role that can:

* Read/write to Firestore
* Invoke Web Build service
