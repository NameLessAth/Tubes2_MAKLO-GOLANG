# WikiRace with BFS and IDS
Tugas Besar 2 IF2211 Strategi Algoritma -  Breadth First Search (BFS) dan Depth First Search (DFS)

<p align="center">
  <img height="360px" src="https://i.ibb.co/9NcPQYb/maklobg.png" alt="preview web"/>
  <br>
  <a><i><sup>Preview Web kelompok "MAKLO GOLANG"</sup></i></a>
</p>

## Anggota 
1. Maria Flora Renata Siringoringo (13522010)
2. M Athaullah Daffa Kusuma M (13522044)
3. Dzaky Satrio Nugroho (13522059)

## Deskripsi Singkat
Program ini merupakan tugas besar 2 dari mata kuliah IF2211 Strategi Algoritma. Program ini berfungsi untuk Menyelesaikan permainan WikiRace, yaitu mencari shortest path dari sebuah artikel wiki ke sebuah artikel wiki lainnya. 

## Informasi Tambahan
Program dibuat dengan : go1.22.2 windows/amd64, HTML, dan CSS.
IDE yang digunakan : Visual Studio Code dengan banyak extension lainnya
Laporan dibuat dengan : Google Docs 

## Petunjuk Cara Menjalankan Program dan lainnya

### Cara Menjalankan Program
Ada 2 langkah untuk menjalankan program, yaitu:
1. Jalankan backend berupa main.go dengan
```
cd src/backend
go run main.go ids.go bfs.go
```
2. Jalankan frontend dengan mendeploy localhost (disarankan menginstall "Live Server" by Ritwick Dey extension di VSCode dan tekan tombol go live di pojok kanan bawah)
### Input
Pengguna nantinya akan mengisi 2 input, yaitu judul artikel mulai dan judul artikel tujuan. Kemudian pengguna menjalankan program dengan menekan tombol IDS/BFS sesuai dengan keinginan pengguna. 


## Commit Messages

Untuk type mengikuti semantic berikut.

- `feat`: (new feature for the user, not a new feature for build script)
- `fix`: (bug fix for the user, not a fix to a build script)
- `docs`: (changes to the documentation)
- `style`: (formatting, missing semi colons, etc; no production code change)
- `refactor`: (refactoring production code, eg. renaming a variable)
- `test`: (adding missing tests, refactoring tests; no production code change)
- `chore`: (updating grunt tasks etc; no production code change)