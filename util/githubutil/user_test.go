package githubutil

import "testing"

func TestEmailForUserLogin(t *testing.T) {
	email, err := EmailForUserLogin("Focinfi")
	if err != nil {
		t.Fatal(err)
	}

	if email != "focinfi@gmail.com" {
		t.Errorf("can not get email for Focinfi, result: %s", email)
	}
}
