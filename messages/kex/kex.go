package kex

import "fmt"
import "bytes"
import "time"
import "github.com/chris-wood/spud/util"
import "github.com/chris-wood/spud/messages/name"
// import "github.com/chris-wood/spud/codec"
// import "github.com/chris-wood/spud/messages/name"
// import typedhash "github.com/chris-wood/spud/messages/hash"
// import "github.com/chris-wood/spud/messages/link"
// import "github.com/chris-wood/spud/messages/payload"
// import "github.com/chris-wood/spud/messages/validation"

import "crypto/hmac"
import "crypto/sha256"
// import "encoding/base64"
import "encoding/binary"

// XXX: use something in the standard crypto library
// import "golang.org/x/crypto/curve25519"

// Extension keys into the encoder dictionary
const _kName = "Name"
const _kSourceChallenge = "SourceChallenge"
const _kSourceToken = "SourceToken"
const _kSourceProof = "SourceProof"
const _kTimestamp = "Timestamp"
const _kMoveChallenge = "MoveChallenge"
const _kMoveProof = "MoveProof"
const _kMoveToken = "MoveToken"
const _kSessionID = "SessionID"

// Digest size(s)
const _sourceChallengeSize = 16

type KEX struct {
    messageType int 
    extensionMap map[string]interface{} // string keys to generic types
}

type kexError struct {
    prob string
}

func (e kexError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

func createChallenge(input []byte) []byte {
    h := sha256.New()
    h.Write(input)
    challenge := h.Sum(nil)
    return challenge
}

func createToken(key []byte, inputs... []byte) []byte {
    h := hmac.New(sha256.New, key)
    for _, v := range inputs {
        h.Write(v)
    }
    tag := h.Sum(nil)
    return tag
}

// Constructors

func KEXHello(prefix name.Name) *KEX {
    emap := make(map[string]interface{})
    emap[_kName] = prefix

    // Generate some random bytes
    bytes, err := util.GenerateRandomBytes(_sourceChallengeSize)
    if err != nil {
        return nil
    }
    emap[_kSourceProof] = bytes
    emap[_kSourceChallenge] = createChallenge(bytes)

    return &KEX{emap}
}

func KEXHelloReject(hello *KEX, macKey []byte) *KEX {
    emap := make(map[string]interface{})

    // Timestamp
    now := time.Now().UnixNano() / int64(time.Millisecond)
    nowBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(nowBytes, uint64(now))
    emap[_kTimestamp] = nowBytes

    // Source challenge
    challenge := hello.extensionMap[_kSourceChallenge].([]byte)
    emap[_kSourceChallenge] = challenge
    emap[_kSourceToken] = createToken(macKey, challenge, nowBytes)

    return &KEX{emap}
}

func KEXFullHello(bare, reject *KEX) *KEX {
    emap := make(map[string]interface{})

    // key share
    // curve25519 stuff

    // source token
    sourceToken := reject.extensionMap[_kSourceToken].([]byte)
    emap[_kSourceToken] = sourceToken

    // source proof
    sourceProof := bare.extensionMap[_kSourceProof]
    emap[_kSourceProof] = sourceProof

    // move challenge
    moveProof, _ := util.GenerateRandomBytes(_sourceChallengeSize)
    emap[_kMoveProof] = moveProof
    emap[_kMoveChallenge] = createChallenge(moveProof)

    // timestamp
    timestamp := reject.extensionMap[_kTimestamp].([]byte)
    emap[_kTimestamp] = timestamp

    return &KEX{emap}
}

func KEXHelloAccept(bare, reject, hello *KEX, macKey, encKey []byte) *KEX {
    emap := make(map[string]interface{})

    // re-compute the token to verify that it's correct
    sourceProof := bare.extensionMap[_kSourceProof].([]byte)
    sourceChallenge := createChallenge(sourceProof)
    timestamp := hello.extensionMap[_kTimestamp].([]byte)
    sourceToken := createToken(macKey, sourceChallenge, timestamp)

    givenToken := hello.extensionMap[_kSourceToken].([]byte)
    if bytes.Compare(sourceToken, givenToken) != 0 {
        // XXX: error here
        return nil
    }

    // key share
    // XXX: curve25519

    // certificate
    // XXX

    // move token
    moveChallenge := hello.extensionMap[_kMoveChallenge].([]byte)
    moveToken := createToken(encKey, moveChallenge)
    emap[_kMoveToken] = moveToken

    // source token
    sessionID, _ := util.GenerateRandomBytes(_sourceChallengeSize)
    emap[_kSessionID] = sessionID

    return &KEX{emap}
}
