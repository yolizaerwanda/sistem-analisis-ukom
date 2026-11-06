# Sistem Analisis UKOM

Sistem analisis peserta UKOM yang dibangun dengan **Go** dan **MySQL**, terintegrasi dengan **Google Spreadsheet** sebagai sumber data.

---

## Persyaratan Sistem

Pastikan software berikut sudah terpasang sebelum memulai:

| Software | Versi / Keterangan |
|----------|--------------------|
| [Go] | **1.26.0** |
| MySQL | Bisa menggunakan [XAMPP] (MySQL/MariaDB) atau server lain yang menyediakan MySQL |
| Akun Google | Untuk membuat Service Account dan mengakses Google Spreadsheet |

Cek versi Go yang terpasang:

```bash
go version
```

---

## Langkah Instalasi dan Menjalankan Sistem

### 1. Konfigurasi Integrasi Google Spreadsheet

Sistem menggunakan **Service Account** agar aplikasi Go dapat membaca Google Spreadsheet.

1. **Buat Service Account** melalui [Google Cloud Console](https://console.cloud.google.com/) (menu *IAM & Admin* → *Service Accounts*).
2. **Aktifkan Google Sheets API** pada project yang sama (menu *APIs & Services* → *Library* → cari *Google Sheets API* → *Enable*).
3. **Download JSON credential** dari Service Account yang sudah dibuat (tab *Keys* → *Add Key* → *Create new key* → *JSON*).
4. **Simpan file tersebut sebagai `credentials.json`** di root project.
5. **Share spreadsheet** ke email Service Account (alamat email berformat `nama@project-id.iam.gserviceaccount.com`) dengan akses minimal *Viewer*.
6. Setelah sistem berjalan, **uji kembali tombol "Proses Analisis"** untuk memastikan integrasi berhasil.

---

### 2. Konfigurasi File `.env`

Buat file `.env` di root project dengan isi berikut:

```env
SESSION_SECRET=isi_sendiri
GOOGLE_SPREADSHEET_ID=url_terakhir_spreadsheet
GOOGLE_CREDENTIALS_FILE=credentials.json
```

| Variabel | Keterangan |
|----------|------------|
| `SESSION_SECRET` | String acak rahasia untuk keamanan session. Isi sendiri dengan nilai yang sulit ditebak. |
| `GOOGLE_SPREADSHEET_ID` | ID spreadsheet, yaitu bagian terakhir dari URL spreadsheet. Contoh: `https://docs.google.com/spreadsheets/d/`**`ID_SPREADSHEET`**`/edit` |
| `GOOGLE_CREDENTIALS_FILE` | Nama file credential Service Account (default: `credentials.json`). |

---

### 3. Persiapan Database MySQL

1. Jalankan MySQL (misalnya lewat XAMPP) lalu buat database bernama **`db-sistem-analisis-ukom`**:

   ```sql
   CREATE DATABASE `db-sistem-analisis-ukom`;
   USE `db-sistem-analisis-ukom`;
   ```

   > Nama database boleh diubah sesuai kebutuhan. Jika diubah, sesuaikan juga nilai **DSN** pada file `config/database.go` **baris ke-10**.

2. Jalankan query berikut untuk membuat seluruh tabel:

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE peserta (
    id_peserta INT AUTO_INCREMENT PRIMARY KEY,
    nama_peserta VARCHAR(150) NOT NULL,
    batch VARCHAR(50) NOT NULL,
    persentase_kehadiran DECIMAL(5,2),
    rata_rata_to DECIMAL(5,2)
);

CREATE TABLE analisis (
    id_analisis INT AUTO_INCREMENT PRIMARY KEY,
    periode INT NOT NULL,
    tahun YEAR NOT NULL,
    institusi VARCHAR(150) NOT NULL,
    batch VARCHAR(50) NOT NULL
);

CREATE TABLE `pengerjaan-to` (
    id_to INT AUTO_INCREMENT PRIMARY KEY,
    id_analisis INT NOT NULL,
    id_peserta INT NOT NULL,
    total_pengerjaan_to INT NOT NULL DEFAULT 0,
    total_belum_pengerjaan_to INT NOT NULL DEFAULT 0,
    total_to INT NOT NULL DEFAULT 0,

    CONSTRAINT fk_to_analisis
        FOREIGN KEY (id_analisis)
        REFERENCES analisis(id_analisis)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_to_peserta
        FOREIGN KEY (id_peserta)
        REFERENCES peserta(id_peserta)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE zona (
    id_zona INT AUTO_INCREMENT PRIMARY KEY,
    id_analisis INT NOT NULL UNIQUE,
    jumlah_merah INT NOT NULL DEFAULT 0,
    jumlah_kuning INT NOT NULL DEFAULT 0,
    jumlah_hijau INT NOT NULL DEFAULT 0,
    total_peserta INT NOT NULL DEFAULT 0,

    CONSTRAINT fk_zona_analisis
        FOREIGN KEY (id_analisis)
        REFERENCES analisis(id_analisis)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE absensi (
    id_absensi INT AUTO_INCREMENT PRIMARY KEY,
    id_analisis INT NOT NULL,
    id_peserta INT NOT NULL,
    persentase_kehadiran DECIMAL(5,2),
    kehadiran_kurang_dari_40 BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_absensi_analisis
        FOREIGN KEY (id_analisis)
        REFERENCES analisis(id_analisis)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_absensi_peserta
        FOREIGN KEY (id_peserta)
        REFERENCES peserta(id_peserta)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

ALTER TABLE analisis
ADD UNIQUE KEY unique_analisis (
    periode,
    tahun,
    institusi,
    batch
);

ALTER TABLE peserta
ADD COLUMN id_analisis INT NOT NULL AFTER id_peserta;

ALTER TABLE peserta
ADD CONSTRAINT fk_peserta_analisis
FOREIGN KEY (id_analisis)
REFERENCES analisis(id_analisis)
ON DELETE CASCADE
ON UPDATE CASCADE;
```

3. Jika nama database, user, atau password MySQL Anda berbeda, ubah **DSN** di `config/database.go` (baris ke-10). Contoh format DSN:

   ```go
   dsn := "root:@tcp(127.0.0.1:3306)/db-sistem-analisis-ukom?charset=utf8mb4&parseTime=True&loc=Local"
   ```

---

### 4. Membuat Akun User

Jalankan perintah berikut untuk membuat akun login:

```bash
go run . create-user
```

Ikuti instruksi yang muncul di terminal untuk mengisi username dan password.

---

### 5. Menjalankan Sistem

```bash
go run .
```

Setelah berjalan, buka alamat yang tampil di terminal (umumnya `http://localhost:<port>`) lalu login menggunakan akun yang dibuat pada langkah sebelumnya.

---

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| Tombol **Proses Analisis** gagal / data tidak terbaca | Pastikan spreadsheet sudah di-share ke email Service Account dan Google Sheets API sudah aktif. |
| `credentials.json` tidak ditemukan | Pastikan file berada di root project dan namanya sesuai `GOOGLE_CREDENTIALS_FILE` di `.env`. |
| Gagal koneksi database | Pastikan MySQL sudah berjalan dan DSN di `config/database.go` (baris ke-10) sudah benar. |
| Error saat `CREATE TABLE` | Pastikan database sudah dibuat dan dipilih (`USE ...`), serta tabel dijalankan sesuai urutan di atas. |
