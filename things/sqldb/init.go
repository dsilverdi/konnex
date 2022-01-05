package sqldb

import (
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	// SSLMode     string
	// SSLCert     string
	// SSLKey      string
	// SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sqlx.DB, error) {
	// url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)
	var db *sqlx.DB
	var err error
	url := cfg.User + ":" + cfg.Pass + "@tcp(" + "things-db" + ":" + cfg.Port + ")/" + cfg.Name + "?parseTime=true&clientFoundRows=true"

	for {
		db, err = sqlx.Connect("mysql", url)
		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "connect: connection refused") {
			return nil, err
		}

		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "things_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS things (
						id       	VARCHAR(255),
						owner    	VARCHAR(255) NOT NULL,
						channel_id	VARCHAR(255), 
						name     	VARCHAR(255) NOT NULL,
						metadata 	TEXT,
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;;`,

					`CREATE TABLE IF NOT EXISTS channels (
						id       VARCHAR(255),
						owner    VARCHAR(255) NOT NULL,
						name     VARCHAR(255) NOT NULL,
						type	 VARCHAR(255) NOT NULL,
						metadata TEXT,
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;;`,

					`CREATE TABLE IF NOT EXISTS connections (
						id	VARCHAR(255) NOT NULL,
						channel_id    VARCHAR(255),
						thing_id      VARCHAR(255),
						protocol	  VARCHAR(255),
						status		  VARCHAR(255),
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;;`,
				},
				Down: []string{
					"DROP TABLE connections",
					"DROP TABLE things",
					"DROP TABLE channels",
				},
			},
			//{
			// 	Id: "things_1",
			// 	Up: []string{
			// 		`CREATE TABLE IF NOT EXISTS things (
			// 			id       UUID,
			// 			owner    VARCHAR(254),
			// 			key      VARCHAR(4096) UNIQUE NOT NULL,
			// 			name     VARCHAR(1024),
			// 			metadata JSON,
			// 			PRIMARY KEY (id, owner)
			// 		)`,
			// 		`CREATE TABLE IF NOT EXISTS channels (
			// 			id       UUID,
			// 			owner    VARCHAR(254),
			// 			name     VARCHAR(1024),
			// 			metadata JSON,
			// 			PRIMARY KEY (id, owner)
			// 		)`,
			// 		`CREATE TABLE IF NOT EXISTS connections (
			// 			channel_id    UUID,
			// 			channel_owner VARCHAR(254),
			// 			thing_id      UUID,
			// 			thing_owner   VARCHAR(254),
			// 			FOREIGN KEY (channel_id, channel_owner) REFERENCES channels (id, owner) ON DELETE CASCADE ON UPDATE CASCADE,
			// 			FOREIGN KEY (thing_id, thing_owner) REFERENCES things (id, owner) ON DELETE CASCADE ON UPDATE CASCADE,
			// 			PRIMARY KEY (channel_id, channel_owner, thing_id, thing_owner)
			// 		)`,
			// 	},
			// 	Down: []string{
			// 		"DROP TABLE connections",
			// 		"DROP TABLE things",
			// 		"DROP TABLE channels",
			// 	},
			// },
			// {
			// 	Id: "things_2",
			// 	Up: []string{
			// 		`ALTER TABLE IF EXISTS things ALTER COLUMN
			// 		 metadata TYPE JSONB using metadata::text::jsonb`,
			// 	},
			// },
			// {
			// 	Id: "things_3",
			// 	Up: []string{
			// 		`ALTER TABLE IF EXISTS channels ALTER COLUMN
			// 		 metadata TYPE JSONB using metadata::text::jsonb`,
			// 	},
			// },
			// {
			// 	Id: "things_4",
			// 	Up: []string{
			// 		`ALTER TABLE IF EXISTS things ADD CONSTRAINT things_id_key UNIQUE (id)`,
			// 	},
			// },
		},
	}

	_, err := migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	return err
}
