module github.com/go-helium/postgres

go 1.13

replace mellium.im/sasl v0.2.1 => github.com/mellium/sasl v0.2.1

require (
	github.com/go-pg/pg/v9 v9.0.1
	github.com/im-kulikov/helium v0.12.2
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.12.0
)
