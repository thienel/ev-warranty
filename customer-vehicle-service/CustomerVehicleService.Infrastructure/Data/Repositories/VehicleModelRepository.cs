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
    public class VehicleModelRepository : IVehicleModelRepository
    {
        private readonly CustomerVehicleDbContext _context;
        private readonly DbSet<VehicleModel> _dbSet;
        
        public VehicleModelRepository(CustomerVehicleDbContext context)
        {
            _context = context;
            _dbSet = context.Set<VehicleModel>();
        }

        // IRepository implementations
        public async Task<VehicleModel> GetByIdAsync(Guid id)
        {
            return await _dbSet.FindAsync(id);
        }

        public async Task<IEnumerable<VehicleModel>> GetAllAsync()
        {
            return await _dbSet.ToListAsync();
        }

        public async Task<IEnumerable<VehicleModel>> FindAsync(Expression<Func<VehicleModel, bool>> predicate)
        {
            return await _dbSet.Where(predicate).ToListAsync();
        }

        public async Task<VehicleModel> FirstOrDefaultAsync(Expression<Func<VehicleModel, bool>> predicate)
        {
            return await _dbSet.FirstOrDefaultAsync(predicate);
        }

        public async Task<bool> ExistsAsync(Expression<Func<VehicleModel, bool>> predicate)
        {
            return await _dbSet.AnyAsync(predicate);
        }

        public IQueryable<VehicleModel> Query()
        {
            return _dbSet.AsQueryable();
        }

        public async Task AddAsync(VehicleModel entity)
        {
            await _dbSet.AddAsync(entity);
        }

        public async Task AddRangeAsync(IEnumerable<VehicleModel> entities)
        {
            await _dbSet.AddRangeAsync(entities);
        }

        public void Update(VehicleModel entity)
        {
            _dbSet.Update(entity);
        }

        public void Remove(VehicleModel entity)
        {
            _dbSet.Remove(entity);
        }

        // IVehicleModelRepository specific implementations
        public async Task<VehicleModel> GetByBrandModelYearAsync(string brand, string modelName, int year)
        {
            return await _dbSet
                .FirstOrDefaultAsync(vm =>
                    vm.Brand == brand &&
                    vm.ModelName == modelName &&
                    vm.Year == year);
        }

        public async Task<bool> ExistsByBrandModelYearAsync(string brand, string modelName, int year, Guid? excludeModelId = null)
        {
            var query = _dbSet.Where(vm =>
                vm.Brand == brand &&
                vm.ModelName == modelName &&
                vm.Year == year);

            if (excludeModelId.HasValue)
            {
                query = query.Where(vm => vm.Id != excludeModelId.Value);
            }

            return await query.AnyAsync();
        }

        public async Task<IEnumerable<VehicleModel>> GetByBrandAsync(string brand)
        {
            return await _dbSet
                .Where(vm => vm.Brand == brand)
                .OrderBy(vm => vm.Year)
                .ThenBy(vm => vm.ModelName)
                .ToListAsync();
        }

        public async Task<IEnumerable<VehicleModel>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
            {
                return await GetAllAsync();
            }

            var lowerSearchTerm = searchTerm.ToLower();

            return await _dbSet
                .Where(vm =>
                    vm.Brand.ToLower().Contains(lowerSearchTerm) ||
                    vm.ModelName.ToLower().Contains(lowerSearchTerm))
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.ModelName)
                .ThenBy(vm => vm.Year)
                .ToListAsync();
        }

        public Task<IEnumerable<string>> GetAllBrandsAsync()
        {
            throw new NotImplementedException();
        }

        public void UpdateRange(IEnumerable<VehicleModel> entities)
        {
            throw new NotImplementedException();
        }

        public void RemoveRange(IEnumerable<VehicleModel> entities)
        {
            throw new NotImplementedException();
        }
    }
}
