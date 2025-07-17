package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

// Go'da type anahtar kelimesi ile yeni veri yapıları tanımlanır. ScanReport, tarama sonucunu özetleyen bir yapıdır. İçinde hedef IP, port aralıkları, tarama süresi, açık portlar gibi bilgiler tutulur.
type ScanReport struct {
	Host        string       `json:"host"`
	StartPort   int          `json:"start_port"`
	EndPort     int          `json:"end_port"`
	TimeoutSec  int          `json:"timeout_seconds"`
	WorkerCount int          `json:"worker_count"`
	DurationSec float64      `json:"duration_seconds"`
	Timestamp   string       `json:"timestamp"`
	OpenPorts   []PortResult `json:"open_ports"`
}

// Bu yapı ise her bir portla ilgili sonucu tutar. Port numarası, açık olup olmadığı (Open boolean tipi) ve eğer biliniyorsa hangi servise ait olduğu (Service) bilgisi.
type PortResult struct {
	Port    int    `json:"port"`
	Open    bool   `json:"-"`
	Service string `json:"service"`
}

// serviceMap: JSON dosyasından gelen port numarası - servis adı eşleşmeleri burada tutulur.
// results: Tüm portlar için elde edilen sonuçlar burada saklanır.
// mu: Eşzamanlılık (concurrency) durumlarında results listesine erişimde veri tutarlılığı için Mutex (kilitleme) kullanılır.
var (
	serviceMap map[string]string
	results    []PortResult
	mu         sync.Mutex
)

// Bu fonksiyon, services.json adlı dosyamı okur. Dosya içeriği bir JSON nesnesidir.
// Okunduktan sonra yukarıda ki serviceMap adlı map veri yapısına yüklenir. Amaç, hangi portun hangi servise ait olduğunu yazdırabilmek.
func loadServiceMap(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &serviceMap)
}

// IPv4 ve IPv6 adreslerini doğru şekilde biçimlendirmek için:
func formatAddress(host string, port int) string {
	if net.ParseIP(host).To4() == nil {
		return fmt.Sprintf("[%s]:%d", host, port) // IPv6
	}
	return fmt.Sprintf("%s:%d", host, port) //IPv4
}

// Bu fonksiyon, her bir iş parçacığının (thread/worker) ne yapacağını tanımlar. Her worker kendisine verilen portları kontrol eder.
func worker(wg *sync.WaitGroup, host string, timeout time.Duration, jobs <-chan int) {
	defer wg.Done()

	for port := range jobs {
		address := formatAddress(host, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		//Burada TCP üzerinden portun açık olup olmadığı test edilir. Eğer bağlantı kurulabiliyorsa, port açıktır. Sonuç results dizisine eklenir.
		isOpen := false
		if err == nil {
			isOpen = true
			conn.Close()
		}
		//

		mu.Lock()
		results = append(results, PortResult{
			Port:    port,
			Open:    isOpen,
			Service: portName(port),
		})
		mu.Unlock()
	}
}

// Portun numarasına göre bilinen bir servis varsa, adını döner. Yoksa “Bilinmeyen Port” yazar.
func portName(port int) string {
	key := fmt.Sprintf("%d", port)
	if name, ok := serviceMap[key]; ok {
		return name
	}
	return "Bilinmeyen Port"
}

/*
	main() Fonksiyonu içinde:

Komut satırından girilen bilgiler alınır: host, start, end, timeout, workers.

Tarama süresi başlangıç zamanı alınır: startTime := time.Now()

Worker sayısı kadar go fonksiyonu başlatılır.

Her port için görev kanalına (jobs) iş verilir.

Tüm işler bitene kadar WaitGroup ile beklenir.

Sonuçlar sıralanır ve sadece açık portlar yazdırılır.

Sonuçlar open_ports.jsonl dosyasına JSON formatında eklenir.
*/

func main() {
	// Flags
	host := flag.String("host", "scanme.nmap.org", "Hedef IP adresi veya domain")
	startPort := flag.Int("start", 1, "Başlangıç Portu")
	endPort := flag.Int("end", 6000, "Bitiş Portu")
	timeout := flag.Int("timeout", 1, "Timeout süresi (saniye)")
	numWorkers := flag.Int("workers", 100, "Aynı anda çalışan worker sayısı")

	flag.Parse()

	startTime := time.Now()

	err := loadServiceMap("services.json")
	if err != nil {
		fmt.Println("Servis listesi yüklenemedi:", err)
		return
	}

	jobs := make(chan int, 100)
	var wg sync.WaitGroup

	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, *host, time.Duration(*timeout)*time.Second, jobs)
	}

	for port := *startPort; port <= *endPort; port++ {
		jobs <- port
	}
	close(jobs)

	wg.Wait()

	duration := time.Since(startTime).Seconds()

	sort.Slice(results, func(i, j int) bool {
		if results[i].Open == results[j].Open {
			return results[i].Port < results[j].Port
		}
		return results[i].Open && !results[j].Open
	})

	var openPorts []PortResult
	for _, res := range results {
		if res.Open {
			openPorts = append(openPorts, res)
		}
	}

	report := ScanReport{
		Host:        *host,
		StartPort:   *startPort,
		EndPort:     *endPort,
		TimeoutSec:  *timeout,
		WorkerCount: *numWorkers,
		DurationSec: duration,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		OpenPorts:   openPorts,
	}

	for _, res := range results {
		if res.Open {
			fmt.Printf("Port %d (%s) AÇIK\n", res.Port, res.Service)
		} else {
			fmt.Printf("Port %d (%s) KAPALI\n", res.Port, res.Service)
		}
	}

	//os.O_APPEND: Varsa dosyanın sonuna ekle
	// os.O_CREATE: Dosya yoksa oluştur
	// 0644: Dosya izinleri (okuma-yazma)
	// Sonuçlar json.NewEncoder(file).Encode(report) ile yazılır.
	file, err := os.OpenFile("open_ports.jsonl",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(report)
	if err != nil {
		fmt.Println("JSON yazma hatası:", err)
		return
	}

	fmt.Println("Tarama raporu open_ports.jsonl dosyasına eklendi!")

}
