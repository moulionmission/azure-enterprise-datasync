package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/moulionmission/azure-enterprise-datasync/pkg/sync"
)

func main() {
	port := flag.Int("port", 8085, "HTTP server port")
	workers := flag.Int("workers", 32, "Parallel worker pool size")
	flag.Parse()

	engine := sync.NewDataSyncEngine(*workers)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":             "healthy",
			"engine":             "azure-enterprise-datasync",
			"worker_pool_size":   *workers,
			"runtime":            "Go 1.22",
			"timestamp":          time.Now().UTC(),
		})
	})

	http.HandleFunc("/api/v1/sync", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		fileName := r.URL.Query().Get("file_name")
		if fileName == "" {
			fileName = "uploaded_enterprise_file.dat"
		}
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default-enterprise-tenant"
		}

		record, err := engine.ProcessAndSyncFile(r.Context(), r.Body, fileName, tenantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(record)
	})

	http.HandleFunc("/api/v1/benchmark", func(w http.ResponseWriter, r *http.Request) {
		totalFiles := 10000
		rate, elapsed := engine.BenchmarkDailyThroughput(totalFiles)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"benchmark":          "500,000 daily file throughput stress test",
			"files_processed":    totalFiles,
			"elapsed_duration":   elapsed.String(),
			"files_per_second":   rate,
			"projected_24h_cap":  fmt.Sprintf("%.0f files/day", rate*86400),
		})
	})

	log.Printf("Starting Azure Enterprise DataSync engine (Go) on port %d...", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

