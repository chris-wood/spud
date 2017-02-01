package flic

import "fmt"

import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/codec"

const hashGroupDataPointerType uint8 = 0
const hashGroupManifestPointerType uint8 = 1

type SizedDataPointer struct {
    size uint64
    ptrHash hash.Hash
}

func (p SizedDataPointer) GetPointerType() uint8 {
    return hashGroupDataPointerType
}

func (p SizedDataPointer) GetSize() uint64 {
    return p.size
}

func (p SizedDataPointer) GetPointer() hash.Hash {
    return p.ptrHash
}

type SizedManifestPointer struct {
    size uint64
    ptrHash hash.Hash
}

func (p SizedManifestPointer) GetPointerType() uint8 {
    return hashGroupManifestPointerType
}

func (p SizedManifestPointer) GetSize() uint64 {
    return p.size
}

func (p SizedManifestPointer) GetPointer() hash.Hash {
    return p.ptrHash
}

type SizedPointer interface {
    GetPointerType() uint8
    GetSize() uint64
    GetPointer() hash.Hash
}

type HashGroupMetadata struct {
    Locator name.Name
    OverallByteCount uint64
    OverallDataDigest hash.Hash
}

type HashGroup struct {
    metadata *HashGroupMetadata
    Pointers []SizedPointer
}

type BlockHashGroup struct {
    metadata *HashGroupMetadata
}

type FLIC struct {
    groups []HashGroup
}

type flicError struct {
    prob string
}

func (e flicError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func CreateFromChunker(chunk chunker.Chunker) *FLIC {
    return &FLIC{make([]HashGroup, 0)}
}

func CreateFromTLV(topLevelTLV codec.TLV) (*FLIC, error) {
    // var result FLIC
    // var err error

    // containers := make([]codec.TLV, 0)
    // for _, tlv := range(topLevelTLV.Children()) {
        // if tlv.Type() == codec.T_NAME {
        //     contentName, err = name.CreateFromTLV(tlv)
        //     if err != nil {
        //         return &result, err
        //     }
        // } else if tlv.Type() == codec.T_PAYLOAD {
        //     dataPayload = payload.Create(tlv.Value())
        // } else if tlv.Type() == codec.T_KEX {
        //     kex, err := kex.CreateFromTLV(tlv)
        //     if err != nil {
        //         return nil, contentError{"Unable to decode the KEX TLV"}
        //     }
        //     containers = append(containers, kex)
        // } else {
        //     fmt.Printf("Unable to parse content TLV type: %d\n", tlv.Type())
        // }
    // }

    return &FLIC{make([]HashGroup, 0)}, nil
}

// Containers

func (f *FLIC) AddContainer(container codec.TLV) {
    // no-op
}

func (f *FLIC) GetContainer(containerType uint16) (codec.TLV, error) {
    var container codec.TLV
    return container, flicError{"FLIC types do not support containers"}
}

// TLV functions

func (f FLIC) Type() uint16 {
    return uint16(codec.T_MANIFEST)
}

func (f FLIC) TypeString() string {
    return "FLIC"
}

func (f FLIC) Length() uint16 {
    length := uint16(0)

    // if c.name.Length() > 0 {
    //     length += c.name.Length() + 4
    // }
    //
    // if c.dataPayload.Length() > 0 {
    //     length += c.dataPayload.Length() + 4
    // }
    //
    // for _, container := range(c.containers) {
    //     length += container.Length() + 4
    // }

    return length
}

func (f FLIC) Value() []byte {
    // e := codec.Encoder{}
    value := make([]byte, 0)

    // if c.name.Length() > 0 {
    //     value = append(value, e.EncodeTLV(c.name)...)
    // }
    //
    // if c.dataPayload.Length() > 0 {
    //     value = append(value, e.EncodeTLV(c.dataPayload)...)
    // }
    //
    // for _, container := range(c.containers) {
    //     value = append(value, e.EncodeTLV(container)...)
    // }

    return value
}

func (f FLIC) Children() []codec.TLV {
    children := []codec.TLV{}
    // XXX: append all the inner TLVs
    return children
}

func (f FLIC) String() string {
    // return Identifier(c)
    // return c.name.String()
}

// Message functions

func (f FLIC) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(c)
    return bytes
}

func (f FLIC) Name() name.Name {
    return f.name
}

func (f FLIC) GetPacketType() uint16 {
    return codec.T_OBJECT
}

func (f FLIC) Payload() payload.Payload {
    var data payload.Payload
    return data // empty payload
}

func (f FLIC) PayloadType() uint16 {
    return codec.T_PAYLOADTYPE_MANIFEST
}
