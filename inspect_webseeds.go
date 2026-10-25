package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type parser struct {
	data     []byte
	pos      int
	infoData []byte
}

type torrentFile struct {
	path       string
	size       int64
	attributes string
	piecesRoot string
}

func (p *parser) parse() (any, error) {
	return p.parseValue(0)
}

func (p *parser) parseValue(depth int) (any, error) {
	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unexpected end of file")
	}

	switch p.data[p.pos] {
	case 'i':
		return p.parseInt()
	case 'l':
		return p.parseList(depth)
	case 'd':
		return p.parseDict(depth)
	default:
		return p.parseBytes()
	}
}

func (p *parser) parseInt() (any, error) {
	p.pos++
	end := p.findByte('e')
	if end < 0 {
		return nil, fmt.Errorf("unterminated integer")
	}

	text := string(p.data[p.pos:end])
	n, err := strconv.ParseInt(text, 10, 64)
	if err != nil || text == "" || (len(text) > 1 && text[0] == '0') || strings.HasPrefix(text, "-0") {
		return nil, fmt.Errorf("invalid integer %q", text)
	}
	p.pos = end + 1
	return n, nil
}

func (p *parser) parseBytes() ([]byte, error) {
	colon := p.findByte(':')
	if colon < 0 {
		return nil, fmt.Errorf("invalid byte string")
	}

	lengthText := string(p.data[p.pos:colon])
	length, err := strconv.Atoi(lengthText)
	if err != nil || length < 0 || (len(lengthText) > 1 && lengthText[0] == '0') {
		return nil, fmt.Errorf("invalid byte-string length %q", lengthText)
	}

	start := colon + 1
	if length > len(p.data)-start {
		return nil, fmt.Errorf("byte string extends past end of file")
	}
	end := start + length
	p.pos = end
	return p.data[start:end], nil
}

func (p *parser) parseList(depth int) ([]any, error) {
	p.pos++
	var items []any
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		item, err := p.parseValue(depth + 1)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unterminated list")
	}
	p.pos++
	return items, nil
}

func (p *parser) parseDict(depth int) (map[string]any, error) {
	p.pos++
	dict := make(map[string]any)
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		keyBytes, err := p.parseBytes()
		if err != nil {
			return nil, fmt.Errorf("invalid dictionary key: %w", err)
		}
		key := string(keyBytes)
		valueStart := p.pos
		value, err := p.parseValue(depth + 1)
		if err != nil {
			return nil, err
		}
		dict[key] = value
		if depth == 0 && key == "info" {
			p.infoData = p.data[valueStart:p.pos]
		}
	}
	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unterminated dictionary")
	}
	p.pos++
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

func decodeTorrent(data []byte) (map[string]any, []byte, error) {
	p := parser{data: data}
	decoded, err := p.parse()
	if err != nil {
		return nil, nil, err
	}
	if p.pos != len(data) {
		return nil, nil, fmt.Errorf("unexpected data after top-level value")
	}
	torrent, ok := decoded.(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("torrent file does not contain a top-level dictionary")
	}
	if _, ok := torrent["info"].(map[string]any); !ok {
		return nil, nil, fmt.Errorf("torrent file does not contain an info dictionary")
	}
	return torrent, p.infoData, nil
}

func torrentReport(torrent map[string]any, infoData []byte) string {
	info := torrent["info"].(map[string]any)
	var report strings.Builder

	fmt.Fprintf(&report, "Torrent: %s\n", printable(preferredBytes(info, "name.utf-8", "name"), "(unknown)"))
	_, hasFileTree := info["file tree"]
	isV2 := intValue(info["meta version"]) == 2 || hasFileTree
	if _, hasV1Pieces := info["pieces"].([]byte); hasV1Pieces || !isV2 {
		sha1Hash := sha1.Sum(infoData)
		fmt.Fprintf(&report, "Info hash (v1/SHA-1): %x\n", sha1Hash)
	}
	if isV2 {
		sha256Hash := sha256.Sum256(infoData)
		fmt.Fprintf(&report, "Info hash (v2/SHA-256): %x\n", sha256Hash)
	}

	if created, ok := torrent["creation date"].(int64); ok {
		fmt.Fprintf(&report, "Created: %s\n", time.Unix(created, 0).UTC().Format(time.RFC3339))
	}
	writeStringField(&report, "Created by", torrent["created by"])
	writeStringField(&report, "Comment", torrent["comment"])
	writeStringField(&report, "Encoding", torrent["encoding"])
	writeStringField(&report, "Source", info["source"])
	if private, ok := info["private"]; ok {
		fmt.Fprintf(&report, "Private: %s\n", yesNo(intValue(private) != 0))
	}

	pieceLength := intValue(info["piece length"])
	if pieceLength > 0 {
		fmt.Fprintf(&report, "Piece length: %d bytes\n", pieceLength)
	}
	if pieces, ok := info["pieces"].([]byte); ok {
		if len(pieces)%sha1.Size == 0 {
			fmt.Fprintf(&report, "Pieces: %d (20-byte SHA-1 hashes)\n", len(pieces)/sha1.Size)
		} else {
			fmt.Fprintf(&report, "Pieces: malformed hash data (%d bytes)\n", len(pieces))
		}
	}

	files := torrentFiles(info)
	var total int64
	for _, file := range files {
		total += file.size
	}
	fmt.Fprintf(&report, "\nContent (%d files, %d bytes):\n", len(files), total)
	if len(files) == 0 {
		fmt.Fprintln(&report, "  (file list unavailable)")
	}
	for _, file := range files {
		fmt.Fprintf(&report, "  %s — %d bytes", file.path, file.size)
		if file.attributes != "" {
			fmt.Fprintf(&report, " [attributes: %s]", file.attributes)
		}
		if file.piecesRoot != "" {
			fmt.Fprintf(&report, " [pieces root: %s]", file.piecesRoot)
		}
		report.WriteByte('\n')
	}

	writeTrackers(&report, torrent)
	writeURLField(&report, "Web seeds (url-list)", torrent["url-list"])
	writeURLField(&report, "HTTP seeds", torrent["httpseeds"])
	writeNodes(&report, torrent["nodes"])
	if layers, ok := torrent["piece layers"].(map[string]any); ok {
		layerCount := 0
		for _, layer := range layers {
			if data, ok := layer.([]byte); ok {
				layerCount += len(data) / sha256.Size
			}
		}
		fmt.Fprintf(&report, "Piece layers: %d SHA-256 hashes across %d layers\n", layerCount, len(layers))
	}
	return report.String()
}

