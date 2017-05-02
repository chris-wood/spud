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
	splits := strings.Split(nameString, "/")
	for index, segmentString := range splits {
		nextSegment, err := name_segment.Parse(segmentString)
		if err != nil {
			return nil, err
		}

		if !(index == len(splits)-1 && nextSegment.Length() == 0) {
			segments = append(segments, nextSegment)
		}
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

func Parse(nameString string) (*Name, error) {
	parsedSegments, err := parseNameString(nameString)
	if err != nil {
		// TODO: what to return here?
		return nil, nameError{"rawr"}
	}
	return &Name{Segments: parsedSegments}, nil
}

func New(segments []name_segment.NameSegment) *Name {
	return &Name{Segments: segments}
}

func CreateFromTLV(nameTlv codec.TLV) (*Name, error) {
	children := make([]name_segment.NameSegment, 0)
	for _, child := range nameTlv.Children() {
		segment, err := name_segment.CreateFromTLV(child)
		if err != nil {
			return nil, err
		}
		children = append(children, segment)
	}
	return &Name{Segments: children}, nil
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
	for _, ns := range n.Segments {
		length += ns.Length() + 4
	}
	return length
}

func (n Name) Value() []byte {
	value := make([]byte, 0)

	e := codec.Encoder{}
	for _, segment := range n.Segments {
		value = append(value, e.EncodeTLV(segment)...)
	}

	return value
}

func (n Name) Children() []codec.TLV {
	children := make([]codec.TLV, 0)
	for _, child := range n.Segments {
		children = append(children, child)
	}
	return children
}

// API

func (n Name) Prefix(num int) string {
	if num >= len(n.Segments) {
		num = len(n.Segments)
	}

	prefix := "/"
	for i := 0; i < num-1; i++ {
		prefix += n.Segments[i].String() + "/"
	}
	prefix += n.Segments[num-1].String()

	return prefix
}

func (n Name) IsPrefix(other *Name) bool {
	if other == nil {
		return false
	}

	if len(other.Segments) < len(n.Segments) {
		return false
	}

	for i := len(n.Segments) - 1; i >= 0; i-- {
		if n.Segments[i] != other.Segments[i] {
			return false
		}
	}

	return true
}

func (n Name) SegmentStrings() []string {
	segments := make([]string, 0)
	for _, v := range n.Segments {
		segments = append(segments, v.String())
	}
	return segments
}

func (n Name) AppendComponent(component string) (*Name, error) {
	segment, err := name_segment.Parse(component)
	if err != nil {
		return nil, err
	}
	return &Name{Segments: append(n.Segments[:], segment)}, nil
}

func (n *Name) DropSuffix() {
	n.Segments = n.Segments[0: len(n.Segments) - 1]
}

// String functions

func (n Name) String() string {
	return "ccnx:" + n.Prefix(len(n.Segments))
}
