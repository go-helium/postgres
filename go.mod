module github.com/go-helium/postgres

go 1.13

replace mellium.im/sasl v0.2.1 => github.com/mellium/sasl v0.2.1

require (
	github.com/go-pg/pg/v10 v10.0.7
	github.com/im-kulikov/helium v0.13.6
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
)
