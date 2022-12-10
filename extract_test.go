package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type extractTestcase struct {
	title string
	opt   []Option
	input []string
	want  string
}

func (e *extractTestcase) test(t *testing.T) {
	t.Run(e.title, func(t *testing.T) {
		var buf bytes.Buffer
		assert.Nil(
			t,
			Extract(
				&buf,
				bytes.NewBufferString(strings.Join(e.input, "\n")),
				e.opt...,
			),
		)
		assert.Equal(t, e.want, buf.String())
	})
}

func TestExtract(t *testing.T) {
	for _, tc := range []extractTestcase{
		{
			title: "empty input",
		},
		{
			title: "ignore function",
			input: []string{
				"# Function",
				"func() {",
				"  echo",
				"}",
			},
		},
		{
			title: "function",
			opt:   []Option{WithNeedFunc(true)},
			input: []string{
				"# Function comment",
				"func() {",
				"  echo",
				"}",
			},
			want: `Function:func()
Function comment

`,
		},
		{
			title: "ignore variable",
			input: []string{
				"# Variable comment",
				"varname=varval",
			},
		},
		{
			title: "variable",
			opt:   []Option{WithNeedVar(true)},
			input: []string{
				"# Variable comment",
				"varname=varval",
			},
			want: `Variable:varname=varval
Variable comment

`,
		},
		{
			title: "ignore comment",
			input: []string{
				"# 1",
				"# 2",
			},
		},
		{
			title: "comment",
			input: []string{
				"# 1",
				"# 2",
				"# 3",
			},
			want: `1
2
3

`,
		},
	} {
		tc := tc
		tc.test(t)
	}
}
