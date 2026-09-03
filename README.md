# Tool Rename File Sertifikat

Aplikasi CLI berbasis Go untuk memudahkan rename file sertifikat hasil unduhan Canva secara massal (PDF, JPG, PNG, JPEG) berdasarkan data peserta di file Excel (.xlsx), tanpa perlu mengubah nama file satu per satu secara manual.

Hasil unduhan Canva yang biasanya hanya berupa nomor halaman (contoh: 1.pdf, 2.png) atau nama desain akan otomatis dicocokkan dengan baris data pada Excel.

## Fitur Utama

- Memudahkan penamaan sertifikat hasil ekspor massal dari Canva.
- Deteksi otomatis file Excel (.xlsx) di folder kerja.
- Deteksi otomatis kolom Nama Mahasiswa, NIM, dan Nomor Urut.
- Dukungan format sertifikat: PDF, JPG, JPEG, dan PNG (format asli tetap dipertahankan).
- Pilihan format nama baru:
  1. NIM dan Nama (contoh: Nim_Nama.pdf)
  2. Nama saja (contoh: Nama.pdf)
- Preview perubahan nama sebelum diproses.
- Pencegahan file tertimpa (overwrite protection).

## Persyaratan

- Go versi 1.21 atau lebih baru.

## Cara Penggunaan

### 1. Persiapan File

1. Unduh sertifikat dari Canva (format PDF, PNG, atau JPG). Jika berupa file ZIP, ekstrak terlebih dahulu.
2. Letakkan file sertifikat tersebut bersama file Excel (.xlsx) atau di dalam folder tersendiri (misal folder `sertifikat`).

Tool dapat mengenali file secara otomatis jika dinamai:
- Nomor urut halaman Canva (contoh: 1.png, 02.jpg, Desain_1.pdf)
- Nama peserta (contoh: Nama.pdf)
- NIM peserta (contoh: Nim.pdf)

### 2. Menjalankan Program

Jika file sertifikat berada di folder yang sama dengan aplikasi:
```bash
go run main.go
```

Jika file sertifikat berada di folder khusus (contoh: folder `sertifikat`):
```bash
go run main.go ./sertifikat
```

### 3. Alur Proses

1. Program membaca file Excel dan mendeteksi file sertifikat di folder.
2. Pilih format nama file baru:
   - Tekan 1 (atau langsung Enter): format NIM_Nama
   - Tekan 2: format Nama saja
3. Program menampilkan preview perubahan nama.
4. Tekan y (atau Enter) untuk memproses rename.

## Contoh Hasil Rename

File Asli Canva -> Opsi 1 (NIM dan Nama):
- 1.pdf -> Nim_Nama.pdf
- 2.jpg -> Nim_Nama.jpg
- 3.png -> Nim_Nama.png

File Asli Canva -> Opsi 2 (Nama Saja):
- 1.pdf -> Nama.pdf
- 2.jpg -> Nama.jpg
- 3.png -> Nama.png
