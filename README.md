# Web API

This service reads/writes data to Google Cloud Firestore.

* Environment-based configurations are defined in the `config` dir.

For running locally, will also need to set credentials for Firestore:

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

The following env vars are also used (and set to dummy values by default):

```
USERNAME          # for auth
PASSWORD          # for auth
TOKEN_SECRET      
```
