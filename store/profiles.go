package store

import (
	"database/sql"
)

// ListProfiles prints the list of available AWS profiles in the config.
func (db *store) ListProfiles() *sql.Rows {
	rows, _ := db.Query("SELECT name FROM credentials")
	return rows
}

// DeleteProfile deletes the given AWS profile from the database, if it's not the currently selected profile.
func (db *store) DeleteProfile(profile string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("DELETE FROM credentials where name = $1", profile)
	return tx.Commit()
}

// UseProfile loads the given profile to the current config in the database.
func (db *store) UseProfile() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("UPDATE current_config SET profile = $1, region = $2", Profile, Region)
	return tx.Commit()
}
