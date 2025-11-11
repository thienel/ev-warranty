using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Infrastructure.Data.Context;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class WorkOrderRepository : BaseRepository<WorkOrder>, IWorkOrderRepository
    {
        public WorkOrderRepository(DbContext context) : base(context) { }

        public async Task<bool> ClaimHasWorkOrderAsync(Guid claimId)
        {
            return await _dbSet
                .AnyAsync(w => w.ClaimId == claimId);
        }

        public async Task<WorkOrder?> GetByClaimIdAsync(Guid claimId)
        {
            return await _dbSet
                .FirstOrDefaultAsync(w => w.ClaimId == claimId);
        }

        public async Task<IEnumerable<WorkOrder>> GetByTechnicianIdAsync(Guid technicianId)
        {
            return await _dbSet
               .Where(w => w.AssignedTechnicianId == technicianId)
               .OrderByDescending(w => w.Status)
               .ThenBy(w => w.ScheduledDate)
               .ToListAsync();
        }

        public async Task<IEnumerable<WorkOrder>> GetByStatusAsync(WorkOrderStatus status)
        {
            return await _dbSet
                .Where(w => w.Status == status)
                .OrderBy(w => w.ScheduledDate)
                .ToListAsync();
        }

        public async Task<IEnumerable<WorkOrder>> GetByTechnicianIdAndStatusAsync(Guid technicianId, WorkOrderStatus status)
        {
            return await _dbSet
                .Where(w => w.AssignedTechnicianId == technicianId && w.Status == status)
                .OrderBy(w => w.ScheduledDate)
                .ToListAsync();
        }
    }
}
