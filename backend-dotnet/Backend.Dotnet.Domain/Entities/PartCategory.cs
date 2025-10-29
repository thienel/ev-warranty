using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public class PartCategory : BaseEntity
    {
        public string CategoryName { get; private set; }
        public string Description { get; private set; }
        public Guid? ParentCategoryId { get; private set; }

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
        }

        // BEHAVIOUR
        public void UpdateDetails(string categoryName, string description)
        {
            SetCategoryName(categoryName);
            SetDescription(description);
            SetUpdatedAt();
        }

        public void ChangeParent(Guid? newParentCategoryId)
        {
            if (newParentCategoryId.HasValue && newParentCategoryId.Value == Id)
                throw new BusinessRuleViolationException("Category cannot be its own parent");

            ParentCategoryId = newParentCategoryId;
            SetUpdatedAt();
        }

        // QUERY
        public bool HasActiveParts() => _parts.Any(p => p.Status == PartStatus.Available || p.Status == PartStatus.Reserved);

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
    }
}
