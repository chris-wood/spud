package kex

import "fmt"
import "bytes"
import "time"
import "github.com/chris-wood/spud/util"
import "github.com/chris-wood/spud/codec"

import "crypto/rand"
import "crypto/hmac"
import "crypto/sha256"
import "encoding/binary"

// XXX: wrap this up
import "golang.org/x/crypto/nacl/box"

// Extension keys into the encoder dictionary
const _kMessageType = "MessageType"
const _kSourceChallenge = "SourceChallenge"
const _kSourceToken = "SourceToken"
const _kSourceProof = "SourceProof"
const _kTimestamp = "Timestamp"
const _kMoveChallenge = "MoveChallenge"
const _kMoveProof = "MoveProof"
const _kMoveToken = "MoveToken"
const _kSessionID = "SessionID"
const _kPrivateKeyShare = "PrivateKeyShare"
const _kPublicKeyShare = "PublicKeyShare"

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

    typeContainer := make([]byte, 2)
    binary.BigEndian.PutUint16(typeContainer, codec.T_KEX_BAREHELLO)
    emap[_kMessageType] = KEXExtension{codec.T_KEX_MESSAGE_TYPE, typeContainer}

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

    typeContainer := make([]byte, 2)
    binary.BigEndian.PutUint16(typeContainer, codec.T_KEX_REJECT)
    emap[_kMessageType] = KEXExtension{codec.T_KEX_MESSAGE_TYPE, typeContainer}

    return &KEX{codec.T_KEX_REJECT, emap}
}

func KEXFullHello(bare, reject *KEX) *KEX {
    emap := make(map[string]KEXExtension)

    // key share
    publicKey, privateKey, err := box.GenerateKey(rand.Reader)
    if err != nil {
        return nil
    }
    emap[_kPublicKeyShare] = KEXExtension{codec.T_KEX_KEYSHARE, publicKey[:]}
    emap[_kPrivateKeyShare] = KEXExtension{0, privateKey[:]}

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

    typeContainer := make([]byte, 2)
    binary.BigEndian.PutUint16(typeContainer, codec.T_KEX_HELLO)
    emap[_kMessageType] = KEXExtension{codec.T_KEX_MESSAGE_TYPE, typeContainer}

    return &KEX{codec.T_KEX_HELLO, emap}
}

func KEXHelloAccept(hello *KEX, macKey, encKey []byte) (*KEX, error) {
    emap := make(map[string]KEXExtension)

    // re-compute the token to verify that it's correct
    sourceProof := hello.extensionMap[_kSourceProof]
    sourceChallenge := createChallenge(sourceProof.ExtValue)
    timestamp := hello.extensionMap[_kTimestamp]

    sourceToken := createToken(macKey, sourceChallenge, timestamp.ExtValue)
    givenToken := hello.extensionMap[_kSourceToken]
    if bytes.Compare(sourceToken, givenToken.ExtValue) != 0 {
        return nil, kexError{"Token verification failed."}
    }

    // key share
    publicKey, privateKey, err := box.GenerateKey(rand.Reader)
    if err != nil {
        return nil, kexError{"Key share generatio failed."}
    }
    emap[_kPublicKeyShare] = KEXExtension{codec.T_KEX_KEYSHARE, publicKey[:]}
    emap[_kPrivateKeyShare] = KEXExtension{0, privateKey[:]}

    // move token
    moveChallenge := hello.extensionMap[_kMoveChallenge]
    moveToken := createToken(encKey, moveChallenge.ExtValue)
    emap[_kMoveToken] = KEXExtension{codec.T_KEX_MOVE_TOKEN, moveToken}

    // source token
    sessionID, _ := util.GenerateRandomBytes(_sourceChallengeSize)
    emap[_kSessionID] = KEXExtension{codec.T_KEX_SESSION_ID, sessionID}

    typeContainer := make([]byte, 2)
    binary.BigEndian.PutUint16(typeContainer, codec.T_KEX_ACCEPT)
    emap[_kMessageType] = KEXExtension{codec.T_KEX_MESSAGE_TYPE, typeContainer}

    return &KEX{codec.T_KEX_ACCEPT, emap}, nil
}

// Container functions
func (k KEX) GetContainerType() uint16 {
    return codec.T_KEX
}

func (k KEX) GetContainerValue() interface{} {
    return k
}

// TLV functions

func (kex KEX) Type() uint16 {
    return codec.T_KEX
}

