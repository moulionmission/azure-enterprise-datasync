using System.IO;
using System.Security.Claims;
using System.Text;
using System.Threading.Tasks;
using DataSync.Core.Models;
using DataSync.Core.Services;
using DataSync.Security.Auth;
using Xunit;

namespace DataSync.Tests
{
    public class SyncEngineTests
    {
        [Fact]
        public async Task ProcessAndSyncFile_ShouldComputeChecksumAndSucceed()
        {
            var engine = new DataSyncEngine();
            byte[] fileBytes = Encoding.UTF8.GetBytes("Header,Metric,Value\n101,Temperature,98.6");
            using var stream = new MemoryStream(fileBytes);

            var result = await engine.ProcessAndSyncFileAsync(stream, "metrics.csv", "tenant_alpha");

            Assert.NotNull(result);
            Assert.Equal("metrics.csv", result.FileName);
            Assert.Equal(ProcessingStatus.Synchronized, result.Status);
            Assert.False(string.IsNullOrEmpty(result.ChecksumSHA256));
        }

        [Fact]
        public async Task BatchSyncFiles_ShouldProcessHighThroughput()
        {
            var engine = new DataSyncEngine();
            int count = 1000;

            int processed = await engine.BatchSyncFilesAsync(count);

            Assert.Equal(count, processed);
        }

        [Fact]
        public void EntraIdTokenValidator_ShouldRejectUnauthorizedTenant()
        {
            var validator = new EntraIdTokenValidator();
            var claims = new[]
            {
                new Claim("tid", "tenant-unauthorized"),
                new Claim(ClaimTypes.Role, "DataSync.Operator")
            };
            var identity = new ClaimsIdentity(claims, "TestAuth");
            var principal = new ClaimsPrincipal(identity);

            bool authorized = validator.ValidateCustomClaims(principal, "DataSync.Operator", "tenant-authorized");

            Assert.False(authorized);
        }
    }
}

