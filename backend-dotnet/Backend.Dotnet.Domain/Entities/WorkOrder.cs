using Backend.Dotnet.Domain.Abstractions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Domain.Entities
{
    public enum WorkOrderStatus
    {
        Pending,
        InProgress,
        Completed
    }
    public class WorkOrder : BaseEntity, IStatus<WorkOrderStatus>
    {
        public WorkOrderStatus Status { get; private set; }

        public void ChangeStatus(WorkOrderStatus newStatus)
        {
            throw new NotImplementedException();
        }
    }
}
