using CustomerVehicleService.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Application.DTOs
{
    public class VehicleModelDto
    {
        /// <summary>
        /// Request DTO for adding a new vehicle model to the catalog
        /// Used when staff adds a new EV model that the company will warranty
        /// </summary>
        public class CreateVehicleModelRequest
        {
            [Required(ErrorMessage = "Brand is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Brand must be between 1 and 100 characters")]
            public string Brand { get; set; } = string.Empty;

            [Required(ErrorMessage = "Model name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Model name must be between 1 and 100 characters")]
            public string ModelName { get; set; } = string.Empty;

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
            [Required(ErrorMessage = "Brand is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Brand must be between 1 and 100 characters")]
            public string Brand { get; set; } = string.Empty;

            [Required(ErrorMessage = "Model name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Model name must be between 1 and 100 characters")]
            public string ModelName { get; set; } = string.Empty;

            [Required(ErrorMessage = "Year is required")]
            [Range(2000, 2100, ErrorMessage = "Year must be between 2000 and 2100")]
            public int Year { get; set; }
        }

        /// <summary>
        /// Response DTO for vehicle model basic information
        /// Used in dropdowns, lists, and vehicle registration forms
        /// </summary>
        public class VehicleModelResponse
        {
            public Guid Id { get; set; }
            public string Brand { get; set; } = string.Empty;
            public string ModelName { get; set; } = string.Empty;
            public int Year { get; set; }
            public DateTime CreatedAt { get; set; }
            public DateTime? UpdatedAt { get; set; }

            // Display for UI
            public string DisplayName => $"{Brand} {ModelName} {Year}";
        }

        /// <summary>
        /// Response DTO for vehicle model with usage statistics
        /// Used when staff wants to see how many vehicles use this model
        /// </summary>
        public class VehicleModelWithStatsResponse
        {
            public Guid Id { get; set; }
            public string Brand { get; set; } = string.Empty;
            public string ModelName { get; set; } = string.Empty;
            public int Year { get; set; }
            public DateTime CreatedAt { get; set; }
            public DateTime? UpdatedAt { get; set; }

            public int VehicleCount { get; set; }

            // Display helpers
            public string DisplayName => $"{Brand} {ModelName} {Year}";
            public bool CanBeDeleted => VehicleCount == 0; // yet supported
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
            entity.UpdateModel(request.Brand, request.ModelName, request.Year);
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
