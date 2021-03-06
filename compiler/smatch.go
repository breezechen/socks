package compiler

import (
	"errors"
	"regexp"
	"strings"
)

type SMatch struct {
	template string
	regexp   *regexp.Regexp
}

var (
	ErrUnrecognizedSMatch = errors.New("smatch resolution fails, unrecognized smatch expression")
)

func NewSMatch(rule string) (match SMatch, err error) {

	if rule[0] != 's' && rule[0] != 'S' {
		return match, ErrUnrecognizedSMatch // errors.New("invalid rule head: " + rule)
	}

	if rule[1] != '@' && rule[0] != '|' {
		return match, ErrUnrecognizedSMatch // errors.New("invalid character segmentation: " + rule)
	}

	split := strings.Split(rule, rule[1:2])

	if len(split) != 4 {
		return match, ErrUnrecognizedSMatch // errors.New("rule string incomplete or invalid: " + rule)
	}

	match.regexp, err = regexp.Compile("(?" + split[3] + ")" + split[1])

	if err != nil {
		return match, ErrUnrecognizedSMatch // err
	}

	match.template = split[2]

	return match, nil
}

func (this *SMatch) Replace(src string) (string, error) {
	var dst []byte

	submatch := this.regexp.FindStringSubmatchIndex(src)

	if len(submatch) == 0 {
		return src, errors.New("regular expression does not match")
	}

	return string(this.regexp.ExpandString(dst, this.template, src, submatch)), nil
}
