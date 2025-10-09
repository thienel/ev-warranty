using CustomerVehicleService.Application.Interfaces.Data;
using CustomerVehicleService.Domain.Entities;
using CustomerVehicleService.Infrastructure.Data.Context;
using Microsoft.EntityFrameworkCore;

namespace CustomerVehicleService.Infrastructure.Data.Repositories
{
    public class CustomerRepository : BaseRepository<Customer>, ICustomerRepository
    {
        public CustomerRepository(DbContext context) : base(context) { }

        public async Task<Customer?> GetByEmailAsync(string email)
        {
            return await _dbSet
                .Where(c => !c.IsDeleted)
                .FirstOrDefaultAsync(c => c.Email.ToLower() == email.ToLower());
        }

        public async Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null)
        {
            var query = _dbSet.Where(c => !c.IsDeleted && c.Email.ToLower() == email.ToLower());

            if (excludeCustomerId.HasValue)
            {
                query = query.Where(c => c.Id != excludeCustomerId.Value);
            }

            return await query.AnyAsync();
        }

        public async Task<Customer?> GetWithVehiclesAsync(Guid customerId)
        {
            return await _dbSet
                .Where(c => !c.IsDeleted)
                .Include(c => c.Vehicles)
                    .ThenInclude(v => v.Model)
                .FirstOrDefaultAsync(c => c.Id == customerId);
        }

        public async Task<IEnumerable<Customer>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
                return await _dbSet.Where(c => !c.IsDeleted).ToListAsync();

            var term = searchTerm.ToLower();
            return await _dbSet
                .Where(c => !c.IsDeleted &&
                       (c.FirstName.ToLower().Contains(term) ||
                        c.LastName.ToLower().Contains(term) ||
                        c.Email.ToLower().Contains(term) ||
                        c.PhoneNumber.Contains(term)))
                .ToListAsync();
        }

        public async Task<Customer?> GetByIdIncludingDeletedAsync(Guid id)
        {
            return await _dbSet
                .IgnoreQueryFilters()
                .FirstOrDefaultAsync(c => c.Id == id);
        }

        public async Task<IEnumerable<Customer>> GetDeletedCustomersAsync(Guid id)
        {
            return await _dbSet
                .IgnoreQueryFilters()
                .Where(c => c.IsDeleted)
                .ToListAsync();
        }

        // Override to exclude soft-deleted customers by default
        public override async Task<Customer?> GetByIdAsync(Guid id)
        {
            return await _dbSet
                .Where(c => !c.IsDeleted)
                .FirstOrDefaultAsync(c => c.Id == id);
        }

        public override async Task<IEnumerable<Customer>> GetAllAsync()
        {
            return await _dbSet
                .Where(c => !c.IsDeleted)
                .ToListAsync();
        }
    }
}
