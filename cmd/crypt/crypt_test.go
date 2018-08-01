package crypt_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bearguy/saito-go/cmd/crypt"
	"github.com/btcsuite/btcd/btcec"
)

func TestToBase58(t *testing.T) {
	hexString := "22a47fa09a223f2aa079edf85a7c2d4f87"
	encodeToBase58String := crypt.ToBase58(hexString)
	correctBase58String := "Kt7xi8ujBNVfdD7Pxs8X22A"
	if encodeToBase58String != correctBase58String {
		t.Errorf("ToBase58 was incorrect, got: %s, want: %s.", encodeToBase58String, correctBase58String)
	}
}

func TestFromBase58(t *testing.T) {
	base58String := "Kt7xi8ujBNVfdD7Pxs8X22A"
	decodeFromBase58String := crypt.FromBase58(base58String)
	correctHexString := "22a47fa09a223f2aa079edf85a7c2d4f87"
	if decodeFromBase58String != correctHexString {
		t.Errorf("FromBase58 was incorrect, got: %s, want: %s.", decodeFromBase58String, correctHexString)
	}
}

func TestDoubleHashB(t *testing.T) {
	message := "this is a test message"
	messageHash := crypt.DoubleHashB([]byte(message))
	correctMessageHash := "a949c222c34b9a59b35214d5a926c406f21b0be5b8311f8dcdbf09dd7b84c242"
	if messageHash != correctMessageHash {
		t.Errorf("DoubleHashB was incorrect, got: %s, want: %s.", messageHash, correctMessageHash)
	}
}

func TestCompressPublicKey(t *testing.T) {
	hexString := "04115c42e757b2efb7671c578530ec191a1" +
		"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
		"21010db67ac11b1b51b651953d2"

		// "02a673638cb9587cb68ea08dbef685c" +
		// "6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5"
	compressedPubKey, err := crypt.CompressPubKey(hexString)
	if err != nil {
		fmt.Println(err)
	}
	correctBase58String := "02115c42e757b2efb7671c578530ec191a1359381e6a71127a9d37c486fd30dae5"
	if compressedPubKey != correctBase58String {
		t.Errorf("TestCompressPublicKey was incorrect, got: %s, want: %s.", compressedPubKey, correctBase58String)
	}
}

func TestUncompressPublicKey(t *testing.T) {
	base58String := "02115c42e757b2efb7671c578530ec191a1359381e6a71127a9d37c486fd30dae5"
	decodeFromBase58String, err := crypt.UncompressPubKey(base58String)
	if err != nil {
		fmt.Println(err)
	}
	correctHexString := "04115c42e757b2efb7671c578530ec191a1" +
		"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
		"21010db67ac11b1b51b651953d2"
	if decodeFromBase58String != correctHexString {
		t.Errorf("TestUncompressPublicKey was incorrect, got: %s, want: %s.", decodeFromBase58String, correctHexString)
	}
}

func TestGenerateKeys(t *testing.T) {
	privateKey, _ := crypt.GenerateKeys()
	privateKeyString := crypt.ConvertPrivKeyToString(privateKey)
	if len(privateKeyString) != 64 {
		t.Errorf("TestGenerateKeys private key length was incorrect, got: %d, want: 64", len(privateKeyString))
	}
}

func TestReturnPublicKey(t *testing.T) {
	pkBytes, err := hex.DecodeString("b5378ce7d7a4e571b9984158e0e913ecd5bb3c931a0d5e7fe220f0588de14da8")
	if err != nil {
		fmt.Println(err)
		return
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	pubKey := crypt.ReturnPublicKey(privKey)
	publicKeyString := crypt.ConvertPubKeyToString(pubKey)
	if len(publicKeyString) != 66 {
		t.Errorf("TestGenerateKeys private key length was incorrect, got: %d, want: 66", len(publicKeyString))
	}
}

func TestSignMessage(t *testing.T) {
	privKey, _ := crypt.GenerateKeys()
	pubKey := crypt.ReturnPublicKey(privKey)
	msg := "this is it"
	signature, messageHash := crypt.SignMessage(msg, privKey)

	//verified := signature.Verify(messageHash, pubKey)
	verified := crypt.VerifyMessage(messageHash, pubKey, signature)
	fmt.Printf("Signature Verified? %v\n", verified)
	if verified != true {
		t.Errorf("TestSignMessage verified value was false")
	}
}
