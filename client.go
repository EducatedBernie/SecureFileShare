package client

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"
	// _ "testing/quick"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	"strconv"

	"bytes"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username          string
	PrivateDecKey     userlib.PrivateKeyType
	PrivateDigitalSig userlib.DSSignKey

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

type File struct {
	Owner []byte
	// contentKey   []byte
	UUIDContents     uuid.UUID
	SharedUsersUUID  uuid.UUID
	InvitedUsersUUID uuid.UUID
	NextNode         uuid.UUID
	IsDummy          bool
	LastNode         uuid.UUID
	NumAppends		 int
}
type Hash struct {
	Hash []byte
}

type Mailbox struct {
	FileStructID uuid.UUID
	FileKey      []byte
}

type HMAC struct {
	Hmac       []byte
	Encryption []byte
}

type Invitation struct {
	IntNodeUUID uuid.UUID
	Key         []byte
}

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdataptr *User, err error) {
	//Generate userStruct
	var userdata User
	userdata.Username = username

	//Generate username hash to get the userStruct
	usernameBytes := []byte(username)
	hashed := userlib.Hash(usernameBytes)[:16]
	userUUID, err := uuid.FromBytes(hashed)
	if err != nil {
		return nil, errors.New("InitUser: Unable to get uuid of usernameHash")
	}

	//Check if username is empty
	if username == "" {
		return nil, errors.New("Username is empty")
	}
	//Check if user is already created (same username)
	_, ok := userlib.DatastoreGet(userUUID)
	if ok {
		return nil, errors.New("username has already been created")
	}
	// Create public and private key
	var publicEncKey, privateEncKey, check = userlib.PKEKeyGen()
	userdata.PrivateDecKey = privateEncKey
	if check != nil {
		return nil, errors.New("InitUser: User public key generation returned an error")
	}

	var privateDigitKey, publicDigitKey, check1 = userlib.DSKeyGen()
	userdata.PrivateDigitalSig = privateDigitKey
	if check1 != nil {
		return nil, errors.New("InitUser: User public key generation returned an error")
	}

	// Store public encryption and public digital key under username in keyStore
	userlib.KeystoreSet(userdata.Username+"a", publicEncKey)
	userlib.KeystoreSet(userdata.Username+"b", publicDigitKey)
	bytesPassword := []byte(password)

	// Password Based Deterministic Generator creates the key to encrypt/decrypt userStruct
	userPassKey := userlib.Argon2Key(bytesPassword, usernameBytes, 16)

	//Create random IV (not saved)
	iv := userlib.RandomBytes(16)
	//Marshal the userStruct
	marshaledUserStruct, check2 := json.Marshal(userdata)
	if check2 != nil {
		return nil, errors.New("InitUser: User json not marshalled properly")
	}

	// Encrypt the user struct, using the Passkey + IV combo we generated
	userCipher := userlib.SymEnc(userPassKey, iv, marshaledUserStruct)

	// Create a signature HMAC(user's priv key, ciphertext)
	userHMAC, check3 := userlib.HMACEval(userPassKey, userCipher)
	if check3 != nil {
		return nil, errors.New("InitUser: userHmac signing not done properly")
	}

	// Store userCipher||userHmac @ UUID(username)

	// Create hmac struct
	var hmacStruct HMAC
	hmacStruct.Hmac = userHMAC
	hmacStruct.Encryption = userCipher

	// Marshal and encrypt hmacStruct
	marshalHmacStruct, hmacCheck := json.Marshal(hmacStruct)
	if hmacCheck != nil {
		return nil, errors.New("InitUser: hmacStruct not done properly")
	}

	userlib.DatastoreSet(userUUID, marshalHmacStruct)

	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata

	// Turn the password and usernames
	bytesPassword := []byte(password)
	usernameBytes := []byte(username)

	//Create password Key (PBKGen)
	userPassKey := userlib.Argon2Key(bytesPassword, usernameBytes, 16)

	//Find the UUID we're interested in using UUID(username)
	userUUID, check0 := uuid.FromBytes(userlib.Hash(usernameBytes)[:16])
	if check0 != nil {
		return nil, errors.New("GetUser: weren't able to find the UUID or some error")
	}
	uuidValue, check1 := userlib.DatastoreGet(userUUID)
	if !check1 {
		return nil, errors.New("GetUser: No UUID with the corresponding username")
	}

	var unMarshaledHMACStruct HMAC
	unmarshalCheck := json.Unmarshal(uuidValue, &unMarshaledHMACStruct)
	if unmarshalCheck != nil {
		return nil, errors.New("GetUser: unmarshaled goes wrong")
	}

	userCipher := unMarshaledHMACStruct.Encryption

	userHMAC, check3 := userlib.HMACEval(userPassKey, userCipher)
	if check3 != nil {
		return nil, errors.New("GetUser: userHmac signing not done properly")
	}

	givenHMAC := unMarshaledHMACStruct.Hmac
	hmacEqual := userlib.HMACEqual(givenHMAC, userHMAC)

	if !hmacEqual {
		return nil, errors.New("getUser: Incorrect Password")
	}

	userStructDec := userlib.SymDec(userPassKey, userCipher)
	err = json.Unmarshal(userStructDec, userdataptr)
	if err != nil {
		return nil, err
	}

	//Decrypt that shit with the password key we just generated at the start
	return userdataptr, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	storageKey, err := userdata.GetMailboxUUID(filename)
	if err != nil {
		return err
	}
	// contentBytes, err := json.Marshal(content)
	// if err != nil {
	// 	return err
	// }

	contentBytes := content

	//Check if filename exists already in user's space
	mailboxValue, check := userlib.DatastoreGet(storageKey)

	var iv []byte = userlib.RandomBytes(16)
	switch {
	case check:
		//At uuid = storageKey, check signature
		userPublicVerifyKey, ok := userlib.KeystoreGet(userdata.Username + "b")
		if !ok {
			return errors.New("can't find signature key of user")
		}

		//Separate signed bytes from ciphertext
		signedBytes, cipheredMailBox := mailboxValue[:256], mailboxValue[256:]

		err := userlib.DSVerify(userPublicVerifyKey, cipheredMailBox, signedBytes)
		if err != nil {
			return errors.New("mailbox has been compromised")
		}

		//Decrypt mailbox
		plaintext, err := userlib.PKEDec(userdata.PrivateDecKey, cipheredMailBox)
		if err != nil {
			return err
		}

		var firstMailBox Mailbox
		err = json.Unmarshal(plaintext, &firstMailBox)
		if err != nil {
			return err
		}

		keyToDecryptIntNode, uuidIntNode := firstMailBox.FileKey, firstMailBox.FileStructID

		intNode, err := FetchIntNode(uuidIntNode, keyToDecryptIntNode)
		if err != nil {
			return err
		}

		hmacStruct, err := FetchHmacStruct(intNode.FileStructID)
		if err != nil {
			return err
		}

		if !checkHmacStruct(hmacStruct, intNode.FileKey) {
			return errors.New("file has been breached/No integrity by hmac")
		}

		//Decrypt hmacStruct's encryption
		headFileStruct, err := FetchHeadFileStruct(hmacStruct.Encryption, intNode.FileKey)
		if err != nil {
			return err
		}

		//Encrypt the contents
		var encryptedContent []byte = symEnc(intNode.FileKey, iv, contentBytes)

		//Overwrite hmac of contents (headFile.contentUUID = hmac || encryptedContent)
		newHmac, err := calculateHMAC(intNode.FileKey, encryptedContent)
		if err != nil {
			return err
		}


		//Set into datastore at the original UUID
		datastoreSet(headFileStruct.UUIDContents, append(newHmac, encryptedContent...))


		//Overwrite any appends the file struct may contain
		headFileStruct.NextNode = uuid.Nil
		headFileStruct.LastNode = uuid.Nil

		//marshal the changed fileStruct
		newMarshaledFileStruct, err := json.Marshal(headFileStruct)
		if err != nil {
			return err
		}

		iv = userlib.RandomBytes(16)
		encryptedMarshaledFileStruct := symEnc(intNode.FileKey, iv, newMarshaledFileStruct)

		//Create new hmacStruct
		newHmacStruct, err := createHmacStruct(intNode.FileKey, encryptedMarshaledFileStruct)
		if err != nil {
			return err
		}
		//Marshal
		newMarshaledHmacStruct, err := json.Marshal(newHmacStruct)
		if err != nil {
			return err
		}
		//Place marshaled hmacStruct into the old uuid
		datastoreSet(intNode.FileStructID, newMarshaledHmacStruct)

		return nil

		//Now start encrypting fileStruct

	// Else, if the file don't exist then...
	case !check:
		// Generate a random Symkey (filekey), encrypt tthe contents of the file
		// Create a fileStruct (Sentinel Node)
		//Initialize all the data to appropriate values (owner, next etc)
		// Store ciphertext + HMAC(Key, Ciphertext) at the appropriate HMAC struct
		var fileKey = userlib.RandomBytes(16)

		//Create dummyFileStruct and that is where the fileStruct is located at
		headFile, err := createHeadFileStruct(userdata.Username, content, fileKey)
		if err != nil {
			return err
		}
		headFile.InvitedUsersUUID = uuid.New()
		headFile.SharedUsersUUID = uuid.New()

		// err = encMarshalList(headFile.InvitedUsersUUID, nil, fileKey, "invitedUsers")
		// if err != nil {
		// 	return err
		// }
		err = encMarshalList(headFile.SharedUsersUUID, nil, fileKey, "sharedUsers")
		if err != nil {
			return err
		}
		marshaledHeadFileStruct, err := json.Marshal(headFile)
		if err != nil {
			return err
		}
		iv := userlib.RandomBytes(16)
		//encrypt marshaled fileStruct
		encryptMarshaledHeadFileStruct := userlib.SymEnc(fileKey, iv, marshaledHeadFileStruct)

		//create Hmac struct
		var hmacStruct HMAC
		hmacStruct, err = createHmacStruct(fileKey, encryptMarshaledHeadFileStruct)
		if err != nil {
			return err
		}

		//marshal the hmacStruct
		marshaledHmacStruct, err := json.Marshal(hmacStruct)
		if err != nil {
			return err
		}

		//Create UUID to store HMACStruct

		uuidHmacStruct := uuid.New()
		userlib.DatastoreSet(uuidHmacStruct, marshaledHmacStruct)

		uuidIntNode := uuid.New()
		intNode := createMailBox(uuidHmacStruct, fileKey)

		marshaledIntNode, err := json.Marshal(intNode)
		if err != nil {
			return err
		}

		keyToEncryptIntNode := userlib.RandomBytes(16)

		//Encrypt marshaled IntNode

		cipherMarshaledIntNode := userlib.SymEnc(keyToEncryptIntNode, iv, marshaledIntNode)

		hmac, err := calculateHMAC(keyToEncryptIntNode, cipherMarshaledIntNode)
		if err != nil {
			return err
		}

		userlib.DatastoreSet(uuidIntNode, append(hmac, cipherMarshaledIntNode...))

		//Create mailbox
		mailBox := createMailBox(uuidIntNode, keyToEncryptIntNode)

		//Marshal mailbox
		marshaledMailBox, err := json.Marshal(mailBox)
		if err != nil {
			return err
		}

		//Public-private key encryption

		userPublicKey, ok := userlib.KeystoreGet(userdata.Username + "a")
		if !ok {
			return errors.New("can't get key")
		}

		cipherMarshaledMailBox, err := userlib.PKEEnc(userPublicKey, marshaledMailBox)
		if err != nil {
			return err
		}
		//Sign encryption

		signedCipherMailbox, err := userlib.DSSign(userdata.PrivateDigitalSig, cipherMarshaledMailBox)
		if err != nil {
			return err
		}
		uuidMailBox := storageKey //StorageKey was calculated at the very first line of func
		datastoreSet(uuidMailBox, append(signedCipherMailbox, cipherMarshaledMailBox...))

	}

	return err
}

