package saito

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cbergoon/merkletree"
	"github.com/mr-tron/base58"
)

// ToBase58 converts hex string to base58 string
func ToBase58(s string) string {
	keyString, _ := hex.DecodeString(s)
	return base58.Encode(keyString)
}

// FromBase58 converts base58 string to hex string
func FromBase58(s string) string {
	base58String, _ := base58.Decode(s)
	return hex.EncodeToString(base58String)
}

// DoubleHashB implements sha256 hashing twice as per bitcoin's standard
func DoubleHashB(b []byte) string {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return hex.EncodeToString(second[:])
}

// CompressPubKey compresses a public key to a 33 byte byte array
func CompressPubKey(pk string) (string, error) {
	// pubKeyBytes, err := hex.DecodeString()
	pubKeyBytes, err := hex.DecodeString(pk)
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return "Error thrown in Compress PublicKey", err
	}
	return hex.EncodeToString(pubKey.SerializeCompressed()), nil
}

// UncompressPubKey uncompresses a 33 byte public key to a 65 byte array
func UncompressPubKey(pk string) (string, error) {
	pubKeyBytes, _ := hex.DecodeString(pk)
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return "Error thrown in UncompressPubKey", err
	}
	return hex.EncodeToString(pubKey.SerializeUncompressed()), nil
}

//////////////////
// generateKeys //
//////////////////
//
// creates a public/private keypair. returns the string
// of the private key from which the public key can be
// re-generated.
//
// @returns {string} private key
//
// Crypt.prototype.generateKeys = function generateKeys() {
// 	let privateKey;
// 	do { privateKey = randomBytes(32) } while (!secp256k1.privateKeyVerify(privateKey, false))
// 	return privateKey.toString('hex');
// }

// GenerateKeys creates a new private key
func GenerateKeys() (*btcec.PrivateKey, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println(err)
	}
	return privKey, err
	// privKeyBytes := privKey.Serialize()
	// privKeyString := hex.EncodeToString(privKeyBytes)
	// return privKeyString, nil
}

/////////////////////
// returnPublicKey //
/////////////////////
//
// returns the public key associated with a private key
//
// @params {string} private key (hex)
// @returns {string} public key (hex)
//
//   Crypt.prototype.returnPublicKey = function returnPublicKey(privkey) {
// 		return this.compressPublicKey(secp256k1.publicKeyCreate(Buffer.from(privkey,'hex'), false).toString('hex'));
//   }

// func ReturnPublickKey(privateKey []byte) *PublicKey {
// 	// assume crypt struct possess context
// 	return secp256k1.EcPubkeyCreate(c.ctx, []byte(hex.DecodeString(privatekey)))
// }

// ReturnPublicKey returns the public key assigned to a private key
func ReturnPublicKey(privateKey *btcec.PrivateKey) *btcec.PublicKey {
	// pkBytes, err := hex.DecodeString(privateKey)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privateKey.Serialize())
	//return hex.EncodeToString(pubKey.SerializeCompressed())
	return pubKey
}

// ConvertPubKeyToString converts btcec.PublicKey to string
func ConvertPubKeyToString(pubKey *btcec.PublicKey) string {
	return hex.EncodeToString(pubKey.SerializeCompressed())
}

// ConvertPrivKeyToString converts btcec.PrivateKey to string
func ConvertPrivKeyToString(privKey *btcec.PrivateKey) string {
	return hex.EncodeToString(privKey.Serialize())
}

// SignMessage produces a signature from a string message and private key
func SignMessage(msg string, privKey *btcec.PrivateKey) (*btcec.Signature, []byte) {
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKey.Serialize())
	messageHash := chainhash.DoubleHashB([]byte(msg))

	signature, err := privateKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
	}

	return signature, messageHash
}

