using System.IO;
using System.Threading;
using System.Threading.Tasks;
using DataSync.Core.Models;

namespace DataSync.Core.Services
{
    public interface IDataSyncEngine
    {
        Task<EnterpriseFileRecord> ProcessAndSyncFileAsync(
            Stream fileStream,
            string fileName,
            string tenantId,
            CancellationToken cancellationToken = default);

        Task<int> BatchSyncFilesAsync(
            int fileCount,
            CancellationToken cancellationToken = default);
    }
}

