package sync

import (
	"bytes"
	"context"
	"testing"

	"github.com/moulionmission/azure-enterprise-datasync/pkg/auth"
)

func TestProcessAndSyncFile(t *testing.T) {
	engine := NewDataSyncEngine(8)
	content := []byte("Column1,Column2,Value\nDataA,DataB,100.5")
	reader := bytes.NewReader(content)

	record, err := engine.ProcessAndSyncFile(context.Background(), reader, "enterprise_sample.csv", "tenant_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if record.FileName != "enterprise_sample.csv" {
		t.Errorf("expected filename 'enterprise_sample.csv', got '%s'", record.FileName)
	}
	if record.FileSizeBytes != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), record.FileSizeBytes)
	}
	if record.Status != StatusSynchronized {
		t.Errorf("expected status SYNCHRONIZED, got %s", record.Status)
	}
	if len(record.ChecksumSHA256) != 64 {
		t.Errorf("expected 64-char SHA256 checksum, got %s", record.ChecksumSHA256)
	}
}

func TestBenchmarkDailyThroughput(t *testing.T) {
	engine := NewDataSyncEngine(16)
	totalFiles := 10000

	rate, elapsed := engine.BenchmarkDailyThroughput(totalFiles)
	t.Logf("Processed %d files in %v (Rate: %.2f files/sec)", totalFiles, elapsed, rate)

	// In an asynchronous environment, 10k files should finish in < 2 seconds
	if rate < 2000 {
		t.Errorf("expected throughput > 2000 files/sec, got %.2f", rate)
	}
}

func TestEntraIDValidator(t *testing.T) {
	validator := auth.NewEntraIDValidator("tenant-prod-100")

	validClaims := &auth.Claims{
		TenantID: "tenant-prod-100",
		Roles:    []string{"DataSync.Operator", "DataSync.Reader"},
	}

	err := validator.ValidateToken(validClaims, "DataSync.Operator")
	if err != nil {
		t.Fatalf("expected valid token to pass, got: %v", err)
	}

	invalidTenantClaims := &auth.Claims{
		TenantID: "tenant-unauthorized",
		Roles:    []string{"DataSync.Operator"},
	}
	err = validator.ValidateToken(invalidTenantClaims, "DataSync.Operator")
	if err == nil {
		t.Fatal("expected unauthorized tenant to fail")
	}
}

