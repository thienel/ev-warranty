using Backend.Dotnet.Domain.Entities;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Backend.Dotnet.Application.DTOs
{
    public class VehicleModelDto
    {
        /// <summary>
        /// Request DTO for adding a new vehicle model to the catalog
        /// Used when staff adds a new EV model that the company will warranty
        /// </summary>
        public class CreateVehicleModelRequest
        {
            [JsonPropertyName("brand")]
            [Required(ErrorMessage = "Brand is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Brand must be between 1 and 100 characters")]
            public string Brand { get; set; } = string.Empty;

            [JsonPropertyName("model_name")]
            [Required(ErrorMessage = "Model name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Model name must be between 1 and 100 characters")]
            public string ModelName { get; set; } = string.Empty;

            [JsonPropertyName("year")]
            [Required(ErrorMessage = "Year is required")]
            [Range(2000, 2100, ErrorMessage = "Year must be between 2000 and 2100")]
            public int Year { get; set; }
        }

        /// <summary>
        /// Request DTO for updating vehicle model information
        /// Used when staff corrects model details (e.g., typo in model name)
        /// </summary>
        public class UpdateVehicleModelRequest
        {
            [JsonPropertyName("brand")]
            [Required(ErrorMessage = "Brand is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Brand must be between 1 and 100 characters")]
            public string Brand { get; set; } = string.Empty;

            [JsonPropertyName("model_name")]
            [Required(ErrorMessage = "Model name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Model name must be between 1 and 100 characters")]
            public string ModelName { get; set; } = string.Empty;

            [JsonPropertyName("year")]
            [Required(ErrorMessage = "Year is required")]
            [Range(2000, 2100, ErrorMessage = "Year must be between 2000 and 2100")]
            public int Year { get; set; }

            [JsonPropertyName("updateMask")]
            [Required(ErrorMessage = "updateMask is required")]
            public List<string> UpdateMask { get; set; } = new();
        }

        /// <summary>
        /// Response DTO for vehicle model basic information
        /// Used in dropdowns, lists, and vehicle registration forms
        /// </summary>
        public class VehicleModelResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("brand")]
            public string Brand { get; set; } = string.Empty;

            [JsonPropertyName("model_name")]
            public string ModelName { get; set; } = string.Empty;

            [JsonPropertyName("year")]
            public int Year { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }
        }

        /// <summary>
        /// Response DTO for vehicle model with usage statistics
        /// Used when staff wants to see how many vehicles use this model
        /// </summary>
        public class VehicleModelWithStatsResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("brand")]
            public string Brand { get; set; } = string.Empty;

            [JsonPropertyName("model_name")]
            public string ModelName { get; set; } = string.Empty;

            [JsonPropertyName("year")]
            public int Year { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("vehicle_count")]
            public int VehicleCount { get; set; }
        }
    }

    public static class VehicleModelMapper
    {
        public static VehicleModel ToEntity(this VehicleModelDto.CreateVehicleModelRequest request)
        {
            return new VehicleModel(request.Brand, request.ModelName, request.Year);
        }

        public static void ApplyToEntity(this VehicleModelDto.UpdateVehicleModelRequest request, VehicleModel entity)
        {
            if (request.ModelName != null || !request.UpdateMask.Any())
            {
                entity.UpdateModel(request.Brand, request.ModelName, request.Year);
                return;
            }

            if (request.UpdateMask.Contains("brand", StringComparer.OrdinalIgnoreCase))
                entity.ChangeBrand(request.Brand ?? string.Empty);

            if (request.UpdateMask.Contains("model_name", StringComparer.OrdinalIgnoreCase))
                entity.ChangeModelName(request.ModelName ?? string.Empty);

            if (request.UpdateMask.Contains("year", StringComparer.OrdinalIgnoreCase))
                entity.ChangeYear(request.Year);
        }

        public static VehicleModelDto.VehicleModelResponse ToResponse(this VehicleModel entity)
        {
            return new VehicleModelDto.VehicleModelResponse
            {
                Id = entity.Id,
                Brand = entity.Brand,
                ModelName = entity.ModelName,
                Year = entity.Year,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt
            };
        }

        public static VehicleModelDto.VehicleModelWithStatsResponse ToStatsResponse(this VehicleModel entity)
        {
            return new VehicleModelDto.VehicleModelWithStatsResponse
            {
                Id = entity.Id,
                Brand = entity.Brand,
                ModelName = entity.ModelName,
                Year = entity.Year,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt,
                VehicleCount = entity.Vehicles?.Count ?? 0
            };
        }
    }

}
