package token

import (
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	secret := "sqs.secret"
	params := map[string]interface{}{
		"userID": "1",
	}
	exp := time.Second
	tkn := NewJWT()
	code, err := tkn.Make(secret, params, exp)
	if err != nil {
		t.Fatal(err)
	}

	p, err := tkn.Verify(code, secret)
	if err != nil {
		t.Fatal(err)
	}

	if p["userID"] != "1" {
		t.Errorf("can not parse a correct token, params is %v\n", p)
	}

	// wait for check the expiration
	time.Sleep(time.Millisecond * 1001)

	_, err = tkn.Verify(code, secret)
	if err == nil {
		t.Error("can not detect the expiration")
	}
}
