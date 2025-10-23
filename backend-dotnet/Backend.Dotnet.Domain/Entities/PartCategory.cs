using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public enum PartCategoryStatus
    {
        Active,
        ReadOnly,   //cannot add any childPart further
        Archived
    }

    public class PartCategory : BaseEntity, IStatus<PartCategoryStatus>
    {
        public string CategoryName { get; private set; }
        public string Description { get; private set; }
        public Guid? ParentCategoryId { get; private set; }
        public PartCategoryStatus Status { get; private set; }

        // Navigation properties
        public PartCategory ParentCategory { get; private set; }
        private readonly List<PartCategory> _childCategories = new();
        public IReadOnlyCollection<PartCategory> ChildCategories => _childCategories.AsReadOnly();

        private readonly List<Part> _parts = new();
        public IReadOnlyCollection<Part> Parts => _parts.AsReadOnly();

        private readonly List<PolicyCoveredPart> _policyCoveredParts = new();
        public IReadOnlyCollection<PolicyCoveredPart> PolicyCoveredParts => _policyCoveredParts.AsReadOnly();

        private PartCategory() { } 

        public PartCategory(string categoryName, string description, Guid? parentCategoryId = null)
        {
            if (string.IsNullOrWhiteSpace(categoryName))
                throw new BusinessRuleViolationException("Category name cannot be empty");

            CategoryName = categoryName;
            Description = description;
            ParentCategoryId = parentCategoryId;
            Status = PartCategoryStatus.Active;
        }

        public void UpdateDetails(string categoryName, string description)
        {
            if (Status == PartCategoryStatus.Archived)
                throw new BusinessRuleViolationException("Cannot update an archived category");

            if (string.IsNullOrWhiteSpace(categoryName))
                throw new BusinessRuleViolationException("Category name cannot be empty");

            CategoryName = categoryName;
            Description = description;
            SetUpdatedAt();
        }

        public void ChangeParent(Guid? newParentCategoryId)
        {
            if (Status == PartCategoryStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change parent of an archived category");

            if (newParentCategoryId == Id)
                throw new BusinessRuleViolationException("A category cannot be its own parent");

            ParentCategoryId = newParentCategoryId;
            SetUpdatedAt();
        }

        public void ChangeStatus(PartCategoryStatus newStatus)
        {
            if (Status == newStatus)
                return;

            // Business rules for status transitions
            if (newStatus == PartCategoryStatus.Archived && _childCategories.Exists(c => c.Status != PartCategoryStatus.Archived))
                throw new BusinessRuleViolationException("Cannot archive a category with active child categories");

            Status = newStatus;
            SetUpdatedAt();
        }

        public void MakeReadOnly()
        {
            ChangeStatus(PartCategoryStatus.ReadOnly);
        }

        public void Archive()
        {
            ChangeStatus(PartCategoryStatus.Archived);
        }

        public void Activate()
        {
            if (ParentCategoryId.HasValue && ParentCategory?.Status != PartCategoryStatus.Active)
                throw new BusinessRuleViolationException("Cannot activate a category whose parent is not active");

            ChangeStatus(PartCategoryStatus.Active);
        }

        public bool CanBeUsedForNewParts()
        {
            return Status == PartCategoryStatus.Active;
        }

    }
}
