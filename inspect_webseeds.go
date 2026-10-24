package main

import (
	"fmt"
	"os"
	"strconv"
)

type parser struct {
	data []byte
	pos  int
}

func (p *parser) parse() (any, error) {
	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unexpected end of file")
	}

	switch p.data[p.pos] {
	case 'i':
		return p.parseInt()
	case 'l':
		return p.parseList()
	case 'd':
		return p.parseDict()
	default:
		return p.parseBytes()
	}
}

func (p *parser) parseInt() (any, error) {
	p.pos++ // Skip 'i'
	end := p.findByte('e')
	if end < 0 {
		return nil, fmt.Errorf("unterminated integer")
	}

	n, err := strconv.ParseInt(string(p.data[p.pos:end]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer: %w", err)
	}
	p.pos = end + 1
	return n, nil
}

func (p *parser) parseBytes() (any, error) {
	colon := p.findByte(':')
	if colon < 0 {
		return nil, fmt.Errorf("invalid byte string")
	}

	length, err := strconv.Atoi(string(p.data[p.pos:colon]))
	if err != nil || length < 0 {
		return nil, fmt.Errorf("invalid byte-string length")
	}

	start := colon + 1
	end := start + length
	if end > len(p.data) {
		return nil, fmt.Errorf("byte string extends past end of file")
	}

	p.pos = end
	return p.data[start:end], nil
}

func (p *parser) parseList() (any, error) {
	p.pos++ // Skip 'l'
	var items []any

	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		item, err := p.parse()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unterminated list")
	}
	p.pos++ // Skip 'e'
	return items, nil
}

func (p *parser) parseDict() (any, error) {
	p.pos++ // Skip 'd'
	dict := make(map[string]any)

	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		keyValue, err := p.parseBytes()
		if err != nil {
			return nil, err
		}
		key := string(keyValue.([]byte))

		value, err := p.parse()
		if err != nil {
			return nil, err
		}
		dict[key] = value
	}

	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unterminated dictionary")
	}
	p.pos++ // Skip 'e'
	return dict, nil
}

func (p *parser) findByte(b byte) int {
	for i := p.pos; i < len(p.data); i++ {
		if p.data[i] == b {
			return i
		}
	}
	return -1
}

func printURLs(key string, value any) {
	switch v := value.(type) {
	case []byte:
		fmt.Printf("%s: %s\n", key, v)
	case []any:
		for _, item := range v {
			if url, ok := item.([]byte); ok {
				fmt.Printf("%s: %s\n", key, url)
			}
		}
	default:
		fmt.Printf("%s: found, but value has an unexpected format\n", key)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s file.torrent\n", os.Args[0])
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Read error:", err)
		os.Exit(1)
	}

	p := parser{data: data}
	decoded, err := p.parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Torrent parse error:", err)
		os.Exit(1)
	}

	torrent, ok := decoded.(map[string]any)
	if !ok {
		fmt.Fprintln(os.Stderr, "Torrent file does not contain a top-level dictionary")
		os.Exit(1)
	}

	found := false
	for _, key := range []string{"url-list", "httpseeds"} {
		if value, ok := torrent[key]; ok {
			printURLs(key, value)
			found = true
		}
	}

	if !found {
		fmt.Println("No url-list or httpseeds fields found.")
	}
}
