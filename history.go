package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	historyPath := filepath.Join(homeDir, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")
	tempHistory := "stats_history.db"

	// Dosyayı kopyala (önceki örnekteki aynı mantık)
	_ = copyFile(historyPath, tempHistory)
	defer os.Remove(tempHistory)

	db, err := sql.Open("sqlite3", tempHistory)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQL: En çok ziyaret edilen ilk 10 siteyi getir
	// visit_count sütunu Chrome tarafından otomatik tutulur
	query := `
		SELECT url, visit_count, title 
		FROM urls 
		ORDER BY visit_count DESC 
		LIMIT 20`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// İstatistikleri tutmak için bir map
	domainStats := make(map[string]int)

	fmt.Println("📊 --- EN ÇOK ZİYARET EDİLEN ALAN ADLARI ---")
	
	for rows.Next() {
		var rawUrl, title string
		var count int
		rows.Scan(&rawUrl, &count, &title)

		// URL'den sadece domain kısmını çek (örn: google.com)
		parsedUrl, err := url.Parse(rawUrl)
		if err == nil {
			domainStats[parsedUrl.Host] += count
		}
	}

	// Sonuçları ekrana bas
	for domain, totalCount := range domainStats {
		fmt.Printf("🌐 %-25s | 🖱️ %d kez tıklandı\n", domain, totalCount)
	}
}

// copyFile fonksiyonunu buraya eklemeyi unutma (önceki cevaptakiyle aynı)