///////////////////
// verifyMessage //
///////////////////
//
// confirms that a message was signed by the private
// key associated with a provided public key
//
// @params {string} message to check
// @params {string} signature to confirm
// @params {string} public key of alleged signer
// @returns {boolean} is signature valid?
//
// Crypt.prototype.verifyMessage = function verifyMessage(msg, sig, pubkey) {
// 	try {
// 		return secp256k1.verify(Buffer.from(this.hash(Buffer.from(msg, 'utf-8').toString('base64')),'hex'), Buffer.from(this.fromBase58(sig),'hex'), Buffer.from(this.uncompressPublicKey(pubkey),'hex'));
// 	} catch (err) {
// 		return false;
// 	}
// }

//VerifyMessage returns boolean from
func VerifyMessage(messageHash []byte, pubKey *btcec.PublicKey, signature *btcec.Signature) bool {
	return signature.Verify(messageHash, pubKey)
}

// func VerifyMessage(ctx *Context, sig *EcdsaSignature, msg32 []byte, pk *PublicKey) (int, error) {
// 	return secp256k1.EcSeckeyVerify(ctx, sig, msg32, pk)
// }

/////////////////
// isPublicKey //
/////////////////
//
// finds out if we have a public key
//
// @params {string} publickey?
//
// Crypt.prototype.isPublicKey = function isPublicKey(publickey) {
// 	if (publickey.length == 44 || publickey.length == 45) {
// 		if (publickey.indexOf("@") > 0) {} else {
// 			return 1;
// 		}
// 	}
// 	return 0;
// }

//////////////////
// MERKLE TREES //
//////////////////

//////////////////////
// returnMerkleTree //
//////////////////////
//
// takes an array of strings and converts them into a merkle tree
// of SHA256 hashes.
//
// @params {array} array of strings
// @returns {merkle-tree}
//
// func returnMerkleTree(inarray) MerkleTree{
// 	var mt   = null;
// 	var args = { array: inarray, hashalgo: 'sha256', hashlist: false };
// 	merkle.fromArray(args, function (err, tree) { mt = tree; });
// 	return mt;
// }

// MerkleContent interface used to prepare BuildMerkleTree for transaction data
// type MerkleContent interface {
// 	CalculateHash() []byte
// 	Equals(merkletree.Content) bool
// }

//BuildMerkleTree returns a merkletree object
func BuildMerkleTree(input []Transaction) *merkletree.MerkleTree {
	// Transaction instead?
	//Create new Merkle Tree
	var list []merkletree.Content

	if len(input) != 0 {
		for _, content := range input {
			list = append(list, content)
		}
	}
	// need to provide some content, even if empty string
	t, _ := merkletree.NewTree(list)

	return t
}

// ReturnMerkleTreeRoot returns the root of the merkle tree
func ReturnMerkleTreeRoot(input []Transaction) []byte {
	// transaction instead?
	return BuildMerkleTree(input).MerkleRoot()
}

////////////////////
// DIFFIE HELLMAN //
////////////////////
//
// The DiffieHellman process allows two people to generate a shared
// secret in an environment where all information exchanged between
// the two can be observed by others.
//
// It is used by our encryption module to generate shared secrets,
// but is generally useful enough that we include it in our core
// cryptography class
//
// see the "encryption" module for an example of how to generate
// a shared secret using these functions
//

/////////////////////////
// createDiffieHellman //
/////////////////////////
//
// @params {string} public key
// @params {string} private key
// @returns {DiffieHellman object} ecdh
//

// mabye should be []bytes?

// func createDiffieHellman(pubkey string, privkey string) string {
// 	var ecdh   = crypto.createECDH("secp256k1");
// 	ecdh.generateKeys();
// 	if (pubkey != "")  { ecdh.setPublicKey(pubkey); }
// 	if (privkey != "") { ecdh.setPrivateKey(privkey); }
// 	return ecdh;
// }

//CreateSharedSecret generates a shared secret from pubkey and privkey, returns byte array
func CreateSharedSecret(privkey *btcec.PrivateKey, pubkey *btcec.PublicKey) []byte {
	return btcec.GenerateSharedSecret(privkey, pubkey)
}

