using CustomerVehicleService.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static CustomerVehicleService.Application.DTOs.CustomerDto;
using static CustomerVehicleService.Application.DTOs.VehicleModelDto;

namespace CustomerVehicleService.Application.DTOs
{
    public class VehicleDto
    {
        /// <summary>
        /// Request DTO for registering a new vehicle in the warranty system
        /// Used when staff registers a customer's vehicle with warranty coverage
        /// </summary>
        public class CreateVehicleRequest
        {
            [Required(ErrorMessage = "VIN is required")]
            [StringLength(17, MinimumLength = 17, ErrorMessage = "VIN must be exactly 17 characters")]
            [RegularExpression(@"^[A-HJ-NPR-Z0-9]{17}$",
                ErrorMessage = "Invalid VIN format. Must be 17 characters (letters and numbers, excluding I, O, Q)")]
            public string Vin { get; set; } = string.Empty;

            [StringLength(20, ErrorMessage = "License plate cannot exceed 20 characters")]
            public string? LicensePlate { get; set; }

            [Required(ErrorMessage = "Customer ID is required")]
            public Guid CustomerId { get; set; }

            [Required(ErrorMessage = "Vehicle model ID is required")]
            public Guid ModelId { get; set; }

            [DataType(DataType.Date)]
            public DateTime? PurchaseDate { get; set; }
        }

        /// <summary>
        /// Request DTO for updating vehicle information
        /// Used when staff corrects vehicle details or updates registration info
        /// </summary>
        public class UpdateVehicleRequest
        {
            [Required(ErrorMessage = "VIN is required")]
            [StringLength(17, MinimumLength = 17, ErrorMessage = "VIN must be exactly 17 characters")]
            [RegularExpression(@"^[A-HJ-NPR-Z0-9]{17}$",
                ErrorMessage = "Invalid VIN format. Must be 17 characters (letters and numbers, excluding I, O, Q)")]
            public string Vin { get; set; } = string.Empty;

            [StringLength(20, ErrorMessage = "License plate cannot exceed 20 characters")]
            public string? LicensePlate { get; set; }

            [Required(ErrorMessage = "Customer ID is required")]
            public Guid CustomerId { get; set; }

            [Required(ErrorMessage = "Vehicle model ID is required")]
            public Guid ModelId { get; set; }

            [DataType(DataType.Date)]
            public DateTime? PurchaseDate { get; set; }
        }

        /// <summary>
        /// Command DTO for transferring vehicle ownership
        /// Used when vehicle is sold/transferred to another customer
        /// Separate from full update to make intent clear and simplify the operation
        /// </summary>
        public class TransferVehicleCommand
        {
            [Required(ErrorMessage = "New customer ID is required")]
            public Guid NewCustomerId { get; set; }
        }

        /// <summary>
        /// Command DTO for updating only the license plate
        /// Common operation when customer gets new plates - no need for full update
        /// </summary>
        public class UpdateLicensePlateCommand
        {
            [StringLength(20, ErrorMessage = "License plate cannot exceed 20 characters")]
            public string? LicensePlate { get; set; }
        }

        /// <summary>
        /// Response DTO for vehicle basic information
        /// Used in vehicle lists and simple lookups
        /// </summary>
        public class VehicleResponse
        {
            public Guid Id { get; set; }
            public string Vin { get; set; } = string.Empty;
            public string? LicensePlate { get; set; }
            public Guid CustomerId { get; set; }
            public Guid ModelId { get; set; }
            public DateTime? PurchaseDate { get; set; }
            public DateTime CreatedAt { get; set; }
            public DateTime? UpdatedAt { get; set; }
        }

        /// <summary>
        /// Response DTO for vehicle with full related information
        /// Used when processing warranty claims - staff needs owner and model details
        /// </summary>
        public class VehicleDetailResponse
        {
            public Guid Id { get; set; }
            public string Vin { get; set; } = string.Empty;
            public string? LicensePlate { get; set; }
            public DateTime? PurchaseDate { get; set; }
            public DateTime CreatedAt { get; set; }
            public DateTime? UpdatedAt { get; set; }

            // Related entities - needed for warranty processing
            public CustomerResponse Owner { get; set; } = null!;
            public VehicleModelResponse Model { get; set; } = null!;

            // Display helpers for UI
            public string DisplayName => $"{Model.DisplayName} - {Vin}";
            public string OwnerName => Owner.FullName;
            public int? VehicleAgeYears => PurchaseDate.HasValue
                ? DateTime.UtcNow.Year - PurchaseDate.Value.Year
                : null;
        }
    }

    // !!! should unify to Mapper(suggest)
    public static class VehicleMapper
    {
        public static Vehicle ToEntity(this VehicleDto.CreateVehicleRequest request)
        {
            return new Vehicle(
                request.Vin,
                request.CustomerId,
                request.ModelId,
                request.LicensePlate,
                request.PurchaseDate
            );
        }

        public static void ApplyToEntity(this VehicleDto.UpdateVehicleRequest request, Vehicle entity)
        {
            entity.UpdateVehicle(
                request.Vin,
                request.CustomerId,
                request.ModelId,
                request.LicensePlate,
                request.PurchaseDate
            );
        }

        public static void ApplyTransfer(this VehicleDto.TransferVehicleCommand request, Vehicle entity)
        {
            entity.TransferToCustomer(request.NewCustomerId);
        }

        public static void ApplyLicensePlateUpdate(this VehicleDto.UpdateLicensePlateCommand request, Vehicle entity)
        {
            entity.ChangeLicensePlate(request.LicensePlate);
        }

        public static VehicleDto.VehicleResponse ToResponse(this Vehicle entity)
        {
            return new VehicleDto.VehicleResponse
            {
                Id = entity.Id,
                Vin = entity.Vin,
                LicensePlate = entity.LicensePlate,
                CustomerId = entity.CustomerId,
                ModelId = entity.ModelId,
                PurchaseDate = entity.PurchaseDate,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt
            };
        }

        public static VehicleDto.VehicleDetailResponse ToDetailResponse(this Vehicle entity)
        {
            return new VehicleDto.VehicleDetailResponse
            {
                Id = entity.Id,
                Vin = entity.Vin,
                LicensePlate = entity.LicensePlate,
                PurchaseDate = entity.PurchaseDate,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt,
                Owner = entity.Customer.ToResponse(),   // mapping CustomerResponse
                Model = entity.Model.ToResponse()      // mapping VehicleModelResponse
            };
        }
    }
}