func symEnc(key []byte, iv []byte, plaintext []byte) []byte {
	return userlib.SymEnc(key, iv, plaintext)
}

func datastoreSet(uuid uuid.UUID, bytes []byte) {
	userlib.DatastoreSet(uuid, bytes)
}

func checkHmacStruct(hmac HMAC, key []byte) bool {
	//Check hmac with just calculated hmac
	checkHmac, err := userlib.HMACEval(key, hmac.Encryption)
	if err != nil {
		return false
	}

	check := userlib.HMACEqual(checkHmac, hmac.Hmac)
	return check
}

func createHmacStruct(key []byte, ciphertext []byte) (HMAC, error) {
	var hmacData HMAC
	hmacData.Encryption = ciphertext
	hmac, err := calculateHMAC(key, ciphertext)
	if err != nil {
		return hmacData, err
	}
	hmacData.Hmac = hmac
	return hmacData, err
}

func calculateHMAC(key []byte, ciphertext []byte) ([]byte, error) {
	hmac, err := userlib.HMACEval(key, ciphertext)
	if err != nil {
		return nil, errors.New("calculateHMAC: I'm in pain")
	}
	return hmac, nil
}

func checkHMACFileContents(fileData File, key []byte) error { //Given fileStruct check the HMAC of file Contents
	fileContents, ok := userlib.DatastoreGet(fileData.UUIDContents)
	if !ok {

		return errors.New("cannot get object")
	}

	hmacCipher, cipher := fileContents[:64], fileContents[64:]
	calculated, err := calculateHMAC(key, cipher)
	if err != nil {
		return err
	}

	ok = userlib.HMACEqual(calculated, hmacCipher)
	if !ok {
		return errors.New("problem with hmac of contents")
	}

	return nil

}

