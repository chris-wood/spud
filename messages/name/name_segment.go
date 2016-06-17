package name

type NameSegment struct {
    SegmentValue string `json:"segment"`
}

// TLV interface functions

func (ns NameSegment) Type() uint16 {
    return 1
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

// String functions

func (ns NameSegment) String() string {
    return ns.SegmentValue
}
