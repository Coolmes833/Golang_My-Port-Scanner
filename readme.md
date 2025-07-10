# Simple TCP Port Scanner (Go)

## About this project

This is a basic TCP port scanner written with the Go programming language.  
The main goal is to check which ports are open on a target IP or domain within a selected port range.  
The scanner uses **goroutines** and a **WaitGroup** for concurrency. This makes the scan faster by running multiple checks at the same time.

I built this to practice Go concurrency concepts and see how real TCP connect scanning works, which is one of the first steps in basic reconnaissance in cybersecurity.

---

## Key Features

- Uses Go's built-in **flag** package for easy CLI parameter management.
- Each port scan runs in a separate **goroutine**, so multiple ports are scanned in parallel.
- All goroutines are synchronized with a **sync.WaitGroup**.
- Timeout can be changed by the user to adjust connection attempts.
- Both IPv4 and IPv6 targets are supported with custom address formatting.
- Open ports and known services are matched with a `services.json` file.
- Results are saved in `open_ports.jsonl` as JSON lines, like a mini log archive.
- Shows clear info about which ports are open or closed during the scan.

---

## How to Use

Run this tool from terminal with Go.

### Example Command
go run scanner.go --host [TARGET_ADDRESS] --start [START_PORT] --end [END_PORT] --timeout [SECONDS] --workers [NUMBER]

### Example
go run scanner.go --host scanme.nmap.org --start 1 --end 1024 --timeout 1 --workers 200



### Parameters

- `--host` : Target IP or domain name (Default: scanme.nmap.org)
- `--start` : Port number to start scanning (Default: 1)
- `--end` : Port number to stop scanning (Default: 1024)
- `--timeout` : Timeout in seconds per port (Default: 1)
- `--workers` : Number of concurrent workers (Default: 100)

---

## Technical Details

- **Goroutines:** Each port is handled in its own goroutine, which gives real parallelism and much faster results than sequential scans.
- **WaitGroup:** Makes sure the main program waits until all scans finish.
- **Mutex:** Used for safe access when workers write results to the same slice.
- **IPv6 Handling:** Uses `[` and `]` for IPv6 literal addresses when needed.
- **Timeout:** Prevents hanging forever if a firewall drops packets instead of rejecting.

---

## How It Works

This scanner does a simple TCP **connect scan**:
- If a port is open and listening, the TCP handshake (`SYN`, `SYN-ACK`) will succeed and `net.Dial` will connect.
- If the port is closed, it fails fast or times out.
- This shows basic open port info, which is a typical first step during network reconnaissance in security assessments.

---

## Legal Note

This tool is for personal learning and testing only.  
Do not scan networks you do not own or do not have clear permission for.  
Unauthorized port scanning can be illegal in many countries.

---

## Author

Berkay KÖSEOĞLU



----------------------------------------------------------------------------------------------------
Türkçe

# Basit TCP Port Tarayıcı (Go)

## Proje Hakkında

Bu proje, Go programlama dili kullanılarak yazılmış basit bir TCP port tarayıcıdır.  
Amacı, verilen bir IP adresi veya alan adı üzerinde belirlenen port aralığında hangi portların açık olduğunu tespit etmektir.  
Tarama işlemi, Go'nun **goroutine** ve **WaitGroup** yapıları ile eşzamanlı çalıştırılır, bu sayede çok daha hızlı sonuç alınır.

Bu projeyi geliştirmemin nedeni, Go dilinde **concurrency**, **goroutine**, **mutex** ve temel TCP port tarama mantığını öğrenmek istememdi.  
Ayrıca siber güvenlik dünyasında ağ keşif (reconnaissance) aşamalarının nasıl ilerlediğini pratikte görmek istedim.

---

## Temel Özellikler

- Komut satırından parametre almak için Go'nun **flag** paketi kullanılmıştır.
- Her port, bağımsız bir **goroutine** içinde taranır, bu da taramayı hızlandırır.
- Tüm goroutine işlemleri **sync.WaitGroup** ile senkronize edilir.
- **Mutex**, aynı anda birden fazla goroutine sonucu yazarken veri çakışmasını engeller.
- IPv4 ve IPv6 adres formatları desteklenir.
- Tarama sırasında bağlantı süresi **timeout** parametresi ile ayarlanabilir.
- Açık portlar ve bilinen servis isimleri `services.json` dosyasından eşleştirilir.
- Tüm tarama sonuçları `open_ports.jsonl` dosyasına JSON Lines formatında kaydedilir.

---

## Nasıl Çalışır

Bu tarayıcı, temel bir TCP **connect scan** uygular:
- Hedef port açık ve dinliyorsa TCP el sıkışma (`SYN`, `SYN-ACK`) başarılı olur, `net.Dial` bağlantı kurar.
- Port kapalıysa veya güvenlik duvarı engelliyorsa bağlantı isteği başarısız olur veya zaman aşımına uğrar.
- Bu, siber güvenlikte ağ keşif (recon) aşamasında ilk adımlardan biridir.

---

## Kullanım

Uygulamayı terminal üzerinden Go ile çalıştırabilirsiniz.

### Örnek Komut
go run scanner.go --host [HEDEF_ADRES] --start [BAŞLANGIÇ_PORT] --end [BİTİŞ_PORT] --timeout [SÜRE_SN] --workers [SAYI]

### Örnek Kullanım
go run scanner.go --host scanme.nmap.org --start 1 --end 1024 --timeout 1 --workers 200


### Parametreler

- `--host` : Hedef IP adresi veya domain (Varsayılan: scanme.nmap.org)
- `--start` : Taramaya başlanacak port numarası (Varsayılan: 1)
- `--end` : Taramayı bitirecek port numarası (Varsayılan: 1024)
- `--timeout` : Her port için bağlantı zaman aşımı (Varsayılan: 1 saniye)
- `--workers` : Aynı anda çalışan işçi (goroutine) sayısı (Varsayılan: 100)

---

## Teknik Detaylar

- **Goroutine Kullanımı:** Her port ayrı bir goroutine ile taranır, bu da taramayı klasik tekli döngüye göre çok daha hızlı yapar.
- **WaitGroup:** Tüm goroutinelerin bitmesi beklenir, bu sayede işlem tamamlanmadan program kapanmaz.
- **Mutex:** Aynı slice'a birden fazla goroutine yazacağı için veri bütünlüğü `sync.Mutex` ile korunur.
- **IPv6 Formatlama:** IPv6 hedeflerinde doğru adres formatı sağlanır (`[host]:port`).
- **Timeout:** Firewall’un cevapsız bırakması durumunda sonsuza kadar beklememek için timeout uygulanır.

---

## Hukuki Uyarı

Bu proje sadece kişisel öğrenme ve test amacıyla hazırlanmıştır.  
Yetkisiz sistemlerde ağ taraması yapmak, birçok ülkede yasa dışıdır.  
Lütfen sadece kendi cihazlarınızda veya açıkça izin verilmiş test sistemlerinde kullanın.

---

## Hazırlayan

Berkay KÖSEOĞLU  

