using CustomerVehicleService.Application.DTOs;
using static CustomerVehicleService.Application.DTOs.VehicleModelDto;

namespace CustomerVehicleService.Application.Interfaces
{
    public interface IVehicleModelService
    {
        //  CREATE 

        /// <summary>
        /// Add a new vehicle model to the catalog
        /// </summary>
        Task<VehicleModelResponse> CreateAsync(CreateVehicleModelRequest request);

        //  READ 

        /// <summary>
        /// Get vehicle model by ID
        /// </summary>
        Task<VehicleModelResponse?> GetByIdAsync(Guid id);

        /// <summary>
        /// Get vehicle model with statistics
        /// </summary>
        Task<VehicleModelWithStatsResponse?> GetWithStatsAsync(Guid id);

        /// <summary>
        /// Get all vehicle models
        /// </summary>
        Task<List<VehicleModelResponse>> GetAllAsync();

        /// <summary>
        /// Get models by brand
        /// </summary>
        Task<List<VehicleModelResponse>> GetByBrandAsync(string brand);

        /// <summary>
        /// Get models by year
        /// </summary>
        Task<List<VehicleModelResponse>> GetByYearAsync(int year);

        /// <summary>
        /// Search models by brand or model name
        /// </summary>
        Task<List<VehicleModelResponse>> SearchAsync(string searchTerm);

        //  UPDATE 

        /// <summary>
        /// Update vehicle model information
        /// </summary>
        Task<VehicleModelResponse> UpdateAsync(Guid id, UpdateVehicleModelRequest request);

        //  DELETE 

        /// <summary>
        /// Delete a vehicle model (only if no vehicles use it)
        /// </summary>
        Task DeleteAsync(Guid id);

        //  VALIDATION 

        /// <summary>
        /// Check if vehicle model exists
        /// </summary>
        Task<bool> ExistsAsync(Guid id);

        /// <summary>
        /// Check if brand/model/year combination exists
        /// </summary>
        Task<bool> ModelCombinationExistsAsync(string brand, string modelName, int year, Guid? excludeModelId = null);

        /// <summary>
        /// Check if any vehicles use this model
        /// </summary>
        Task<bool> IsUsedByVehiclesAsync(Guid id);
    }
}
