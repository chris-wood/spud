package messages

// import "encoding/json"

type NameSegment struct {
    // TOOD: this really has a type and whatnot
    SegmentValue string `json:"segment"`
}

func (ns NameSegment) Type() uint16 {
    return 1
}

func (ns NameSegment) Length() uint16 {
    return uint16(len(ns.SegmentValue))
}

func (ns NameSegment) Value() []byte {
    return []byte(ns.SegmentValue)
}
