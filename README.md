# Profile micro-service

[![pipeline status](https://gitlab.slade360emr.com/go/profile/badges/develop/pipeline.svg)](https://gitlab.slade360emr.com/go/profile/-/commits/develop)
[![coverage report](https://gitlab.slade360emr.com/go/profile/badges/develop/coverage.svg)](https://gitlab.slade360emr.com/go/profile/-/commits/develop)

## Environment variables

For local development, you need to *export* the following env vars:

```bash
# Google Cloud Settings
export GOOGLE_APPLICATION_CREDENTIALS="<a path to a Google service account JSON file>"
export GOOGLE_CLOUD_PROJECT="<the name of the project that the service account above belongs to>"
export FIREBASE_WEB_API_KEY="<an API key from the Firebase console for the project mentioned above>"

# Mailgun settings
export MAILGUN_API_KEY=key="<an API key>"
export MAILGUN_DOMAIN=app.healthcloud.co.ke
export MAILGUN_FROM=hello@app.healthcloud.co.ke

# AfricasTalking SMS API settings
export AIT_API_KEY="<an API key>"
export AIT_USERNAME=sandbox
export AIT_SENDER_ID=HealthCloud
export AIT_ENVIRONMENT=sandbox

# Go private modules
export GOPRIVATE="gitlab.slade360emr.com/go/*,gitlab.slade360emr.com/optimalhealth/*"
```

The server deploys to Google Cloud Run. The environment variables defined above
should also be set on Google Cloud Run.

