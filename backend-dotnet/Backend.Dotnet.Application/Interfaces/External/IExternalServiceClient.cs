using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.IntermediaryDto;

namespace Backend.Dotnet.Application.Interfaces.External
{
    public interface IExternalServiceClient
    {
        // Warranty Claim Service
        Task<ClaimInfo?> GetClaimAsync(Guid claimId);
        Task<List<ClaimItem>?> GetClaimItemsAsync(Guid claimId);
        Task<TechnicianInfo?> GetTechnicianAsync(Guid technicianId);
        Task<bool> CompleteClaimAsync(Guid claimId);
    }
}
