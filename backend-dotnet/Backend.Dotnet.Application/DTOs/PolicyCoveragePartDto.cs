using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;

namespace Backend.Dotnet.Application.DTOs
{
    public class PolicyCoveragePartDto
    {
        /// <summary>
        /// Request DTO for adding a part category to a warranty policy's coverage
        /// Used when staff configures which part categories are covered under warranty
        /// </summary>
        public class CreatePolicyCoveragePartRequest
        {
            [JsonPropertyName("policy_id")]
            [Required(ErrorMessage = "Policy ID is required")]
            public Guid PolicyId { get; set; }

            [JsonPropertyName("part_category_id")]
            [Required(ErrorMessage = "Part category ID is required")]
            public Guid PartCategoryId { get; set; }

            [JsonPropertyName("coverage_conditions")]
            [StringLength(1000, ErrorMessage = "Coverage conditions cannot exceed 1000 characters")]
            public string? CoverageConditions { get; set; }
        }

        /// <summary>
        /// Request DTO for updating coverage conditions of an existing policy-category mapping
        /// Used when staff modifies warranty terms for a specific part category
        /// </summary>
        public class UpdatePolicyCoveragePartRequest
        {
            [JsonPropertyName("coverage_conditions")]
            [StringLength(1000, ErrorMessage = "Coverage conditions cannot exceed 1000 characters")]
            public string? CoverageConditions { get; set; }
        }
        
        /// <summary>
        /// Response DTO for policy coverage part basic information
        /// Used in lists and simple queries
        /// </summary>
        public class PolicyCoveragePartResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("policy_id")]
            public Guid PolicyId { get; set; }

            [JsonPropertyName("policy_name")]
            public string? PolicyName { get; set; }

            [JsonPropertyName("part_category_id")]
            public Guid PartCategoryId { get; set; }

            [JsonPropertyName("category_name")]
            public string? CategoryName { get; set; }

            [JsonPropertyName("coverage_conditions")]
            public string? CoverageConditions { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }
        }

        /// <summary>
        /// Response DTO for policy coverage part with full related entities
        /// Used when detailed information about policy and category is needed
        /// </summary>
        public class PolicyCoveragePartDetailResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("policy_id")]
            public Guid PolicyId { get; set; }

            [JsonPropertyName("part_category_id")]
            public Guid PartCategoryId { get; set; }

            [JsonPropertyName("coverage_conditions")]
            public string? CoverageConditions { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            // Related entities
            [JsonPropertyName("policy")]
            public WarrantyPolicyDto.WarrantyPolicyResponse Policy { get; set; } = null!;

            [JsonPropertyName("part_category")]
            public PartCategoryDto.PartCategoryResponse PartCategory { get; set; } = null!;
        }

        public class CoverageDetailsResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("policy_id")]
            public Guid PolicyId { get; set; }

            [JsonPropertyName("part_category_id")]
            public Guid PartCategoryId { get; set; }

            [JsonPropertyName("part_category_name")]
            public string PartCategoryName { get; set; }

            [JsonPropertyName("coverage_conditions")]
            public string CoverageConditions { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }
        }

    }

    public static class PolicyCoveragePartMapper
    {
        public static PolicyCoveragePart ToEntity(this PolicyCoveragePartDto.CreatePolicyCoveragePartRequest request)
        {
            return new PolicyCoveragePart(
                request.PolicyId,
                request.PartCategoryId,
                request.CoverageConditions
            );
        }

        public static void ApplyToEntity(this PolicyCoveragePartDto.UpdatePolicyCoveragePartRequest request, PolicyCoveragePart entity)
        {
            entity.UpdateCoverageConditions(request.CoverageConditions);
        }

        public static PolicyCoveragePartDto.PolicyCoveragePartResponse ToResponse(this PolicyCoveragePart entity)
        {
            return new PolicyCoveragePartDto.PolicyCoveragePartResponse
            {
                Id = entity.Id,
                PolicyId = entity.PolicyId,
                PolicyName = entity.Policy?.PolicyName,
                PartCategoryId = entity.PartCategoryId,
                CategoryName = entity.PartCategory?.CategoryName,
                CoverageConditions = entity.CoverageConditions,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt
            };
        }

        public static PolicyCoveragePartDto.PolicyCoveragePartDetailResponse ToDetailResponse(this PolicyCoveragePart entity)
        {
            return new PolicyCoveragePartDto.PolicyCoveragePartDetailResponse
            {
                Id = entity.Id,
                PolicyId = entity.PolicyId,
                PartCategoryId = entity.PartCategoryId,
                CoverageConditions = entity.CoverageConditions,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt,
                Policy = entity.Policy?.ToResponse(),
                PartCategory = entity.PartCategory?.ToResponse()
            };
        }

        public static CoverageDetailsResponse ToCoverageDetailsResponse(this PolicyCoveragePart entity)
        {
            return new CoverageDetailsResponse
            {
                Id = entity.Id,
                PolicyId = entity.PolicyId,
                PartCategoryId = entity.PartCategoryId,
                PartCategoryName = entity.PartCategory.CategoryName,
                CoverageConditions = entity.CoverageConditions,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt
            };
        }
    }
}
