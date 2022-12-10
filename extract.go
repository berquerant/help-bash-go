package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/berquerant/logger"
)

//go:generate go run github.com/berquerant/goconfig@latest -type "CommentThreshold uint,NeedFunc bool,NeedVar bool" -option -configOption Option -output extract_config_generated.go

// Extract finds comments to be written to `w` as a documentation from `r`.
//
// # Option:
//
//   - WithCommentThreshold
//
// Number of rows of which threshold to determine whether to find the top-level (outside of any statements) comments.
// e.g. the value is 3 then finds 3 or more lines of top-level comments.
//
//   - WithNeedFunc
//
// Find the top-level function documentations.
//
//   - WithNeedVar
//
// Find the top-level variable documentations.
func Extract(w io.Writer, r io.Reader, opt ...Option) error {
	extractor := NewExtractor(w, opt...)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := extractor.Read(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return extractor.Flush()
}

type Extractor interface {
	Read(string) error
	Flush() error
}

func NewExtractor(w io.Writer, opt ...Option) Extractor {
	config := NewConfigBuilder().
		NeedFunc(false).
		NeedVar(false).
		CommentThreshold(3).
		Build()
	config.Apply(opt...)
	return &extractor{
		config: config,
		w:      w,
	}
}

type extractor struct {
	config *Config
	buf    []string
	w      io.Writer
}

func (e *extractor) Read(line string) error {
	if e.readComment(line) {
		return nil
	}
	if err := e.readFunction(line); err != nil {
		return err
	}
	if err := e.readVariable(line); err != nil {
		return err
	}
	return e.writeComment(line)
}

func (e *extractor) writeComment(line string) error {
	if IsTopLevelComment(line) {
		return nil
	}

	logger.G().Trace("writeComment()")
	if e.size() >= int(e.config.CommentThreshold.Get()) {
		if err := e.writeBuf(); err != nil {
			return fmt.Errorf("WriteComment: %w", err)
		}
	}
	e.reset()
	return nil
}

func (e *extractor) readVariable(line string) error {
	if e.size() == 0 {
		return nil
	}
	head, ok := ExtractVariableDeclHead(line)
	if !ok {
		return nil
	}

	logger.G().Trace("readVariable()")
	if e.config.NeedVar.Get() {
		if _, err := fmt.Fprintf(e.w, "Variable:%s\n", head); err != nil {
			return fmt.Errorf("WriteVariable: %w", err)
		}
		if err := e.writeBuf(); err != nil {
			return fmt.Errorf("WriteVariable: %w", err)
		}
	}
	e.reset()
	return nil
}

func (e *extractor) readFunction(line string) error {
	if e.size() == 0 {
		return nil
	}
	head, ok := ExtractFunctionHead(line)
	if !ok {
		return nil
	}

	logger.G().Trace("readFunction()")
	if e.config.NeedFunc.Get() {
		if _, err := fmt.Fprintf(e.w, "Function:%s\n", head); err != nil {
			return fmt.Errorf("WriteFunction: %w", err)
		}
		if err := e.writeBuf(); err != nil {
			return fmt.Errorf("WriteFunction: %w", err)
		}
	}
	e.reset()
	return nil
}

func (e *extractor) readComment(line string) bool {
	x, ok := ExtractTopLevelComment(line)
	if ok {
		logger.G().Trace("readComment()")
		e.add(x)
	}
	e.debug(line)
	return ok
}

func (e *extractor) Flush() error {
	defer e.reset()
	if e.size() >= int(e.config.CommentThreshold.Get()) {
		return e.writeBuf() // last comments
	}
	return nil
}

func (e *extractor) debug(line string) {
	logger.G().Debug("[read] %s", line)
	logger.G().Debug("[buf][%d]%v", e.size(), e.buf)
}

func (e *extractor) writeBuf() error {
	if e.size() == 0 {
		return nil
	}
	for _, x := range e.buf {
		if _, err := fmt.Fprintln(e.w, x); err != nil {
			return fmt.Errorf("WriteBuffer: %w", err)
		}
	}
	if _, err := fmt.Fprintln(e.w); err != nil {
		return fmt.Errorf("WriteBuffer: %w", err)
	}
	return nil
}

func (e *extractor) add(line string) { e.buf = append(e.buf, line) }
func (e *extractor) reset()          { e.buf = nil }
func (e *extractor) size() int       { return len(e.buf) }
