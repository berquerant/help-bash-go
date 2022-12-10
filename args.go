package main

import (
	"fmt"
	"strings"
)

// FixArgs translates `-abc` style arguments into `-a -b -c` style.
//
// flagChars is a set of a boolean flag names.
// flagWithArgChars is a set of names of flags with an argument.
func FixArgs(args []string, flagChars, flagWithArgChars string) []string {
	return NewArgsFixer(flagChars, flagWithArgChars).FixArgs(args)
}

func NewArgsFixer(flagChars, flagWithArgChars string) ArgsFixer {
	return &argsFixer{
		flagChars:        flagChars,
		flagWithArgChars: flagWithArgChars,
	}
}

type ArgsFixer interface {
	FixArgs(args []string) []string
}

type argsFixer struct {
	flagChars        string
	flagWithArgChars string
}

func (a *argsFixer) FixArgs(args []string) []string {
	var ret []string
	for _, arg := range args {
		ret = append(ret, a.fixArg(arg)...)
	}
	return ret
}

func (a *argsFixer) fixArg(x string) []string {
	if x == "" {
		return nil
	}
	if x == "-" || !strings.HasPrefix(x, "-") {
		return []string{x}
	}

	x = x[1:] // remove head -
	var (
		ret []string
		add = func(s string) {
			ret = append(ret, s)
		}
		addFlag = func(r rune) {
			ret = append(ret, fmt.Sprintf("-%s", string(r)))
		}
		arg    []rune
		addArg = func(r rune) {
			arg = append(arg, r)
		}
		flushArg = func() {
			if len(arg) == 0 {
				return
			}
			add(string(arg))
			arg = nil
		}
	)
	for _, c := range x {
		switch {
		case strings.ContainsRune(a.flagWithArgChars, c):
			flushArg()
			addFlag(c)
		case strings.ContainsRune(a.flagChars, c):
			flushArg()
			addFlag(c)
		default:
			addArg(c)
		}
	}
	flushArg()
	return ret
}
