# go-app-scaffolder

CLI sederhana untuk membuat kerangka proyek Go berdasarkan struktur `go-app-scaffolds`.

## Fitur

- Membuat proyek baru dari template yang di-embed ke dalam binary.
- Mengganti placeholder nama proyek dan module path.
- Mempertahankan struktur Ports and Adapters / Hexagonal Architecture.
- Menolak menimpa direktori tujuan yang sudah ada, kecuali `--force` diberikan.

## Menjalankan

```sh
go run ./cmd/scaffold \
  --module github.com/username/myapp \
  --name myapp \
  --out ../myapp
```

Opsi:
- `--module`: module path untuk `go.mod` dan dokumentasi (wajib).
- `--name`: nama proyek (default: nama direktori output).
- `--out`: direktori tujuan (wajib).
- `--force`: izinkan menulis ke direktori tujuan yang sudah ada. File dengan nama sama akan ditimpa.

Template ditanam di binary menggunakan `embed`, sehingga generator tidak membutuhkan direktori template terpisah ketika dijalankan.

## Struktur proyek yang dihasilkan

```text
myapp/
├── internal/
│   ├── adapter/
│   │   ├── inbound/
│   │   └── outbound/
│   ├── core/
│   └── port/
│       ├── inbound/
│       └── outbound/
├── .gitignore
├── go.mod
├── LICENSE
├── main.go
└── README.md
```

Direktori `internals` mengikuti struktur ZIP sumber. File implementasi domain dan adapter sengaja tidak dibuat otomatis karena ZIP sumber belum berisi implementasi di direktori tersebut.

## Catatan

`main.go` pada ZIP sumber berisi logika khusus aplikasi lelang dan mengimpor package serta file schema yang tidak ada di ZIP. Karena itu, template awal generator menggunakan `main.go` minimal yang dapat dikompilasi dan dijalankan, sambil mempertahankan struktur direktori. Logika aplikasi khusus bisa ditambahkan sebagai template lanjutan.
