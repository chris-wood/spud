package kex

import "fmt"
import "bytes"
import "time"
import "github.com/chris-wood/spud/util"
// import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/codec"
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
    messageType uint16
    extensionMap map[string]KEXExtension
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

func KEXHello() *KEX {
    emap := make(map[string]KEXExtension)

    // Generate some random bytes
    bytes, err := util.GenerateRandomBytes(_sourceChallengeSize)
    if err != nil {
        return nil
    }
    emap[_kSourceProof] = KEXExtension{codec.T_KEX_SOURCE_PROOF, bytes}
    emap[_kSourceChallenge] = KEXExtension{codec.T_KEX_SOURCE_CHALLENGE, createChallenge(bytes)}

    return &KEX{codec.T_KEX_BAREHELLO, emap}
}

func KEXHelloReject(hello *KEX, macKey []byte) *KEX {
    emap := make(map[string]KEXExtension)

    // Timestamp
    now := time.Now().UnixNano() / int64(time.Millisecond)
    nowBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(nowBytes, uint64(now))
    emap[_kTimestamp] = KEXExtension{codec.T_KEX_TIMESTAMP, nowBytes}

    // Source challenge
    challenge := hello.extensionMap[_kSourceChallenge]
    emap[_kSourceChallenge] = challenge
    emap[_kSourceToken] = KEXExtension{codec.T_KEX_SOURCE_TOKEN, createToken(macKey, challenge.ExtValue, nowBytes)}

    return &KEX{codec.T_KEX_REJECT, emap}
}

func KEXFullHello(bare, reject *KEX) *KEX {
    emap := make(map[string]KEXExtension)

    // key share
    // curve25519 stuff

    // source token
    emap[_kSourceToken] = reject.extensionMap[_kSourceToken]

    // source proof
    emap[_kSourceProof] = bare.extensionMap[_kSourceProof]

    // move challenge
    moveProof, _ := util.GenerateRandomBytes(_sourceChallengeSize)
    emap[_kMoveProof] = KEXExtension{codec.T_KEX_MOVE_PROOF, moveProof}
    emap[_kMoveChallenge] = KEXExtension{codec.T_KEX_MOVE_CHALLENGE, createChallenge(moveProof)}

    // timestamp
    emap[_kTimestamp] = reject.extensionMap[_kTimestamp]

    return &KEX{codec.T_KEX_HELLO, emap}
}

func KEXHelloAccept(bare, reject, hello *KEX, macKey, encKey []byte) *KEX {
    emap := make(map[string]KEXExtension)

    // re-compute the token to verify that it's correct
    sourceProof := bare.extensionMap[_kSourceProof]
    sourceChallenge := createChallenge(sourceProof.ExtValue)
    timestamp := hello.extensionMap[_kTimestamp]
    sourceToken := createToken(macKey, sourceChallenge, timestamp.ExtValue)

    givenToken := hello.extensionMap[_kSourceToken]
    if bytes.Compare(sourceToken, givenToken.ExtValue) != 0 {
        // XXX: error here
        return nil
    }

    // key share
    // XXX: curve25519

    // certificate
    // XXX

    // move token
    moveChallenge := hello.extensionMap[_kMoveChallenge]
    moveToken := createToken(encKey, moveChallenge.ExtValue)
    emap[_kMoveToken] = KEXExtension{codec.T_KEX_MOVE_TOKEN, moveToken}

    // source token
    sessionID, _ := util.GenerateRandomBytes(_sourceChallengeSize)
    emap[_kSessionID] = KEXExtension{codec.T_KEX_SESSION_ID, sessionID}

    return &KEX{codec.T_KEX_ACCEPT, emap}
}

// TLV functions

func (kex KEX) Type() uint16 {
    return kex.messageType
}

func (kex KEX) TypeString() string {
    return "KEX"
}

func (kex KEX) bareHelloValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceChallenge])...)

    return value
}

func (kex KEX) rejectValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    value = append(value, e.EncodeTLV(kex.extensionMap[_kTimestamp])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceToken])...)

    return value
}

func (kex KEX) helloValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    // XXX: key share!
    value = append(value, e.EncodeTLV(kex.extensionMap[_kTimestamp])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceToken])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceProof])...)

    return value
}

func (kex KEX) acceptValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    // XXX: key share!
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSessionID])...)

    return value
}

func (kex KEX) Value() []byte  {
    switch (kex.messageType) {
    case codec.T_KEX_BAREHELLO:
        return kex.bareHelloValue()
    case codec.T_KEX_REJECT:
        return kex.rejectValue()
    case codec.T_KEX_HELLO:
        return kex.helloValue()
    case codec.T_KEX_ACCEPT:
        return kex.acceptValue()
    }
    return nil
}

func (kex KEX) bareHelloLength() uint16 {
    helloValue := kex.bareHelloValue()
    return uint16(len(helloValue))
}

func (kex KEX) rejectLength() uint16 {
    rejectValue := kex.rejectValue()
    return uint16(len(rejectValue))
}

func (kex KEX) helloLength() uint16 {
    helloValue := kex.helloValue()
    return uint16(len(helloValue))
}

func (kex KEX) acceptLength() uint16 {
    acceptValue := kex.acceptValue()
    return uint16(len(acceptValue))
}

func (kex KEX) Length() uint16 {
    switch (kex.messageType) {
    case codec.T_KEX_BAREHELLO:
        return kex.bareHelloLength()
    case codec.T_KEX_REJECT:
        return kex.rejectLength()
    case codec.T_KEX_HELLO:
        return kex.helloLength()
    case codec.T_KEX_ACCEPT:
        return kex.acceptLength()
    }
    return 0
}

func (kex KEX) bareHelloChildren() []codec.TLV {
    return make([]codec.TLV, 0)
}

func (kex KEX) rejectChidlren() []codec.TLV {
    return make([]codec.TLV, 0)
}

func (kex KEX) helloChildren() []codec.TLV {
    return make([]codec.TLV, 0)
}

func (kex KEX) acceptChildren() []codec.TLV {
    return make([]codec.TLV, 0)
}

func (kex KEX) Children() []codec.TLV {
    children := make([]codec.TLV, 0)

    // for _, child := range(n.Segments) {
    //     children = append(children, child)
    // }

    return children
}
