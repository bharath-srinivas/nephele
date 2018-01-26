package store

import (
	"fmt"
	"os"
	"os/user"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
)

func newDBSession() *sqlx.DB {
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

	db := sqlx.MustConnect("sqlite3", configFile)

	if f, err := os.Stat(configFile); !os.IsNotExist(err) {
		if f.Size() < 1 {
			db.MustExec(schema)
			db.MustExec("INSERT INTO current_config (profile, region) VALUES (NULL, NULL)")
		}
	}

	return db
}