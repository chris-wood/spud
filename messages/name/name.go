package name

import "github.com/chris-wood/spud/codec"
import "strings"
import "fmt"

// import "encoding/json"

type Name struct {
    Segments []NameSegment `json:"segments"`
}

type NameError struct {
    prob string
}


func (e NameError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Name parsing functions

func parseNameStringWithoutSchema(nameString string) ([]NameSegment, error) {
    segments := make([]NameSegment, 0)
    for _, segmentString := range(strings.Split(nameString, "/")) {
        segments = append(segments, NameSegment{segmentString})
    }
    return segments, nil
}

func parseNameString(nameString string) ([]NameSegment, error) {
    if index := strings.Index(nameString, "ccnx:/"); index == 0 {
        return parseNameStringWithoutSchema(nameString[6:])
    } else if index := strings.Index(nameString, "/"); index == 0 {
        return parseNameStringWithoutSchema(nameString[1:])
    } else {
        return nil, NameError{"rawr"}
    }
}

// Constructor functions

func New(nameString string) (*Name, error) {
    parsedSegments, err := parseNameString(nameString)
    if (err != nil) {
        // TODO: what to return here?
        return nil, NameError{"rawr"}
    }
    return &Name{Segments: parsedSegments}, nil
}

// TLV interface functions

func (n Name) Type() uint16 {
    return uint16(1)
}

func (n Name) TypeString() string {
    return "Name"
}

func (n Name) Length() uint16 {
    length := uint16(0)
    for _, ns := range(n.Segments) {
        length += ns.Length() + 4
    }
    return length
}

func (n Name) Value() []byte  {
    value := make([]byte, 0)

    e := codec.Encoder{}
    for _, segment := range(n.Segments) {
        value = append(value, e.Encode(segment)...)
    }

    return value
}

// String functions

func (n Name) String() string {
    segmentStrings := make([]string, len(n.Segments))
    for index, segment := range(n.Segments) {
        segmentStrings[index] = segment.String()
    }
    return "ccnx:/" + strings.Join(segmentStrings, "/")
}
