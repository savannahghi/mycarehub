module github.com/savannahghi/mycarehub

go 1.18

require (
	cloud.google.com/go/pubsub v1.26.0
	firebase.google.com/go v3.13.0+incompatible
	github.com/99designs/gqlgen v0.17.26
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.33.0
	github.com/alexedwards/scs/v2 v2.5.1
	github.com/basgys/goxml2json v1.1.0
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/getsentry/sentry-go v0.15.0
	github.com/go-testfixtures/testfixtures/v3 v3.8.1
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/imroc/req v0.3.2
	github.com/jackc/pgtype v1.12.0
	github.com/jarcoal/httpmock v1.2.0
	github.com/kevinburke/twilio-go v0.0.0-20220922200631-8f3f155dfe1f
	github.com/lib/pq v1.10.7
	github.com/mailgun/mailgun-go v2.0.0+incompatible
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/ory/fosite v0.44.0
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.3.0
	github.com/savannahghi/converterandformatter v0.0.11
	github.com/savannahghi/enumutils v0.0.3
	github.com/savannahghi/errorcodeutil v0.0.6
	github.com/savannahghi/feedlib v0.0.6
	github.com/savannahghi/firebasetools v0.0.19
	github.com/savannahghi/interserviceclient v0.0.18
	github.com/savannahghi/profileutils v0.0.27
	github.com/savannahghi/pubsubtools v0.0.3
	github.com/savannahghi/scalarutils v0.0.4
	github.com/savannahghi/serverutils v0.0.7
	github.com/savannahghi/silcomms v0.0.3
	github.com/segmentio/ksuid v1.0.4
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.2
	github.com/tj/assert v0.0.3
	github.com/vektah/gqlparser/v2 v2.5.1
	github.com/xdg-go/pbkdf2 v1.0.0
	go.opencensus.io v0.24.0
	golang.org/x/crypto v0.0.0-20221012134737-56aed061732a
	golang.org/x/text v0.8.0
	google.golang.org/api v0.103.0
	gopkg.in/go-playground/validator.v9 v9.31.0
	gorm.io/driver/postgres v1.4.5
	gorm.io/gorm v1.24.2
)

require (
	cloud.google.com/go v0.105.0 // indirect
	cloud.google.com/go/compute v1.12.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.1 // indirect
	cloud.google.com/go/container v1.6.0 // indirect
	cloud.google.com/go/errorreporting v0.2.0 // indirect
	cloud.google.com/go/firestore v1.6.1 // indirect
	cloud.google.com/go/iam v0.6.0 // indirect
	cloud.google.com/go/logging v1.4.2 // indirect
	cloud.google.com/go/longrunning v0.1.1 // indirect
	cloud.google.com/go/monitoring v1.7.0 // indirect
	cloud.google.com/go/profiler v0.2.0 // indirect
	cloud.google.com/go/storage v1.27.0 // indirect
	cloud.google.com/go/trace v1.3.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.6 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.37.0 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cristalhq/jwt/v4 v4.0.2 // indirect
	github.com/dave/jennifer v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/ecordell/optgen v0.0.6 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/gobuffalo/envy v1.10.2 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20220113144219-d25a53d42d00 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.8 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kevinburke/go-types v0.0.0-20210723172823-2deba1f80ba7 // indirect
	github.com/kevinburke/rest v0.0.0-20210506044642-5611499aa33c // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lithammer/shortuuid v3.0.0+incompatible // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/goveralls v0.0.6 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.22.1 // indirect
	github.com/ory/go-acc v0.2.6 // indirect
	github.com/ory/go-convenience v0.1.0 // indirect
	github.com/ory/viper v1.7.5 // indirect
	github.com/ory/x v0.0.214 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/afero v1.3.2 // indirect
	github.com/spf13/cast v1.3.2-0.20200723214538-8d17101741c8 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.2.1 // indirect
	github.com/urfave/cli/v2 v2.24.4 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.opentelemetry.io/contrib v0.21.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.21.0 // indirect
	go.opentelemetry.io/otel v1.0.0-RC1 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC1 // indirect
	go.opentelemetry.io/otel/internal/metric v0.21.0 // indirect
	go.opentelemetry.io/otel/metric v0.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.0.0-RC1 // indirect
	go.opentelemetry.io/otel/trace v1.0.0-RC1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
	google.golang.org/grpc v1.50.1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/square/go-jose.v2 v2.5.2-0.20210529014059-a5c7eec3c614 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
