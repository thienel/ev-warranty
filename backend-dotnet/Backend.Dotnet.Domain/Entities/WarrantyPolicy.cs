using Backend.Dotnet.Domain.Abstractions;
using Backend.Dotnet.Domain.Exceptions;

namespace Backend.Dotnet.Domain.Entities
{
    public enum WarrantyPolicyStatus
    {
        Draft,      // Can be edited freely
        Active,     // In use, cannot be edited
        Expired,    // Past validity period
        Superseded, // Replaced by newer policy
        Archived    // Removed from active use
    }

    public class WarrantyPolicy : BaseEntity, IStatus<WarrantyPolicyStatus>
    {
        public string PolicyName { get; private set; }
        public int WarrantyDurationMonths { get; private set; }
        public int? KilometerLimit { get; private set; }
        public string TermsAndConditions { get; private set; }
        public WarrantyPolicyStatus Status { get; private set; }

        // Navigation properties
        public virtual VehicleModel? AssignedModel { get; private set; }

        private readonly List<PolicyCoveragePart> _coverageParts = new();
        public virtual IReadOnlyCollection<PolicyCoveragePart> CoverageParts => _coverageParts.AsReadOnly();

        private WarrantyPolicy() { }

        public WarrantyPolicy(string policyName, int warrantyDurationMonths, int? kilometerLimit, string termsAndConditions)
        {
            SetPolicyName(policyName);
            SetWarrantyDuration(warrantyDurationMonths);
            SetKilometerLimit(kilometerLimit);
            SetTermsAndConditions(termsAndConditions);
            Status = WarrantyPolicyStatus.Draft;
        }

        // BEHAVIOUR
        public void UpdateDetails(string policyName, int warrantyDurationMonths, int? kilometerLimit, string termsAndConditions)
        {
            if (Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Only draft policies can be edited");

            SetPolicyName(policyName);
            SetWarrantyDuration(warrantyDurationMonths);
            SetKilometerLimit(kilometerLimit);
            SetTermsAndConditions(termsAndConditions);
            SetUpdatedAt();
        }

        // Status
        public void ChangeStatus(WarrantyPolicyStatus newStatus)
        {
            if (Status == newStatus)
                return;

            ValidateStatusTransition(Status, newStatus);
            Status = newStatus;
            SetUpdatedAt();
        }

        public void Activate()
        {
            if (_coverageParts.Count == 0)
                throw new BusinessRuleViolationException("Cannot activate policy without coverage parts");

            if (Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Only draft policies can be activated");

            Status = WarrantyPolicyStatus.Active;
            SetUpdatedAt();
        }

        public void Expire()
        {
            if (Status != WarrantyPolicyStatus.Active)
                throw new BusinessRuleViolationException("Only active policies can be expired");

            Status = WarrantyPolicyStatus.Expired;
            SetUpdatedAt();
        }

        public void Supersede()
        {
            if (Status != WarrantyPolicyStatus.Active)
                throw new BusinessRuleViolationException("Only active policies can be superseded");

            Status = WarrantyPolicyStatus.Superseded;
            SetUpdatedAt();
        }

        public void Archive()
        {
            if (Status == WarrantyPolicyStatus.Active)
                throw new BusinessRuleViolationException("Cannot archive active policy");

            Status = WarrantyPolicyStatus.Archived;
            SetUpdatedAt();
        }

        // Vehicle Model Assignment
        public void AssignToVehicleModel(Guid vehicleModelId)
        {
            if (vehicleModelId == Guid.Empty)
                throw new BusinessRuleViolationException("Vehicle model ID cannot be empty");

            if (AssignedModel != null)
                throw new BusinessRuleViolationException("This policy is already assigned to a vehicle model");

            // Note: The actual assignment will be handled by the repository/service layer
            SetUpdatedAt();
        }

        public void RemoveVehicleModelAssignment()
        {
            if (AssignedModel == null)
                throw new BusinessRuleViolationException("This policy is not assigned to any vehicle model");

            SetUpdatedAt();
        }

        // Coverage Parts
        public void AddCoveragePart(PolicyCoveragePart coveragePart)
        {
            if (Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Can only add coverage parts to draft policies");

            if (coveragePart == null)
                throw new ArgumentNullException(nameof(coveragePart));

            if (_coverageParts.Any(cp => cp.PartCategoryId == coveragePart.PartCategoryId))
                throw new BusinessRuleViolationException("Part category already coverage by this policy");

            _coverageParts.Add(coveragePart);
            SetUpdatedAt();
        }

        public void RemoveCoveragePart(Guid partCategoryId)
        {
            if (Status != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Can only remove coverage parts from draft policies");

            var coveragePart = _coverageParts.FirstOrDefault(cp => cp.PartCategoryId == partCategoryId);
            if (coveragePart == null)
                throw new BusinessRuleViolationException("Part category not found in policy coverage");

            _coverageParts.Remove(coveragePart);
            SetUpdatedAt();
        }

        // QUERY
        public bool CanBeAssignedToVehicles() => Status == WarrantyPolicyStatus.Active;
        public bool IsEditable() => Status == WarrantyPolicyStatus.Draft;
        public bool IsPartCategoryCoverage(Guid partCategoryId) =>
            _coverageParts.Any(cp => cp.PartCategoryId == partCategoryId);

        // PRIVATE SETTERS
        private void SetPolicyName(string policyName)
        {
            if (string.IsNullOrWhiteSpace(policyName))
                throw new BusinessRuleViolationException("Policy name cannot be empty");

            if (policyName.Length > 255)
                throw new BusinessRuleViolationException("Policy name cannot exceed 255 characters");

            PolicyName = policyName.Trim();
        }

        private void SetWarrantyDuration(int months)
        {
            if (months <= 0)
                throw new BusinessRuleViolationException("Warranty duration must be greater than zero");

            if (months > 600)
                throw new BusinessRuleViolationException("Warranty duration cannot exceed 600 months");

            WarrantyDurationMonths = months;
        }

        private void SetKilometerLimit(int? limit)
        {
            if (limit.HasValue && limit.Value <= 0)
                throw new BusinessRuleViolationException("Kilometer limit must be greater than zero");

            KilometerLimit = limit;
        }

        private void SetTermsAndConditions(string terms)
        {
            if (string.IsNullOrWhiteSpace(terms))
                throw new BusinessRuleViolationException("Terms and conditions cannot be empty");

            if (terms.Length > 5000)
                throw new BusinessRuleViolationException("Terms and conditions cannot exceed 5000 characters");

            TermsAndConditions = terms.Trim();
        }

        private void ValidateStatusTransition(WarrantyPolicyStatus from, WarrantyPolicyStatus to)
        {
            if (from == WarrantyPolicyStatus.Archived)
                throw new BusinessRuleViolationException("Cannot change status of archived policy");

            if (to == WarrantyPolicyStatus.Draft && from != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Cannot revert to draft status");

            if (to == WarrantyPolicyStatus.Active && from != WarrantyPolicyStatus.Draft)
                throw new BusinessRuleViolationException("Only draft policies can be activated");
        }
    }
}
