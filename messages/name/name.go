package name

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name_segment"
import "strings"
import "fmt"

// import "encoding/json"

type Name struct {
    Segments []*name_segment.NameSegment `json:"segments"`
}

type nameError struct {
    prob string
}

func (e nameError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Name parsing functions

func parseNameStringWithoutSchema(nameString string) ([]*name_segment.NameSegment, error) {
    segments := make([]*name_segment.NameSegment, 0)
    for _, segmentString := range(strings.Split(nameString, "/")) {
        segments = append(segments, name_segment.Parse(segmentString))
    }
    return segments, nil
}

func parseNameString(nameString string) ([]*name_segment.NameSegment, error) {
    if index := strings.Index(nameString, "ccnx:/"); index == 0 {
        return parseNameStringWithoutSchema(nameString[6:])
    } else if index := strings.Index(nameString, "/"); index == 0 {
        return parseNameStringWithoutSchema(nameString[1:])
    } else {
        return nil, nameError{"rawr"}
    }
}

// Constructor functions

func Parse(nameString string) (*Name, error) {
    parsedSegments, err := parseNameString(nameString)
    if (err != nil) {
        // TODO: what to return here?
        return nil, nameError{"rawr"}
    }
    return &Name{Segments: parsedSegments}, nil
}

// TLV interface functions

func (n Name) Type() uint16 {
    return uint16(codec.T_NAME)
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
