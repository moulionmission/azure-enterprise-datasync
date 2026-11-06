package sync

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// ProcessingStatus tracks state of file synchronization.
type ProcessingStatus string

const (
	StatusPending      ProcessingStatus = "PENDING"
	StatusSynchronized ProcessingStatus = "SYNCHRONIZED"
	StatusFailed       ProcessingStatus = "FAILED"
)

// EnterpriseFileRecord contains metadata for processed enterprise files.
type EnterpriseFileRecord struct {
	FileID            string           `json:"file_id"`
	FileName          string           `json:"file_name"`
	TenantID          string           `json:"tenant_id"`
	FileSizeBytes     int64            `json:"file_size_bytes"`
	ChecksumSHA256    string           `json:"checksum_sha256"`
	Status            ProcessingStatus `json:"status"`
	IngestedTimestamp time.Time        `json:"ingested_timestamp"`
}

// DataSyncEngine coordinates high-throughput file streaming and validation.
type DataSyncEngine struct {
	workerPoolSize int
	totalProcessed uint64
}

// NewDataSyncEngine initializes an engine configured for high concurrency.
func NewDataSyncEngine(workerPoolSize int) *DataSyncEngine {
	if workerPoolSize <= 0 {
		workerPoolSize = 32
	}
	return &DataSyncEngine{
		workerPoolSize: workerPoolSize,
	}
}

// ProcessAndSyncFile streams the file content, calculates cryptographic checksum, and commits metadata.
func (e *DataSyncEngine) ProcessAndSyncFile(ctx context.Context, reader io.Reader, fileName, tenantID string) (*EnterpriseFileRecord, error) {
	hasher := sha256.New()
	written, err := io.Copy(hasher, reader)
	if err != nil {
		return nil, fmt.Errorf("failed streaming file content: %w", err)
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	record := &EnterpriseFileRecord{
		FileID:            fmt.Sprintf("file_%d", time.Now().UnixNano()),
		FileName:          fileName,
		TenantID:          tenantID,
		FileSizeBytes:     written,
		ChecksumSHA256:    checksum,
		Status:            StatusSynchronized,
		IngestedTimestamp: time.Now().UTC(),
	}

	atomic.AddUint64(&e.totalProcessed, 1)
	return record, nil
}

// BenchmarkDailyThroughput simulates parallel file ingestion across the worker pool.
func (e *DataSyncEngine) BenchmarkDailyThroughput(totalFiles int) (float64, time.Duration) {
	start := time.Now()
	var wg sync.WaitGroup
	fileChan := make(chan int, 1000)

	// Spawn worker pool
	for w := 0; w < e.workerPoolSize; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileIdx := range fileChan {
				// 1KB payload per file
				buf := make([]byte, 1024)
				r := &mockReader{data: buf}
				_, _ = e.ProcessAndSyncFile(context.Background(), r, fmt.Sprintf("file_%d.dat", fileIdx), "tenant-enterprise")
			}
		}()
	}

	for i := 0; i < totalFiles; i++ {
		fileChan <- i
	}
	close(fileChan)
	wg.Wait()

	elapsed := time.Since(start)
	rate := float64(totalFiles) / elapsed.Seconds()
	return rate, elapsed
}

type mockReader struct {
	data []byte
	read bool
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.read {
		return 0, io.EOF
	}
	n = copy(p, m.data)
	m.read = true
	return n, nil
}

