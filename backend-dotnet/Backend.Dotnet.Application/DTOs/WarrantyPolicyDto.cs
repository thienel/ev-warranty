using Backend.Dotnet.Domain.Entities;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Backend.Dotnet.Application.DTOs
{
    public class WarrantyPolicyDto
    {
        public class CreateWarrantyPolicyRequest
        {
            [JsonPropertyName("policy_name")]
            [Required(ErrorMessage = "Policy name is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Policy name must be between 1 and 255 characters")]
            public string PolicyName { get; set; } = string.Empty;

            [JsonPropertyName("warranty_duration_months")]
            [Required(ErrorMessage = "Warranty duration is required")]
            [Range(1, 600, ErrorMessage = "Warranty duration must be between 1 and 600 months")]
            public int WarrantyDurationMonths { get; set; }

            [JsonPropertyName("kilometer_limit")]
            [Range(1, 9999999, ErrorMessage = "Kilometer limit must be between 1 and 9,999,999")]
            public int? KilometerLimit { get; set; }

            [JsonPropertyName("terms_and_conditions")]
            [Required(ErrorMessage = "Terms and conditions are required")]
            [StringLength(5000, MinimumLength = 1, ErrorMessage = "Terms must be between 1 and 5000 characters")]
            public string TermsAndConditions { get; set; } = string.Empty;

            [JsonPropertyName("vehicle_model_id")]
            public Guid? VehicleModelId { get; set; }
        }

        public class UpdateWarrantyPolicyRequest
        {
            [JsonPropertyName("policy_name")]
            [Required(ErrorMessage = "Policy name is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Policy name must be between 1 and 255 characters")]
            public string PolicyName { get; set; } = string.Empty;

            [JsonPropertyName("warranty_duration_months")]
            [Required(ErrorMessage = "Warranty duration is required")]
            [Range(1, 600, ErrorMessage = "Warranty duration must be between 1 and 600 months")]
            public int WarrantyDurationMonths { get; set; }

            [JsonPropertyName("kilometer_limit")]
            [Range(1, 9999999, ErrorMessage = "Kilometer limit must be between 1 and 9,999,999")]
            public int? KilometerLimit { get; set; }

            [JsonPropertyName("terms_and_conditions")]
            [Required(ErrorMessage = "Terms and conditions are required")]
            [StringLength(5000, MinimumLength = 1, ErrorMessage = "Terms must be between 1 and 5000 characters")]
            public string TermsAndConditions { get; set; } = string.Empty;

            [JsonPropertyName("vehicle_model_id")]
            public Guid? VehicleModelId { get; set; }
        }

        public class ChangeStatusRequest
        {
            [JsonPropertyName("status")]
            [Required(ErrorMessage = "Status is required")]
            [RegularExpression("^(Draft|Active|Expired|Superseded|Archived)$",
                ErrorMessage = "Status must be Draft, Active, Expired, Superseded, or Archived")]
            public string Status { get; set; } = string.Empty;
        }

        public class WarrantyPolicyResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("policy_name")]
            public string PolicyName { get; set; } = string.Empty;

            [JsonPropertyName("warranty_duration_months")]
            public int WarrantyDurationMonths { get; set; }

            [JsonPropertyName("kilometer_limit")]
            public int? KilometerLimit { get; set; }

            [JsonPropertyName("terms_and_conditions")]
            public string TermsAndConditions { get; set; } = string.Empty;

            [JsonPropertyName("status")]
            public string Status { get; set; } = string.Empty;

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("vehicle_model")]
            public VehicleModelDto.VehicleModelResponse? AssignedModel { get; set; }
        }

        public class WarrantyPolicyWithDetailsResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("policy_name")]
            public string PolicyName { get; set; } = string.Empty;

            [JsonPropertyName("warranty_duration_months")]
            public int WarrantyDurationMonths { get; set; }

            [JsonPropertyName("kilometer_limit")]
            public int? KilometerLimit { get; set; }

            [JsonPropertyName("terms_and_conditions")]
            public string TermsAndConditions { get; set; } = string.Empty;

            [JsonPropertyName("status")]
            public string Status { get; set; } = string.Empty;

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("vehicle_model")]
            public VehicleModelDto.VehicleModelResponse? AssignedModel { get; set; }

            [JsonPropertyName("covered_parts")]
            public List<PolicyCoveragePartDto.PolicyCoveragePartResponse> CoveredParts { get; set; } = new();
        }
    }

    public static class WarrantyPolicyMapper
    {
        public static WarrantyPolicy ToEntity(this WarrantyPolicyDto.CreateWarrantyPolicyRequest request)
        {
            return new WarrantyPolicy(
                request.PolicyName,
                request.WarrantyDurationMonths,
                request.KilometerLimit,
                request.TermsAndConditions
            );
        }

        public static void ApplyToEntity(this WarrantyPolicyDto.UpdateWarrantyPolicyRequest request, WarrantyPolicy policy)
        {
            policy.UpdateDetails(
                request.PolicyName,
                request.WarrantyDurationMonths,
                request.KilometerLimit,
                request.TermsAndConditions
            );
        }

        public static WarrantyPolicyDto.WarrantyPolicyResponse ToResponse(this WarrantyPolicy policy)
        {
            return new WarrantyPolicyDto.WarrantyPolicyResponse
            {
                Id = policy.Id,
                PolicyName = policy.PolicyName,
                WarrantyDurationMonths = policy.WarrantyDurationMonths,
                KilometerLimit = policy.KilometerLimit,
                TermsAndConditions = policy.TermsAndConditions,
                Status = policy.Status.ToString(),
                CreatedAt = policy.CreatedAt,
                UpdatedAt = policy.UpdatedAt,
                AssignedModel = null
            };
        }

        public static WarrantyPolicyDto.WarrantyPolicyWithDetailsResponse ToWithDetailsResponse(this WarrantyPolicy policy)
        {
            return new WarrantyPolicyDto.WarrantyPolicyWithDetailsResponse
            {
                Id = policy.Id,
                PolicyName = policy.PolicyName,
                WarrantyDurationMonths = policy.WarrantyDurationMonths,
                KilometerLimit = policy.KilometerLimit,
                TermsAndConditions = policy.TermsAndConditions,
                Status = policy.Status.ToString(),
                CreatedAt = policy.CreatedAt,
                UpdatedAt = policy.UpdatedAt,
                AssignedModel = policy.AssignedModel?.ToResponse(),
                CoveredParts = policy.CoverageParts?.Select(cp => cp.ToResponse()).ToList() ?? new()
            };
        }
    }
}
