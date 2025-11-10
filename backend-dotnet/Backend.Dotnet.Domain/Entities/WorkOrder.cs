using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

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
        public Guid ClaimId { get; private set; }
        public Guid AssignedTechnicianId { get; private set; }
        public WorkOrderStatus Status { get; private set; }
        public DateTime ScheduledDate { get; private set; } = DateTime.UtcNow;
        public DateTime? CompletedDate { get; private set; }
        public string Note { get; private set; }

        public WorkOrder() { }

        public WorkOrder(Guid claimId, Guid assignedTechnicianId, DateTime scheduledDate, DateTime? completedDate = null, string note = null)
        {
            ClaimId = claimId;
            AssignedTechnicianId = assignedTechnicianId;
            SetScheduledDate(scheduledDate);
            SetNote(note);
            Status = WorkOrderStatus.Pending;
        }

        public void ChangeStatus(WorkOrderStatus newStatus)
        {
            if (Status == WorkOrderStatus.Completed)
                throw new BusinessRuleViolationException("Cannot change status of completed work order");
            
            if (Status == WorkOrderStatus.InProgress && newStatus == WorkOrderStatus.Pending)
                throw new BusinessRuleViolationException("InProgress work order cannot revert to Pending");
            
            if (Status == WorkOrderStatus.Pending && newStatus == WorkOrderStatus.Completed)
                throw new BusinessRuleViolationException("Must be InProgress before Completed");

            if (newStatus == WorkOrderStatus.Completed)
                CompletedDate = DateTime.UtcNow;

            Status = newStatus;
            SetUpdatedAt();
        }

        // Behaviour Methods
        public void UpdateScheduledDate(DateTime scheduledDate)
        {
            if (Status == WorkOrderStatus.Completed)
                throw new BusinessRuleViolationException("Cannot reschedule completed work order");

            SetScheduledDate(scheduledDate);
            SetUpdatedAt();
        }

        public void UpdateNote(string note)
        {
            if (Status == WorkOrderStatus.Completed)
                throw new BusinessRuleViolationException("Cannot update note of completed work order");

            SetNote(note);
            SetUpdatedAt();
        }

        // Private Setter
        private void SetScheduledDate(DateTime scheduledDate)
        {
            if (scheduledDate.Date < DateTime.Now.Date)
                throw new BusinessRuleViolationException("Scheduled date cannot be in the past");

            ScheduledDate = scheduledDate.Date;
        }

        private void SetNote(string note)
        {
            Note = string.IsNullOrWhiteSpace(note) ? null : note.Trim();
        }
    }
}
