package postgres

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewDefaultConfig(t *testing.T) {
	t.Run("Check orm module", func(t *testing.T) {
		l := zap.L()

		t.Run("must fail on empty", func(t *testing.T) {
			v := viper.New()
			cfg, err := NewDefaultConfig(v)
			require.Nil(t, cfg)
			require.Error(t, err)
		})

		t.Run("should be ok", func(t *testing.T) {
			v := viper.New()
			hostname := "localhost"
			v.SetDefault("postgres.hostname", hostname)

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, hostname, cfg.Hostname)
		})

		t.Run("should fail for empty config", func(t *testing.T) {
			con, err := NewConnection(nil, l)
			require.Nil(t, con)
			require.EqualError(t, err, ErrEmptyConfig.Error())
		})

		t.Run("should fail for empty logger", func(t *testing.T) {
			v := viper.New()
			hostname := "localhost"
			v.SetDefault("postgres.hostname", hostname)

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, hostname, cfg.Hostname)

			cli, err := NewConnection(cfg, nil)
			require.Nil(t, cli)
			require.EqualError(t, err, ErrEmptyLogger.Error())
		})

		t.Run("should not connect", func(t *testing.T) {
			v := viper.New()
			v.SetDefault("postgres.username", "unknown")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)

			cli, err := NewConnection(cfg, l)
			require.Error(t, err)
			require.Nil(t, cli)
		})

		t.Run("should connect", func(t *testing.T) {
			v := viper.New()
			v.SetDefault("postgres.debug", true)
			v.SetDefault("postgres.username", "postgres")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, "postgres", cfg.Username)
			require.Equal(t, "postgres", cfg.Password)
			require.Equal(t, "postgres", cfg.Database)

			cli, err := NewConnection(cfg, l)
			require.NoError(t, err)
			require.NotNil(t, cli)

			_, err = cli.ExecOne("SELECT 1")
			require.NoError(t, err)

			err = cli.Close()
			require.NoError(t, err)
		})

		t.Run("should verify all ssl-modes", func(t *testing.T) {
			v := viper.New()
			v.SetDefault("postgres.debug", true)
			v.SetDefault("postgres.username", "postgres")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")

			modes := []string{
				"disable",
				"require",
				"unknown",
				"verify-ca",
				"verify-full",
			}

			for _, mode := range modes {
				if mode == "require" {
					v.Set("postgres.options.sslrootcert", mode)
				}

				v.Set("postgres.options.sslmode", mode)

				cfg, err := NewDefaultConfig(v)
				require.NoError(t, err)
				require.Equal(t, "postgres", cfg.Username)
				require.Equal(t, "postgres", cfg.Password)
				require.Equal(t, "postgres", cfg.Database)

				_, _ = NewConnection(cfg, l)
			}
		})

		t.Run("should not fail on couldn't parse pem in sslrootcert", func(t *testing.T) {
			v := viper.New()

			file, err := ioutil.TempFile("/tmp", "something")
			require.NoError(t, err)

			require.NoError(t, file.Chmod(0701))

			v.SetDefault("postgres.debug", true)
			v.SetDefault("postgres.username", "postgres")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")
			v.Set("postgres.options.sslmode", "verify-full")
			v.Set("postgres.options.sslrootcert", file.Name())

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, "postgres", cfg.Username)
			require.Equal(t, "postgres", cfg.Password)
			require.Equal(t, "postgres", cfg.Database)

			cli, err := NewConnection(cfg, l)
			require.Nil(t, cli)
			require.Error(t, err)

			require.NoError(t, file.Close())
			require.NoError(t, os.Remove(file.Name()))
		})

		t.Run("should not fail on file sslrootcert not found", func(t *testing.T) {
			v := viper.New()

			v.SetDefault("postgres.debug", true)
			v.SetDefault("postgres.username", "postgres")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")
			v.Set("postgres.options.sslmode", "verify-full")
			v.Set("postgres.options.sslrootcert", "file-not-exists")

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, "postgres", cfg.Username)
			require.Equal(t, "postgres", cfg.Password)
			require.Equal(t, "postgres", cfg.Database)

			cli, err := NewConnection(cfg, l)
			require.Nil(t, cli)
			require.Error(t, err)
		})

		t.Run("should not fail on bad PEM data", func(t *testing.T) {
			v := viper.New()

			file, err := ioutil.TempFile("/tmp", "something")
			require.NoError(t, err)

			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			require.NoError(t, err)

			tmpl := &x509.Certificate{
				SerialNumber:          big.NewInt(1),
				NotBefore:             time.Now(),
				NotAfter:              time.Now().Add(24 * time.Hour),
				BasicConstraintsValid: true,
				KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			}

			cert, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
			require.NoError(t, err)

			_, err = file.Write(cert)
			require.NoError(t, err)

			v.SetDefault("postgres.debug", true)
			v.SetDefault("postgres.username", "postgres")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")
			v.Set("postgres.options.sslmode", "verify-full")
			v.Set("postgres.options.sslcert", file.Name())
			v.Set("postgres.options.sslkey", file.Name())

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, "postgres", cfg.Username)
			require.Equal(t, "postgres", cfg.Password)
			require.Equal(t, "postgres", cfg.Database)

			modes := []os.FileMode{
				0600,
				0701,
			}

			for _, mode := range modes {
				require.NoError(t, file.Chmod(mode))

				cli, err := NewConnection(cfg, l)
				require.Nil(t, cli)
				require.Error(t, err)
			}

			require.NoError(t, file.Close())
			require.NoError(t, os.Remove(file.Name()))
		})

		t.Run("should connect with before/after hooks", func(t *testing.T) {
			v := viper.New()
			v.SetDefault("postgres.debug", true)
			v.SetDefault("postgres.username", "postgres")
			v.SetDefault("postgres.password", "postgres")
			v.SetDefault("postgres.database", "postgres")

			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.Equal(t, "postgres", cfg.Username)
			require.Equal(t, "postgres", cfg.Password)
			require.Equal(t, "postgres", cfg.Database)

			cli, err := NewConnection(cfg, l)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.AddQueryHook(&Hook{
				Before: func(ctx context.Context, e *pg.QueryEvent) (context.Context, error) { return ctx, e.Err },
				After:  nil,
			})

			_, err = cli.ExecOne("SELECT 1")
			require.NoError(t, err)

			err = cli.Close()
			require.NoError(t, err)
		})
	})
}
