# Web API

For running locally, you will need to set a few env vars.

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

Environment-based configurations are defined in `config` dir. The docker container sets `NODE_ENV` to production to use production config.
