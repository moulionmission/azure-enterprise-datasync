using System;
using System.IO;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using DataSync.Core.Models;

namespace DataSync.Core.Services
{
    public class DataSyncEngine : IDataSyncEngine
    {
        public async Task<EnterpriseFileRecord> ProcessAndSyncFileAsync(
            Stream fileStream,
            string fileName,
            string tenantId,
            CancellationToken cancellationToken = default)
        {
            if (fileStream == null) throw new ArgumentNullException(nameof(fileStream));
            if (string.IsNullOrWhiteSpace(fileName)) throw new ArgumentException("Filename required", nameof(fileName));

            using var sha256 = SHA256.Create();
            byte[] hashBytes = await sha256.ComputeHashAsync(fileStream, cancellationToken);
            string checksum = BitConverter.ToString(hashBytes).Replace("-", "").ToLowerInvariant();

            var record = new EnterpriseFileRecord
            {
                FileName = fileName,
                TenantId = tenantId,
                FileSizeBytes = fileStream.Length,
                ChecksumSHA256 = checksum,
                Status = ProcessingStatus.Synchronized,
                IngestedTimestamp = DateTime.UtcNow
            };

            return record;
        }

        public async Task<int> BatchSyncFilesAsync(
            int fileCount,
            CancellationToken cancellationToken = default)
        {
            // Simulates parallel worker pool processing enterprise multi-format files
            int completed = 0;
            var parallelOptions = new ParallelOptions
            {
                MaxDegreeOfParallelism = Environment.ProcessorCount * 2,
                CancellationToken = cancellationToken
            };

            await Parallel.ForEachAsync(
                Enumerable.Range(0, fileCount),
                parallelOptions,
                async (i, token) =>
                {
                    byte[] dummyBytes = new byte[1024];
                    using var ms = new MemoryStream(dummyBytes);
                    await ProcessAndSyncFileAsync(ms, $"enterprise_file_{i}.dat", "tenant_global", token);
                    Interlocked.Increment(ref completed);
                });

            return completed;
        }
    }
}

