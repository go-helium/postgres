module github.com/go-helium/postgres

require (
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/im-kulikov/helium v0.11.11
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.3.0
	go.uber.org/zap v1.10.0
	mellium.im/sasl v0.2.1 // indirect
)

replace mellium.im/sasl v0.2.1 => github.com/mellium/sasl v0.2.1
