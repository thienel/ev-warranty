using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IPolicyCoveragePartRepository : IRepository<PolicyCoveragePart>
    {
        Task<IEnumerable<PolicyCoveragePart>> GetByPolicyIdAsync(Guid policyId);
        Task<PolicyCoveragePart?> GetByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId);
        Task<bool> ExistsByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId, Guid? excludeId = null);
        
        Task<PolicyCoveragePart?> GetWithDetailsAsync(Guid id);
        Task<IEnumerable<PolicyCoveragePart>> GetAllWithDetailsAsync();

        Task RemoveByPolicyIdAsync(Guid policyId);

        Task<bool> IsCategoryCoveredByPolicyAsync(Guid policyId, Guid partCategoryId);
        Task<int> GetCountByPolicyIdAsync(Guid policyId);
    }
}