func createHeadFileStruct(owner string, content []byte, fileKey []byte) (File, error) {
	var filedata File
	filedata.Owner = userlib.Hash([]byte(owner))
	filedata.UUIDContents = uuid.New()
	iv := userlib.RandomBytes(16)
	var encryptedContent = userlib.SymEnc(fileKey, iv, content)
	contentHmac, err := calculateHMAC(fileKey, encryptedContent)
	if err != nil {
		return filedata, err
	}
	userlib.DatastoreSet(filedata.UUIDContents, append(contentHmac, encryptedContent...))
	filedata.IsDummy = true
	filedata.LastNode = uuid.Nil
	return filedata, nil
}
func createFileStruct(owner string, content []byte, fileKey []byte) (File, error) {
	// initialize owner's name
	// generate random key with a randomly generated
	// encrypt the contents
	var filedata File
	filedata.Owner = userlib.Hash([]byte(owner))
	filedata.UUIDContents = uuid.New()
	iv := userlib.RandomBytes(16)
	var encryptedContent = userlib.SymEnc(fileKey, iv, content)
	contentHmac, err := calculateHMAC(fileKey, encryptedContent)
	if err != nil {
		return filedata, err
	}
	userlib.DatastoreSet(filedata.UUIDContents, append(contentHmac, encryptedContent...))

	return filedata, nil

}

func createMailBox(uuid uuid.UUID, key []byte) Mailbox {
	var mailbox Mailbox
	mailbox.FileKey = key
	mailbox.FileStructID = uuid
	return mailbox
}

func createInvitationStruct(uuid uuid.UUID, key []byte) Invitation {
	var invitationData Invitation
	invitationData.Key = key
	invitationData.IntNodeUUID = uuid
	return invitationData
}

//Fetches the FileContents (get, decrypt)
func fetchFileContents(uuid uuid.UUID, key []byte) ([]byte, error) {
	userObj, ok := userlib.DatastoreGet(uuid)
	if !ok {
		return nil, errors.New("fetchFileContents: cannot get object")
	}

	hmacFileContents, cipher := userObj[:64], userObj[64:]
	calculatedHmac, err := calculateHMAC(key, cipher)
	if err != nil {
		return nil, err
	}

	ok = userlib.HMACEqual(hmacFileContents, calculatedHmac)
	if !ok {
		return nil, errors.New("fetchFileContents: hmac not equal")
	}

	plaintext := userlib.SymDec(key, cipher)


	return plaintext, nil

}

func (userdata *User) FetchMailBox(uuid uuid.UUID) (Mailbox, error) {
	userObj, ok := userlib.DatastoreGet(uuid)
	var check error = nil
	var maildata Mailbox
	if !ok {
		check = errors.New("datastore can't get value")
		return maildata, check
	}
	verifykey, ok := userdata.getVerifySignatureKey()
	if !ok {
		check = errors.New("verifyKey is hating us")
		return maildata, check
	}
	//Get the signature and encrypted text
	signedBytes, ciphertext := userObj[:256], userObj[256:]
	//verify signature
	check = userlib.DSVerify(verifykey, ciphertext, signedBytes)
	if check != nil {
		return maildata, check
	}

	//
	plaintext, err := userlib.PKEDec(userdata.PrivateDecKey, ciphertext)
	if err != nil {
		return maildata, err
	}

	check = json.Unmarshal(plaintext, &maildata)
	return maildata, check

}

func FetchIntNode(uuid uuid.UUID, key []byte) (Mailbox, error) {
	userObj, ok := userlib.DatastoreGet(uuid)
	var check error = nil
	var intNodeData Mailbox
	if !ok {
		check = errors.New("datastore can't get value")
		return intNodeData, check
	}

	hmac, encMarshaledIntNode := userObj[:64], userObj[64:]
	calculatedHmac, err := calculateHMAC(key, encMarshaledIntNode)
	if err != nil {
		return intNodeData, err
	}

	ok = userlib.HMACEqual(hmac, calculatedHmac)
	if !ok {
		return intNodeData, errors.New("FetchIntNode: intNode has been compromised")
	}

	marshaledPlaintext := userlib.SymDec(key, encMarshaledIntNode)

	check = json.Unmarshal(marshaledPlaintext, &intNodeData)
	return intNodeData, check

}

