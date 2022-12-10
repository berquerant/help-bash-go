package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgsFixer(t *testing.T) {
	const (
		flagChars        = "abc"
		flagWithArgChars = "xyz"
	)
	for _, tc := range []struct {
		title string
		args  []string
		want  []string
	}{
		{
			title: "empty input",
		},
		{
			title: "an empty string",
			args:  []string{""},
		},
		{
			title: "identity -",
			args:  []string{"-"},
			want:  []string{"-"},
		},
		{
			title: "a flag",
			args:  []string{"-a"},
			want:  []string{"-a"},
		},
		{
			title: "flags",
			args:  []string{"-ab"},
			want:  []string{"-a", "-b"},
		},
		{
			title: "a flag with argument",
			args:  []string{"-x10"},
			want:  []string{"-x", "10"},
		},
		{
			title: "flags with argument",
			args:  []string{"-x10yA"},
			want:  []string{"-x", "10", "-y", "A"},
		},
		{
			title: "mixed flags",
			args:  []string{"-x10byA"},
			want:  []string{"-x", "10", "-b", "-y", "A"},
		},
		{
			title: "2 mixed flags",
			args:  []string{"-x10byA", "-c", "-zZa"},
			want:  []string{"-x", "10", "-b", "-y", "A", "-c", "-z", "Z", "-a"},
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			fixer := NewArgsFixer(flagChars, flagWithArgChars)
			got := fixer.FixArgs(tc.args)
			assert.Equal(t, tc.want, got)
		})
	}
}
