module github.com/KejawenLab/skeleton/v3

go 1.16

replace github.com/KejawenLab/bima/v4 v4.0.0 => ../bima

require (
	github.com/KejawenLab/bima/v4 v4.0.0
	github.com/fatih/color v1.13.0
	github.com/gertd/go-pluralize v0.2.1
	github.com/goccy/go-json v0.9.8
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3
	github.com/iancoleman/strcase v0.2.0
	github.com/jinzhu/copier v0.3.5
	github.com/joho/godotenv v1.4.0
	github.com/sarulabs/di/v2 v2.4.2
	github.com/sarulabs/dingo/v4 v4.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/vito/go-interact v1.0.1
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.103.0 // indirect
	cloud.google.com/go/pubsub v1.23.1 // indirect
	github.com/Shopify/sarama v1.34.1 // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama v0.32.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220630143837-2104d58473e0 // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	golang.org/x/tools v0.1.11 // indirect
	google.golang.org/genproto v0.0.0-20220630174209-ad1d48641aa7
	google.golang.org/grpc v1.47.0
	gorm.io/driver/postgres v1.3.8 // indirect
	gorm.io/gorm v1.23.7
)