//Get value at uuid and Unmarshal (DO NOT DECRYPT)
func FetchHmacStruct(uuid uuid.UUID) (HMAC, error) {
	userObj, ok := userlib.DatastoreGet(uuid)
	var check error = nil
	var hmacdata HMAC
	if !ok {
		check = errors.New("datastore can't get value")
		return hmacdata, check
	}
	check = json.Unmarshal(userObj, &hmacdata)
	return hmacdata, check

}
func FetchFileStruct(uuid uuid.UUID, key []byte) (File, error) {
	var filedata File
	userObj, ok := userlib.DatastoreGet(uuid)
	if !ok {
		return filedata, errors.New("datastore can't get value")
	}

	hmacFile, cipher := userObj[:64], userObj[64:]
	//check hmac of file

	calculatedHmac, err := calculateHMAC(key, cipher)
	if err != nil {
		return filedata, err
	}
	ok = userlib.HMACEqual(hmacFile, calculatedHmac)
	if !ok {
		return filedata, errors.New("fetchFileStruct: fileStruct has been compromised")
	}

	plaintext := userlib.SymDec(key, cipher)
	err = json.Unmarshal(plaintext, &filedata)
	if err != nil {
		return filedata, errors.New("fetchFileStruct: unable to unmarshal")
	}

	return filedata, nil
}

func FetchHeadFileStruct(ciphertext []byte, key []byte) (File, error) {
	var check error = nil
	var dummydata File
	plaintext := userlib.SymDec(key, ciphertext)
	check = json.Unmarshal(plaintext, &dummydata)
	return dummydata, check

}

func (userdata *User) getPublicEncKey() (userlib.PKEEncKey, bool) {
	return userlib.KeystoreGet(userdata.Username + "a")
}

func (userdata *User) GetMailboxUUID(filename string) (uuid.UUID, error) {
	// takes in a string, fetches the UUID of the mailbox
	thingToHash := filename + "-" + userdata.Username
	return uuid.FromBytes(userlib.Hash([]byte(thingToHash))[:16])
}

