#!/usr/bin/env sh

set -eux

# Create the namespace
kubectl create namespace $NAMESPACE || true

# Delete Kubernetes secret if exists
kubectl delete secret mycarehub-service-account --namespace $NAMESPACE || true

# Create GCP service account file
cat $GOOGLE_APPLICATION_CREDENTIALS >> ./service-account.json

# Recreate service account file as Kubernetes secret
kubectl create secret generic mycarehub-service-account \
    --namespace $NAMESPACE \
    --from-file=key.json=./service-account.json

helm upgrade \
    --install \
    --debug \
    --create-namespace \
    --namespace "${NAMESPACE}" \
    --set app.replicaCount="${APP_REPLICA_COUNT}" \
    --set service.port="${PORT}"\
    --set app.container.image="${DOCKER_IMAGE_TAG}"\
    --set app.container.env.googleCloudProject="${GOOGLE_CLOUD_PROJECT}"\
    --set app.container.env.firebaseWebApiKey="${FIREBASE_WEB_API_KEY}"\
    --set app.container.env.jwtKey="${JWT_KEY}"\
    --set app.container.env.repository="${REPOSITORY}"\
    --set app.container.env.environment="${ENVIRONMENT}"\
    --set app.container.env.googleProjectNumber="${GOOGLE_PROJECT_NUMBER}"\
    --set app.container.env.sentryDSN="${SENTRY_DSN}"\
    --set app.container.env.serviceHost="${SERVICE_HOST}"\
    --set app.container.env.postgresUser="${POSTGRES_USER}"\
    --set app.container.env.postgresPassword="${POSTGRES_PASSWORD}"\
    --set app.container.env.postgresDB="${POSTGRES_DB}"\
    --set app.container.env.databaseRegion="${DATABASE_REGION}"\
    --set app.container.env.databaseInstance="${DATABASE_INSTANCE}"\
    --set app.container.env.databaseInstanceConnectionName="${DATABASE_INSTANCE_CONNECTION_NAME}"\
    --set app.container.env.defaultOrgID="${DEFAULT_ORG_ID}"\
    --set app.container.env.proInviteLink="${PRO_INVITE_LINK}"\
    --set app.container.env.consumerInviteLink="${CONSUMER_INVITE_LINK}"\
    --set app.container.env.sensitiveContentSecretKey="${SENSITIVE_CONTENT_SECRET_KEY}"\
    --set app.container.env.mailgunAPIKey="${MAILGUN_API_KEY}"\
    --set app.container.env.mailgunDomain="${MAILGUN_DOMAIN}"\
    --set app.container.env.mailgunFrom="${MAILGUN_FROM}"\
    --set app.container.env.djangoAuthorizationToken="${DJANGO_AUTHORIZATION_TOKEN}"\
    --set app.container.env.contentAPIURL="${CONTENT_API_URL}"\
    --set app.container.env.contentServiceBaseURL="${CONTENT_SERVICE_BASE_URL}"\
    --set app.container.env.googleCloudStorageURL="${GOOGLE_CLOUD_STORAGE_URL}"\
    --set app.container.env.invitePinExpiryDays="${INVITE_PIN_EXPIRY_DAYS}"\
    --set app.container.env.pinExpiryDays="${PIN_EXPIRY_DAYS}"\
    --set app.container.env.myCareHubAdminEmail="${MYCAREHUB_ADMIN_EMAIL}"\
    --set app.container.env.surveysSystemEmail="${SURVEYS_SYSTEM_EMAIL}"\
    --set app.container.env.surveysSystemPassword="${SURVEYS_SYSTEM_PASSWORD}"\
    --set app.container.env.surveysBaseURL="${SURVEYS_BASE_URL}"\
    --set app.container.env.consumerAppIdentifier="${CONSUMER_APP_IDENTIFIER}"\
    --set app.container.env.proAppIdentifier="${PRO_APP_IDENTIFIER}"\
    --set app.container.env.consumerAppName="${CONSUMER_APP_NAME}"\
    --set app.container.env.proAppName="${PRO_APP_NAME}"\
    --set app.container.env.silCommsBaseURL="${SIL_COMMS_BASE_URL}"\
    --set app.container.env.silCommsEmail="${SIL_COMMS_EMAIL}"\
    --set app.container.env.silCommsPassword="${SIL_COMMS_PASSWORD}"\
    --set app.container.env.silCommsSenderID="${SIL_COMMS_SENDER_ID}"\
    --set app.container.env.twilioAccountSID="${TWILIO_ACCOUNT_SID}"\
    --set app.container.env.twilioAccountAuthToken="${TWILIO_ACCOUNT_AUTH_TOKEN}"\
    --set app.container.env.twilioSMSNumber="${TWILIO_SMS_NUMBER}"\
    --set app.container.env.defaultProgramID="${DEFAULT_PROGRAM_ID}"\
    --set app.container.env.matrixBaseURL="${MATRIX_BASE_URL}"\
    --set app.container.env.mchMatrixUser="${MCH_MATRIX_USER}"\
    --set app.container.env.mchMatrixPassword="${MCH_MATRIX_PASSWORD}"\
    --set app.container.env.matrixDomain="${MATRIX_DOMAIN}"\
    --set app.container.env.fositeSecret="${FOSITE_SECRET}"\
    --set networking.issuer.name="letsencrypt-prod"\
    --set networking.issuer.privateKeySecretRef="letsencrypt-prod"\
    --set networking.ingress.host="${APPDOMAIN}"\
    --wait \
    --timeout 300s \
    -f ./charts/mycarehub-multitenant/values.yaml \
    $APPNAME \
    ./charts/mycarehub-multitenant
