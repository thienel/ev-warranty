using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class CustomerRepository : BaseRepository<Customer>, ICustomerRepository
    {
        public CustomerRepository(DbContext context) : base(context) { }

        public async Task<Customer?> GetByEmailAsync(string email)
        {
            return await _dbSet
                .Where(c => c.DeletedAt == null)
                .FirstOrDefaultAsync(c => c.Email.ToLower() == email.ToLower());
        }

        public async Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null)
        {
            var query = _dbSet.Where(c => c.DeletedAt == null && c.Email.ToLower() == email.ToLower());

            if (excludeCustomerId.HasValue)
            {
                query = query.Where(c => c.Id != excludeCustomerId.Value);
            }

            return await query.AnyAsync();
        }

        public async Task<Customer?> GetByPhoneAsync(string phone)
        {
            return await _dbSet
                .Where(c => c.DeletedAt == null)
                .FirstOrDefaultAsync(c => c.PhoneNumber == phone);
        }

        public async Task<IEnumerable<Customer>> GetByNameAsync(string firstName, string lastName)
        {
            var query = _dbSet.Where(c => c.DeletedAt == null);

            if (!string.IsNullOrWhiteSpace(firstName))
            {
                query = query.Where(c => c.FirstName.ToLower() == firstName.ToLower());
            }

            if (!string.IsNullOrWhiteSpace(lastName))
            {
                query = query.Where(c => c.LastName.ToLower() == lastName.ToLower());
            }

            return await query.ToListAsync();
        }

        public async Task<Customer?> GetWithVehiclesAsync(Guid customerId)
        {
            return await _dbSet
                .Where(c => c.DeletedAt == null)
                .Include(c => c.Vehicles)
                    .ThenInclude(v => v.Model)
                .FirstOrDefaultAsync(c => c.Id == customerId);
        }

        // Soft Delete
        public async Task<Customer?> GetByIdIncludingDeletedAsync(Guid id)
        {
            return await _dbSet
                .IgnoreQueryFilters()
                .FirstOrDefaultAsync(c => c.Id == id);
        }

        public async Task<IEnumerable<Customer>> GetDeletedCustomersAsync()
        {
            return await _dbSet
                .IgnoreQueryFilters()
                .Where(c => c.DeletedAt != null)
                .ToListAsync();
        }

        // Override to exclude soft-deleted
        public override async Task<Customer?> GetByIdAsync(Guid id)
        {
            return await _dbSet
                .Where(c => c.DeletedAt == null)
                .FirstOrDefaultAsync(c => c.Id == id);
        }

        public override async Task<IEnumerable<Customer>> GetAllAsync()
        {
            return await _dbSet
                .Where(c => c.DeletedAt == null)
                .OrderBy(c => c.LastName)
                .ThenBy(c => c.FirstName)
                .ToListAsync();
        }
    }
}
