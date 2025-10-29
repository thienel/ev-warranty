using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class PolicyCoveragePartRepository : BaseRepository<PolicyCoveragePart>, IPolicyCoveragePartRepository
    {
        public PolicyCoveragePartRepository(DbContext context) : base(context) { }

        public async Task<IEnumerable<PolicyCoveragePart>> GetByPolicyIdAsync(Guid policyId)
        {
            return await _dbSet
                .Include(cp => cp.PartCategory)
                .Where(cp => cp.PolicyId == policyId)
                .OrderBy(cp => cp.PartCategory.CategoryName)
                .ToListAsync();
        }

        public async Task<PolicyCoveragePart?> GetByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId)
        {
            return await _dbSet
                .Include(cp => cp.PartCategory)
                .FirstOrDefaultAsync(cp =>
                    cp.PolicyId == policyId &&
                    cp.PartCategoryId == partCategoryId);
        }

        public async Task<bool> ExistsByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId, Guid? excludeId = null)
        {
            var query = _dbSet.Where(cp =>
                cp.PolicyId == policyId &&
                cp.PartCategoryId == partCategoryId);

            if (excludeId.HasValue)
                query = query.Where(cp => cp.Id != excludeId.Value);

            return await query.AnyAsync();
        }

        public async Task<PolicyCoveragePart?> GetWithDetailsAsync(Guid id)
        {
            return await _dbSet
                .Include(cp => cp.Policy)
                .Include(cp => cp.PartCategory)
                .FirstOrDefaultAsync(cp => cp.Id == id);
        }

        public async Task<IEnumerable<PolicyCoveragePart>> GetAllWithDetailsAsync()
        {
            return await _dbSet
                .Include(cp => cp.Policy)
                .Include(cp => cp.PartCategory)
                .OrderByDescending(cp => cp.CreatedAt)
                .ToListAsync();
        }

        public async Task RemoveByPolicyIdAsync(Guid policyId)
        {
            var coverageParts = await _dbSet
                .Where(cp => cp.PolicyId == policyId)
                .ToListAsync();

            _dbSet.RemoveRange(coverageParts);
        }

        public async Task<bool> IsCategoryCoveredByPolicyAsync(Guid policyId, Guid partCategoryId)
        {
            return await _dbSet.AnyAsync(cp =>
                cp.PolicyId == policyId &&
                cp.PartCategoryId == partCategoryId);
        }

        public async Task<int> GetCountByPolicyIdAsync(Guid policyId)
        {
            return await _dbSet.CountAsync(cp => cp.PolicyId == policyId);
        }
    }
}
