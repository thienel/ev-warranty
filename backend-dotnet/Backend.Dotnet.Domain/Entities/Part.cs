using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public enum PartStatus
    {
        Available,  // In stock, ready for use
        Reserved,   // Allocated for work order
        Installed,  // Used in vehicle
        Defective,  // Damaged, not usable
        Obsolete,   // No longer in use
        Archived    // Removed from inventory
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
        public virtual PartCategory Category { get; private set; }

        private Part() { }

        public Part(
            string serialNumber,
            string partName,
            decimal unitPrice,
            Guid categoryId,
            Guid? officeLocationId = null)
        {
            SetSerialNumber(serialNumber);
            SetPartName(partName);
            SetUnitPrice(unitPrice);
            CategoryId = categoryId;
            OfficeLocationId = officeLocationId;
            Status = PartStatus.Available;
        }

        // BEHAVIOUR
        public void UpdateDetails(string partName, decimal unitPrice, Guid? officeLocationId)
        {
            if (Status == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot update archived part");

            SetPartName(partName);
            SetUnitPrice(unitPrice);
            OfficeLocationId = officeLocationId;
            SetUpdatedAt();
        }

        public void UpdateCategory(Guid newCategoryId)
        {
            if (Status == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change category of archived part");

            if (Status == PartStatus.Installed)
                throw new BusinessRuleViolationException("Cannot change category of installed part");

            CategoryId = newCategoryId;
            SetUpdatedAt();
        }

        public void UpdateOfficeLocation(Guid? officeLocationId)
        {
            if (Status == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot update archived part");

            OfficeLocationId = officeLocationId;
            SetUpdatedAt();
        }

        // Status
        public void ChangeStatus(PartStatus newStatus)
        {
            if (Status == newStatus)
                return;
            
            ValidateStatusTransition(Status, newStatus);
            Status = newStatus;
            SetUpdatedAt();
        }

        public void Reserve()
        {
            if (Status != PartStatus.Available)
                throw new BusinessRuleViolationException($"Cannot reserve part with status {Status}");

            Status = PartStatus.Reserved;
            SetUpdatedAt();
        }

        public void MarkAsInstalled()
        {
            if (Status != PartStatus.Reserved)
                throw new BusinessRuleViolationException("Only reserved parts can be marked as installed");

            Status = PartStatus.Installed;
            SetUpdatedAt();
        }

        public void MarkAsDefective()
        {
            if (Status == PartStatus.Installed || Status == PartStatus.Archived)
                throw new BusinessRuleViolationException($"Cannot mark {Status} part as defective");

            Status = PartStatus.Defective;
            SetUpdatedAt();
        }

        public void MakeObsolete()
        {
            if (Status == PartStatus.Reserved || Status == PartStatus.Installed)
                throw new BusinessRuleViolationException($"Cannot make {Status} part obsolete");

            Status = PartStatus.Obsolete;
            SetUpdatedAt();
        }

        public void Archive()
        {
            if (Status == PartStatus.Reserved || Status == PartStatus.Installed)
                throw new BusinessRuleViolationException($"Cannot archive {Status} part");

            Status = PartStatus.Archived;
            SetUpdatedAt();
        }

        public void MakeAvailable()
        {
            if (Status == PartStatus.Installed || Status == PartStatus.Archived)
                throw new BusinessRuleViolationException($"Cannot make {Status} part available");

            Status = PartStatus.Available;
            SetUpdatedAt();
        }

        // QUERY
        public bool CanBeUsedInWorkOrder() => Status == PartStatus.Available;
        public bool IsInStock() => Status == PartStatus.Available || Status == PartStatus.Reserved;

        // PRIVATE SETTERS
        private void SetSerialNumber(string serialNumber)
        {
            if (string.IsNullOrWhiteSpace(serialNumber))
                throw new BusinessRuleViolationException("Serial number cannot be empty");

            if (serialNumber.Length > 255)
                throw new BusinessRuleViolationException("Serial number cannot exceed 255 characters");

            SerialNumber = serialNumber.Trim().ToUpperInvariant();
        }

        private void SetPartName(string partName)
        {
            if (string.IsNullOrWhiteSpace(partName))
                throw new BusinessRuleViolationException("Part name cannot be empty");

            if (partName.Length > 255)
                throw new BusinessRuleViolationException("Part name cannot exceed 255 characters");

            PartName = partName.Trim();
        }

        private void SetUnitPrice(decimal unitPrice)
        {
            if (unitPrice <= 0)
                throw new BusinessRuleViolationException("Unit price must be greater than zero");

            if (unitPrice > 999999999.99m)
                throw new BusinessRuleViolationException("Unit price cannot exceed 999,999,999.99");

            UnitPrice = unitPrice;
        }

        private void ValidateStatusTransition(PartStatus from, PartStatus to)
        {
            if (from == PartStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change status of archived part");
        }
    }
}
