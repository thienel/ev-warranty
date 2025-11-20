using Backend.Dotnet.Domain.Entities;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Backend.Dotnet.Application.DTOs
{
    public class PartCategoryDto
    {
        public class CreatePartCategoryRequest
        {
            [JsonPropertyName("category_name")]
            [Required(ErrorMessage = "Category name is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Category name must be between 1 and 255 characters")]
            public string CategoryName { get; set; } = string.Empty;

            [JsonPropertyName("description")]
            [StringLength(1000, ErrorMessage = "Description cannot exceed 1000 characters")]
            public string? Description { get; set; }

            [JsonPropertyName("parent_category_id")]
            public Guid? ParentCategoryId { get; set; }
        }

        public class UpdatePartCategoryRequest
        {
            [JsonPropertyName("category_name")]
            [Required(ErrorMessage = "Category name is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Category name must be between 1 and 255 characters")]
            public string CategoryName { get; set; } = string.Empty;

            [JsonPropertyName("description")]
            [StringLength(1000, ErrorMessage = "Description cannot exceed 1000 characters")]
            public string? Description { get; set; }
        }

        public class ChangeParentCategoryRequest
        {
            [JsonPropertyName("new_parent_category_id")]
            public Guid? NewParentCategoryId { get; set; }
        }

        public class PartCategoryResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("category_name")]
            public string CategoryName { get; set; } = string.Empty;

            [JsonPropertyName("description")]
            public string? Description { get; set; }

            [JsonPropertyName("parent_category_id")]
            public Guid? ParentCategoryId { get; set; }

            [JsonPropertyName("parent_category_name")]
            public string? ParentCategoryName { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("can_be_used_for_new_parts")]
            public bool CanBeUsedForNewParts { get; set; }

            [JsonPropertyName("has_active_parts")]
            public bool HasActiveParts { get; set; }
        }

        public class PartCategoryWithHierarchyResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("category_name")]
            public string CategoryName { get; set; } = string.Empty;

            [JsonPropertyName("description")]
            public string? Description { get; set; }

            [JsonPropertyName("parent_category_id")]
            public Guid? ParentCategoryId { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("parent_category")]
            public PartCategoryResponse? ParentCategory { get; set; }

            [JsonPropertyName("child_categories")]
            public List<PartCategoryWithHierarchyResponse> ChildCategories { get; set; } = new();

            [JsonPropertyName("parts_count")]
            public int PartsCount { get; set; }
        }
    }

    public static class PartCategoryMapper
    {
        public static PartCategory ToEntity(this PartCategoryDto.CreatePartCategoryRequest request)
        {
            return new PartCategory(
                request.CategoryName,
                request.Description,
                request.ParentCategoryId
            );
        }

        public static void ApplyToEntity(this PartCategoryDto.UpdatePartCategoryRequest request, PartCategory category)
        {
            category.UpdateDetails(
                request.CategoryName,
                request.Description
            );
        }

        public static PartCategoryDto.PartCategoryResponse ToResponse(this PartCategory category)
        {
            return new PartCategoryDto.PartCategoryResponse
            {
                Id = category.Id,
                CategoryName = category.CategoryName,
                Description = category.Description,
                ParentCategoryId = category.ParentCategoryId,
                ParentCategoryName = category.ParentCategory?.CategoryName,
                CreatedAt = category.CreatedAt,
                UpdatedAt = category.UpdatedAt,
                HasActiveParts = category.HasActiveParts()
            };
        }

        public static PartCategoryDto.PartCategoryWithHierarchyResponse ToWithHierarchyResponse(this PartCategory category)
        {
            return new PartCategoryDto.PartCategoryWithHierarchyResponse
            {
                Id = category.Id,
                CategoryName = category.CategoryName,
                Description = category.Description,
                ParentCategoryId = category.ParentCategoryId,
                CreatedAt = category.CreatedAt,
                UpdatedAt = category.UpdatedAt,
                ParentCategory = category.ParentCategory?.ToResponse(),
                ChildCategories = category.ChildCategories?.Select(c => c.ToWithHierarchyResponse()).ToList() 
                        ?? new List<PartCategoryDto.PartCategoryWithHierarchyResponse>(),
                PartsCount = category.Parts?.Count ?? 0
            };
        }
    }
}
