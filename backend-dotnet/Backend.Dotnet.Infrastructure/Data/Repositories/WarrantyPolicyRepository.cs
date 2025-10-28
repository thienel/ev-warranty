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
    public class WarrantyPolicyRepository : BaseRepository<WarrantyPolicy>, IWarrantyPolicyRepository
    {
        public WarrantyPolicyRepository(DbContext context) : base(context) { }

        public async Task<IEnumerable<WarrantyPolicy>> GetByStatusAsync(WarrantyPolicyStatus status)
        {
            return await _dbSet
                .Where(wp => wp.Status == status)
                .OrderByDescending(wp => wp.CreatedAt)
                .ToListAsync();
        }

        public async Task<WarrantyPolicy?> GetByPolicyNameAsync(string policyName)
        {
            return await _dbSet
                .FirstOrDefaultAsync(wp => wp.PolicyName.ToLower() == policyName.ToLower());
        }

        public async Task<bool> PolicyNameExistsAsync(string policyName, Guid? excludePolicyId = null)
        {
            var query = _dbSet.Where(wp => wp.PolicyName.ToLower() == policyName.ToLower());

            if (excludePolicyId.HasValue)
                query = query.Where(wp => wp.Id != excludePolicyId.Value);

            return await query.AnyAsync();
        }

        public async Task<WarrantyPolicy?> GetWithDetailsAsync(Guid policyId)
        {
            return await _dbSet
                .Include(wp => wp.VehicleModels)
                .Include(wp => wp.CoverageParts)
                    .ThenInclude(cp => cp.PartCategory)
                .FirstOrDefaultAsync(wp => wp.Id == policyId);
        }

        public async Task<IEnumerable<WarrantyPolicy>> GetAllWithDetailsAsync(WarrantyPolicyStatus? status = null)
        {
            IQueryable<WarrantyPolicy> query = _dbSet
               .Include(wp => wp.VehicleModels)
               .Include(wp => wp.CoverageParts)
               .OrderByDescending(wp => wp.CreatedAt);
            //IOrderedQueryable implement IQueryable

            if (status.HasValue)
            {
                query = query.Where(wp => wp.Status == status.Value);
            }

            return await query
               .ToListAsync();
        }

        public async Task<bool> CanBeAssignedToVehiclesAsync(Guid policyId)
        {
            var policy = await _dbSet.FindAsync(policyId);
            return policy?.CanBeAssignedToVehicles() ?? false;
        }

        public async Task<int> GetCoveragePartCountAsync(Guid policyId)
        {
            return await _context.Set<PolicyCoveragePart>()
                .CountAsync(cp => cp.PolicyId == policyId);
        }
    }
}
