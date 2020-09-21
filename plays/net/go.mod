module github.com/wenerme/letsgo/plays/net

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/json-iterator/go v1.1.10
	github.com/pion/turn/v2 v2.0.4
	github.com/pion/webrtc/v3 v3.0.0-beta.5
	github.com/rancher/remotedialer v0.2.5
	github.com/sirupsen/logrus v1.6.0
	github.com/wenerme/letsgo v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
)

replace github.com/wenerme/letsgo => ../..
