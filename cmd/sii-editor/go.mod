module github.com/Luzifer/sii/cmd/sii-editor

go 1.13

replace github.com/Luzifer/sii => ../../

require (
	github.com/Luzifer/go_helpers/v2 v2.9.1
	github.com/Luzifer/rconfig/v2 v2.2.1
	github.com/Luzifer/scs-extract v0.1.0
	github.com/Luzifer/sii v0.0.0-00010101000000-000000000000
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/validator.v2 v2.0.0-20191107172027-c3144fdedc21 // indirect
	gopkg.in/yaml.v2 v2.2.7
)
