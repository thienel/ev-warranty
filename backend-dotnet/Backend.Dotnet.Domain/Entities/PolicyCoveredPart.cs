using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Domain.Entities
{
    public class PolicyCoveredPart : BaseEntity
    {
        public Guid PolicyId { get; private set; }
        public Guid PartCategoryId { get; private set; }
        public string CoverageConditions { get; private set; }

        // Navigation properties
        public WarrantyPolicy Policy { get; private set; }
        public PartCategory PartCategory { get; private set; }

        private PolicyCoveredPart() { } // EF Core

        public PolicyCoveredPart(Guid policyId, Guid partCategoryId, string coverageConditions = null)
        {
            PolicyId = policyId;
            PartCategoryId = partCategoryId;
            CoverageConditions = coverageConditions;
        }

        public void UpdateCoverageConditions(string coverageConditions)
        {
            if (Policy?.Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Cannot update coverage conditions for a non-draft policy");

            CoverageConditions = coverageConditions;
            SetUpdatedAt();
        }

        public void ValidateAgainstPolicy(WarrantyPolicy policy)
        {
            if (policy == null)
                throw new BusinessRuleViolationException("Policy cannot be null");

            if (!policy.IsEditable())
                throw new BusinessRuleViolationException("Cannot modify covered parts for a non-draft policy");
        }

        public void ValidateAgainstCategory(PartCategory category)
        {
            if (category == null)
                throw new BusinessRuleViolationException("Part category cannot be null");

            if (!category.CanBeUsedForNewParts())
                throw new BusinessRuleViolationException("Cannot add an inactive category to policy coverage");
        }
    }
}
