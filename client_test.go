package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "encoding/json" //TODO: comment out when sending to autograder
	_ "errors"

	// "math/rand"
	// "strconv"
	_ "strconv"
	_ "strings"
	"testing"

	// "time"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	"github.com/google/uuid" //TODO: comment out when sending to autograder
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"
const contentDorisChange = "123"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {
	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var doris *client.User
	var eve *client.User
	var frank *client.User
	var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	fooFile := "foo.txt"
	dorisFile := "dorisFile.txt"
	eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	measureBandwidth := func(probe func()) (bandwidth int) {
		before := userlib.DatastoreGetBandwidth()
		probe()
		after := userlib.DatastoreGetBandwidth()
		return after - before
	}

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			userlib.DebugMsg(alice.Username)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("InitUser: The same username shouldn't be allowed", func() {
			userlib.DebugMsg("Creating bernie")
			_, err := client.InitUser("bernie", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Trying to create bernie --should error")
			_, err = client.InitUser("bernie", defaultPassword)
			Expect(err).ToNot(BeNil())

		})

		Specify("InitUser: Username with different cases are unique", func() {
			userlib.DebugMsg("InitUser: lowercase bernie")
			_, err := client.InitUser("bernie", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("InitUser: Uppercase Bernie")
			_, err = client.InitUser("Bernie", defaultPassword)
			Expect(err).To(BeNil())

		})

		Specify("InitUser: Empty String as username not allowed", func() {
			userlib.DebugMsg("InitUser: No username -- Should error")
			_, err := client.InitUser(emptyString, defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("AcceptInvitation: Recepient should be able to accept and load a file", func() {
			userlib.DebugMsg("To test if system can have same filename in the global scheme")
			userlib.DebugMsg("Initializing Alice and Bob")
			alice, err := client.InitUser("alice", "hello") //alicePassword = hello
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", "bye") //bobPassword = bye
			Expect(err).To(BeNil())

			userlib.DebugMsg("Store fooFile under Alice's space")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Store fooFile under Bob's space")
			err = bob.StoreFile(fooFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Check that Alice can access the contents of file")
			data, err := alice.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Check that Bob can access the contents of bob.foofile")
			data, err = bob.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Create invitation for Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to accept but errors")
			err = bob.AcceptInvitation("alice", invite, fooFile)
			Expect(err).ToNot(BeNil())

		})

		Specify("InitUser: Password length is zero", func() {
			userlib.DebugMsg("Password can be zero")
			_, err = client.InitUser("alice", emptyString)
			Expect(err).To(BeNil())
		})

		Specify("Our own test: No user initialized", func() {
			userlib.DebugMsg("No user initialized")
			_, err = client.GetUser("alice", emptyString)
			Expect(err).ToNot(BeNil())

		})

		Specify("GetUser: User credientials invalid", func() {
			userlib.DebugMsg("Initialize alice with default password")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Get alice with wrong password")
			_, err = client.GetUser("alice", emptyString)
			Expect(err).ToNot(BeNil())

		})

		//specify a test with invalid login
		Specify("GetUser: Invalid Username", func() {
			userlib.DebugMsg("Initialize alice with default password")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//get alice with the wrong username
			userlib.DebugMsg("Get alice with wrong username")
			_, err = client.GetUser("bob", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("LoadFile: loading file that exists but user does not have access to", func() {
			userlib.DebugMsg("Initialize alice and frank with default password")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			frank, err = client.InitUser("frank", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Store fooFile under Alice's space")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())
			userlib.DebugMsg("Check that Frank cannot see the file")
			_, err = frank.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())

		})

		Specify("LoadFile: File does not exist", func() {
			userlib.DebugMsg("Initializefrank with default password")
			frank, err = client.InitUser("frank", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Check that Frank cannot see the file")
			_, err = frank.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("Random Test: Large Sharing Tree ", func() {
			userlib.DebugMsg("Create eight users: alice, bob, charles, doris, eve, frank, grace, horace")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			eve, err = client.InitUser("eve", defaultPassword)
			Expect(err).To(BeNil())

			frank, err = client.InitUser("frank", defaultPassword)
			Expect(err).To(BeNil())

			grace, err = client.InitUser("grace", defaultPassword)
			Expect(err).To(BeNil())

			_, err = client.InitUser("horace", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Store fooFile under Alice's space")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Create invitation for Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice and calls it foo.txt")
			err = bob.AcceptInvitation("alice", invite, fooFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Check that Bob can access fooFile")
			content, err := bob.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Create invitation for Charles")
			invite, err = alice.CreateInvitation(fooFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepts invitation from Alice")
			err = charles.AcceptInvitation("alice", invite, fooFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob shares file to Doris and Eve")
			invite, err = bob.CreateInvitation(fooFile, "doris")
			Expect(err).To(BeNil())
			err = doris.AcceptInvitation("bob", invite, fooFile)
			Expect(err).To(BeNil())

			invite, err = bob.CreateInvitation(fooFile, "eve")
			Expect(err).To(BeNil())
			err = eve.AcceptInvitation("bob", invite, fooFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Check that Doris and Eve can see the file")
			content, err = doris.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			content, err = eve.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Charles shares file with frank")
			invite, err = charles.CreateInvitation(fooFile, "frank")
			Expect(err).To(BeNil())
			err = frank.AcceptInvitation("charles", invite, fooFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Check that Frank can see the file")
			content, err = frank.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Ensure that Grace has no access to the file")

			userlib.DebugMsg("2. LoadFile: Grace should not be able to load fooFile")
			_, err = grace.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("3. AppendFile: Grace should not be able to append fooFile")
			err = grace.AppendToFile(fooFile, []byte(contentThree))
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("4. CreateInvitation: Grace should not be able to create invitation for fooFile")
			invite, err = grace.CreateInvitation(fooFile, "horace")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("5. AcceptInvitation: Grace should not be able to accept invitation for fooFile")
			err = grace.AcceptInvitation("grace", invite, "horace")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Bob modifies fooFile")
			err = bob.StoreFile(fooFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Check that all can access the right contents:")
			userlib.DebugMsg("Alice can access")
			content, err = alice.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Bob can access")
			content, err = bob.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Charles can access")
			content, err = charles.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Doris can access")
			content, err = doris.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Eve can access")
			content, err = eve.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Frank can access")
			content, err = frank.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Doris Modifies fooFile")

			err = doris.StoreFile(fooFile, []byte(contentDorisChange))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Doris's changes to the file should be reflected in everyone's file as well")
			content, err = bob.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			content, err = charles.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			content, err = alice.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			content, err = doris.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			content, err = eve.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			content, err = frank.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			//Starting to test revoke
			userlib.DebugMsg("Alice revokes bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob now loads a file he doesn't have access to (fooFile)")
			_, err = bob.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking access of D and E")

			// B, D E shouldnt have access to (foo file)
			userlib.DebugMsg("Doris should not be able to load file")
			_, err = doris.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Eve should not be able to load file")
			_, err = eve.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())

			// A, C, F should have access (foofile)
			userlib.DebugMsg("Alice should have access to fooFile")
			content, err = alice.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			userlib.DebugMsg("Charles should have access to fooFile")
			content, err = charles.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

			userlib.DebugMsg("Frank should have access to fooFile")
			content, err = frank.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentDorisChange)))

		})

		Specify("Store File: Able to store to a file they were invited to", func() {
			userlib.DebugMsg("1. Filename does not exist in the personal file namespace of the caller")
			userlib.DebugMsg("Create Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation from Alice and calls it foo.txt")
			err = bob.AcceptInvitation("alice", invite, fooFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob modifies fooFile")
			err = bob.StoreFile(fooFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Check that all can access the right contents:")
			userlib.DebugMsg("Alice can access")
			content, err := alice.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

			userlib.DebugMsg("Bob can access")
			content, err = bob.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentThree)))

		})
		Specify("AppendToFile: filename does not exist to be appended", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Alice appends to an non existent file")
			err = alice.AppendToFile(fooFile, []byte("hi"))
			Expect(err).ToNot(BeNil())
		})

		Specify("AppendToFile: Append empty string to file", func() {
			userlib.DebugMsg("Initialize Alice")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Creating fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())
			userlib.DebugMsg("Alice appends an empty string")
			err = alice.AppendToFile(fooFile, []byte(""))
			Expect(err).To(BeNil())
		})

		Specify("Create Invitation: Filename doesn't exist in namespace of caller", func() {
			userlib.DebugMsg("1. Filename does not exist in the personal file namespace of the caller")
			userlib.DebugMsg("Create Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice invites Bob")
			_, err = alice.CreateInvitation(fooFile, "bob")
			Expect(err).ToNot(BeNil())
		})

		Specify("Create Invitation: RecipientUsername does not exist", func() {
			//create alice and bob
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			//alice creates file foo.txt
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())
			//alice invites charles to foo.txt
			_, err := alice.CreateInvitation(fooFile, "charles")
			//expect an error
			Expect(err).ToNot(BeNil())
			//done
		})

		Specify("Accept Invitation: Unable to verify invite", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = bob.StoreFile("bar.txt", []byte(contentOne))
			Expect(err).To(BeNil())
			invite, err := alice.CreateInvitation(fooFile, "charles")
			Expect(err).To(BeNil())
			// invite2, err := bob.CreateInvitation("boo.txt", "bob")
			// Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, "trying.txt")
			Expect(err).ToNot(BeNil())

		})

		Specify("Accept Invitation: Swapping invites", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = bob.StoreFile("bar.txt", []byte(contentOne))
			Expect(err).To(BeNil())
			invite, err := alice.CreateInvitation(fooFile, "charles")
			Expect(err).To(BeNil())
			invite2, err := bob.CreateInvitation("bar.txt", "charles")
			Expect(err).To(BeNil())
			err = charles.AcceptInvitation("bob", invite, "trying.txt")
			Expect(err).ToNot(BeNil())
			err = charles.AcceptInvitation("alice", invite2, "trying.txt")
			Expect(err).ToNot(BeNil())

		})

		//specify a test for double active invitations
		Specify("AcceptInvitation: Double active invitations of the same fileName in owner's space", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			//init charles
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			//create file foo.txt
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			//charles create foo.txt
			err = charles.StoreFile(fooFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			//alice invites bob to foo.txt
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			//charles invites bob to foo.txt
			invite2, err := charles.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			//bob accepts charles's invite
			err = bob.AcceptInvitation("charles", invite2, "charlesfoo.txt")
			Expect(err).To(BeNil())

			//bob accepts alice's invite
			err = bob.AcceptInvitation("alice", invite, "alicefoo.txt")
			Expect(err).To(BeNil())

			//bob loads charlesfoo.txt
			data, err := bob.LoadFile("charlesfoo.txt")
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(contentTwo))

			//bob loads alicefoo.txt
			data, err = bob.LoadFile("alicefoo.txt")
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(contentOne))

		})

		Specify("AcceptInvitation: Invitee already has a file with the given filename in their personal file namespace", func() {
			userlib.DebugMsg("Initalize user: Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice stores a fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob stores a fooFile")
			err = bob.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("1. Caller already has a file with the given filename in their personal file namespace")
			err = bob.AcceptInvitation("alice", invite, fooFile)
			Expect(err).ToNot(BeNil())

		})

		Specify("AcceptInvitation: User revokes before invitation Pointer not intended for you", func() {
			userlib.DebugMsg("Initalize user Alice, Bob and Charlie")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice stores a fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).ToNot(BeNil())

		})

		Specify("AcceptInvitation: Accept a a UUID.nil invite", func() {
			userlib.DebugMsg("Initalize user Alice, Bob and Charlie")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice stores a fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite := uuid.New()

			err = charles.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).ToNot(BeNil())

		})

		Specify("AcceptInvitation: Invitation is no longer valid due to Revocation", func() {
			userlib.DebugMsg("Initalize user Alice, Bob and Charlie")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice has fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, "boo.txt")
			Expect(err).ToNot(BeNil())

		})

		Specify("RevokeAccess: Filename does not exist in owner's space", func() {
			userlib.DebugMsg("1. Filename does not exist in owner's space")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob has fooFile but Alice does not")
			err = bob.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes fooFile (does not have it)")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).ToNot(BeNil())

		})

		Specify("RevokeAccess: Filename does not exist in recipient's space", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Filename does not exist in recipient's space")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).ToNot(BeNil())
		})

		Specify("RevokeAccess: Recipient tries to revoke access of a file that they own with a valid user that does not have access", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			// alice creates a file
			userlib.DebugMsg("Alice creates a file")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			//alice revokes from bob
			userlib.DebugMsg("Alice revokes from bob")
			err = alice.RevokeAccess(fooFile, "bob")
			//expect an error
			Expect(err).ToNot(BeNil())
		})

		Specify("RevokeAcess: Recipient tries to revoke access of a file they do not own", func() {
			userlib.DebugMsg("Initiate Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to revoke access on Alice")
			err = bob.RevokeAccess("bar.txt", "alice")
			Expect(err).ToNot(BeNil())

		})

		Specify("RevokeAccess: Recipient tries to load after revocation", func() {
			userlib.DebugMsg("Initiate Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to load file")
			_, err = bob.LoadFile(fooFile)
			Expect(err).ToNot(BeNil())

		})

		Specify("RevokeAccess: Recipient tries to create Invitation after revocation", func() {
			userlib.DebugMsg("Initiate Alice, Bob, Charles")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to create invitation of fooFile")
			_, err = bob.CreateInvitation(fooFile, "charles")
			Expect(err).ToNot(BeNil())

		})
		
		Specify("Revoke Access: Recipient makes same file after being revoked", func() {
			userlib.DebugMsg("Initiate Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to store file")
			err = bob.StoreFile(fooFile, []byte(contentTwo))
			Expect(err).To(BeNil())
		}) 

		Specify("RevokeAccess: Recipient tries to append after revocation", func() {
			userlib.DebugMsg("Initiate Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to append file")
			err = bob.AppendToFile(fooFile, []byte("hi"))
			Expect(err).ToNot(BeNil())

		})

		Specify("RevokeAccess: Recipient cannot revoke after being revoked", func() {
			userlib.DebugMsg("Initiate Alice, Bob & Charles")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to charles")
			invite, err = alice.CreateInvitation(fooFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepts invitation/ Now has access")
			err = charles.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries to revoke file")
			err = bob.RevokeAccess(fooFile, "charles")
			Expect(err).ToNot(BeNil())
		})

		Specify("RevokeAccess: Check revoked users can still access other files", func() {
			userlib.DebugMsg("Initiate Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice makes fooFile")
			err = alice.StoreFile(fooFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creates invitation and shares to Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepts invitation/ Now has access")
			err = bob.AcceptInvitation("alice", invite, "bar.txt")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob can get into fooFile")
			content, err := bob.LoadFile("bar.txt")
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob has more files that he owns: doris.txt, eve.txt")
			err = bob.StoreFile(dorisFile, []byte(contentOne))
			Expect(err).To(BeNil())

			content, err = bob.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			err = bob.StoreFile(eveFile, []byte(contentOne))
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			content, err = bob.LoadFile(eveFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revokes on Bob")
			err = alice.RevokeAccess(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob tries his other files")
			content, err = bob.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

			content, err = bob.LoadFile(eveFile)
			Expect(err).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))

		})

		Specify("Late to the party bug", func() {
			userlib.DebugMsg("Initalize user: Alice and Bob and Charles")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			aliceContent := []byte("chips")

			userlib.DebugMsg("Alice stores fooFile")
			err := alice.StoreFile(fooFile, aliceContent)
			Expect(err).To(BeNil())

			// Alice share to Charles
			userlib.DebugMsg("Alice share to Charles")
			fooInvPtr, err := alice.CreateInvitation(fooFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepts invitation")
			err = charles.AcceptInvitation("alice", fooInvPtr, fooFile)
			Expect(err).To(BeNil())

			//Alice invites Bob
			userlib.DebugMsg("Alice creates invitation for Bob")
			invite, err := alice.CreateInvitation(fooFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revokes Charlie's access to foo.txt")
			userlib.DebugMsg("Charlie's access to foo.txt")
			err = alice.RevokeAccess(fooFile, "charles")
			Expect(err).To(BeNil())

			//Bob accepts Alice's invitation to foo.txt
			userlib.DebugMsg("Bob accepts Alice's invitation to foo.txt")
			err = bob.AcceptInvitation("alice", invite, fooFile)
			Expect(err).To(BeNil())

			//Bob loads foo.txt
			userlib.DebugMsg("Bob loads foo.txt")
			bobContent, err := bob.LoadFile(fooFile)
			Expect(err).To(BeNil())
			Expect(bobContent).To(Equal(aliceContent))

		})

		Specify("Invalid username when creating invitation", func() {
			userlib.DebugMsg("Initalize user: Alice and Bob")
			alice, err = client.InitUser("alice", defaultPassword)
			//alice creates a file
			aliceContent := []byte("chips")
			err = alice.StoreFile("foo.txt", aliceContent)

			//Alice invites Bob
			_, err = alice.CreateInvitation("foo.txt", "notbob")
			Expect(err).ToNot(BeNil())
		})

		//specify a Multiple Session test
		Specify("Multiple Session test: Append to a file, the other should see appended content upon loading", func() {
			userlib.DebugMsg("Initalize user: Alice and Bob")
			aliceLaptop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//alice creates a file
			aliceContent := []byte("chips")
			err = aliceLaptop.StoreFile("foo.txt", aliceContent)
			Expect(err).To(BeNil())

			//alice appends to the file
			aliceContent = []byte("chipschips")
			err = aliceLaptop.AppendToFile("foo.txt", aliceContent)
			Expect(err).To(BeNil())

		})

		Specify("Delete User Struct", func() {
			userlib.DebugMsg("Initalize user: Alice and Bob")
			aliceLaptop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			username := "alice"
			usernameBytes := []byte(username)
			hashed := userlib.Hash(usernameBytes)[:16]
			userUUID, err := uuid.FromBytes(hashed)
			Expect(err).To(BeNil())

			userlib.DatastoreDelete(userUUID)

			_, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())

		})

		// // //TODO Swap test, how do we get the UUID of the blobs?
		// Specify("Swapping two File Structs", func() {
		// 	//create alice and bob
		// 	alice, err = client.InitUser("alice", defaultPassword)
		// 	Expect(err).To(BeNil())
		// 	bob, err = client.InitUser("bob", defaultPassword)
		// 	Expect(err).To(BeNil())
		// 	//create a file
		// 	aliceContent := []byte("chips")
		// 	//alice stores the file
		// 	err = alice.StoreFile("foo.txt", aliceContent)
		// 	Expect(err).To(BeNil())

		// 	//alice appends
		// 	aliceAppendContent := []byte("are")
		// 	aliceAppendContent2 := []byte("awesome")

		// 	err = alice.AppendToFile("foo.txt", aliceAppendContent)
		// 	Expect(err).To(BeNil())
		// 	// firstAppendFileStructUUID =

		// 	err = alice.AppendToFile("foo.txt", aliceAppendContent2)
		// 	Expect(err).To(BeNil())

		// 	// perform some malcious action by swapping
		// 	// TODO Complete the function and find the UUID of the appended content

		// 	thingToHash := "foo.txt" + "-" + "alice"
		// 	aliceMailboxID, err := uuid.FromBytes(userlib.Hash([]byte(thingToHash))[:16])
		// 	Expect(err).To(BeNil())

		// 	//get alice's mailboxUUID for foo.txt

		// 	aliceMailbox, err := alice.FetchMailBox(aliceMailboxID)
		// 	Expect(err).To(BeNil())

		// 	//Now that we have the mailbox, we can find the IntNode
		// 	aliceIntNode, err := client.FetchIntNode(aliceMailbox.FileStructID, aliceMailbox.FileKey)
		// 	Expect(err).To(BeNil())

		// 	//get the UUID of the first FileStruct
		// 	aliceHmacStructUUID := aliceIntNode.FileStructID

		// 	var hmacData client.HMAC
		// 	hmacData, err = client.FetchHmacStruct(aliceHmacStructUUID)
		// 	Expect(err).To(BeNil())

		// 	var headFileData client.File
		// 	headFileData, err = client.FetchHeadFileStruct(hmacData.Encryption, aliceIntNode.FileKey)
		// 	Expect(err).To(BeNil())

		// 	//get the UUID of the second file struct from the first file struct

		// 	// unmarshal the secondFileStruct
		// 	var secondFileStruct client.File
		// 	secondFileStructKey, err := userlib.HashKDF(aliceIntNode.FileKey, []byte("append"+strconv.Itoa(1)))
		// 	Expect(err).To(BeNil())
		// 	secondFileStruct, err = client.FetchFileStruct(headFileData.NextNode, secondFileStructKey[:16])
		// 	Expect(err).To(BeNil())

		// 	// var thirdFileStruct client.File
		// 	// thirdFileStruct, err = client.FetchFileStruct(secondFileStruct.NextNode, aliceIntNode.FileKey)
		// 	// Expect(err).To(BeNil())

		// 	//get the uuid of the third file struct's blob

		// 	//get the dataStore map
		// 	dataStoreMap := userlib.DatastoreGetMap()

		// 	//swap the second and third file structs
		// 	temp := dataStoreMap[headFileData.NextNode]
		// 	dataStoreMap[headFileData.NextNode] = dataStoreMap[secondFileStruct.NextNode]
		// 	dataStoreMap[secondFileStruct.NextNode] = temp

		// 	_, err = alice.LoadFile("foo.txt")
		// 	Expect(err).ToNot(BeNil())

		// 	// TODO This is relevant for the content swap version of the test
		// 	// get the uuid of the the first file's blob
		// 	// firstAppendFileStructUUID := aliceFileStruct.UUIDContents
		// 	// get the uuid of the second file struct's blob
		// 	// secondAppendFileStructUUID := secondFileStruct.UUIDContents
		// 	// get the uuid of the third file's blob
		// 	// thirdAppendFileStructUUID := thirdFileStruct.UUIDContents

		// })

		// //Specify a test for swapping two appended contents
		// Specify("Swapping two appended contents", func() {
		// 	//create alice and bob
		// 	alice, err = client.InitUser("alice", defaultPassword)
		// 	Expect(err).To(BeNil())
		// 	bob, err = client.InitUser("bob", defaultPassword)
		// 	Expect(err).To(BeNil())
		// 	//create a file
		// 	aliceContent := []byte("chips")
		// 	//alice stores the file
		// 	err = alice.StoreFile("foo.txt", aliceContent)
		// 	Expect(err).To(BeNil())

		// 	//alice appends
		// 	aliceAppendContent := []byte("are")
		// 	aliceAppendContent2 := []byte("awesome")

		// 	err = alice.AppendToFile("foo.txt", aliceAppendContent)
		// 	Expect(err).To(BeNil())
		// 	// firstAppendFileStructUUID =

		// 	err = alice.AppendToFile("foo.txt", aliceAppendContent2)
		// 	Expect(err).To(BeNil())

		// 	// perform some malcious action by swapping
		// 	// TODO Complete the function and find the UUID of the appended content

		// 	thingToHash := "foo.txt" + "-" + "alice"
		// 	aliceMailboxID, err := uuid.FromBytes(userlib.Hash([]byte(thingToHash))[:16])
		// 	Expect(err).To(BeNil())

		// 	//get alice's mailboxUUID for foo.txt
		// 	aliceMailbox, err := alice.FetchMailBox(aliceMailboxID)
		// 	Expect(err).To(BeNil())

		// 	//Now that we have the mailbox, we can find the IntNode
		// 	aliceIntNode, err := client.FetchIntNode(aliceMailbox.FileStructID, aliceMailbox.FileKey)
		// 	Expect(err).To(BeNil())

		// 	//get the UUID of the first FileStruct
		// 	aliceHmacStructUUID := aliceIntNode.FileStructID

		// 	var hmacData client.HMAC
		// 	hmacData, err = client.FetchHmacStruct(aliceHmacStructUUID)
		// 	Expect(err).To(BeNil())

		// 	var headFileData client.File
		// 	headFileData, err = client.FetchHeadFileStruct(hmacData.Encryption, aliceIntNode.FileKey)
		// 	Expect(err).To(BeNil())

		// 	//get the UUID of the second file struct from the first file struct

		// 	// unmarshal the secondFileStruct
		// 	var secondFileStruct client.File
		// 	secondFileStructKey, err := userlib.HashKDF(aliceIntNode.FileKey, []byte("append" + strconv.Itoa(1)))
		// 	Expect(err).To(BeNil())
		// 	secondFileStruct, err = client.FetchFileStruct(headFileData.NextNode, secondFileStructKey[:16])
		// 	Expect(err).To(BeNil())

		// 	// var thirdFileStruct client.File
		// 	// thirdFileStruct, err = client.FetchFileStruct(secondFileStruct.NextNode, aliceIntNode.FileKey)
		// 	// Expect(err).To(BeNil())

		// 	//get the uuid of the third file struct's blob

		// 	//get the dataStore map
		// 	dataStoreMap := userlib.DatastoreGetMap()

		// 	//swap the second and third file structs
		// 	temp := dataStoreMap[headFileData.UUIDContents]
		// 	dataStoreMap[headFileData.UUIDContents] = dataStoreMap[secondFileStruct.UUIDContents]
		// 	dataStoreMap[secondFileStruct.UUIDContents] = temp

		// 	_, err = alice.LoadFile("foo.txt")
		// 	Expect(err).ToNot(BeNil())
		// })

		//Specify a test for corruption of the userstruct
		Specify("Corrupt the Userstruct: Append byte, flip bytes, delete entry ", func() {
			//create alice and bob
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			//get the dataStore map
			dataStoreMap := userlib.DatastoreGetMap()
			//get alice's user struct UUID
			aliceUsernameBytes := []byte("alice")
			aliceUsernameBytesHashed := userlib.Hash(aliceUsernameBytes)[:16]
			aliceUUID, err := uuid.FromBytes(aliceUsernameBytesHashed)
			Expect(err).To(BeNil())

			//in dataStore map, get the user struct for alice using aliceUUID
			aliceUserStruct := dataStoreMap[aliceUUID]
			// corrupt the user struct for alice by appending a byte to the end of the user struct
			aliceCorruptedUserStruct := append(aliceUserStruct, []byte("a")...)
			//put it into the actual dataStore
			userlib.DatastoreSet(aliceUUID, aliceCorruptedUserStruct)

			//getUser alice
			alice, err = client.GetUser("alice", defaultPassword)
			//expect an error
			Expect(err).ToNot(BeNil())

			//get bob's user struct UUID
			bobUsernameBytes := []byte("bob")
			bobUsernameBytesHashed := userlib.Hash(bobUsernameBytes)[:16]
			bobUUID, err := uuid.FromBytes(bobUsernameBytesHashed)
			Expect(err).To(BeNil())
			//in dataStore map, get the user struct for bob using bobUUID
			bobUserStruct := dataStoreMap[bobUUID]
			//corrupt bobUserStruct by flipping bytes in the slice in a loop
			for i := 0; i < len(bobUserStruct); i++ {
				bobUserStruct[i] = ^bobUserStruct[i]
			}

			//put bobUserStruct in the actual dataStore
			userlib.DatastoreSet(bobUUID, bobUserStruct)

			bob, err = client.GetUser("bob", defaultPassword)
			//expect an error
			Expect(err).ToNot(BeNil())

			//init charles
			charles, err = client.InitUser("charles", defaultPassword)
			//expect no errors
			Expect(err).To(BeNil())
			//get charles's user struct UUID
			charlesUsernameBytes := []byte("charles")
			charlesUsernameBytesHashed := userlib.Hash(charlesUsernameBytes)[:16]
			charlesUUID, err := uuid.FromBytes(charlesUsernameBytesHashed)
			Expect(err).To(BeNil())

			//corrupt charlesUserStruct by deleting charle's user struct entry
			userlib.DatastoreDelete(charlesUUID)
			//load charles from the dataStore
			charles, err = client.GetUser("charles", defaultPassword)
			//expect an error
			Expect(err).ToNot(BeNil())

		})

		//Specify a test for Swapping User Structs
		Specify("Swap: User Structs", func() {
			//create alice and bob
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			//get Alice's user struct UUID
			aliceUsernameBytes := []byte("alice")
			aliceUsernameBytesHashed := userlib.Hash(aliceUsernameBytes)[:16]
			aliceUUID, err := uuid.FromBytes(aliceUsernameBytesHashed)
			Expect(err).To(BeNil())

			//get Bob's user struct UUID
			bobUsernameBytes := []byte("bob")
			bobUsernameBytesHashed := userlib.Hash(bobUsernameBytes)[:16]
			bobUUID, err := uuid.FromBytes(bobUsernameBytesHashed)
			Expect(err).To(BeNil())

			//get the dataStore map
			//swap the two user structs
			dataStoreMap := userlib.DatastoreGetMap()
			temp := dataStoreMap[aliceUUID]
			dataStoreMap[aliceUUID] = dataStoreMap[bobUUID]
			dataStoreMap[bobUUID] = temp
			//get alice
			alice, err = client.GetUser("alice", defaultPassword)
			//expect an error
			Expect(err).ToNot(BeNil())
			//get bob
			bob, err = client.GetUser("bob", defaultPassword)
			//expect an error
			Expect(err).ToNot(BeNil())

		})

		//specify a corruption test for the dataStore
		Specify("Corrupt the Mailbox: Append byte, flip bytes, delete entry ", func() {

			//create alice and bob and charles
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

		})
		Specify("Efficiency Append", func() {
			userlib.DebugMsg("Efficiency Append: Initializing Alice")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			alice.StoreFile(fooFile, []byte(contentOne))
			// bw := measureBandwidth(func() {
			// 	alice.AppendToFile(fooFile, []byte("hi"))
			// })

			for i := 10000; i > 0; i-- {
				alice.AppendToFile(fooFile, []byte("hi extremely long text"))
			}

			bw1 := measureBandwidth(func() {
				alice.AppendToFile(fooFile, []byte("hi"))
			})

			Expect(bw1).To(BeNumerically("<", 3600))
			

		})

		Specify("Efficiency Append: First append has 1000 words", func() {
			userlib.DebugMsg("First Append: Initializing Alice")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			alice.StoreFile(fooFile, []byte(contentOne))

			var longText []byte
			for i := 10000; i > 0; i-- {
				longText = append(longText, []byte("rhinosaurus")...)
			}

			//  bw := measureBandwidth(func() {
			// 	alice.AppendToFile(fooFile, longText)
			//  })

			bw1 := measureBandwidth(func() {
				alice.AppendToFile(fooFile, []byte("hi"))
			})

			Expect(bw1).To(BeNumerically("<", 3600))

		})

		// //Specify a test for corruption of the mailbox struct
		// Specify("Corrupt the mailbox struct", func() {
		// 	//create alice
		// 	alice, err = client.InitUser("alice", defaultPassword)
		// 	Expect(err).To(BeNil())

		// 	//make a file for alice
		// 	aliceContent := []byte("chips")
		// 	//alice stores the file
		// 	err = alice.StoreFile("foo.txt", aliceContent)
		// 	Expect(err).To(BeNil())

		// 	//find alice's mailbox UUID for foo.txt
		// 	aliceMailboxID, err := alice.GetMailboxUUID("foo.txt")
		// 	Expect(err).To(BeNil())

		// 	//corrupt the mailbox struct
		// 	//get the dataStore map
		// 	dataStoreMap := userlib.DatastoreGetMap()
		// 	//get the mailbox struct for alice
		// 	aliceMailbox := dataStoreMap[aliceMailboxID]
		// 	//append a byte to the end of the mailbox struct
		// 	dataStoreMap[aliceMailboxID] = append([]byte("a"), aliceMailbox...)
		// 	//load the file
		// 	_, err = alice.LoadFile("foo.txt")
		// 	//expect an error
		// 	Expect(err).ToNot(BeNil())

		// })

		// //specify a random swap test on random items
		// Specify("Random Swap test", func() {
		// 	//create alice and bob
		// 	alice, err = client.InitUser("alice", defaultPassword)

		// 	Expect(err).To(BeNil())
		// 	bob, err = client.InitUser("bob", defaultPassword)

		// 	Expect(err).To(BeNil())
		// 	//create a file
		// 	aliceContent := []byte("chips")
		// 	//alice stores the file
		// 	err = alice.StoreFile("foo.txt", aliceContent)
		// 	Expect(err).To(BeNil())
		// 	//alice appends
		// 	aliceAppendContent := []byte("are")
		// 	aliceAppendContent2 := []byte("awesome")

		// 	err = alice.AppendToFile("foo.txt", aliceAppendContent)
		// 	Expect(err).To(BeNil())
		// 	err = alice.AppendToFile("foo.txt", aliceAppendContent2)
		// 	Expect(err).To(BeNil())

		// 	//get the datastore map
		// 	dataStoreMap := userlib.DatastoreGetMap()
		// 	//get all the keys from the map and put it into an array
		// 	var keys_arr []uuid.UUID
		// 	for key := range dataStoreMap {
		// 		keys_arr = append(keys_arr, key)
		// 	}
		// 	//write a for loop to swap two random keys from the array, and swap the corresponding keys in the map
		// 	for i := 0; i < len(keys_arr); i++ {
		// 		rand.Seed(time.Now().UnixNano())
		// 		rand1 := rand.Intn(len(keys_arr))
		// 		rand2 := rand.Intn(len(keys_arr))
		// 		temp := dataStoreMap[keys_arr[rand1]]
		// 		dataStoreMap[keys_arr[rand1]] = dataStoreMap[keys_arr[rand2]]
		// 		dataStoreMap[keys_arr[rand2]] = temp
		// 	}

		// 	//now load the file
		// 	_, err = alice.LoadFile("foo.txt")
		// 	Expect(err).ToNot(BeNil())

		// })

	})
})
