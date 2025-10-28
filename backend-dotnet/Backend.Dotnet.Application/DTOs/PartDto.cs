using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.DTOs
{
    public class PartDto
    {
        public class CreatePartRequest
        {
            [JsonPropertyName("serial_number")]
            [Required(ErrorMessage = "Serial number is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Serial number must be between 1 and 255 characters")]
            public string SerialNumber { get; set; } = string.Empty;

            [JsonPropertyName("part_name")]
            [Required(ErrorMessage = "Part name is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Part name must be between 1 and 255 characters")]
            public string PartName { get; set; } = string.Empty;

            [JsonPropertyName("unit_price")]
            [Required(ErrorMessage = "Unit price is required")]
            [Range(0.01, 999999999.99, ErrorMessage = "Unit price must be between 0.01 and 999,999,999.99")]
            public decimal UnitPrice { get; set; }

            [JsonPropertyName("category_id")]
            [Required(ErrorMessage = "Category ID is required")]
            public Guid CategoryId { get; set; }

            [JsonPropertyName("office_location_id")]
            public Guid? OfficeLocationId { get; set; }
        }

        public class UpdatePartRequest
        {
            [JsonPropertyName("part_name")]
            [Required(ErrorMessage = "Part name is required")]
            [StringLength(255, MinimumLength = 1, ErrorMessage = "Part name must be between 1 and 255 characters")]
            public string PartName { get; set; } = string.Empty;

            [JsonPropertyName("unit_price")]
            [Required(ErrorMessage = "Unit price is required")]
            [Range(0.01, 999999999.99, ErrorMessage = "Unit price must be between 0.01 and 999,999,999.99")]
            public decimal UnitPrice { get; set; }

            [JsonPropertyName("office_location_id")]
            public Guid? OfficeLocationId { get; set; }
        }

        public class ChangePartCategoryRequest
        {
            [JsonPropertyName("category_id")]
            [Required(ErrorMessage = "Category ID is required")]
            public Guid CategoryId { get; set; }
        }

        public class ChangeStatusRequest
        {
            [JsonPropertyName("status")]
            [Required(ErrorMessage = "Status is required")]
            [RegularExpression("^(Available|Reserved|Installed|Defective|Obsolete|Archived)$",
                ErrorMessage = "Status must be Available, Reserved, Installed, Defective, Obsolete, or Archived")]
            public string Status { get; set; } = string.Empty;
        }

        public class PartResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("serial_number")]
            public string SerialNumber { get; set; } = string.Empty;

            [JsonPropertyName("part_name")]
            public string PartName { get; set; } = string.Empty;

            [JsonPropertyName("unit_price")]
            public decimal UnitPrice { get; set; }

            [JsonPropertyName("category_id")]
            public Guid CategoryId { get; set; }

            [JsonPropertyName("category_name")]
            public string? CategoryName { get; set; }

            [JsonPropertyName("office_location_id")]
            public Guid? OfficeLocationId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; } = string.Empty;

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("can_be_used_in_work_order")]
            public bool CanBeUsedInWorkOrder { get; set; }

            [JsonPropertyName("is_in_stock")]
            public bool IsInStock { get; set; }
        }

        public class PartWithDetailsResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("serial_number")]
            public string SerialNumber { get; set; } = string.Empty;

            [JsonPropertyName("part_name")]
            public string PartName { get; set; } = string.Empty;

            [JsonPropertyName("unit_price")]
            public decimal UnitPrice { get; set; }

            [JsonPropertyName("category_id")]
            public Guid CategoryId { get; set; }

            [JsonPropertyName("office_location_id")]
            public Guid? OfficeLocationId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; } = string.Empty;

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("category")]
            public PartCategoryDto.PartCategoryResponse Category { get; set; } = null!;
        }
    }

    public static class PartMapper
    {
        public static Part ToEntity(this PartDto.CreatePartRequest request)
        {
            return new Part(
                request.SerialNumber,
                request.PartName,
                request.UnitPrice,
                request.CategoryId,
                request.OfficeLocationId
            );
        }

        public static void ApplyToEntity(this PartDto.UpdatePartRequest request, Part part)
        {
            part.UpdateDetails(
                request.PartName,
                request.UnitPrice,
                request.OfficeLocationId
            );
        }

        public static PartDto.PartResponse ToResponse(this Part part)
        {
            return new PartDto.PartResponse
            {
                Id = part.Id,
                SerialNumber = part.SerialNumber,
                PartName = part.PartName,
                UnitPrice = part.UnitPrice,
                CategoryId = part.CategoryId,
                CategoryName = part.Category?.CategoryName,
                OfficeLocationId = part.OfficeLocationId,
                Status = part.Status.ToString(),
                CreatedAt = part.CreatedAt,
                UpdatedAt = part.UpdatedAt,
                CanBeUsedInWorkOrder = part.CanBeUsedInWorkOrder(),
                IsInStock = part.IsInStock()
            };
        }

        public static PartDto.PartWithDetailsResponse ToWithDetailsResponse(this Part part)
        {
            return new PartDto.PartWithDetailsResponse
            {
                Id = part.Id,
                SerialNumber = part.SerialNumber,
                PartName = part.PartName,
                UnitPrice = part.UnitPrice,
                CategoryId = part.CategoryId,
                OfficeLocationId = part.OfficeLocationId,
                Status = part.Status.ToString(),
                CreatedAt = part.CreatedAt,
                UpdatedAt = part.UpdatedAt,
                Category = part.Category?.ToResponse()
            };
        }
    }
}
