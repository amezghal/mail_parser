package mail_parser

import (
	"unicode/utf8"
)

const (
	alphanum = iota
	group1

	specialChars1 // [\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]
	specialChars2 // [\x01-\x09\x0b\x0c\x0e-\x7f]
	specialChars3 // [\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5ax53-\x7f]
)

type Validator struct {
	input    []rune
	cursor   int
	max      int
	rollback func(int, bool)
}

func (v *Validator) Compile(email string) *Validator {

	for _, c := range email {
		v.input = append(v.input, c)
	}

	v.input = append(v.input, utf8.RuneError)
	v.max = len(v.input)

	return v
}

func New() *Validator {
	v := &Validator{}
	v.rollback = func(cursor int, valid bool) {
		if !valid {
			v.cursor = cursor
		}
	}
	return v
}

func (v *Validator) peek() rune {
	if v.cursor == v.max-1 {
		return utf8.RuneError
	}
	return v.input[v.cursor+1]
}

func (v *Validator) current() rune {
	if v.cursor > v.max-1 {
		return utf8.RuneError
	}
	return v.input[v.cursor]
}

func (v *Validator) next() bool {
	if v.cursor+1 > v.max-1 {
		return false
	}

	v.cursor++
	return true
}

// Validate test
func (v *Validator) Validate() bool {
	valid := false
	if v.state1BeforeAT() || v.state2BeforeAT() {
		if v.current() == '@' && v.next() {
			if v.state1AfterAT() || v.state2AfterAT() {
				valid = true
			}
		}
	}
	return valid && v.cursor == v.max-1
}

// region beforeAT
func (v *Validator) state1BeforeAT() bool {
	cursor := v.cursor
	if !v.consumeRule1() {
		return false
	}

	valid := true
	if v.current() == '.' {
		for v.current() != utf8.RuneError {
			cc := v.cursor
			if v.current() == '.' && v.next() && v.consumeRule1() {
				valid = true
			} else {
				v.cursor = cc
				if v.current() == '.' {
					valid = false
				}
				break
			}
		}
	}

	if !valid {
		v.cursor = cursor
	}
	return valid
}

func (v *Validator) state2BeforeAT() bool {
	cursor := v.cursor
	if v.current() == '"' && v.next() {

		if v.current() == '"' && v.next() { // enclosing '"'
			return true
		}

		valid := false
		for v.current() != utf8.RuneError {
			if isGroupX(specialChars1, v.current()) {
				v.next()
				valid = true
			} else if v.current() == '\\' && v.next() {

				if isGroupX(specialChars2, v.current()) {
					v.next()
					valid = true
				} else {
					valid = false
					break
				}
			} else {
				break
			}
		}

		if valid && v.current() == '"' && v.next() { // enclosing '"'
			return true
		}
	}
	v.cursor = cursor
	return false
}

// endregion beforeAT

// region afterAT
func (v *Validator) state1AfterAT() bool {
	cursor := v.cursor
	checkBlock := func() bool {
		if isGroupX(alphanum, v.current()) {
			v.next()
			for v.current() != utf8.RuneError {

				if isGroupX(alphanum, v.current()) {
					v.next()
				} else if v.current() == '-' {
					// check trailing -
					for v.current() != utf8.RuneError {
						if v.peek() == '-' {
							v.next()
						} else {
							break
						}
					}
					if isGroupX(alphanum, v.peek()) {
						v.next()
					} else {
						return false
					}
				} else {
					break
				}
			}
			return true
		}
		return false
	}

	valid := false
	for v.cursor < utf8.RuneError {
		c := v.cursor
		if checkBlock() && v.current() == '.' {
			valid = v.next()
		} else {
			v.cursor = c
			break
		}
	}

	if valid && checkBlock() {
		return true
	}

	v.cursor = cursor
	return false
}

