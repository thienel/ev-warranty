using Backend.Dotnet.Domain.Entities;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface ICustomerRepository : IRepository<Customer>
    {
        Task<Customer?> GetByEmailAsync(string email);
        Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null);
        Task<Customer?> GetByPhoneAsync(string phone);
        Task<IEnumerable<Customer>> GetByNameAsync(string firstName, string lastName);
        Task<Customer?> GetWithVehiclesAsync(Guid customerId);

        Task<Customer?> GetByIdIncludingDeletedAsync(Guid id);
        Task<IEnumerable<Customer>> GetDeletedCustomersAsync();
    }
}
