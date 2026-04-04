package cmd

import (
	"context"
	"embed"
	"net/url"
	"os"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/logging"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var EmbeddedMigrations embed.FS

var migrateCmd = cobra.Command{
	Use:  "migrate",
	Long: "Migrate database strucutures. This will create new tables and add missing columns and indexes.",
	Run:  migrate,
}

func migrate(cmd *cobra.Command, args []string) {
	globalConfig := loadGlobalConfig(cmd.Context())
	u, err := url.Parse(globalConfig.DB.URL)
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "parsing db connection url"))
	}

	if globalConfig.DB.Driver == "" && globalConfig.DB.URL != "" {
		globalConfig.DB.Driver = u.Scheme
	}

	log := logrus.StandardLogger()

	pop.Debug = false
	if globalConfig.Logging.Level != "" {
		level, err := logrus.ParseLevel(globalConfig.Logging.Level)
		if err != nil {
			log.Fatalf("Failed to parse log level: %+v", err)
		}
		log.SetLevel(level)
		if level == logrus.DebugLevel {
			// Set to true to display query info
			pop.Debug = true
		}
		if level != logrus.DebugLevel {
			var noopLogger = func(lvl logging.Level, s string, args ...interface{}) {
			}
			// Hide pop migration logging
			pop.SetLogger(noopLogger)
		}
	}

	q := u.Query()
	q.Add("application_name", "auth_migrations")
	u.RawQuery = q.Encode()
	deets := &pop.ConnectionDetails{
		Dialect: globalConfig.DB.Driver,
		URL:     u.String(),
	}
	deets.Options = map[string]string{
		"migration_table_name": "schema_migrations",
		"Namespace":            globalConfig.DB.Namespace,
	}

	db, err := pop.NewConnection(deets)
	if err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "opening db connection"))
	}
	defer db.Close()

	if err := db.Open(); err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "checking database connection"))
	}

	if globalConfig.DB.AutoCreateNamespace {
		log.Infof("Create schema if not exists: %s", globalConfig.DB.Namespace)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		defer cancel()
		_, err = db.Store.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", globalConfig.DB.Namespace))
		if err != nil {
			log.Fatalf("%+v", errors.Wrap(err, "creating namespace"))
		}
	} else {
		log.Infof("Using existing schema: %s", globalConfig.DB.Namespace)
	}

	log.Debugf("Reading migrations from executable")
	box, err := pop.NewMigrationBox(EmbeddedMigrations, db)
	if err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "creating db migrator"))
	}

	mig := box.Migrator

	log.Debugf("before status")

	if log.Level == logrus.DebugLevel {
		err = mig.Status(os.Stdout)
		if err != nil {
			log.Fatalf("%+v", errors.Wrap(err, "migration status"))
		}
	}

	// turn off schema dump
	mig.SchemaPath = ""

	count, err := mig.UpTo(0)
	if err != nil {
		log.Fatalf("%v", errors.Wrap(err, "running db migrations"))
	} else {
		log.WithField("count", count).Infof("GoTrue migrations applied successfully")
	}

	log.Debugf("after status")

	if log.Level == logrus.DebugLevel {
		err = mig.Status(os.Stdout)
		if err != nil {
			log.Fatalf("%+v", errors.Wrap(err, "migration status"))
		}
	}
}