//////////////////////////////
// returnDiffieHellmanKeys //
/////////////////////////////
//
// Given a Diffie-Hellman object, fetch the keys
//
// @params {DiffieHellman object} dh
// @returns {{pubkey:"",privkey:""}} object with keys
//
// func returnDiffieHellmanKeys(dh string) (string, string) {
// 	var keys = {};
// 	keys.pubkey  = dh.getPublicKey(null, "compressed");
// 	keys.privkey = dh.getPrivateKey(null, "compressed");
// 	return keys;
// }

///////////////////////////////
// createDiffieHellmanSecret //
//////////////////////////////
//
// Given your private key and your counterparty's public
// key and an extra piece of information, you can generate
// a shared secret.
//
// @params {DiffieHellman object} counterparty DH
// @params {string} my_publickey
//
// @returns {{pubkey:"",privkey:""}} object with keys
//
// func createDiffieHellmanSecret(a_dh, b_pubkey) {
// 	return a_dh.computeSecret(b_pubkey);
// }

////////////////////////////////
// AES SYMMETRICAL ENCRYPTION //
////////////////////////////////
//
// once we have a shared secret (possibly generated through the
// Diffie-Hellman method above), we can use it to encrypt and
// decrypt communications using a symmetrical encryption method
// like AES.
//

////////////////
// aesEncrypt //
////////////////
//
// @param {string} msg to encrypt
// @param {string} shared secret
// @returns {string} json object
//
// func aesEncrypt(msg string, secret string) {
// 	var rp = new Buffer(secret.toString("hex"), "hex").toString("base64");
// 	var en = CryptoJS.AES.encrypt(msg, rp, { format: JsonFormatter });
// 	return en.toString();
// }

// Encrypt encrypts a message using elliptic compute
func Encrypt(pubKey *btcec.PublicKey, message string) []byte {
	ciphertext, err := btcec.Encrypt(pubKey, []byte(message))
	if err != nil {
		fmt.Println(err)
	}
	return ciphertext
}

// Decrypt decrypts a message using a private key
func Decrypt(privKey *btcec.PrivateKey, ciphertext []byte) []byte {
	plaintext, err := btcec.Decrypt(privKey, ciphertext)
	if err != nil {
		fmt.Println(err)
	}
	return plaintext
}

// example of AES encryption
// func ExampleNewGCM_encrypt() {
// 	// Load your secret key from a safe place and reuse it across multiple
// 	// Seal/Open calls. (Obviously don't use this example key for anything
// 	// real.) If you want to convert a passphrase to a key, use a suitable
// 	// package like bcrypt or scrypt.
// 	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
// 	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
// 	plaintext := []byte("exampleplaintext")

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
// 	nonce := make([]byte, 12)
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		panic(err.Error())
// 	}

// 	aesgcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
// 	fmt.Printf("%x\n", ciphertext)
// }

////////////////
// aesDecrypt //
////////////////
//
// @param {string} encrypted json object from aesEncrypt
// @param {string} shared secret
// @returns {string} unencrypted

// func aesDecrypt(msg, secret) {
// 	var rp = new Buffer(secret.toString("hex"), "hex").toString("base64");
// 	var de = CryptoJS.AES.decrypt(msg, rp, { format: JsonFormatter });
// 	return CryptoJS.enc.Utf8.stringify(de);
// }

// example of AES Decryption

// func ExampleNewGCM_decrypt() {
// 	// Load your secret key from a safe place and reuse it across multiple
// 	// Seal/Open calls. (Obviously don't use this example key for anything
// 	// real.) If you want to convert a passphrase to a key, use a suitable
// 	// package like bcrypt or scrypt.
// 	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
// 	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
// 	ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
// 	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	aesgcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	fmt.Printf("%s\n", plaintext)
// 	// Output: exampleplaintext
// }
