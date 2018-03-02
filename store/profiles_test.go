package store

import "testing"

func TestListProfiles(t *testing.T) {
	db := newTest()
	defer db.Close()

	rows := db.ListProfiles()

	for rows.Next() {
		var got string
		rows.Scan(&got)

		if want := "test"; got != want {
			t.Errorf("ListProfiles returned incorrect env name, got: %s, want: %s", got, want)
		}

		if ok := db.CurrentProfile(got); !ok {
			t.Errorf("Unexpected error: current active profile, got: %t, want: %t", ok, true)
		}
	}
}

func TestUseProfile(t *testing.T) {
	db := newTest()
	defer db.Close()

	Profile = "test1"
	Region = "us-west-1"

	if err := db.SetCredentials("AnotherAccessId", "AnotherSecretKey"); err != nil {
		t.Errorf("Unexpected error: insert credentials: %s", err)
		return
	}

	if ok := db.EntryExists(Profile); !ok {
		t.Errorf("Unexpected error: profile exists: %t", ok)
		return
	}

	if err := db.UseProfile(); err != nil {
		t.Errorf("UseProfile returned an error: %s", err)
	}
}

func TestDeleteProfile(t *testing.T) {
	db := newTest()
	defer func() {
		db.Close()
		cleanup()
	}()

	profile := "test"

	if ok := db.EntryExists(profile); !ok {
		t.Errorf("Unexpected error: profile exists: %t", ok)
		return
	}

	if ok := db.CurrentProfile(profile); ok {
		t.Errorf("Unexpected error: current active profile: %t", ok)
		return
	}

	if err := db.DeleteProfile(profile); err != nil {
		t.Errorf("DeleteProfile returned an error: %s", err)
	}
}
