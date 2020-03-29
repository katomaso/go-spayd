package main

import (
	"strings"
	"testing"
)

var (
	minimalSpayd   = Spayd{Account: "CZ0000000000123456789012", Amount: 12.34, Currency: "CZK"}
	minimalEncoded = "SPD*1.0*ACC:CZ0000000000123456789012*AM:12.34*CC:CZK"
)

func TestEncodeMinimal(t *testing.T) {
	bytes, err := minimalSpayd.Encode()
	if err != nil {
		t.Fatal(err)
	}
	encoded := string(bytes)
	if minimalEncoded != encoded {
		t.Fatalf("Should be \"%s\"; is \"%s\"", minimalEncoded, encoded)
	}
}

func TestEncodeCurrecntyNot3Long(t *testing.T) {
	spayd := minimalSpayd // perform a copy
	spayd.Currency = "CZKXX"
	_, err := spayd.Encode()
	if err == nil {
		t.Fatal("No error was returned")
	}
	if !strings.Contains(err.Error(), "wrong length") {
		t.Fatalf("Error should contain \"wrong length\"\nError: %s", err.Error())
	}
}

func TestEncodeRefTooLong(t *testing.T) {
	spayd := minimalSpayd
	spayd.Ref = "1234567890123456X"
	_, err := spayd.Encode()
	if err == nil {
		t.Fatal("No error was returned")
	}
	if !strings.Contains(err.Error(), "too long") {
		t.Fatalf("Error should contain \"too long\".\nError: %s", err.Error())
	}
}

func TestEncodeAccNotIBAN(t *testing.T) {
	spayd := minimalSpayd
	spayd.Account = "NotIBAN"
	_, err := spayd.Encode()
	if err == nil {
		t.Fatal("No error was returned")
	}
	if !strings.Contains(err.Error(), "IBAN") {
		fmt.Println("Error should contain \"IBAN\".\nError: " + err.Error())
		t.Fail()
	}
}
