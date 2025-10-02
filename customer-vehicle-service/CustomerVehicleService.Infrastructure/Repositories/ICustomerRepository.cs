using CustomerVehicleService.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Infrastructure.Repositories
{
    public interface ICustomerRepository : IRepository<Customer> 
    {
        // Customer-specific queries
        Task<Customer> GetByEmailAsync(string email);
        Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null);
        Task<Customer> GetWithVehiclesAsync(Guid customerId);
        Task<IEnumerable<Customer>> SearchAsync(string searchTerm);

        // Soft delete operations
        Task<Customer> GetByIdIncludingDeletedAsync(Guid id);
        Task<IEnumerable<Customer>> GetDeletedCustomersAsync(Guid id);
    }
}
