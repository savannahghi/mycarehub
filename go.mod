module github.com/savannahghi/mycarehub

go 1.17

require (
	firebase.google.com/go v3.13.0+incompatible
	github.com/99designs/gqlgen v0.14.0
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.27.0
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/imroc/req v0.3.0
	github.com/lib/pq v1.10.3
	github.com/savannahghi/converterandformatter v0.0.11
	github.com/savannahghi/enumutils v0.0.3
	github.com/savannahghi/errorcodeutil v0.0.3
	github.com/savannahghi/firebasetools v0.0.16
	github.com/savannahghi/interserviceclient v0.0.18
	github.com/savannahghi/onboarding v0.0.29
	github.com/savannahghi/serverutils v0.0.6
	github.com/segmentio/ksuid v1.0.4
	github.com/sirupsen/logrus v1.8.1
	github.com/tj/assert v0.0.3
	github.com/vektah/gqlparser/v2 v2.2.0
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.26.1
	gopkg.in/go-playground/validator.v9 v9.31.0
	gorm.io/driver/postgres v1.2.1
	gorm.io/gorm v1.22.2
)

require (
	cloud.google.com/go v0.97.0 // indirect
	cloud.google.com/go/errorreporting v0.1.0 // indirect
	cloud.google.com/go/firestore v1.6.1 // indirect
	cloud.google.com/go/kms v1.1.0 // indirect
	cloud.google.com/go/logging v1.4.2 // indirect
	cloud.google.com/go/monitoring v1.1.0 // indirect
	cloud.google.com/go/profiler v0.1.1 // indirect
	cloud.google.com/go/pubsub v1.12.2 // indirect
	cloud.google.com/go/storage v1.18.2 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.41.19 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20211011173535-cb28da3451f1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/go-control-plane v0.10.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.2 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/getsentry/sentry-go v0.11.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/gofrs/uuid v4.1.0+incompatible // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/pprof v0.0.0-20211104044539-f987b9c94b31 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lithammer/shortuuid v3.0.0+incompatible // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/savannahghi/feedlib v0.0.6 // indirect
	github.com/savannahghi/profileutils v0.0.23 // indirect
	github.com/savannahghi/pubsubtools v0.0.2 // indirect
	github.com/savannahghi/scalarutils v0.0.4 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.2.1 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.opentelemetry.io/contrib v0.21.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.21.0 // indirect
	go.opentelemetry.io/otel v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.1.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.21.0 // indirect
	go.opentelemetry.io/otel/metric v0.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/trace v1.1.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211105183446-c75c47738b0c // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.60.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247 // indirect
	google.golang.org/grpc v1.42.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
