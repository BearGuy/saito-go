package crypt_test

import (
	"testing"
)

func TestToBase58(t *testing.T) {
	hexString := "22a47fa09a223f2aa079edf85a7c2d4f87"
	calcBase58String := ToBase58(hexString)
	correctBase58String := "28pMtVHPTn8uek9kGF4GMweSYzGeePxBCpamk3tj6p8bhAn"
	if calcBase58String != correctBase58String {
		t.Errorf("ToBase58 was incorrect, got: %s, want: %s.", calcBase58String, correctBase58String)
	}
}
