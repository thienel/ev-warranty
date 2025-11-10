using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IWorkOrderRepository : IRepository<WorkOrder>
    {
        Task<WorkOrder?> GetByClaimIdAsync(Guid claimId);
        Task <IEnumerable<WorkOrder>> GetByTechnicianIdAsync(Guid technicianId);
        Task <IEnumerable<WorkOrder>> GetByStatusAsync(WorkOrderStatus status);
        Task <IEnumerable<WorkOrder>> GetByTechnicianIdAndStatusAsync(Guid technicianId, WorkOrderStatus status);
        Task <bool> ClaimHasWorkOrderAsync(Guid claimId);
    }
}
