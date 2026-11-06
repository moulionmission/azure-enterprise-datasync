using System;

namespace DataSync.Core.Models
{
    public enum ProcessingStatus
    {
        Pending,
        Processing,
        Synchronized,
        Failed,
        Quarantined
    }

    public class EnterpriseFileRecord
    {
        public string FileId { get; set; } = Guid.NewGuid().ToString("N");
        public string FileName { get; set; } = string.Empty;
        public string ContentType { get; set; } = "application/octet-stream";
        public long FileSizeBytes { get; set; }
        public DateTime IngestedTimestamp { get; set; } = DateTime.UtcNow;
        public string TenantId { get; set; } = string.Empty;
        public string ChecksumSHA256 { get; set; } = string.Empty;
        public ProcessingStatus Status { get; set; } = ProcessingStatus.Pending;
    }
}
