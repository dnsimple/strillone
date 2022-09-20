// +heroku goVersion go1.19.1
// +heroku install ./cmd/...

module github.com/dnsimple/strillone

go 1.19

require (
	github.com/bluele/slack v0.0.0-20180528010058-b4b4d354a079
	github.com/dnsimple/dnsimple-go v1.0.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/wunderlist/ttlcache v0.0.0-20180801091818-7dbceb0d5094
)

require (
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	google.golang.org/appengine v1.6.1 // indirect
)
