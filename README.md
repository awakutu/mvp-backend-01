# MVP-Frontend-01

# Cara Melakukan Kolaborasi Kode Sumber Terbuka
 
Disarankan telah belum memiliki git, [ install segera ]( https://git-scm.com/downloads ).

## 1. Fork repositori ini

Fork repositori ini dengan cara menekan tombol fork yang ada di bagian atas kanan layar.
Hal tersebut akan membuat salinan repositori ini di akun anda. (Tujuannya agar tidak ada timpang tindih dengan repo utama)

## 2. Klon repositori

Sekarang klon repositori ini ke komputer anda. Tekan tombol clone lalu tekan ikon "copy to clipboard".

Buka terminal dan eksekusi perintah git berikut:

```
git clone "url yang telah anda salin"
```
yang mana "url yang telah anda salin" (tanpa tanda petik) adalah url ke repositori ini. Lihat langkah sebelumnya untuk mendapatkan url.

Contoh:
```
git clone https://github.com/username-anda/first-contributions.git
```
Dimana `username-anda` adalah username Github anda. Pada langkah ini anda menduplikasi konten dari repositori first-contributions di GitHub ke komputer anda.

## 3. Membuat Cabang

Ganti ke direktori repositori di komputer (jika belum ada di sana):

```
cd first-contributions
```
Buat cabang dengan perintah `git checkout`:
```
git checkout -b <add-nama>
```

Contoh:
```
git checkout -b halaman-registrasi
```
(Nama cabang disarankan sesuai laman yang dikerjakan.)

## 4. Buat perubahan yang diperlukan lalu commit perubahan tersebut

Lakukan perubahan di dalam  dan jangan lupa memberikan pesan yang jelas di commit.


## 5. Dorong perubahan ke GitLab

Dorong perubahan menggunakan perintah `git push`:
```
git push origin <halaman-registrasi>
```
Ganti `<halaman-registrasi>` dengan nama cabang yang sebelumnya telah dibuat.

## 6. Submit perubahan untuk diulas

Jika anda membuka repositori anda di GitLab, maka akan ada tombol `Compare & pull request`. Tekan tombol tersebut.

Submit pull request.

Suatu saat nanti saya akan melakukan penggabungan terhadap semua perubahan anda ke cabang master proyek ini. Anda akan mendapatkan pemberitahuan melalui email setelah perubahan tersebut selesai digabungkan. eheh

Cabang master dari fork anda tidak akan memiliki perubahan-perubahan tersebut. Untuk membuat fork anda selaras dengan milik saya, ikuti langkah-langkah berikut.
