using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Domain.Entities
{
    public enum WarrantyPolicyStatus
    {
        Draft,
        Active,
        Expired,
        Superseded,
        Archived
    }

    public class WarrantyPolicy : BaseEntity, IStatus<WarrantyPolicyStatus>
    {
        public string PolicyName { get; private set; }
        public Guid ModelId { get; private set; }
        public int WarrantyDurationMonths { get; private set; }
        public int? KilometerLimit { get; private set; }
        public string TermsAndConditions { get; private set; }
        public WarrantyPolicyStatus Status { get; private set; }

        private readonly List<PolicyCoveredPart> _coveredParts = new();
        public IReadOnlyCollection<PolicyCoveredPart> CoveredParts => _coveredParts.AsReadOnly();

        private WarrantyPolicy() { } // EF Core

        public WarrantyPolicy(
            string policyName,
            Guid modelId,
            int warrantyDurationMonths,
            int? kilometerLimit,
            string termsAndConditions)
        {
            if (string.IsNullOrWhiteSpace(policyName))
                throw new BusinessRuleViolationException("Policy name cannot be empty");

            if (warrantyDurationMonths <= 0)
                throw new BusinessRuleViolationException("Warranty duration must be greater than zero");

            if (kilometerLimit.HasValue && kilometerLimit.Value <= 0)
                throw new BusinessRuleViolationException("Kilometer limit must be greater than zero");

            if (string.IsNullOrWhiteSpace(termsAndConditions))
                throw new BusinessRuleViolationException("Terms and conditions cannot be empty");

            PolicyName = policyName;
            ModelId = modelId;
            WarrantyDurationMonths = warrantyDurationMonths;
            KilometerLimit = kilometerLimit;
            TermsAndConditions = termsAndConditions;
            Status = WarrantyPolicyStatus.Draft;
        }

        public void UpdateDetails(
            string policyName,
            int warrantyDurationMonths,
            int? kilometerLimit,
            string termsAndConditions)
        {
            if (Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Only draft policies can be freely edited");

            if (string.IsNullOrWhiteSpace(policyName))
                throw new BusinessRuleViolationException("Policy name cannot be empty");

            if (warrantyDurationMonths <= 0)
                throw new BusinessRuleViolationException("Warranty duration must be greater than zero");

            if (kilometerLimit.HasValue && kilometerLimit.Value <= 0)
                throw new BusinessRuleViolationException("Kilometer limit must be greater than zero");

            if (string.IsNullOrWhiteSpace(termsAndConditions))
                throw new BusinessRuleViolationException("Terms and conditions cannot be empty");

            PolicyName = policyName;
            WarrantyDurationMonths = warrantyDurationMonths;
            KilometerLimit = kilometerLimit;
            TermsAndConditions = termsAndConditions;
            SetUpdatedAt();
        }

        public void ChangeStatus(WarrantyPolicyStatus newStatus)
        {
            if (Status == newStatus)
                return;

            // Business rules for status transitions
            if (Status == WarrantyPolicyStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change status of an archived policy");

            if (newStatus == WarrantyPolicyStatus.Active && Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Only draft policies can be activated");

            if (newStatus == WarrantyPolicyStatus.Draft && Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Cannot revert to draft status");

            Status = newStatus;
            SetUpdatedAt();
        }

        public void Activate()
        {
            if (_coveredParts.Count == 0)
                throw new BusinessRuleViolationException("Cannot activate a policy with no covered parts");

            ChangeStatus(WarrantyPolicyStatus.Active);
        }

        public void Expire()
        {
            if (Status != WarrantyPolicyStatus.Active)
                throw new BusinessRuleViolationException("Only active policies can be expired");

            ChangeStatus(WarrantyPolicyStatus.Expired);
        }

        public void Supersede()
        {
            if (Status != WarrantyPolicyStatus.Active)
                throw new BusinessRuleViolationException("Only active policies can be superseded");

            ChangeStatus(WarrantyPolicyStatus.Superseded);
        }

        public void Archive()
        {
            if (Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Only draft policies can be archived");

            ChangeStatus(WarrantyPolicyStatus.Archived);
        }

        public bool CanBeAssignedToVehicles()
        {
            return Status == WarrantyPolicyStatus.Active;
        }

        public bool IsEditable()
        {
            return Status == WarrantyPolicyStatus.Draft;
        }
    }
}
