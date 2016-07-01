package name_segment

import "fmt"
import "github.com/chris-wood/spud/codec"

type NameSegment struct {
    segmentType uint16
    SegmentValue string `json:"segment"`
}

type nameSegmentError struct {
    problem string
}

func (e nameSegmentError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

// Constructor functions

func Parse(segmentString string) (NameSegment, error) {
    // TODO: the type needs to be derived or inferred
    return NameSegment{segmentType: 1, SegmentValue: segmentString}, nil
}

func New(segmentType uint16, segmentString string) NameSegment {
    return NameSegment{segmentType: segmentType, SegmentValue: segmentString}
}

func CreateFromTLV(tlv codec.TLV) (NameSegment, error) {
    return NameSegment{segmentType: tlv.Type(), SegmentValue: string(tlv.Value())}, nil
}

// TLV interface functions

func (ns NameSegment) Type() uint16 {
    return ns.segmentType
}

func (ns NameSegment) TypeString() string {
    return "NameSegment"
}

func (ns NameSegment) Length() uint16 {
    return uint16(len(ns.SegmentValue))
}

func (ns NameSegment) Value() []byte {
    return []byte(ns.SegmentValue)
}

func (ns NameSegment) Children() []codec.TLV {
    return nil
}

// String functions

func (ns NameSegment) String() string {
    return ns.SegmentValue
}
