// Package store implements the database operations required for config.
package store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AWS Profile.
var Profile string

// AWS Region.
var Region string

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

// entryExists returns true if an entry exists in database for the given profile, else false.
func entryExists(profile string) bool {
	db := newDBSession()
	defer db.Close()

	var result string
	query := db.QueryRow("SELECT access_id FROM credentials WHERE name = $1", profile)
	query.Scan(&result)

	if result == "" {
		return false
	}

	return true
}

// GetCredentials returns accessId, secretKey and region information for the current AWS profile.
func GetCredentials() (accessId string, secretKey string, region string) {
	db := newDBSession()

	result := db.QueryRow("WITH t1 as (SELECT name FROM credentials), t2 as (SELECT profile FROM current_config) " +
		"SELECT t1.access_id, t1.secret_key, t2.region FROM credentials t1, current_config t2 WHERE t1.name = t2.profile")
	result.Scan(&accessId, &secretKey, &region)

	defer db.Close()
	return
}

// SetCredentials creates a new entry in the database for the current AWS config, if not present, else returns error.
func SetCredentials() {
	db := newDBSession()

	tx, txErr := db.Begin()
	if txErr != nil {
		fmt.Println(txErr.Error())
	}
	defer db.Close()

	if entryExists(Profile) {
		fmt.Println("Error: profile already exists!")
		os.Exit(1)
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("AWS Access Key ID: ")
		accessId, _ := reader.ReadString('\n')

		fmt.Print("AWS Secret Access Key: ")
		secretKey, _ := reader.ReadString('\n')

		accessId = strings.TrimRight(accessId, "\n")
		secretKey = strings.TrimRight(secretKey, "\n")

		tx.Exec("INSERT INTO credentials (name, access_id, secret_key) VALUES ($1, $2, $3)", Profile, accessId, secretKey)
		tx.Exec("UPDATE current_config SET profile = $1, region = $2", Profile, Region)
		tx.Commit()

		fmt.Printf("Successfully added new profile: %s\n", Profile)
		fmt.Printf("Switched to '%s' with region '%s'\n", Profile, Region)
	}
}

// ListProfiles returns the list of available AWS profiles in the config.
func ListProfiles() {
	db := newDBSession()
	defer db.Close()

	rows, _ := db.Query("SELECT name FROM credentials")
	var currentEnv string
	current := db.QueryRow("SELECT profile FROM current_config")
	current.Scan(&currentEnv)

	for rows.Next() {
		var envName string
		rows.Scan(&envName)
		if envName == currentEnv {
			fmt.Printf("\033[1;32m%s* (current)\033[m\n", envName)
		} else {
			fmt.Println(envName)
		}
	}
}

// DeleteProfile deletes the given AWS profile from the database.
func DeleteProfile(profile string) {
	db := newDBSession()
	defer db.Close()

	tx, txErr := db.Begin()
	if txErr != nil {
		fmt.Println(txErr.Error())
	}

	if !entryExists(profile) {
		fmt.Println("Error: no such profile found!")
		os.Exit(1)
	}

	tx.Exec("DELETE FROM credentials where name = $1", profile)
	tx.Commit()

	fmt.Printf("Successfully deleted '%s' from config!\n", profile)
}

// UseProfile loads the given profile to the current config in the database.
func UseProfile() {
	db := newDBSession()

	tx, txErr := db.Begin()
	if txErr != nil {
		fmt.Println(txErr.Error())
	}
	defer db.Close()

	if !entryExists(Profile) {
		fmt.Println("Error: no such profile found!")
		os.Exit(1)
	}

	tx.Exec("UPDATE current_config SET profile = $1, region = $2", Profile, Region)
	tx.Commit()

	fmt.Printf("Switched to '%s' with region '%s'\n", Profile, Region)
}
