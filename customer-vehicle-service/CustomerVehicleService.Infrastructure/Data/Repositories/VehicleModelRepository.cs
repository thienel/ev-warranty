using CustomerVehicleService.Application.Interfaces.Data;
using CustomerVehicleService.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace CustomerVehicleService.Infrastructure.Data.Repositories
{
    public class VehicleModelRepository : BaseRepository<VehicleModel>, IVehicleModelRepository
    {
        public VehicleModelRepository(DbContext context) : base(context) { }

        public async Task<VehicleModel?> GetByBrandModelYearAsync(string brand, string modelName, int year)
        {
            return await _dbSet
                .FirstOrDefaultAsync(vm =>
                    vm.Brand.ToLower() == brand.ToLower() &&
                    vm.ModelName.ToLower() == modelName.ToLower() &&
                    vm.Year == year);
        }

        public async Task<bool> ExistsByBrandModelYearAsync(string brand, string modelName, int year, Guid? excludeModelId = null)
        {
            var query = _dbSet.Where(vm =>
                vm.Brand.ToLower() == brand.ToLower() &&
                vm.ModelName.ToLower() == modelName.ToLower() &&
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
                .Where(vm => vm.Brand.ToLower() == brand.ToLower())
                .OrderBy(vm => vm.Year)
                .ThenBy(vm => vm.ModelName)
                .ToListAsync();
        }

        public async Task<IEnumerable<VehicleModel>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
                return await GetAllAsync();

            var term = searchTerm.ToLower();
            return await _dbSet
                .Where(vm =>
                    vm.Brand.ToLower().Contains(term) ||
                    vm.ModelName.ToLower().Contains(term))
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.ModelName)
                .ThenBy(vm => vm.Year)
                .ToListAsync();
        }

        public async Task<IEnumerable<string>> GetAllBrandsAsync()
        {
            return await _dbSet
                .Select(vm => vm.Brand)
                .Distinct()
                .OrderBy(b => b)
                .ToListAsync();
        }
    }
}
