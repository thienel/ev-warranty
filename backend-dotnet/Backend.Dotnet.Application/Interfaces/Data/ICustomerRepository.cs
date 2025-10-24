using Backend.Dotnet.Domain.Entities;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface ICustomerRepository : IRepository<Customer>
    {
        Task<IEnumerable<Customer>> GetByEmailAsync(string email);
        Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null);
        Task<IEnumerable<Customer>> GetByPhoneAsync(string phone);
        Task<IEnumerable<Customer>> GetByNameAsync(string name);
        Task<Customer?> GetWithVehiclesAsync(Guid customerId);

        Task<Customer?> GetByIdIncludingDeletedAsync(Guid id);
        Task<IEnumerable<Customer>> GetDeletedCustomersAsync();
    }
}
