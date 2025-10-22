using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public enum PartStatus
    {
        Available,
        Reserved,
        Installed,
        Defective,
        Obsolete,
        Archived
    }

    public class Part : BaseEntity, IStatus<PartStatus>
    {
        public string SerialNumber { get; private set; }
        public string PartName { get; private set; }
        public decimal UnitPrice { get; private set; }
        public Guid CategoryId { get; private set; }
        public Guid? OfficeLocationId { get; private set; }
        public PartStatus Status { get; private set; }

        // Navigation properties
        public PartCategory Category { get; private set; }

        private Part() { } // EF Core

        public Part(string serialNumber, string partName, decimal unitPrice, Guid categoryId, Guid? officeLocationId = null)
        {
            if (string.IsNullOrWhiteSpace(serialNumber))
                throw new BusinessRuleViolationException("Serial number cannot be empty");

            if (string.IsNullOrWhiteSpace(partName))
                throw new BusinessRuleViolationException("Part name cannot be empty");

            if (unitPrice <= 0)
                throw new BusinessRuleViolationException("Unit price must be greater than zero");

            SerialNumber = serialNumber;
            PartName = partName;
            UnitPrice = unitPrice;
            CategoryId = categoryId;
            OfficeLocationId = officeLocationId;
            Status = PartStatus.Available;
        }

        public void UpdateDetails(string partName, decimal unitPrice, Guid? officeLocationId)
        {
            if (Status == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot update an archived part");

            if (string.IsNullOrWhiteSpace(partName))
                throw new BusinessRuleViolationException("Part name cannot be empty");

            if (unitPrice <= 0)
                throw new BusinessRuleViolationException("Unit price must be greater than zero");

            PartName = partName;
            UnitPrice = unitPrice;
            OfficeLocationId = officeLocationId;
            SetUpdatedAt();
        }

        public void UpdateCategory(Guid newCategoryId, PartCategory newCategory)
        {
            if (Status == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change category of an archived part");

            if (newCategory == null || !newCategory.CanBeUsedForNewParts())
                throw new BusinessRuleViolationException("Cannot assign part to an inactive category");

            CategoryId = newCategoryId;
            SetUpdatedAt();
        }

        public void ChangeStatus(PartStatus newStatus)
        {
            if (Status == newStatus)
                return;

            // Business rules for status transitions
            if (Status == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change status of an archived part");

            if (newStatus == PartStatus.Available && Status == PartStatus.Reserved)
            {
                // Could add additional checks here (e.g., ensure no active work order)
            }

            Status = newStatus;
            SetUpdatedAt();
        }

        public void Reserve()
        {
            if (Status != PartStatus.Available)
                throw new BusinessRuleViolationException($"Cannot reserve a part with status {Status}");

            ChangeStatus(PartStatus.Reserved);
        }

        public void MarkAsInstalled()
        {
            if (Status != PartStatus.Reserved)
                throw new BusinessRuleViolationException("Only reserved parts can be marked as installed");

            ChangeStatus(PartStatus.Installed);
        }

        public void MarkAsDefective()
        {
            if (Status == PartStatus.Installed || Status == PartStatus.Archived)
                throw new BusinessRuleViolationException($"Cannot mark a part with status {Status} as defective");

            ChangeStatus(PartStatus.Defective);
        }

        public void MakeObsolete()
        {
            if (Status == PartStatus.Reserved || Status == PartStatus.Installed)
                throw new BusinessRuleViolationException($"Cannot make a part with status {Status} obsolete");

            ChangeStatus(PartStatus.Obsolete);
        }

        public void Archive()
        {
            if (Status == PartStatus.Reserved || Status == PartStatus.Installed)
                throw new BusinessRuleViolationException("Cannot archive a part that is reserved or installed");

            ChangeStatus(PartStatus.Archived);
        }

        public void MakeAvailable()
        {
            if (Status == PartStatus.Installed || Status == PartStatus.Archived)
                throw new BusinessRuleViolationException($"Cannot make a part with status {Status} available");

            ChangeStatus(PartStatus.Available);
        }

        public bool CanBeUsedInWorkOrder()
        {
            return Status == PartStatus.Available;
        }
    }
}
