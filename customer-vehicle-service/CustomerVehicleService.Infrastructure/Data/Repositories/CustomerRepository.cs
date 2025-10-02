using CustomerVehicleService.Domain.Entities;
using CustomerVehicleService.Infrastructure.Data.Context;
using CustomerVehicleService.Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Infrastructure.Data.Repositories
{
    public class CustomerRepository : ICustomerRepository
    {
        private readonly CustomerVehicleDbContext _context;
        private readonly DbSet<Customer> _dbSet;

        public CustomerRepository(CustomerVehicleDbContext context)
        {
            _context = context;
            _dbSet = context.Set<Customer>();
        }

        // IRepository implementations
        public async Task<Customer> GetByIdAsync(Guid id)
        {
            return await _dbSet.FindAsync(id);
        }

        public async Task<IEnumerable<Customer>> GetAllAsync()
        {
            return await _dbSet.ToListAsync();
        }

        public async Task<IEnumerable<Customer>> FindAsync(Expression<Func<Customer, bool>> predicate)
        {
            return await _dbSet.Where(predicate).ToListAsync();
        }

        public async Task<Customer> FirstOrDefaultAsync(Expression<Func<Customer, bool>> predicate)
        {
            return await _dbSet.FirstOrDefaultAsync(predicate);
        }

        public async Task<bool> ExistsAsync(Expression<Func<Customer, bool>> predicate)
        {
            return await _dbSet.AnyAsync(predicate);
        }

        public IQueryable<Customer> Query()
        {
            return _dbSet.AsQueryable();
        }

        public async Task AddAsync(Customer entity)
        {
            await _dbSet.AddAsync(entity);
        }

        public async Task AddRangeAsync(IEnumerable<Customer> entities)
        {
            await _dbSet.AddRangeAsync(entities);
        }

        public void Update(Customer entity)
        {
            _dbSet.Update(entity);
        }

        public void Remove(Customer entity)
        {
            _dbSet.Remove(entity);
        }

        // ICustomerRepository specific implementations
        public async Task<Customer> GetByEmailAsync(string email)
        {
            return await _dbSet
                .FirstOrDefaultAsync(c => c.Email == email);
        }

        public async Task<bool> EmailExistsAsync(string email, Guid? excludeCustomerId = null)
        {
            var query = _dbSet.Where(c => c.Email == email);

            if (excludeCustomerId.HasValue)
            {
                query = query.Where(c => c.Id != excludeCustomerId.Value);
            }

            return await query.AnyAsync();
        }

        public async Task<Customer> GetWithVehiclesAsync(Guid customerId)
        {
            return await _dbSet
                .Include(c => c.Vehicles)
                    .ThenInclude(v => v.Model)
                .FirstOrDefaultAsync(c => c.Id == customerId);
        }

        public async Task<IEnumerable<Customer>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
            {
                return await GetAllAsync();
            }

            var lowerSearchTerm = searchTerm.ToLower();

            return await _dbSet
                .Where(c =>
                    c.FirstName.ToLower().Contains(lowerSearchTerm) ||
                    c.LastName.ToLower().Contains(lowerSearchTerm) ||
                    c.Email.ToLower().Contains(lowerSearchTerm) ||
                    c.PhoneNumber.Contains(searchTerm))
                .ToListAsync();
        }

        public Task<Customer> GetByIdIncludingDeletedAsync(Guid id)
        {
            throw new NotImplementedException();
        }

        public Task<IEnumerable<Customer>> GetDeletedCustomersAsync(Guid id)
        {
            throw new NotImplementedException();
        }

        public void UpdateRange(IEnumerable<Customer> entities)
        {
            throw new NotImplementedException();
        }

        public void RemoveRange(IEnumerable<Customer> entities)
        {
            throw new NotImplementedException();
        }
    }
}
