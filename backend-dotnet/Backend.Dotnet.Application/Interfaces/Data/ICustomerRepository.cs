using Backend.Dotnet.Domain.Entities;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface ICustomerRepository : IRepository<Customer>
    {
        // Customer-specific queries
        Task<Customer?> GetByEmailAsync(string email);
        Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null);
        Task<Customer?> GetWithVehiclesAsync(Guid customerId);
        Task<IEnumerable<Customer>> SearchAsync(string searchTerm);

        // Soft delete operations
        Task<Customer?> GetByIdIncludingDeletedAsync(Guid id);
        Task<IEnumerable<Customer>> GetDeletedCustomersAsync();
    }
}
