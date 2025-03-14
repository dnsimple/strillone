// +heroku goVersion go1.24.1
// +heroku install ./cmd/...

module github.com/dnsimple/strillone

go 1.24.1

require (
	github.com/dnsimple/dnsimple-go v1.7.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/slack-go/slack v0.16.0
	github.com/stretchr/testify v1.10.0
	github.com/wunderlist/ttlcache v0.0.0-20180801091818-7dbceb0d5094
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/oauth2 v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
