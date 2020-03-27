package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestEncodeMinimal(t *testing.T) {
	spayd := Spayd{Account: "Hello", Amount: 12.34, Currency: "CZK"}
	bytes, err := spayd.Encode()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	encoded := string(bytes)
	if "SPD*1.0*ACC:Hello*AM:12.34*CC:CZK" != encoded {
		fmt.Println("Should be \"SPD*1.0*ACC:Hello*AM:12.34*CC:CZK\"; is " + encoded)
		t.Fail()
	}
}

func TestEncodeCurrecntyNot3Long(t *testing.T) {
	spayd := Spayd{Account: "Hello", Amount: 12.34, Currency: "CZKX"}
	_, err := spayd.Encode()
	if err == nil {
		t.Fail()
	}
	if !strings.Contains(err.Error(), "wrong length") {
		fmt.Println("Error should contain \"wrong length\"\nError: " + err.Error())
		t.Fail()
	}
}

func TestEncodeRefTooLong(t *testing.T) {
	spayd := Spayd{Account: "Hello", Amount: 12.34, Currency: "CZK", Ref: "1234567890123456X"}
	_, err := spayd.Encode()
	if err == nil {
		t.Fail()
	}
	if !strings.Contains(err.Error(), "too long") {
		fmt.Println("Error should contain \"too long\".\nError: " + err.Error())
		t.Fail()
	}
}
