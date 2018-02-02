package store

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bharath-srinivas/aws-go/internal/colors"
)

// ListProfiles prints the list of available AWS profiles in the config.
func ListProfiles() {
	db := newSession()
	defer db.Close()

	rows := db.listProfiles()

	for rows.Next() {
		var profile string
		rows.Scan(&profile)

		if db.currentProfile(profile) {
			fmt.Printf("%s[1;%dm%s* (current)%s[m\n", escape, colors.Green, profile, escape)
		} else {
			fmt.Println(profile)
		}
	}
}

// DeleteProfile deletes the given AWS profile from the database, if it's not the currently selected profile.
func DeleteProfile(profile string) {
	db := newSession()
	defer db.Close()

	if !db.entryExists(profile) {
		fmt.Println("Error: no such profile found!")
		os.Exit(1)
	} else if db.currentProfile(profile) {
		fmt.Println("Error: cannot delete currently active profile!")
		os.Exit(1)
	}

	db.deleteProfile(profile)
	fmt.Printf("Successfully deleted '%s' from config!\n", profile)
}

// UseProfile loads the given profile to the current config in the database.
func UseProfile() {
	db := newSession()
	defer db.Close()

	if !db.entryExists(Profile) {
		fmt.Println("Error: no such profile found!")
		os.Exit(1)
	}

	db.useProfile()
	fmt.Printf("Switched to '%s' with region '%s'\n", Profile, Region)
}

func (db *store) listProfiles() *sql.Rows {
	rows, _ := db.Query("SELECT name FROM credentials")
	return rows
}

func (db *store) deleteProfile(profile string) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err.Error())
	}

	tx.Exec("DELETE FROM credentials where name = $1", profile)

	return tx.Commit()
}

func (db *store) useProfile() error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err.Error())
	}

	tx.Exec("UPDATE current_config SET profile = $1, region = $2", Profile, Region)

	return tx.Commit()
}
