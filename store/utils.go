package store

// entryExists is a helper function that returns true if an entry exists in database for the given profile, else false.
func (db *store) entryExists(profile string) bool {
	var credentials Credentials

	db.Get(&credentials, "SELECT access_id FROM credentials WHERE name = $1", profile)

	if credentials.AccessId == "" {
		return false
	}

	return true
}

// currentProfile is a helper function which will return true if the specified profile is the current active profile,
// else false.
func (db *store) currentProfile(profile string) bool {
	var currentProfile string
	current := db.QueryRow("SELECT profile FROM current_config")
	current.Scan(&currentProfile)

	if profile == currentProfile {
		return true
	}
	return false
}
