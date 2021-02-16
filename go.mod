module github.com/nickysemenza/gourd

go 1.14

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.16.0
	github.com/Masterminds/squirrel v1.5.0
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/buckket/go-blurhash v1.1.0
	github.com/charmbracelet/glamour v0.2.0
	github.com/cosmtrek/air v1.15.1
	github.com/davecgh/go-spew v1.1.1
	github.com/deepmap/oapi-codegen v1.5.1
	github.com/dgraph-io/ristretto v0.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getkin/kin-openapi v0.39.0
	github.com/getsentry/sentry-go v0.9.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golangci/golangci-lint v1.36.0
	github.com/gphotosuploader/google-photos-api-client-go v1.1.5
	github.com/gphotosuploader/googlemirror v0.5.0
	github.com/jmoiron/sqlx v1.3.1
	github.com/kjk/notionapi v0.0.0-20201230072046-b69038831038
	github.com/labstack/echo/v4 v4.2.0
	github.com/lib/pq v1.9.0
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/ory/go-acc v0.2.6
	github.com/sirupsen/logrus v1.7.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.22.6
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.16.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.16.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.16.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.16.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	golang.org/x/oauth2 v0.0.0-20210216194517-16ff1888fd2e
	golang.org/x/tools v0.1.0
	google.golang.org/api v0.40.0
	gopkg.in/guregu/null.v3 v3.5.0
	sigs.k8s.io/yaml v1.2.0
)