func (userdata *User) getVerifySignatureKey() (userlib.PKEEncKey, bool) {
	return userlib.KeystoreGet(userdata.Username + "b")
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	// Find the mailbox @ UUID(Hash(myname)||Hash(filename))
	mailboxID, check1 := userdata.GetMailboxUUID(filename)
	if check1 != nil {
		return errors.New("AppendToFile: did not get MailboxUUID successfully")
	}
	// mailbox
	mailboxData, check2 := userdata.FetchMailBox(mailboxID)
	if check2 != nil {
		return errors.New("AppendToFile: did not fetchWithPKE mailbox @ UUID successfully")
	}

	intNodeData, err := FetchIntNode(mailboxData.FileStructID, mailboxData.FileKey)
	if err != nil {
		return err
	}

	hmacData, err := FetchHmacStruct(intNodeData.FileStructID)
	if err != nil {
		return err
	}

	//Verify Hmac(ciphertext, key)
	if !checkHmacStruct(hmacData, intNodeData.FileKey) {
		return errors.New("hmac of fileStruct (not contents) has been compromised")
	}

	//Decrypt headFileStruct
	headFileStructData, err := FetchHeadFileStruct(hmacData.Encryption, intNodeData.FileKey)
	if err != nil {
		return err
	}
	//Change the number of appends here
	headFileStructData.NumAppends += 1
	appendFileContentKey, err := userlib.HashKDF(intNodeData.FileKey, []byte("content" + strconv.Itoa(headFileStructData.NumAppends)))
	if err != nil {
		return err
	}

	//Create new appendNode (File)
	appendNode, err := createFileStruct("", content, appendFileContentKey[:16])
	if err != nil {
		return err
	}
	
	appendNode.Owner = headFileStructData.Owner

	//Marshal appendNode
	marshaledAppendNode, err := json.Marshal(appendNode)
	if err != nil {
		return err
	}
	iv := userlib.RandomBytes(16)

	appendFileKey, err := userlib.HashKDF(intNodeData.FileKey, []byte("append" + strconv.Itoa(headFileStructData.NumAppends)))
	if err != nil {
		return err
	}
	encryptMarshaledAppendNode := symEnc(appendFileKey[:16], iv, marshaledAppendNode)

	//Create hmac value of the appendNode
	hmacValueAppendNode, err := calculateHMAC(appendFileKey[:16], encryptMarshaledAppendNode)
	if err != nil {
		return err
	}
	newAppendUUID := uuid.New()
	datastoreSet(newAppendUUID, append(hmacValueAppendNode, encryptMarshaledAppendNode...))

	// Case 1: Only node is head node
	pointerLastNode := &headFileStructData
	var lastIsHead bool = true
	var pointerLastNodeUUID uuid.UUID
	//Case 2: There's a bunch of nodes, we need the last one
	if pointerLastNode.LastNode != uuid.Nil {
		lastIsHead = false
		nodeFileKey, err := userlib.HashKDF(intNodeData.FileKey, []byte("append" + strconv.Itoa(headFileStructData.NumAppends - 1)))
		if err != nil {
			return err
		}
		lastNodeData, err := FetchFileStruct(headFileStructData.LastNode, nodeFileKey[:16])
		if err != nil {
			return err
		}
		pointerLastNodeUUID = headFileStructData.LastNode
		pointerLastNode = &lastNodeData
	}
	//The old node's next should be pointing to the newAppend
	pointerLastNode.NextNode = newAppendUUID

	//Change the headNode's lastNode to point to the new append
	headFileStructData.LastNode = newAppendUUID
	
	//marshal and encrypt (NOW second to) lastNode
	marshaledLastNode, err := json.Marshal(pointerLastNode)
	if err != nil {
		return err
	}
	

	if !lastIsHead {
		iv = userlib.RandomBytes(16)
		nodeFileKey, err := userlib.HashKDF(intNodeData.FileKey, []byte("append" + strconv.Itoa(headFileStructData.NumAppends - 1)))
		if err != nil {
			return err
		}
		encMarshaledLastNode := symEnc(nodeFileKey[:16], iv, marshaledLastNode)
		hmacLastNode, err := calculateHMAC(nodeFileKey[:16], encMarshaledLastNode)
		if err != nil {
			return err
		}

		datastoreSet(pointerLastNodeUUID, append(hmacLastNode, encMarshaledLastNode...))
		iv = userlib.RandomBytes(16)

		marshaledLastNode, err = json.Marshal(headFileStructData)
		if err != nil {
			return err
		}
	}

	encMarshaledLastNode := symEnc(intNodeData.FileKey, iv, marshaledLastNode)
	//Marshal and encrypt headFileStruct
	newHmacStruct, err := createHmacStruct(intNodeData.FileKey, encMarshaledLastNode)
	if err != nil {
		return err
	}
	marshaledNewHmacStruct, err := json.Marshal(newHmacStruct)
	if err != nil {
		return err
	}
	datastoreSet(intNodeData.FileStructID, marshaledNewHmacStruct)

	return err

}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	storageKey, err := userdata.GetMailboxUUID(filename)
	if err != nil {
		return nil, err
	}
	_, ok := userlib.DatastoreGet(storageKey)
	if !ok {
		return nil, errors.New(strings.ToTitle("file not found"))
	}
	var mailBoxData Mailbox
	mailBoxData, err = userdata.FetchMailBox(storageKey)
	if err != nil {
		return nil, err
	}

	var intNodeData Mailbox
	intNodeData, err = FetchIntNode(mailBoxData.FileStructID, mailBoxData.FileKey)
	if err != nil {
		return nil, err
	}

	var hmacData HMAC
	hmacData, err = FetchHmacStruct(intNodeData.FileStructID)
	if err != nil {
		return nil, err
	}

	//Check hmacStruct
	if !checkHmacStruct(hmacData, intNodeData.FileKey) {
		return nil, errors.New(strings.ToTitle("Hmac compromised"))
	}

	var headFile File
	headFile, err = FetchHeadFileStruct(hmacData.Encryption, intNodeData.FileKey)
	if err != nil {
		return nil, err
	}

	//Check hmac of fileContents
	err = checkHMACFileContents(headFile, intNodeData.FileKey)
	if err != nil {
		return nil, err
	}
	//Get fileContents
	var fileContents []byte
	fileContents, err = fetchFileContents(headFile.UUIDContents, intNodeData.FileKey)
	if err != nil {
		return nil, err
	}

	content = fileContents
	if headFile.NextNode != uuid.Nil {

		var node uuid.UUID = headFile.NextNode
		// var prevNodeData File = headFile
		for i := 1; i <= headFile.NumAppends; i++ {
			var nodeData File
			nodeStructKey, err := userlib.HashKDF(intNodeData.FileKey, []byte("append" + strconv.Itoa(i)))
			if err != nil {
				return nil, err
			}

			nodeData, err = FetchFileStruct(node, nodeStructKey[:16])
			if err != nil {
				return nil, err
			}
			fileContentKey, err := userlib.HashKDF(intNodeData.FileKey, []byte("content" + strconv.Itoa(i)))
			if err != nil {
				return nil, err
			}

			//Get fileContents of nodeData
			fileContents, err = fetchFileContents(nodeData.UUIDContents, fileContentKey[:16])
			if err != nil {
				return nil, err
			}

			content = append(content, fileContents...)
			// prevNodeData = nodeData
			node = nodeData.NextNode

		}
	}

	return content, err
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	// Create a new mailbox encrypted with sharee's key stored at
	// UUID(Hash(ShareeName)||Filename)
	// Put the key to the intNode inside of said mailbox
	//Create a new intermediary node?

	// Store the intNode UUID at that mailbox

	//Find deterministic mailbox UUID using hash(filename) || hash(user)
	mailboxID, check1 := userdata.GetMailboxUUID(filename)
	if check1 != nil {
		return uuid.Nil, errors.New("CreateInvitation: did not get MailboxUUID successfully")
		// return uuid.Nil, check1

	}

	//Check that the recipientUsername exists by looking if Keystore has
	recipientPublicKey, ok := userlib.KeystoreGet(recipientUsername + "a")
	if !ok {
		return uuid.Nil, errors.New("recipient does not exist")
	}

	//Get intNode
	mailbox, err := userdata.FetchMailBox(mailboxID)
	if err != nil {
		return uuid.Nil, err
	}

	intNode, err := FetchIntNode(mailbox.FileStructID, mailbox.FileKey)
	if err != nil {
		return uuid.Nil, err
	}
	hmacData, err := FetchHmacStruct(intNode.FileStructID)
	if err != nil {
		return uuid.Nil, err
	}

	headFile, err := FetchHeadFileStruct(hmacData.Encryption, intNode.FileKey)
	if err != nil {
		return uuid.Nil, err
	}

	//Check if user is owner of the file. If yes, create a new intNode. Otherwise use it's own intNode
	userHash := userlib.Hash([]byte(userdata.Username))
	var yesOwner bool = bytes.Equal(userHash, headFile.Owner)

	//Add recipient to invited user list
	_, ok = userlib.DatastoreGet(headFile.InvitedUsersUUID)
	if !ok {
		err = encMarshalList(headFile.InvitedUsersUUID, nil, intNode.FileKey, "invitedUsers")
		if err != nil {
			return uuid.Nil, err
		}
	}
	
	invitedUserList, _, err := fetchUserList(headFile.InvitedUsersUUID, intNode.FileKey, "invitedUsers")
	if err != nil {
		return uuid.Nil, err
	}
	var recipientHashStruct Hash
	recipientHashStruct.Hash = userlib.Hash([]byte(recipientUsername))

	invitedUserList = append(invitedUserList, recipientHashStruct)
	err = encMarshalList(headFile.InvitedUsersUUID, invitedUserList, intNode.FileKey, "invitedUsers")
	if err != nil {
		return uuid.Nil, err
	}
	//Marshal & Encrypt senderIntNode
	var invite Invitation
	if yesOwner {
		recipientUsernameHash := userlib.Hash([]byte(recipientUsername))
		marshaledSenderIntNode, err := json.Marshal(intNode)
		if err != nil {
			return uuid.Nil, err
		}

		recipientFileKey, err := userlib.HashKDF(mailbox.FileKey, recipientUsernameHash)
		if err != nil {
			return uuid.Nil, err
		}
		iv := userlib.RandomBytes(16)

		cipherMarshaledIntNode := symEnc(recipientFileKey[:16], iv, marshaledSenderIntNode)

		//Append hmac to intNode
		hmacIntNode, err := calculateHMAC(recipientFileKey[:16], cipherMarshaledIntNode)
		if err != nil {
			return uuid.Nil, err
		}

		//Create intNode for recipient
		intNodePtr, err := userdata.fetchIntNodePtr(filename, recipientUsernameHash)
		if err != nil {
			return uuid.Nil, err
		}

		//Marshal newIntNode
		datastoreSet(intNodePtr, append(hmacIntNode, cipherMarshaledIntNode...))

		invite = createInvitationStruct(intNodePtr, recipientFileKey[:16])
	} else {
		invite = createInvitationStruct(mailbox.FileStructID, mailbox.FileKey)
	}
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	// //Marshal invitationStruct
	marshaledInvitationStruct, err := json.Marshal(invite)
	if err != nil {
		return uuid.Nil, err
	}

	//Encrypt invitationStruct
	encryptedMarshaledInvitationStruct, err := userlib.PKEEnc(recipientPublicKey, marshaledInvitationStruct)
	if err != nil {
		return uuid.Nil, errors.New("problem with encrypting invitation struct")
	}

	// Get sender's private key

	signedBytes, err := userlib.DSSign(userdata.PrivateDigitalSig, encryptedMarshaledInvitationStruct)
	if err != nil {
		return uuid.Nil, err
	}

	//Hash the outcomes of get uuID mailbox, as well as every where we created a string of this kind
	hashedMailboxString := userlib.Hash([]byte(filename + "+" + recipientUsername + "+" + userdata.Username))
	invitationPtr, err = uuid.FromBytes(hashedMailboxString[:16])
	if err != nil {
		return uuid.Nil, err
	}
	datastoreSet(invitationPtr, append(signedBytes, encryptedMarshaledInvitationStruct...))

	return invitationPtr, nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	//TODO: Support the feature where, if a user revokes before the recipient accepted. We need to delete that invitation
	// the invitation pointer or something
	//Check if filename is in user's existing filespace
	uuidMailBox, err := userdata.GetMailboxUUID(filename)
	if err != nil {
		return err
	}
	// Attempt to fetch it, if we find that particular mailbox (ok = true)
	// then the user must already have this file in its file space
	_, ok := userlib.DatastoreGet(uuidMailBox)
	if ok {
		mailboxData, err := userdata.FetchMailBox(uuidMailBox)
		if err != nil {
			return nil
		}

	
		_, ok = userlib.DatastoreGet(mailboxData.FileStructID)
		if ok {
			return errors.New("file already exists in recipient's filespace")
		}
		
	}
	invitationContents, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New("cannot get datastore")
	}

	senderVerifyKey, ok := userlib.KeystoreGet(senderUsername + "b")
	if !ok {
		return errors.New("something wrong with verification key")
	}

	signedBytes, encMarshaledInvStruct := invitationContents[:256], invitationContents[256:]

	//Verify signedBytes
	err = userlib.DSVerify(senderVerifyKey, encMarshaledInvStruct, signedBytes)
	if err != nil {
		return err
	}

	marshaledInvStruct, err := userlib.PKEDec(userdata.PrivateDecKey, encMarshaledInvStruct)
	if err != nil {
		return err
	}

	var invitationData Invitation
	err = json.Unmarshal(marshaledInvStruct, &invitationData)
	if err != nil {
		return err
	}

	keyToDecryptIntNode, uuidIntNode := invitationData.Key, invitationData.IntNodeUUID

	var intNodeData Mailbox
	intNodeData, err = FetchIntNode(uuidIntNode, keyToDecryptIntNode)
	if err != nil {
		return err
	}

	//Create mailbox

	mailBoxStruct := createMailBox(uuidIntNode, keyToDecryptIntNode)

	//Marshal mailbox
	marshaledMailBox, err := json.Marshal(mailBoxStruct)
	if err != nil {
		return err
	}

	publicEncKey, ok := userdata.getPublicEncKey()
	if !ok {
		return errors.New("publicEncKey has problems")
	}
	encryptMarshaledMailBox, err := userlib.PKEEnc(publicEncKey, marshaledMailBox)
	if err != nil {
		return err
	}

	//Sign encryptedMarshalMailBox
	signedMailBox, err := userlib.DSSign(userdata.PrivateDigitalSig, encryptMarshaledMailBox)
	if err != nil {
		return err
	}
	datastoreSet(uuidMailBox, append(signedMailBox, encryptMarshaledMailBox...))

	//Delete invitation
	userlib.DatastoreDelete(invitationPtr)

	//Check if sender is owner of file
	//1. Decrypt encryptedMarshaledIntNode
	//2. Unmarshal encryptedMarshaledIntNode
	//3. Get the HMAC


	var hmacData HMAC
	hmacData, err = FetchHmacStruct(intNodeData.FileStructID)
	if err != nil {
		return err
	}

	if !checkHmacStruct(hmacData, intNodeData.FileKey) {
		return errors.New("hmac has been compromised")
	}

	var fileData File
	fileData, err = FetchHeadFileStruct(hmacData.Encryption, intNodeData.FileKey)
	if err != nil {
		return err
	}

	//If sender is owner, then we need to place user in the sharee's list of direct descendents
	//Regardless of who is owner, we need to remove recipient from the list of invited but not accepted users
	senderHash := userlib.Hash([]byte(senderUsername))
	var directDescendent bool = bytes.Equal(fileData.Owner, senderHash)
	sharedUserHash := userlib.Hash([]byte(userdata.Username))
	var sharedUserHashStruct Hash
	sharedUserHashStruct.Hash = sharedUserHash

	//Fetch sharedUserList
	sharedUsersList, _, err := fetchUserList(fileData.SharedUsersUUID, intNodeData.FileKey, "sharedUsers")
	if err != nil {
		return err
	}

	if directDescendent {
		//Add recipient to list of shared users
		if sharedUsersList == nil {
			newSharedUsers := []Hash{sharedUserHashStruct}
			sharedUsersList = newSharedUsers
		} else {
			sharedUsersList = append(sharedUsersList, sharedUserHashStruct)
		}

		//Store changed list into sharedUser's uuid
		encMarshalList(fileData.SharedUsersUUID, sharedUsersList, intNodeData.FileKey, "sharedUsers")
	}

	//Fetch invitedUserList

	invitedUsersList, ok, err := fetchUserList(fileData.InvitedUsersUUID, intNodeData.FileKey, "invitedUsers")
	if ok {
		if err != nil {
			return err
		}
	
		newInvitedUsers := []Hash{}
		for _, v := range invitedUsersList {
			if !bytes.Equal(v.Hash, sharedUserHash) {
				newInvitedUsers = append(newInvitedUsers, sharedUserHashStruct)
			}
		}
		encMarshalList(fileData.InvitedUsersUUID, newInvitedUsers, intNodeData.FileKey, "invitedUsers")
	}
	

	//Store changed list into invitedUser's uuid

	// marshaledFileData, err := json.Marshal(fileData)
	// if err != nil {
	// 	return err
	// }
	// iv := userlib.RandomBytes(16)
	// encryptMarshaledFileData := symEnc(intNodeData.FileKey, iv, marshaledFileData)

	// hmacData, err = createHmacStruct(intNodeData.FileKey, encryptMarshaledFileData)
	// if err != nil {
	// 	return err
	// }
	// marshalHmacData, err := json.Marshal(hmacData)
	// if err != nil {
	// 	return err
	// }
	// userlib.DatastoreSet(intNodeData.FileStructID, marshalHmacData)

	return nil

}

