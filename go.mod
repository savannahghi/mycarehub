module github.com/savannahghi/onboarding-service

go 1.16

require (
	cloud.google.com/go/errorreporting v0.1.0 // indirect
	cloud.google.com/go/firestore v1.5.0
	cloud.google.com/go/kms v1.0.0 // indirect
	cloud.google.com/go/monitoring v1.0.0 // indirect
	cloud.google.com/go/profiler v0.1.0 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	firebase.google.com/go v3.13.0+incompatible
	github.com/99designs/gqlgen v0.13.0
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.26.0
	github.com/casbin/casbin/v2 v2.31.3
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/imroc/req v0.3.0
	github.com/lib/pq v1.10.3
	github.com/savannahghi/converterandformatter v0.0.9
	github.com/savannahghi/enumutils v0.0.3
	github.com/savannahghi/feedlib v0.0.4
	github.com/savannahghi/firebasetools v0.0.15
	github.com/savannahghi/interserviceclient v0.0.13
	github.com/savannahghi/onboarding v0.0.22
	github.com/savannahghi/profileutils v0.0.17
	github.com/savannahghi/scalarutils v0.0.4
	github.com/savannahghi/serverutils v0.0.6
	github.com/segmentio/ksuid v1.0.4
	github.com/sirupsen/logrus v1.8.1
	github.com/tj/assert v0.0.3
	github.com/vektah/gqlparser/v2 v2.1.0
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	gorm.io/datatypes v1.0.2
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.16
)
