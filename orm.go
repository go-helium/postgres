package postgres

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/im-kulikov/helium/module"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	// Config alias
	Config struct {
		Hostname string
		Username string
		Password string
		Database string
		Debug    bool
		PoolSize int
		Options  map[string]string
	}

	// Hook is a simple implementation of pg.QueryHook
	Hook struct {
		StartAt time.Time
		Before  func(context.Context, *pg.QueryEvent) (context.Context, error)
		After   func(context.Context, *pg.QueryEvent) error
	}

	// Error is constant error
	Error string
)

const (
	// ErrPemParse when couldn't parse pem in sslrootcert
	ErrPemParse = Error("couldn't parse pem in sslrootcert")
	// ErrEmptyConfig when given empty options
	ErrEmptyConfig = Error("database empty config")
	// ErrEmptyLogger when logger not initialized
	ErrEmptyLogger = Error("database empty logger")
	// ErrSSLKeyHasWorldPermissions when pk permissions no u=rw (0600) or less
	ErrSSLKeyHasWorldPermissions = Error("private key file has group or world access. Permissions should be u=rw (0600) or less")

	errUnsupportedSSLMode = `unsupported sslmode %q; only "require" (default), "verify-full", "verify-ca", and "disable" supported`
)

var (
	// Module is default connection to PostgreSQL
	Module = module.Module{
		{Constructor: NewDefaultConfig},
		{Constructor: NewConnection},
	}
)

// Error implementation
func (e Error) Error() string {
	return string(e)
}

// BeforeQuery callback
func (h *Hook) BeforeQuery(ctx context.Context, e *pg.QueryEvent) (context.Context, error) {
	h.StartAt = time.Now()

	if h.Before == nil {
		return ctx, e.Err
	}

	return h.Before(ctx, e)
}

// AfterQuery callback
func (h Hook) AfterQuery(ctx context.Context, e *pg.QueryEvent) error {
	if h.After == nil {
		return e.Err
	}

	return h.After(ctx, e)
}

// NewDefaultConfig returns connection config
func NewDefaultConfig(v *viper.Viper) (*Config, error) {
	if !v.IsSet("postgres") {
		return nil, ErrEmptyConfig
	}

	// v.SetDefault("postgres.hostname", "localhost")
	v.SetDefault("postgres.options.sslmode", "disable")

	// re-fetch by full key
	options := v.GetStringMapString("postgres.options")
	if len(options) > 0 {
		for opt := range options {
			options[opt] = v.GetString("postgres.options." + opt)
		}
	}

	return &Config{
		Hostname: v.GetString("postgres.hostname"),
		Username: v.GetString("postgres.username"),
		Password: v.GetString("postgres.password"),
		Database: v.GetString("postgres.database"),
		Debug:    v.GetBool("postgres.debug"),
		PoolSize: v.GetInt("postgres.pool_size"),
		Options:  options,
	}, nil
}

// NewConnection returns database connection
func NewConnection(cfg *Config, l *zap.Logger) (db *pg.DB, err error) {
	if cfg == nil {
		err = ErrEmptyConfig
		return
	}

	if l == nil {
		err = ErrEmptyLogger
		return
	}

	opts := &pg.Options{
		Addr:     cfg.Hostname,
		User:     cfg.Username,
		Password: cfg.Password,
		Database: cfg.Database,
		PoolSize: cfg.PoolSize,
	}

	if cfg.Debug {
		l.Debug("Connect to PostgreSQL",
			zap.String("hostname", cfg.Hostname),
			zap.String("username", cfg.Username),
			zap.String("password", cfg.Password),
			zap.String("database", cfg.Database),
			zap.Int("pool_size", cfg.PoolSize),
			zap.Any("options", cfg.Options))
	}

	if opts.TLSConfig, err = ssl(cfg.Options); err != nil {
		return nil, err
	}

	db = pg.Connect(opts)
	if _, err = db.ExecOne("SELECT 1"); err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	if cfg.Debug {
		h := new(Hook)
		h.After = func(ctx context.Context, e *pg.QueryEvent) error {
			query, qErr := e.FormattedQuery()
			l.Debug("pg query",
				zap.String("query", query),
				zap.Duration("query_time", time.Since(h.StartAt)),
				zap.Any("params", e.Params),
				zap.NamedError("format_error", qErr),
				zap.Error(e.Err))

			return e.Err
		}
		db.AddQueryHook(h)
	}

	return
}
