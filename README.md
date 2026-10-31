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

## Cross-platform build archives

Use the provided Makefile to build for every GOOS/GOARCH target supported by
your local Go toolchain and produce one ZIP file per target.

```sh
make package-all
```

Some targets may require external cross-compilers/cgo and cannot be built in a
given environment. Those targets are reported and skipped. Windows `386` and
`amd64` are required and the build fails if either does not compile.

Artifacts are written to:

- `dist/bin/` for compiled binaries
- `dist/zip/` for per-target zip files (each zip includes the binary and `README.md`)
