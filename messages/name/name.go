package name

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name_segment"
import "strings"
import "fmt"

// import "encoding/json"

type Name struct {
    Segments []name_segment.NameSegment `json:"segments"`
}

type nameError struct {
    prob string
}

func (e nameError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Name parsing functions

func parseNameStringWithoutSchema(nameString string) ([]name_segment.NameSegment, error) {
    segments := make([]name_segment.NameSegment, 0)
    for _, segmentString := range(strings.Split(nameString, "/")) {
        nextSegment, err := name_segment.Parse(segmentString)
        if err != nil {
            return nil, err
        }
        segments = append(segments, nextSegment)
    }
    return segments, nil
}

func parseNameString(nameString string) ([]name_segment.NameSegment, error) {
    if index := strings.Index(nameString, "ccnx:/"); index == 0 {
        return parseNameStringWithoutSchema(nameString[6:])
    } else if index := strings.Index(nameString, "/"); index == 0 {
        return parseNameStringWithoutSchema(nameString[1:])
    } else {
        return nil, nameError{"rawr"}
    }
}

// Constructor functions

func Parse(nameString string) (Name, error) {
    var result Name
    parsedSegments, err := parseNameString(nameString)
    if (err != nil) {
        // TODO: what to return here?
        return result, nameError{"rawr"}
    }
    return Name{Segments: parsedSegments}, nil
}

func New(segments []name_segment.NameSegment) Name {
    return Name{Segments: segments}
}

func CreateFromTLV(nameTlv codec.TLV) (Name, error) {
    var result Name

    children := make([]name_segment.NameSegment, 0)

    fmt.Println(len(nameTlv.Children()))
    for _, child := range(nameTlv.Children()) {
        segment, err := name_segment.CreateFromTLV(child)
        fmt.Println(segment.String())
        if err != nil {
            return result, nil
        }
        children = append(children, segment)
    }
    return Name{Segments: children}, nil
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
        value = append(value, e.EncodeTLV(segment)...)
    }

    return value
}

func (n Name) Children() []codec.TLV {
    children := make([]codec.TLV, 0)
    for _, child := range(n.Segments) {
        children = append(children, child)
    }
    return children
}

// String functions

func (n Name) String() string {
    segmentStrings := make([]string, len(n.Segments))
    for index, segment := range(n.Segments) {
        segmentStrings[index] = segment.String()
    }
    return "ccnx:/" + strings.Join(segmentStrings, "/")
}
