using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IWarrantyPolicyRepository : IRepository<WarrantyPolicy>
    {
        Task<IEnumerable<WarrantyPolicy>> GetByModelIdAsync(Guid modelId, WarrantyPolicyStatus? status = null);

        // Status   
        Task<IEnumerable<WarrantyPolicy>> GetByStatusAsync(WarrantyPolicyStatus status);

        Task<WarrantyPolicy?> GetByPolicyNameAsync(string policyName);
        Task<bool> PolicyNameExistsAsync(string policyName, Guid? excludePolicyId = null);

        Task<WarrantyPolicy?> GetWithDetailsAsync(Guid policyId);
        Task<IEnumerable<WarrantyPolicy>> GetAllWithDetailsAsync();

        Task<bool> CanBeAssignedToVehiclesAsync(Guid policyId);
        Task<bool> HasCoveragePartsAsync(Guid policyId);
        Task<int> GetCoveragePartCountAsync(Guid policyId);
        Task<bool> IsPartCategoryCoveredAsync(Guid policyId, Guid partCategoryId);
    }
}
