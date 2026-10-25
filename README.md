# torrent-411

A dependency-free Go command-line tool that reads a `.torrent` file and prints
its metadata in a human-readable report.

```sh
go run inspect_torrent.go file.torrent
(or)
torrent-411 torrent-411/bt-hybrid/debian-live-13.7.0-amd64-lxde.iso.torrent
Torrent: debian-live-13.7.0-amd64-lxde.iso
Info hash (v1/SHA-1): 78fb2a1b08069f2dd415302b7ab517bb82eb3619
Created: 2026-09-12T11:35:50Z
Created by: mktorrent 1.1
Comment: Debian CD from cdimage.debian.org
Piece length: 262144 bytes
Pieces: 14500 (20-byte SHA-1 hashes)

Content (1 files, 3800989696 bytes):
  debian-live-13.7.0-amd64-lxde.iso — 3800989696 bytes

Trackers:
  Primary: http://bttracker.debian.org:6969/announce

Web seeds (url-list):
  https://cdimage.debian.org/cdimage/release/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://cdimage.debian.org/cdimage/archive/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
```

The report includes v1 and (when applicable) v2 info hashes, torrent name,
creation and creator metadata, privacy/source fields, piece length and piece
hash counts, file paths and sizes, tracker tiers, web/HTTP seeds, DHT nodes,
and v2 piece-layer/file-root metadata where present. The info hashes are
computed from the original bencoded `info` dictionary, as required by the
torrent formats.
