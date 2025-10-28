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
        Task<int> GetCountByPolicyIdAsync(Guid policyId);

        Task<IEnumerable<PolicyCoveragePart>> GetByPartCategoryIdAsync(Guid partCategoryId);
        Task<int> GetCountByPartCategoryIdAsync(Guid partCategoryId);

        Task<PolicyCoveragePart?> GetByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId);
        Task<bool> ExistsByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId, Guid? excludeId = null);

        Task<PolicyCoveragePart?> GetWithDetailsAsync(Guid id);
        Task<PolicyCoveragePart?> GetWithPolicyAsync(Guid id);
        Task<PolicyCoveragePart?> GetWithCategoryAsync(Guid id);
        Task<IEnumerable<PolicyCoveragePart>> GetAllWithDetailsAsync();

        Task<bool> IsCategoryCoveredByPolicyAsync(Guid policyId, Guid partCategoryId);
        Task<IEnumerable<Guid>> GetCoveredCategoryIdsByPolicyAsync(Guid policyId);
        Task<IEnumerable<Guid>> GetPolicyIdsCoveringCategoryAsync(Guid partCategoryId);

        Task<IEnumerable<PolicyCoveragePart>> GetByCategoryIdsAsync(Guid policyId, IEnumerable<Guid> categoryIds);
        Task RemoveByPolicyIdAsync(Guid policyId);
    }
}
