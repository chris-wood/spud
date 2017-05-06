package name

import (
	"strings"
	"regexp"
)

type RegularName struct {
	pattern string
	regexer *regexp.Regexp
}

func CreateFromString(name string) (*RegularName, error) {
	if len(name) == 0 {
		return nil, nil
	}

	newName := strings.Replace(name, "><", "/", -1)
	newName = strings.Replace(newName, "<", "/", -1)
	newName = strings.Replace(newName, ">", "", -1)

	matcher, err := regexp.Compile(newName)
	if err != nil {
		return nil, err
	}

	return &RegularName{newName, matcher}, nil
}

// API

func (r RegularName) Accepts(other *Name) bool {
	nameString := other.String()
	matched := r.regexer.MatchString(nameString)
	return matched
}

func (n RegularName) Prefix(num int) string {
	return ""
}

func (n RegularName) IsPrefix(other *RegularName) bool {
	return false
}

func (n RegularName) SegmentStrings() []string {
	return nil
}

func (n RegularName) AppendComponent(component string) (*RegularName, error) {
	return nil, nil
}

func (n *RegularName) DropSuffix() {

}
