using System.IO;
using System.Threading.Tasks;
using DataSync.Core.Services;

namespace DataSync.Functions
{
    public class BlobSyncTrigger
    {
        private readonly IDataSyncEngine _syncEngine;

        public BlobSyncTrigger(IDataSyncEngine syncEngine)
        {
            _syncEngine = syncEngine;
        }

        public async Task RunAsync(Stream blobStream, string name)
        {
            // Simulated Azure Blob Trigger invocation
            await _syncEngine.ProcessAndSyncFileAsync(blobStream, name, "enterprise-tenant-100");
        }
    }
}
