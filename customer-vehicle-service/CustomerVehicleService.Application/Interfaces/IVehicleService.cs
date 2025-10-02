using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static CustomerVehicleService.Application.DTOs.VehicleDto;

namespace CustomerVehicleService.Application.Interfaces
{
    public interface IVehicleService
    {
        //  CREATE 

        /// <summary>
        /// Register a new vehicle in the warranty system
        /// </summary>
        Task<VehicleDetailResponse> CreateAsync(CreateVehicleRequest request);

        //  READ 

        /// <summary>
        /// Get vehicle by ID with basic information
        /// </summary>
        Task<VehicleResponse?> GetByIdAsync(Guid id);

        /// <summary>
        /// Get vehicle with full details (owner, model)
        /// </summary>
        Task<VehicleDetailResponse?> GetDetailsAsync(Guid id);

        /// <summary>
        /// Get vehicle by VIN (common lookup)
        /// </summary>
        Task<VehicleDetailResponse?> GetByVinAsync(string vin);

        /// <summary>
        /// Get vehicle by license plate
        /// </summary>
        Task<VehicleDetailResponse?> GetByLicensePlateAsync(string licensePlate);

        /// <summary>
        /// Get all vehicles
        /// </summary>
        Task<List<VehicleResponse>> GetAllAsync();

        /// <summary>
        /// Get all vehicles owned by a customer
        /// </summary>
        Task<List<VehicleDetailResponse>> GetByCustomerIdAsync(Guid customerId);

        /// <summary>
        /// Get all vehicles of a specific model
        /// </summary>
        Task<List<VehicleDetailResponse>> GetByModelIdAsync(Guid modelId);

        /// <summary>
        /// Search vehicles by VIN, license plate, or owner name
        /// </summary>
        Task<List<VehicleDetailResponse>> SearchAsync(string searchTerm);

        //  UPDATE 

        /// <summary>
        /// Update vehicle information
        /// </summary>
        Task<VehicleDetailResponse> UpdateAsync(Guid id, UpdateVehicleRequest request);

        /// <summary>
        /// Update only the license plate
        /// </summary>
        Task<VehicleDetailResponse> UpdateLicensePlateAsync(Guid id, UpdateLicensePlateCommand command);

        /// <summary>
        /// Transfer vehicle to another customer
        /// </summary>
        Task<VehicleDetailResponse> TransferOwnershipAsync(Guid id, TransferVehicleCommand command);

        //  DELETE 

        /// <summary>
        /// Delete a vehicle from the system
        /// </summary>
        Task DeleteAsync(Guid id);

        //  VALIDATION 

        /// <summary>
        /// Check if vehicle exists
        /// </summary>
        Task<bool> ExistsAsync(Guid id);

        /// <summary>
        /// Check if VIN is already registered
        /// </summary>
        Task<bool> VinExistsAsync(string vin, Guid? excludeVehicleId = null);
    }
}
