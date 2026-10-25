# go-my-webseed-lister

A dependency-free Go command-line tool that reads a `.torrent` file and prints
its metadata in a human-readable report.

```sh
go run inspect_torrent.go file.torrent
```

The report includes v1 and (when applicable) v2 info hashes, torrent name,
creation and creator metadata, privacy/source fields, piece length and piece
hash counts, file paths and sizes, tracker tiers, web/HTTP seeds, DHT nodes,
and v2 piece-layer/file-root metadata where present. The info hashes are
computed from the original bencoded `info` dictionary, as required by the
torrent formats.
