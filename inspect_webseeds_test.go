package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestDecodeTorrentAndReport(t *testing.T) {
	info := "d6:lengthi123e4:name4:test12:piece lengthi16e6:pieces40:" + strings.Repeat("a", 40) + "e"
	data := []byte("d8:announce14:http://tracker13:announce-listll14:http://tracker" +
		"ee7:comment4:test13:creation datei1700000000e10:created by3:app8:url-list12:https://seed4:info" + info + "e")

	torrent, rawInfo, err := decodeTorrent(data)
	if err != nil {
		t.Fatalf("decodeTorrent() error = %v", err)
	}
	if string(rawInfo) != info {
		t.Fatalf("captured info dictionary = %q, want %q", rawInfo, info)
	}

	report := torrentReport(torrent, rawInfo)
	hash := sha1.Sum([]byte(info))
	for _, expected := range []string{
		"Torrent: test",
		fmt.Sprintf("Info hash (v1/SHA-1): %x", hash),
		"Created by: app",
		"Comment: test",
		"Piece length: 16 bytes",
		"Pieces: 2 (20-byte SHA-1 hashes)",
		"Content (1 files, 123 bytes):",
		"test — 123 bytes",
		"Primary: http://tracker",
		"Tier 1: http://tracker",
		"https://seed",
	} {
		if !strings.Contains(report, expected) {
			t.Errorf("report does not contain %q:\n%s", expected, report)
		}
	}
}

func TestDecodeTorrentRejectsTrailingData(t *testing.T) {
	_, _, err := decodeTorrent([]byte("d4:infod4:name1:xeextra"))
	if err == nil || !strings.Contains(err.Error(), "unexpected data") {
		t.Fatalf("decodeTorrent() error = %v, want trailing-data error", err)
	}
}

func TestV2OnlyReport(t *testing.T) {
	infoData := []byte("original bencoded info")
	torrent := map[string]any{
		"info": map[string]any{
			"name":         []byte("dir"),
			"meta version": int64(2),
			"file tree": map[string]any{
				"file": map[string]any{
					"": map[string]any{
						"length":      int64(5),
						"pieces root": []byte(strings.Repeat("b", 32)),
					},
				},
			},
		},
	}

	report := torrentReport(torrent, infoData)
	hash := sha256.Sum256(infoData)
	if !strings.Contains(report, fmt.Sprintf("Info hash (v2/SHA-256): %x", hash)) {
		t.Errorf("report is missing the v2 info hash:\n%s", report)
	}
	if strings.Contains(report, "Info hash (v1/SHA-1)") {
		t.Errorf("v2-only report incorrectly includes a v1 info hash:\n%s", report)
	}
	if !strings.Contains(report, "file — 5 bytes") || !strings.Contains(report, "pieces root: "+strings.Repeat("62", 32)) {
		t.Errorf("report is missing v2 file details:\n%s", report)
	}
}
