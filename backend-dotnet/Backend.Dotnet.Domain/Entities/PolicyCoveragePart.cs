using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public class PolicyCoveragePart : BaseEntity
    {
        public Guid PolicyId { get; private set; }
        public Guid PartCategoryId { get; private set; }
        public string CoverageConditions { get; private set; }

        // Navigation properties
        public virtual WarrantyPolicy Policy { get; private set; }
        public virtual PartCategory PartCategory { get; private set; }

        private PolicyCoveragePart() { }

        public PolicyCoveragePart(Guid policyId, Guid partCategoryId, string coverageConditions = null)
        {
            PolicyId = policyId;
            PartCategoryId = partCategoryId;
            SetCoverageConditions(coverageConditions);
        }

        // BEHAVIOUR
        public void UpdateCoverageConditions(string coverageConditions)
        {
            // Validate for warranty policy in draft status
            if (Policy != null && !Policy.IsEditable())
                throw new BusinessRuleViolationException("Cannot update coverage for non-draft policy");

            SetCoverageConditions(coverageConditions);
            SetUpdatedAt();
        }

        // VALIDATION
        public void ValidateAgainstPolicy(WarrantyPolicy policy)
        {
            if (policy == null)
                throw new ArgumentNullException(nameof(policy));

            if (policy.Id != PolicyId)
                throw new BusinessRuleViolationException("Policy ID mismatch");

            if (!policy.IsEditable())
                throw new BusinessRuleViolationException("Cannot modify coverage for non-draft policy");
        }

        public void ValidateAgainstCategory(PartCategory category)
        {
            if (category == null)
                throw new ArgumentNullException(nameof(category));

            if (category.Id != PartCategoryId)
                throw new BusinessRuleViolationException("Part category ID mismatch");

            if (!category.CanBeUsedForNewParts())
                throw new BusinessRuleViolationException("Cannot add inactive category to policy coverage");
        }

        // PRIVATE SETTERS
        private void SetCoverageConditions(string conditions)
        {
            if (!string.IsNullOrWhiteSpace(conditions) && conditions.Length > 1000)
                throw new BusinessRuleViolationException("Coverage conditions cannot exceed 1000 characters");

            CoverageConditions = conditions?.Trim();
        }
    }
}
