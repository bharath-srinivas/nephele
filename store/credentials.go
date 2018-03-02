package store

type Credentials struct {
	AccessId  string `db:"access_id"`  // AWS Access Key ID
	SecretKey string `db:"secret_key"` // AWS Secret Access Key
	Region    string `db:"region"`     // AWS Region
}

// GetCredentials returns accessId, secretKey and region information for the current AWS profile.
func GetCredentials() *Credentials {
	db := NewSession()
	defer db.Close()

	return db.getCredentials()
}

func (db *store) getCredentials() *Credentials {
	var credentials Credentials
	db.Get(&credentials, "WITH t1 as (SELECT name FROM credentials), t2 as (SELECT profile FROM current_config) "+
		"SELECT t1.access_id, t1.secret_key, t2.region FROM credentials t1, current_config t2 WHERE t1.name = t2.profile")
	return &credentials
}

// SetCredentials creates a new entry in the database for the current AWS config, if not present, else returns error.
func (db *store) SetCredentials(accessId string, secretKey string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tx.Exec("INSERT INTO credentials (name, access_id, secret_key) VALUES ($1, $2, $3)", Profile, accessId, secretKey)
	tx.Exec("UPDATE current_config SET profile = $1, region = $2", Profile, Region)
	return tx.Commit()
}
