package store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Credentials struct {
	AccessId  string `db:"access_id"`  // AWS Access Key ID
	SecretKey string `db:"secret_key"` // AWS Secret Access Key
	Region    string `db:"region"`     // AWS Region
}

// GetCredentials returns accessId, secretKey and region information for the current AWS profile.
func GetCredentials() *Credentials {
	db := newSession()
	defer db.Close()

	return db.getCredentials()
}

// SetCredentials creates a new entry in the database for the current AWS config, if not present, else returns error.
func SetCredentials() {
	db := newSession()
	defer db.Close()

	if db.entryExists(Profile) {
		fmt.Println("Error: profile already exists!")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("AWS Access Key ID: ")
	accessId, _ := reader.ReadString('\n')

	fmt.Print("AWS Secret Access Key: ")
	secretKey, _ := reader.ReadString('\n')

	accessId = strings.TrimRight(accessId, "\n")
	secretKey = strings.TrimRight(secretKey, "\n")

	db.setCredentials(accessId, secretKey)
	fmt.Printf("Successfully added new profile: %s\n", Profile)
	fmt.Printf("Switched to '%s' with region '%s'\n", Profile, Region)

}

func (db *store) getCredentials() *Credentials {
	var credentials Credentials

	db.Get(&credentials, "WITH t1 as (SELECT name FROM credentials), t2 as (SELECT profile FROM current_config) "+
		"SELECT t1.access_id, t1.secret_key, t2.region FROM credentials t1, current_config t2 WHERE t1.name = t2.profile")

	return &credentials
}

func (db *store) setCredentials(accessId string, secretKey string) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err.Error())
	}

	tx.Exec("INSERT INTO credentials (name, access_id, secret_key) VALUES ($1, $2, $3)", Profile, accessId, secretKey)
	tx.Exec("UPDATE current_config SET profile = $1, region = $2", Profile, Region)

	return tx.Commit()
}
