package store

import (
	"testing"
)

func TestSetCredentials(t *testing.T) {
	db := newTest()
	defer db.Close()

	Profile = "test"
	Region = "us-east-1"

	if ok := db.EntryExists(Profile); ok {
		t.Errorf("Unexpected error: profile exists: %t", ok)
		return
	}

	if err := db.SetCredentials("ThisIsATestAccessId123", "ThisIsATestSecretKey456"); err != nil {
		t.Errorf("SetCredentials returned an error: %s", err)
	}
}

func TestGetCredentials(t *testing.T) {
	db := newTest()
	defer db.Close()

	want := &Credentials{
		AccessId:  "ThisIsATestAccessId123",
		SecretKey: "ThisIsATestSecretKey456",
		Region:    "us-east-1",
	}

	got := db.getCredentials()

	if *got != *want {
		t.Errorf("GetCredentials returned incorrect output, got: %+v, want: %+v", *got, *want)
	}
}