func (v *Validator) state2AfterAT() bool {
	cursor := v.cursor
	valid := true
	defer func() {
		if !valid {
			v.cursor = cursor
		}
	}()

	checkBlock1 := func() bool {
		valid := true
		c := v.cursor
		switch {
		case v.current() == '2' && v.peek() == '5':
			v.next()
			v.next()
			if !(v.current() >= '0' && v.current() <= '5') {
				valid = false
				break
			} else {
				v.next()
			}
			break
		case v.current() == '2' && v.peek() >= '0' && v.peek() <= '4':
			v.next()
			v.next()
			if !(v.current() >= '0' && v.current() <= '9') {
				valid = false
				break
			} else {
				v.next()
			}
			break
		case v.current() == '0' || v.current() == '1':
			v.next()
			if v.current() >= '0' && v.current() <= '9' {
				v.next()
				if v.current() >= '0' && v.current() <= '9' {
					v.next()
				}
			} else {
				valid = false
				break
			}
			break
		case v.current() >= '0' && v.current() <= '9':
			v.next()
			if v.current() >= '0' && v.current() <= '9' {
				v.next()
			}
			break
		default:
			valid = false
		}
		if !valid {
			v.cursor = c
		}
		return valid
	}

	checkBlock2 := func() bool {
		valid := true
		cursor := v.cursor

		for v.current() != utf8.RuneError {

			if isGroupX(alphanum, v.current()) {
				v.next()
			} else if v.current() == '-' {
				// check trailing -
				for v.current() != utf8.RuneError {
					if v.peek() == '-' {
						v.next()
					} else {
						break
					}
				}
				if isGroupX(alphanum, v.peek()) {
					v.next()
				} else {
					valid = false
				}
			} else {
				break
			}
		}
		if valid && v.cursor > cursor {
			return true
		}
		v.cursor = cursor
		return false
	}
	checkBlock3 := func() bool {
		valid := false
		cursor = v.cursor
		if v.current() == '\\' {
			for v.cursor < v.max-2 {
				if isGroupX(specialChars2, v.current()) {
					v.next()
					valid = true
				} else {
					if isGroupX(specialChars3, v.current()) {
						valid = false
					}
					break
				}
			}
		}

		if !valid {
			cursor = v.cursor // rollback
			for v.cursor < v.max-2 {
				if isGroupX(specialChars3, v.current()) {
					v.next()
					valid = true
				} else {
					break
				}
			}
		}

		if !valid {
			v.cursor = cursor
		}
		return valid
	}

	if v.current() == '[' {
		v.next()
		for i := 0; i < 3; i++ {
			if checkBlock1() && v.current() == '.' {
				v.next()
			} else {
				valid = false
				return false
			}
		}
		if checkBlock1() {

		} else if checkBlock2() {
			if v.current() == ':' && v.next() && checkBlock3() {
				valid = true
			} else {
				valid = false
			}
		} else {
			valid = false
		}
	} else {
		valid = false
	}

	if v.current() == ']' {
		v.next()
	} else {
		valid = false
	}

	return valid
}

// endregion afterAT

func (v *Validator) consumeRule1() bool {
	c := v.cursor
	for v.current() != utf8.RuneError {
		if isGroupX(group1, v.current()) {
			v.cursor++
		} else {
			break
		}
	}
	return v.cursor > c
}

func isGroupX(groupID int, c rune) bool {
	if c == utf8.RuneError {
		return false
	}
	switch groupID {
	case alphanum:
		return c >= 'a' && c <= 'z' || c >= '0' && c <= '9'
	case group1:
		// check alpha num
		if isGroupX(alphanum, c) {
			return true
		}
		// check special chars
		for _, cc := range "!#$%&'*+/=?^_`{|}~-" {
			if cc == c {
				return true
			}
		}
	case specialChars1:
		// [\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]
		cc := byte(c)
		return cc >= '\x01' && cc <= '\x08' || cc >= '\x0e' && cc <= '\x1f' || cc >= '\x23' && cc <= '\x5b' || cc >= '\x5d' && cc <= '\x7f' ||
			cc == '\x0b' || cc == '\x0c' || cc == '\x21'
	case specialChars2:
		// [\x01-\x09\x0b\x0c\x0e-\x7f]
		cc := byte(c)
		return cc >= '\x01' && cc <= '\x09' || cc >= '\x0e' && cc <= '\x7f' || cc == '\x0b' || cc == '\x0c'
	case specialChars3:
		// [\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]
		cc := byte(c)
		return cc >= '\x01' && cc <= '\x08' || cc >= '\x0e' && cc <= '\x1f' || cc >= '\x21' && cc <= '\x5a' || cc >= '\x53' && cc <= '\x7f' ||
			cc == '\x0b' || cc == '\x0c'
	}
	return false
}
