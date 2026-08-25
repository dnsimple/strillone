// +heroku goVersion go1.24.4
// +heroku install ./cmd/...

module github.com/dnsimple/strillone

go 1.25.2

require (
	github.com/caarlos0/env/v11 v11.4.1
	github.com/dnsimple/dnsimple-go/v9 v9.1.1
	github.com/slack-go/slack v0.27.0
	github.com/stretchr/testify v1.12.1
	github.com/wunderlist/ttlcache v0.0.0-20180801091818-7dbceb0d5094
)

require (
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
)
