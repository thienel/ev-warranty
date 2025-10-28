using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public enum PartCategoryStatus
    {
        Active,     // Can be used for new parts
        ReadOnly,   // Cannot add new parts, existing parts remain
        Archived    // Completely inactive
    }

    public class PartCategory : BaseEntity, IStatus<PartCategoryStatus>
    {
        public string CategoryName { get; private set; }
        public string Description { get; private set; }
        public Guid? ParentCategoryId { get; private set; }
        public PartCategoryStatus Status { get; private set; }

        // Navigation properties
        public virtual PartCategory ParentCategory { get; private set; }

        private readonly List<PartCategory> _childCategories = new();
        public virtual IReadOnlyCollection<PartCategory> ChildCategories => _childCategories.AsReadOnly();

        private readonly List<Part> _parts = new();
        public virtual IReadOnlyCollection<Part> Parts => _parts.AsReadOnly();

        private readonly List<PolicyCoveragePart> _policyCoverageParts = new();
        public virtual IReadOnlyCollection<PolicyCoveragePart> PolicyCoverageParts => _policyCoverageParts.AsReadOnly();

        private PartCategory() { }

        public PartCategory(string categoryName, string description = null, Guid? parentCategoryId = null)
        {
            SetCategoryName(categoryName);
            SetDescription(description);
            ParentCategoryId = parentCategoryId;
            Status = PartCategoryStatus.Active;
        }

        // BEHAVIOUR
        public void UpdateDetails(string categoryName, string description)
        {
            if (Status == PartCategoryStatus.Archived)
                throw new BusinessRuleViolationException("Cannot update archived category");

            SetCategoryName(categoryName);
            SetDescription(description);
            SetUpdatedAt();
        }

        public void ChangeParent(Guid? newParentCategoryId)
        {
            if (Status == PartCategoryStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change parent of archived category");

            if (newParentCategoryId.HasValue && newParentCategoryId.Value == Id)
                throw new BusinessRuleViolationException("Category cannot be its own parent");

            ParentCategoryId = newParentCategoryId;
            SetUpdatedAt();
        }

        // Status
        public void ChangeStatus(PartCategoryStatus newStatus)
        {
            if (Status == newStatus)
                return;

            ValidateStatusTransition(Status, newStatus);
            Status = newStatus;
            SetUpdatedAt();
        }

        public void MakeReadOnly()
        {
            if (Status == PartCategoryStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change status of archived category");

            Status = PartCategoryStatus.ReadOnly;
            SetUpdatedAt();
        }

        public void Archive()
        {
            if (_childCategories.Any(c => c.Status != PartCategoryStatus.Archived))
                throw new BusinessRuleViolationException("Cannot archive category with active child categories");

            Status = PartCategoryStatus.Archived;
            SetUpdatedAt();
        }

        public void Activate()
        {
            if (ParentCategoryId.HasValue && ParentCategory?.Status != PartCategoryStatus.Active)
                throw new BusinessRuleViolationException("Cannot activate category with non-active parent");

            Status = PartCategoryStatus.Active;
            SetUpdatedAt();
        }

        // QUERY
        public bool CanBeUsedForNewParts() => Status == PartCategoryStatus.Active;
        public bool HasActiveParts() => _parts.Any(p => p.Status == PartStatus.Available || p.Status == PartStatus.Reserved);
        public bool HasActiveChildren() => _childCategories.Any(c => c.Status == PartCategoryStatus.Active);

        // PRIVATE SETTERS
        private void SetCategoryName(string categoryName)
        {
            if (string.IsNullOrWhiteSpace(categoryName))
                throw new BusinessRuleViolationException("Category name cannot be empty");

            if (categoryName.Length > 255)
                throw new BusinessRuleViolationException("Category name cannot exceed 255 characters");

            CategoryName = categoryName.Trim();
        }

        private void SetDescription(string description)
        {
            if (!string.IsNullOrWhiteSpace(description) && description.Length > 1000)
                throw new BusinessRuleViolationException("Description cannot exceed 1000 characters");

            Description = description?.Trim();
        }

        private void ValidateStatusTransition(PartCategoryStatus from, PartCategoryStatus to)
        {
            if (to == PartCategoryStatus.Archived && _childCategories.Any(c => c.Status != PartCategoryStatus.Archived))
                throw new BusinessRuleViolationException("Cannot archive category with active child categories");
        }
    }
}