func preferredBytes(dict map[string]any, preferred, fallback string) []byte {
	if value, ok := dict[preferred].([]byte); ok && len(value) > 0 {
		return value
	}
	value, _ := dict[fallback].([]byte)
	return value
}

func printable(value []byte, fallback string) string {
	if len(value) == 0 {
		return fallback
	}
	return string(value)
}

func writeStringField(report *strings.Builder, label string, value any) {
	if text, ok := value.([]byte); ok && len(text) > 0 {
		fmt.Fprintf(report, "%s: %s\n", label, text)
	}
}

func intValue(value any) int64 {
	n, _ := value.(int64)
	return n
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func torrentFiles(info map[string]any) []torrentFile {
	if entries, ok := info["files"].([]any); ok {
		files := make([]torrentFile, 0, len(entries))
		for _, entry := range entries {
			fileInfo, ok := entry.(map[string]any)
			if !ok {
				continue
			}
			path := pathFromList(fileInfo["path.utf-8"])
			if path == "" {
				path = pathFromList(fileInfo["path"])
			}
			files = append(files, torrentFile{
				path:       path,
				size:       intValue(fileInfo["length"]),
				attributes: stringValue(fileInfo["attr"]),
			})
		}
		return files
	}
	if length, ok := info["length"]; ok {
		return []torrentFile{{
			path:       printable(preferredBytes(info, "name.utf-8", "name"), "(unnamed)"),
			size:       intValue(length),
			attributes: stringValue(info["attr"]),
		}}
	}
	if tree, ok := info["file tree"].(map[string]any); ok {
		var files []torrentFile
		walkFileTree(tree, nil, &files)
		return files
	}
	return nil
}

func walkFileTree(tree map[string]any, path []string, files *[]torrentFile) {
	names := make([]string, 0, len(tree))
	for name := range tree {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		value := tree[name]
		node, ok := value.(map[string]any)
		if !ok {
			continue
		}
		if fileInfo, ok := node[""].(map[string]any); ok {
			*files = append(*files, torrentFile{
				path:       strings.Join(append(append([]string(nil), path...), name), "/"),
				size:       intValue(fileInfo["length"]),
				attributes: stringValue(fileInfo["attr"]),
				piecesRoot: hexBytes(fileInfo["pieces root"]),
			})
		}
		if name == "" {
			continue
		}
		childPath := append(append([]string(nil), path...), name)
		walkFileTree(node, childPath, files)
	}
}

func pathFromList(value any) string {
	parts, ok := value.([]any)
	if !ok {
		return ""
	}
	var path []string
	for _, part := range parts {
		if text, ok := part.([]byte); ok {
			path = append(path, string(text))
		}
	}
	return strings.Join(path, "/")
}

func stringValue(value any) string {
	text, _ := value.([]byte)
	return string(text)
}

func hexBytes(value any) string {
	data, _ := value.([]byte)
	if len(data) == 0 {
		return ""
	}
	return hex.EncodeToString(data)
}

func writeTrackers(report *strings.Builder, torrent map[string]any) {
	fmt.Fprintln(report, "\nTrackers:")
	tierCount := 0
	if tiers, ok := torrent["announce-list"].([]any); ok {
		for i, tier := range tiers {
			urls, ok := tier.([]any)
			if !ok {
				continue
			}
			for _, item := range urls {
				if url, ok := item.([]byte); ok {
					fmt.Fprintf(report, "  Tier %d: %s\n", i+1, url)
					tierCount++
				}
			}
		}
	}
	if announce, ok := torrent["announce"].([]byte); ok && len(announce) > 0 {
		fmt.Fprintf(report, "  Primary: %s\n", announce)
		tierCount++
	}
	if tierCount == 0 {
		fmt.Fprintln(report, "  (none listed)")
	}
}

func writeURLField(report *strings.Builder, label string, value any) {
	var urls [][]byte
	switch item := value.(type) {
	case []byte:
		urls = append(urls, item)
	case []any:
		for _, entry := range item {
			if url, ok := entry.([]byte); ok {
				urls = append(urls, url)
			}
		}
	}
	if len(urls) == 0 {
		return
	}
	fmt.Fprintf(report, "\n%s:\n", label)
	for _, url := range urls {
		fmt.Fprintf(report, "  %s\n", url)
	}
}

func writeNodes(report *strings.Builder, value any) {
	nodes, ok := value.([]any)
	if !ok || len(nodes) == 0 {
		return
	}
	fmt.Fprintln(report, "\nDHT nodes:")
	for _, entry := range nodes {
		node, ok := entry.([]any)
		if !ok || len(node) != 2 {
			continue
		}
		host, hostOK := node[0].([]byte)
		port, portOK := node[1].(int64)
		if hostOK && portOK {
			fmt.Fprintf(report, "  %s:%d\n", host, port)
		}
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
	torrent, infoData, err := decodeTorrent(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Torrent parse error:", err)
		os.Exit(1)
	}
	fmt.Print(torrentReport(torrent, infoData))
}
