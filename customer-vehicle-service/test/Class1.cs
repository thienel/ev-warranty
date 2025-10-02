namespace test
{
    // ============================================================================
    // CUSTOMER SERVICE INTERFACE
    // ============================================================================

    namespace CustomerVehicleService.Application.Interfaces;

    /// <summary>
    /// Service interface for managing vehicle owners in the warranty system
    /// Used by staff to register, update, and retrieve customer information
    /// </summary>
    public interface ICustomerService
    {
        // ==================== CREATE ====================

        /// <summary>
        /// Register a new customer (vehicle owner) in the system
        /// </summary>
        Task<CustomerResponse> CreateAsync(CreateCustomerRequest request, CancellationToken cancellationToken = default);

        // ==================== READ ====================

        /// <summary>
        /// Get customer by ID
        /// </summary>
        Task<CustomerResponse?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get customer with all their vehicles (useful for warranty overview)
        /// </summary>
        Task<CustomerWithVehiclesResponse?> GetByIdWithVehiclesAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get customer by email (common lookup when customer calls support)
        /// </summary>
        Task<CustomerResponse?> GetByEmailAsync(string email, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all customers with optional filter for deleted records
        /// </summary>
        Task<List<CustomerResponse>> GetAllAsync(bool includeDeleted = false, CancellationToken cancellationToken = default);

        /// <summary>
        /// Search customers by name, email, or phone (for staff lookup)
        /// </summary>
        Task<List<CustomerResponse>> SearchAsync(string searchTerm, CancellationToken cancellationToken = default);

        // ==================== UPDATE ====================

        /// <summary>
        /// Update customer information (corrections or changes)
        /// </summary>
        Task<CustomerResponse> UpdateAsync(Guid id, UpdateCustomerRequest request, CancellationToken cancellationToken = default);

        // ==================== DELETE ====================

        /// <summary>
        /// Soft delete a customer (mark as deleted, keep for historical records)
        /// </summary>
        Task SoftDeleteAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Restore a soft-deleted customer
        /// </summary>
        Task RestoreAsync(Guid id, CancellationToken cancellationToken = default);

        // ==================== VALIDATION ====================

        /// <summary>
        /// Check if customer exists by ID
        /// </summary>
        Task<bool> ExistsAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Check if email is already registered (for duplicate prevention)
        /// </summary>
        Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null, CancellationToken cancellationToken = default);
    }

    // ============================================================================
    // VEHICLE MODEL SERVICE INTERFACE
    // ============================================================================

    /// <summary>
    /// Service interface for managing vehicle model catalog
    /// Used by staff to maintain the list of EV models covered by warranty
    /// </summary>
    public interface IVehicleModelService
    {
        // ==================== CREATE ====================

        /// <summary>
        /// Add a new vehicle model to the catalog
        /// </summary>
        Task<VehicleModelResponse> CreateAsync(CreateVehicleModelRequest request, CancellationToken cancellationToken = default);

        // ==================== READ ====================

        /// <summary>
        /// Get vehicle model by ID
        /// </summary>
        Task<VehicleModelResponse?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get vehicle model with statistics (how many vehicles use this model)
        /// </summary>
        Task<VehicleModelWithStatsResponse?> GetByIdWithStatsAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get vehicle model by brand, model name, and year (unique combination)
        /// </summary>
        Task<VehicleModelResponse?> GetByBrandModelYearAsync(string brand, string modelName, int year, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all vehicle models (for dropdown lists and catalogs)
        /// </summary>
        Task<List<VehicleModelResponse>> GetAllAsync(CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all models from a specific brand
        /// </summary>
        Task<List<VehicleModelResponse>> GetByBrandAsync(string brand, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all models from a specific year
        /// </summary>
        Task<List<VehicleModelResponse>> GetByYearAsync(int year, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all unique brands in the catalog
        /// </summary>
        Task<List<string>> GetAllBrandsAsync(CancellationToken cancellationToken = default);

        /// <summary>
        /// Search models by brand or model name
        /// </summary>
        Task<List<VehicleModelResponse>> SearchAsync(string searchTerm, CancellationToken cancellationToken = default);

        // ==================== UPDATE ====================

        /// <summary>
        /// Update vehicle model information (corrections)
        /// </summary>
        Task<VehicleModelResponse> UpdateAsync(Guid id, UpdateVehicleModelRequest request, CancellationToken cancellationToken = default);

        // ==================== DELETE ====================

        /// <summary>
        /// Delete a vehicle model (only if no vehicles are using it)
        /// </summary>
        Task DeleteAsync(Guid id, CancellationToken cancellationToken = default);

        // ==================== VALIDATION ====================

        /// <summary>
        /// Check if vehicle model exists by ID
        /// </summary>
        Task<bool> ExistsAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Check if this brand/model/year combination already exists (for duplicate prevention)
        /// </summary>
        Task<bool> ModelCombinationExistsAsync(string brand, string modelName, int year, Guid? excludeModelId = null, CancellationToken cancellationToken = default);

        /// <summary>
        /// Check if any vehicles are using this model (blocks deletion if true)
        /// </summary>
        Task<bool> IsUsedByVehiclesAsync(Guid id, CancellationToken cancellationToken = default);
    }

    // ============================================================================
    // VEHICLE SERVICE INTERFACE
    // ============================================================================

    /// <summary>
    /// Service interface for managing vehicles in the warranty system
    /// Used by staff to register vehicles and track ownership/warranty information
    /// </summary>
    public interface IVehicleService
    {
        // ==================== CREATE ====================

        /// <summary>
        /// Register a new vehicle in the warranty system
        /// </summary>
        Task<VehicleDetailResponse> CreateAsync(CreateVehicleRequest request, CancellationToken cancellationToken = default);

        // ==================== READ ====================

        /// <summary>
        /// Get vehicle by ID with basic information
        /// </summary>
        Task<VehicleResponse?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get vehicle with full details (owner, model) - used for warranty processing
        /// </summary>
        Task<VehicleDetailResponse?> GetDetailsByIdAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get vehicle by VIN (common lookup when customer brings vehicle for service)
        /// </summary>
        Task<VehicleDetailResponse?> GetByVinAsync(string vin, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get vehicle by license plate (alternative lookup method)
        /// </summary>
        Task<VehicleDetailResponse?> GetByLicensePlateAsync(string licensePlate, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all vehicles
        /// </summary>
        Task<List<VehicleResponse>> GetAllAsync(CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all vehicles owned by a specific customer
        /// </summary>
        Task<List<VehicleDetailResponse>> GetByCustomerIdAsync(Guid customerId, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get all vehicles of a specific model (useful for recalls or model-specific issues)
        /// </summary>
        Task<List<VehicleDetailResponse>> GetByModelIdAsync(Guid modelId, CancellationToken cancellationToken = default);

        /// <summary>
        /// Search vehicles by VIN, license plate, or owner name
        /// </summary>
        Task<List<VehicleDetailResponse>> SearchAsync(string searchTerm, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get vehicles purchased within a date range (useful for warranty batch operations)
        /// </summary>
        Task<List<VehicleDetailResponse>> GetByPurchaseDateRangeAsync(DateTime startDate, DateTime endDate, CancellationToken cancellationToken = default);

        // ==================== UPDATE ====================

        /// <summary>
        /// Update vehicle information (corrections or changes)
        /// </summary>
        Task<VehicleDetailResponse> UpdateAsync(Guid id, UpdateVehicleRequest request, CancellationToken cancellationToken = default);

        /// <summary>
        /// Update only the license plate (quick common operation)
        /// </summary>
        Task<VehicleDetailResponse> UpdateLicensePlateAsync(Guid id, UpdateLicensePlateCommand command, CancellationToken cancellationToken = default);

        /// <summary>
        /// Transfer vehicle ownership to another customer (warranty transfer)
        /// </summary>
        Task<VehicleDetailResponse> TransferOwnershipAsync(Guid id, TransferVehicleCommand command, CancellationToken cancellationToken = default);

        // ==================== DELETE ====================

        /// <summary>
        /// Delete a vehicle from the system (hard delete)
        /// </summary>
        Task DeleteAsync(Guid id, CancellationToken cancellationToken = default);

        // ==================== VALIDATION ====================

        /// <summary>
        /// Check if vehicle exists by ID
        /// </summary>
        Task<bool> ExistsAsync(Guid id, CancellationToken cancellationToken = default);

        /// <summary>
        /// Check if VIN is already registered (for duplicate prevention)
        /// </summary>
        Task<bool> VinExistsAsync(string vin, Guid? excludeVehicleId = null, CancellationToken cancellationToken = default);

        /// <summary>
        /// Get count of vehicles owned by a customer (useful for validation and stats)
        /// </summary>
        Task<int> GetCustomerVehicleCountAsync(Guid customerId, CancellationToken cancellationToken = default);
    }
}

namespace test2
{
    
}