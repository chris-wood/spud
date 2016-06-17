package name_segment

type NameSegment struct {
    segmentType uint16
    SegmentValue string `json:"segment"`
}

// Constructor functions

func Parse(segmentString string) *NameSegment {
    // TODO: the type needs to be derived or inferred
    return &NameSegment{segmentType: 0, SegmentValue: segmentString}
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

// String functions

func (ns NameSegment) String() string {
    return ns.SegmentValue
}
