module github.com/reconinfosec/usp-adapters

go 1.20

require (
	cloud.google.com/go/pubsub v1.30.0
	cloud.google.com/go/storage v1.30.1
	github.com/Azure/azure-event-hubs-go/v3 v3.5.0
	github.com/aws/aws-sdk-go v1.44.240
	github.com/duosecurity/duo_api_golang v0.0.0-20230203160531-b221c950c2b0
	github.com/nxadm/tail v1.4.8
	github.com/refractionPOINT/evtx v0.0.0-20221217013001-bbcbd9938c6e
	github.com/refractionPOINT/go-uspclient v0.0.0-20230227011818-1c8752a42005
	github.com/refractionPOINT/usp-adapters v0.0.0-20230511222808-d88d9bcb44b1
	github.com/vmihailenco/msgpack/v5 v5.3.5
	golang.org/x/oauth2 v0.7.0
	golang.org/x/sync v0.1.0
	golang.org/x/text v0.9.0
	google.golang.org/api v0.117.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.19.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.0.0 // indirect
	github.com/Azure/azure-amqp-common-go/v4 v4.1.0 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/go-amqp v0.19.1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Velocidex/json v0.0.0-20220224052537-92f3c0326e5a // indirect
	github.com/Velocidex/ordereddict v0.0.0-20221110130714-6a7cb85851cd // indirect
	github.com/Velocidex/pkcs7 v0.0.0-20230220112103-d4ed02e1862a // indirect
	github.com/Velocidex/yaml/v2 v2.2.8 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/s2a-go v0.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/refractionPOINT/go-limacharlie/limacharlie v0.0.0-20230215163839-251d0a4a5966 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	www.velocidex.com/golang/binparsergen v0.1.1-0.20220107080050-ae6122c5ed14 // indirect
	www.velocidex.com/golang/go-pe v0.1.1-0.20220107093716-e91743c801de // indirect
)

replace github.com/nxadm/tail => github.com/refractionPOINT/tail v0.0.0-20211216163028-4472660a31a6
