using System;
using System.Security.Claims;

namespace DataSync.Security.Auth
{
    public interface ITokenValidator
    {
        bool ValidateCustomClaims(ClaimsPrincipal principal, string requiredRole, string tenantId);
    }

    public class EntraIdTokenValidator : ITokenValidator
    {
        public bool ValidateCustomClaims(ClaimsPrincipal principal, string requiredRole, string tenantId)
        {
            if (principal == null) return false;

            // Enforce tenant boundary isolation
            var tidClaim = principal.FindFirst("tid")?.Value;
            if (string.IsNullOrEmpty(tidClaim) || !string.Equals(tidClaim, tenantId, StringComparison.OrdinalIgnoreCase))
            {
                return false;
            }

            // Role-Based Access Control (RBAC) verification
            bool hasRole = principal.HasClaim(ClaimTypes.Role, requiredRole);
            return hasRole;
        }
    }
}
