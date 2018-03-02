// Package store implements the database operations required for config.
package store

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// AWS Profile.
	Profile string

	// AWS Region.
	Region string
)

// schema to be executed during config creation.
var schema = `CREATE TABLE credentials (
	name varchar(20) NOT NULL,
	access_id varchar(100) NOT NULL,
	secret_key varchar(150) NOT NULL,
	token varchar(255) NULL,
	UNIQUE (name)
);
	CREATE TABLE current_config (
	profile varchar(20) NULL,
	region varchar(20) NULL
);`

type store struct {
	*sqlx.DB // an instance of sqlx DB
}

// NewSession creates the config file required by aws-go to function, if not present and creates a new database
// connection and returns a new store.
func NewSession() *store {
	usr, _ := user.Current()
	homePath := usr.HomeDir

	configPath := path.Join(homePath, ".aws")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join(homePath, ".aws"), 0755)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	configFileName := "aws_go.credentials"
	configFile := path.Join(configPath, configFileName)

	return &store{
		newDBSession(configFile),
	}
}

// newDBSession opens a new database connection with the specified config file and returns an instance of sqlx DB.
func newDBSession(configFile string) *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", configFile)

	if f, err := os.Stat(configFile); !os.IsNotExist(err) {
		if f.Size() < 1 {
			db.MustExec(schema)
			db.MustExec("INSERT INTO current_config (profile, region) VALUES (NULL, NULL)")
		}
	}
	return db
}

// newTest creates a new database connection for testing purposes and returns a new store.
func newTest() *store {
	usr, _ := user.Current()
	homePath := usr.HomeDir

	configPath := path.Join(homePath, ".aws")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join(homePath, ".aws"), 0755)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	configFileName := "aws_go_test.credentials"
	configFile := path.Join(configPath, configFileName)

	return &store{
		newDBSession(configFile),
	}
}

// cleanup removes the testing config file created during the invocation of newTest function.
func cleanup() error {
	usr, _ := user.Current()
	homePath := usr.HomeDir

	configPath := path.Join(homePath, ".aws")

	configFileName := "aws_go_test.credentials"
	configFile := path.Join(configPath, configFileName)
	return os.Remove(configFile)
}
