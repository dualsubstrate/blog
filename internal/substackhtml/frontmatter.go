package substackhtml

import (
	"bytes"
	"errors"
)

// StripFrontMatter removes a leading YAML block delimited by --- lines. If no
// opening --- is found, the input is returned unchanged.
func StripFrontMatter(src []byte) ([]byte, error) {
	s := bytes.TrimPrefix(src, []byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
	// Allow leading whitespace/newlines before front matter.
	s = bytes.TrimLeft(s, " \t\r\n")
	if !bytes.HasPrefix(s, []byte("---")) {
		return bytes.TrimSpace(s), nil
	}
	_, rest, ok := cutLine(s)
	if !ok {
		return nil, errors.New("substackhtml: front matter opening --- not terminated")
	}
	for {
		line, after, ok := cutLine(rest)
		if !ok {
			return nil, errors.New("substackhtml: front matter not closed with ---")
		}
		if bytes.Equal(bytes.TrimSpace(line), []byte("---")) {
			return bytes.TrimSpace(after), nil
		}
		rest = after
	}
}

// FrontMatterBlock returns the leading --- YAML --- block (including delimiters)
// and the remainder of the file. If there is no front matter, ok is false.
func FrontMatterBlock(src []byte) (block, rest []byte, ok bool) {
	s := bytes.TrimPrefix(src, []byte{0xEF, 0xBB, 0xBF})
	leadingSkip := len(src) - len(s)
	s = bytes.TrimLeft(s, " \t\r\n")
	leadingSkip += len(src) - leadingSkip - len(s)
	if !bytes.HasPrefix(s, []byte("---")) {
		return nil, src, false
	}
	_, restBytes, lineOK := cutLine(s)
	if !lineOK {
		return nil, src, false
	}
	for {
		line, after, lineOK := cutLine(restBytes)
		if !lineOK {
			return nil, src, false
		}
		if bytes.Equal(bytes.TrimSpace(line), []byte("---")) {
			used := len(s) - len(after)
			block = s[:used]
			rest = src[leadingSkip+used:]
			return block, rest, true
		}
		restBytes = after
	}
}

func cutLine(b []byte) (line, after []byte, ok bool) {
	idx := bytes.IndexByte(b, '\n')
	if idx < 0 {
		return nil, nil, false
	}
	line = b[:idx]
	after = b[idx+1:]
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line, after, true
}
