// +heroku goVersion go1.14.5
// +heroku install ./cmd/...

module github.com/dnsimple/strillone

go 1.14

require (
	github.com/bluele/slack v0.0.0-20180528010058-b4b4d354a079
	github.com/dnsimple/dnsimple-go v0.71.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/wunderlist/ttlcache v0.0.0-20180801091818-7dbceb0d5094
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	google.golang.org/appengine v1.6.1 // indirect
)
