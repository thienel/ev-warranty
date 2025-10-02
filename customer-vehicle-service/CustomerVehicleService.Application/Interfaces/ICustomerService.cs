using static CustomerVehicleService.Application.DTOs.CustomerDto;

namespace CustomerVehicleService.Application.Interfaces
{
    public interface ICustomerService
    {
        //  CREATE 

        /// <summary>
        /// Register a new customer (vehicle owner) in the system
        /// </summary>
        Task<CustomerResponse> CreateAsync(CreateCustomerRequest request);

        //  READ 

        /// <summary>
        /// Get customer by ID
        /// </summary>
        Task<CustomerResponse?> GetByIdAsync(Guid id);

        /// <summary>
        /// Get customer with all their vehicles (useful for warranty overview)
        /// </summary>
        Task<CustomerWithVehiclesResponse?> GetWithVehiclesAsync(Guid id);

        /// <summary>
        /// Get customer by email (common lookup when customer calls support)
        /// </summary>
        Task<CustomerResponse?> GetByEmailAsync(string email);

        /// <summary>
        /// Get all customers with optional filter for deleted records
        /// </summary>
        Task<List<CustomerResponse>> GetAllAsync(bool includeDeleted = false);

        /// <summary>
        /// Search customers by name, email, or phone
        /// </summary>
        Task<List<CustomerResponse>> SearchAsync(string searchTerm);

        //  UPDATE 

        /// <summary>
        /// Update customer information
        /// </summary>
        Task<CustomerResponse> UpdateAsync(Guid id, UpdateCustomerRequest request);

        //  DELETE 

        /// <summary>
        /// Soft delete a customer (keep for historical records)
        /// </summary>
        Task SoftDeleteAsync(Guid id);

        /// <summary>
        /// Restore a soft-deleted customer
        /// </summary>
        Task RestoreAsync(Guid id);

        //  VALIDATION 

        /// <summary>
        /// Check if customer exists
        /// </summary>
        Task<bool> ExistsAsync(Guid id);

        /// <summary>
        /// Check if email is already registered
        /// </summary>
        Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null);
    }
}