//Encrypt, marshal and store list
func encMarshalList(uuidList uuid.UUID, data []Hash, notDerivedKey []byte, typeOfList string) error {
	//Get derived key
	derivedKey, err := userlib.HashKDF(notDerivedKey, []byte(typeOfList))
	if err != nil {
		return err
	}

	//Marshal then encrypt
	marshaledList, err := json.Marshal(data)
	if err != nil {
		return err
	}

	iv := userlib.RandomBytes(16)
	encMarshaledList := userlib.SymEnc(derivedKey[:16], iv, marshaledList)

	//Add hmac
	hmacOfList, err := calculateHMAC(derivedKey[:16], encMarshaledList)
	if err != nil {
		return err
	}

	//Append to encryption
	userlib.DatastoreSet(uuidList, append(hmacOfList, encMarshaledList...))
	return nil

}
func fetchUserList(uuidList uuid.UUID, notDerivedKey []byte, typeOfList string) ([]Hash, bool,  error) {
	dataOfList, ok := userlib.DatastoreGet(uuidList)
	if !ok {
		return nil, false, nil
	}

	derivedKey, err := userlib.HashKDF(notDerivedKey, []byte(typeOfList))
	if err != nil {
		return nil, true, err
	}

	//Get Hmac and data
	hmacOfList, encMarshaledOfList := dataOfList[:64], dataOfList[64:]

	//Check Hmac
	calculatedHmac, err := calculateHMAC(derivedKey[:16], encMarshaledOfList)
	if err != nil {
		return nil, true, err
	}
	ok = userlib.HMACEqual(calculatedHmac, hmacOfList)
	if !ok {
		return nil, true, errors.New("list modified")
	}

	//Decrypt and unmarshal
	marshaledOfList := userlib.SymDec(derivedKey[:16], encMarshaledOfList)
	var hashData []Hash
	err = json.Unmarshal(marshaledOfList, &hashData)
	if err != nil {
		return nil, true, err
	}

	return hashData, true, nil
}
func (userdata *User) fetchIntNodePtr(filename string, recipientUsernameHash []byte) (uuid.UUID, error) {
	filenameBytes := []byte(filename + "-")
	usernameBytes := []byte("-" + userdata.Username)
	intNodeString := userlib.Hash(append(append(filenameBytes, recipientUsernameHash...), usernameBytes...))
	intNodePtr, err := uuid.FromBytes(intNodeString[:16])
	return intNodePtr, err
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	// First ensure that filename does not exist in caller's personal file namespace
	// Set up
	// Getting the mailbox
	// Go form mailbox -> intnode
	// Check if the user's we're revoking already have permissions.
	// TODO: Overcome the late to the party bug
	//Revocation:
	//

	// Generate the mailBox, then check if it exist's in the caller's personal fileSpace
	mailboxID, check1 := userdata.GetMailboxUUID(filename)
	if check1 != nil {
		return check1

	}
	_, ok := userlib.DatastoreGet(mailboxID)
	if !ok {
		return errors.New("file does not exist in user's filespace")
	}

	//If recipient did not accept invitation (and not shared duh), delete the invitation
	//Find the invitation pointer

	invitationString := userlib.Hash([]byte(filename + "+" + recipientUsername + "+" + userdata.Username))
	invitationPtr, err := uuid.FromBytes(invitationString[:16])
	if err != nil {
		return err
	}
	_, ok = userlib.DatastoreGet(invitationPtr)
	if ok {
		recipientHash := userlib.Hash([]byte(recipientUsername))
		intNodePtr, err := userdata.fetchIntNodePtr(filename, recipientHash)
		if err != nil {
			return err
		}
		userlib.DatastoreDelete(intNodePtr)
		userlib.DatastoreDelete(invitationPtr)
		//return error because the file was not shared with recipient
		return nil
	}

	//Error if given filename is not currently shared with recipientUsername
	recipientUsernameHash := userlib.Hash([]byte(recipientUsername))
	intNodePtr, err := userdata.fetchIntNodePtr(filename, recipientUsernameHash)
	if err != nil {
		return err
	}

	_, ok = userlib.DatastoreGet(intNodePtr)
	if !ok {
		return errors.New("recipient does not have access to file")
	}

	//Get intNodeKey through username's own intNode
	var mailboxData Mailbox
	mailboxData, err = userdata.FetchMailBox(mailboxID)
	if err != nil {
		return err
	}
	var intNodeData Mailbox
	intNodeData, err = FetchIntNode(mailboxData.FileStructID, mailboxData.FileKey)
	if err != nil {
		return err
	}

	var hmacData HMAC
	hmacData, err = FetchHmacStruct(intNodeData.FileStructID)
	if err != nil {
		return err
	}
	var headFile File
	headFile, err = FetchHeadFileStruct(hmacData.Encryption, intNodeData.FileKey)
	if err != nil {
		return err
	}
	if !bytes.Equal(headFile.Owner, userlib.Hash([]byte(userdata.Username))) {
		return errors.New("user does not have the power to revoke the file")
	}

	content, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}

	// contentBytes, err := json.Marshal(content)
	// if err != nil {
	// 	return err
	// }

	//Encrypt contentBytes with new fileKey
	newFileKey := userlib.RandomBytes(16)
	newHeadFileStruct, err := createHeadFileStruct(userdata.Username, content, newFileKey)
	if err != nil {
		return err
	}
	
	//Fetch sharedUsersList
	headSharedUsersList, _, err := fetchUserList(headFile.SharedUsersUUID, intNodeData.FileKey, "sharedUsers")
	if err != nil {
		return err
	}
	//Remove revoked user from direct shared users list (headFileStruct.SharedUsers slice)
	recipientHash := userlib.Hash([]byte(recipientUsername))
	var newSharedUserString []Hash
	for _, s := range headSharedUsersList {
		if !bytes.Equal(s.Hash, recipientHash) {
			newSharedUserString = append(newSharedUserString, s)
		} else {
			intNodePtr, err := userdata.fetchIntNodePtr(filename, recipientHash)
			if err != nil {
				return err
			}
			userlib.DatastoreDelete(intNodePtr)
		}
	}
	//Create new uuid for sharedUsers
	newUUIDSharedUsers := uuid.New()
	encMarshalList(newUUIDSharedUsers, newSharedUserString, newFileKey, "sharedUsers")
	newHeadFileStruct.SharedUsersUUID = newUUIDSharedUsers
	marshalNewHeadFileStruct, err := json.Marshal(newHeadFileStruct)
	if err != nil {
		return err
	}
	iv := userlib.RandomBytes(16)
	encryptMarshalNewHeadFileStruct := symEnc(newFileKey, iv, marshalNewHeadFileStruct)

	var newHmacStruct HMAC
	newHmacStruct, err = createHmacStruct(newFileKey, encryptMarshalNewHeadFileStruct)
	if err != nil {
		return err
	}

	marshaledNewHmacStruct, err := json.Marshal(newHmacStruct)
	if err != nil {
		return err
	}

	newHmacStructUUID := uuid.New()
	datastoreSet(newHmacStructUUID, marshaledNewHmacStruct)
	newIntNode := createMailBox(newHmacStructUUID, newFileKey)

	//Marshal and encrypt newIntNode
	marshaledNewIntNode, err := json.Marshal(newIntNode)
	if err != nil {
		return err
	}
	iv = userlib.RandomBytes(16)
	encMarshaledNewIntNode := symEnc(mailboxData.FileKey, iv, marshaledNewIntNode)
	//Get hmac of newIntNode
	hmacNewIntNode, err := calculateHMAC(mailboxData.FileKey, encMarshaledNewIntNode)
	if err != nil {
		return err
	}

	x := append(hmacNewIntNode, encMarshaledNewIntNode...)
	datastoreSet(mailboxData.FileStructID, x)

	//Loop through all the un-revoked users to update their mailboxes
	newSharedUsersList, _, err := fetchUserList(newHeadFileStruct.SharedUsersUUID, newFileKey, "sharedUsers")
	if err != nil {
		return err
	}
	for _, s := range newSharedUsersList {
		uuidIntNode, err := userdata.fetchIntNodePtr(filename, s.Hash)
		if err != nil {
			return err
		}
		recipientFileKey, err := userlib.HashKDF(mailboxData.FileKey, s.Hash)
		if err != nil {
			return err
		}
		encMarshaledNewIntNodeShared := symEnc(recipientFileKey[:16], iv, marshaledNewIntNode)
		hmacNewIntNodeShared, err := calculateHMAC(recipientFileKey[:16], encMarshaledNewIntNodeShared)
		if err != nil {
			return err
		}
		datastoreSet(uuidIntNode, append(hmacNewIntNodeShared, encMarshaledNewIntNodeShared...))
	}

	

	invitedUsersList, ok, err := fetchUserList(newHeadFileStruct.InvitedUsersUUID, newFileKey, "invitedUsers")
	if ok {
		if err != nil {
			return err
		}
		
		for _, s := range invitedUsersList {
			uuidIntNode, err := userdata.fetchIntNodePtr(filename, s.Hash)
			if err != nil {
				return err
			}
			recipientFileKey, err := userlib.HashKDF(mailboxData.FileKey, s.Hash)
			if err != nil {
				return err
			}
			encMarshaledNewIntNodeShared := symEnc(recipientFileKey[:16], iv, marshaledNewIntNode)
			hmacNewIntNodeShared, err := calculateHMAC(recipientFileKey[:16], encMarshaledNewIntNodeShared)
			if err != nil {
				return err
			}
			datastoreSet(uuidIntNode, append(hmacNewIntNodeShared, encMarshaledNewIntNodeShared...))
		}
		
	}
	

	return nil
}
