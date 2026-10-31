# torrent-411 sponsored by https://linuxtorrents.com

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
```
supaplex@sky:~/usr/src/github-by-user/supaplextor$ torrent-411/torrent-411 torrent-411/bt-hybrid/debian-live-13.7.0-amd64-lxde.iso.67webseeds.torrent 
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
  https://cdimage.debian.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://debian-cd.mirror.web4u.cz/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.as3741.net/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.connesi.it/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.cs.nycu.edu.tw/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.envisagecloud.net.za/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.ethz.ch/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.inf.tu-dresden.de/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.interhost.co.il/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.koyanet.lv/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.lagoon.nc/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.mirror.digitalpacific.com.au/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.mirror.digitalpacific.com.au/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.mirror.vu.lt/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.nautile.nc/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.netcologne.de/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.nic.cz/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.obspm.fr/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.osuosl.org/debian-cdimage/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://debian.saix.net/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.snt.utwente.nl/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://debian.superhosting.cz/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://debian.uni-due.de/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://debian.xfree.com.ar/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp-stud.hs-esslingen.de/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.acc.umu.se/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.bme.hu/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.caliu.cat/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.cica.es/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.crifo.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.fau.de/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.hosteurope.de/mirror/ftp.debian.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.jaist.ac.jp/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.kaist.ac.kr/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.lanet.kr/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.riken.jp/Linux/debian/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://ftp.ticklers.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://ftp.uio.no/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.uni-hannover.de/debian/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://ftp.yz.yamagata-u.ac.jp/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://kartolo.sby.datautama.net.id/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://lidsol.fi-b.unam.mx/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://miroir.univ-lorraine.fr/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.141networks.com/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.as35701.net/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.cogentco.com/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.corbina.net/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.dkm.cz/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.ihost.md/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.lzu.edu.cn/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.maeen.sa/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://mirror.mephi.ru/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  http://mirror.overthewire.com.au/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.raiolanetworks.com/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.siwoo.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.steadfast.net/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.truenetwork.ru/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.twds.com.tw/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.us.leaseweb.net/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirror.xtom.com.hk/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirrors.cloudforest.co.th/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirrors.lug.mtu.edu/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirrors.ocf.berkeley.edu/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirrors.pidginhost.com/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirrors.up.pt/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://mirrorservice.org/sites/cdimage.debian.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
  https://cdimage.debian.org/debian-cd/13.7.0-live/amd64/iso-hybrid/debian-live-13.7.0-amd64-lxde.iso
```