func (kex KEX) TypeString() string {
    return "KEX"
}

func (kex KEX) String() string {
    return kex.TypeString()
}

func CreateFromTLV(kexTLV codec.TLV) (*KEX, error) {
    emap := make(map[string]KEXExtension)

    messageType := codec.T_KEX_BAREHELLO
    for _, child := range(kexTLV.Children()) {
        extension, err := CreateExtensionFromTLV(child)
        if err != nil {
            return nil, err
        }

        // Drop the extension into the right slot
        switch (extension.ExtType) {
        case codec.T_KEX_MESSAGE_TYPE:
            emap[_kMessageType] = extension
            messageType = binary.BigEndian.Uint16(extension.Value())
            break
        case codec.T_KEX_SOURCE_CHALLENGE:
            emap[_kSourceChallenge] = extension
            break
        case codec.T_KEX_SOURCE_TOKEN:
            emap[_kSourceToken] = extension
            break
        case codec.T_KEX_SOURCE_PROOF:
            emap[_kSourceProof] = extension
            break
        case codec.T_KEX_MOVE_CHALLENGE:
            emap[_kMoveChallenge] = extension
            break
        case codec.T_KEX_MOVE_TOKEN:
            emap[_kMoveProof] = extension
            break
        case codec.T_KEX_MOVE_PROOF:
            emap[_kMoveProof] = extension
            break
        case codec.T_KEX_SESSION_ID:
            emap[_kSessionID] = extension
            break
        case codec.T_KEX_TIMESTAMP:
            emap[_kTimestamp] = extension
            break
        }
    }

    return &KEX{messageType, emap}, nil
}

func (kex KEX) bareHelloValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    value = append(value, e.EncodeTLV(kex.extensionMap[_kMessageType])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceChallenge])...)

    return value
}

func (kex KEX) rejectValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    value = append(value, e.EncodeTLV(kex.extensionMap[_kMessageType])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kTimestamp])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceToken])...)

    return value
}

func (kex KEX) helloValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    value = append(value, e.EncodeTLV(kex.extensionMap[_kMessageType])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kTimestamp])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceToken])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSourceProof])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kPublicKeyShare])...)

    return value
}

func (kex KEX) acceptValue() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    value = append(value, e.EncodeTLV(kex.extensionMap[_kMessageType])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kSessionID])...)
    value = append(value, e.EncodeTLV(kex.extensionMap[_kPublicKeyShare])...)

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
    children := make([]codec.TLV, 0)

    children = append(children, kex.extensionMap[_kMessageType])
    children = append(children, kex.extensionMap[_kSourceChallenge])

    return children
}

func (kex KEX) rejectChidlren() []codec.TLV {
    children := make([]codec.TLV, 0)

    children = append(children, kex.extensionMap[_kMessageType])
    children = append(children, kex.extensionMap[_kTimestamp])
    children = append(children, kex.extensionMap[_kSourceToken])

    return children
}

func (kex KEX) helloChildren() []codec.TLV {
    children := make([]codec.TLV, 0)

    children = append(children, kex.extensionMap[_kMessageType])
    children = append(children, kex.extensionMap[_kTimestamp])
    children = append(children, kex.extensionMap[_kSourceToken])
    children = append(children, kex.extensionMap[_kSourceProof])

    return children
}

func (kex KEX) acceptChildren() []codec.TLV {
    children := make([]codec.TLV, 0)

    children = append(children, kex.extensionMap[_kMessageType])
    children = append(children, kex.extensionMap[_kSessionID])

    return children
}

func (kex KEX) Children() []codec.TLV {
    switch (kex.messageType) {
    case codec.T_KEX_BAREHELLO:
        return kex.bareHelloChildren()
    case codec.T_KEX_REJECT:
        return kex.rejectChidlren()
    case codec.T_KEX_HELLO:
        return kex.helloChildren()
    case codec.T_KEX_ACCEPT:
        return kex.acceptChildren()
    }
    return make([]codec.TLV, 0)
}

// KEX API
func (k KEX) GetPublicKeyShare() []byte {
    return k.extensionMap[_kPublicKeyShare].ExtValue
}

func (k KEX) GetPrivateKeyShare() []byte {
    return k.extensionMap[_kPrivateKeyShare].ExtValue
}

func (k KEX) GetSessionID() string {
    return string(k.extensionMap[_kSessionID].ExtValue)
}

func (k KEX) GetMessageType() uint16 {
    return k.messageType
}
