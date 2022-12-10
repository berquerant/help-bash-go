package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTopLevelComment(t *testing.T) {
	for _, tc := range []struct {
		title   string
		input   string
		matched bool
		want    string
	}{
		{
			title:   "empty string",
			input:   "",
			matched: false,
			want:    "",
		},
		{
			title:   "empty comment",
			input:   "#",
			matched: true,
			want:    "",
		},
		{
			title:   "not comment",
			input:   "usage() {",
			matched: false,
			want:    "",
		},
		{
			title:   "space only",
			input:   "# ",
			matched: true,
			want:    "",
		},
		{
			title:   "spaces only",
			input:   "#  ",
			matched: true,
			want:    " ",
		},
		{
			title:   "comment",
			input:   "#Comment",
			matched: true,
			want:    "Comment",
		},
		{
			title:   "comment with space",
			input:   "# Comment",
			matched: true,
			want:    "Comment",
		},
		{
			title:   "comment with spaces",
			input:   "#  Comment",
			matched: true,
			want:    " Comment",
		},
		{
			title:   "sharps",
			input:   "## Comment",
			matched: true,
			want:    "Comment",
		},
		{
			title:   "commented out func",
			input:   "#func() {",
			matched: true,
			want:    "func() {",
		},
		{
			title:   "spaced sharps",
			input:   "# # comment",
			matched: true,
			want:    "# comment",
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, matched := ExtractTopLevelComment(tc.input)
			assert.Equal(t, tc.matched, matched)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestExtractFunctionHead(t *testing.T) {
	for _, tc := range []struct {
		title   string
		input   string
		matched bool
		want    string
	}{
		{
			title:   "empty string",
			input:   "",
			matched: false,
			want:    "",
		},
		{
			title:   "not function head",
			input:   "x=y",
			matched: false,
			want:    "",
		},
		{
			title:   "function head",
			input:   "func(){",
			matched: true,
			want:    "func()",
		},
		{
			title:   "comment",
			input:   "#func(){",
			matched: false,
			want:    "",
		},
		{
			title: "ref var",
			input: "$(varname)",
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, matched := ExtractFunctionHead(tc.input)
			assert.Equal(t, tc.matched, matched)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestExtractVariableDeclHead(t *testing.T) {
	for _, tc := range []struct {
		title   string
		input   string
		matched bool
		want    string
	}{
		{
			title:   "empty string",
			input:   "",
			matched: false,
			want:    "",
		},
		{
			title:   "not variable decl",
			input:   "func() {",
			matched: false,
			want:    "",
		},
		{
			title:   "variable decl",
			input:   "x=y",
			matched: true,
			want:    "x=y",
		},
		{
			title:   "readonly var",
			input:   "readonly x=y",
			matched: true,
			want:    "readonly x=y",
		},
		{
			title:   "comment",
			input:   "#x=y",
			matched: false,
			want:    "",
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, matched := ExtractVariableDeclHead(tc.input)
			assert.Equal(t, tc.matched, matched)
			assert.Equal(t, tc.want, got)
		})
	}
}
