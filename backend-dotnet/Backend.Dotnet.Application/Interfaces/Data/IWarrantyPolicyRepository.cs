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
        // Status   
        Task<IEnumerable<WarrantyPolicy>> GetByStatusAsync(WarrantyPolicyStatus status);

        Task<WarrantyPolicy?> GetByPolicyNameAsync(string policyName);
        Task<bool> PolicyNameExistsAsync(string policyName, Guid? excludePolicyId = null);

        Task<WarrantyPolicy?> GetWithDetailsAsync(Guid policyId);
        Task<IEnumerable<WarrantyPolicy>> GetAllWithDetailsAsync(WarrantyPolicyStatus? status = null);

        Task<VehicleModel?> GetAssignedModelAsync(Guid policyId);


        Task<bool> IsPolicyAssignedAsync(Guid policyId, Guid? excludeModelId = null);
        Task<bool> CanBeAssignedToVehiclesAsync(Guid policyId);
        Task<int> GetCoveragePartCountAsync(Guid policyId);
    }
}
