using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class VehicleModelRepository : BaseRepository<VehicleModel>, IVehicleModelRepository
    {
        public VehicleModelRepository(DbContext context) : base(context) { }

        public async Task<IEnumerable<VehicleModel>> GetByBrandModelYearAsync(string brand, string modelName, int year)
        {
            return await _dbSet
                .Where(vm =>
                vm.Brand.ToLower().Contains(brand.ToLower()) &&
                vm.ModelName.ToLower().Contains(modelName.ToLower()) &&
                    vm.Year == year)
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.ModelName)
                .ToListAsync();
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
                .Where(vm => vm.Brand.ToLower().Contains(brand.ToLower()))
                .OrderBy(vm => vm.Year)
                .ThenBy(vm => vm.ModelName)
                .ToListAsync();
        }
        public async Task<IEnumerable<VehicleModel>> GetByModelNameAsync(string modelName)
        {
            return await _dbSet
                .Where(vm => vm.ModelName.ToLower().Contains(modelName.ToLower()))
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.Year)
                .ToListAsync();
        }

        public async Task<IEnumerable<VehicleModel>> GetByYearAsync(int year)
        {
            return await _dbSet
                .Where(vm => vm.Year == year)
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.ModelName)
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

        public async Task<bool> HasActiveVehiclesAsync(Guid modelId)
        {
            return await _context.Set<Vehicle>()
                .Where(v => v.DeletedAt == null)
                .AnyAsync(v => v.ModelId == modelId);
        }

        public async Task<int> GetActiveVehicleCountAsync(Guid modelId)
        {
            return await _context.Set<Vehicle>()
                .Where(v => v.DeletedAt == null)
                .CountAsync(v => v.ModelId == modelId);
        }

        public async Task<IEnumerable<VehicleModel>> GetByPolicyIdAsync(Guid policyId)
        {
            return await _dbSet
                .Where(vm => vm.PolicyId == policyId)
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.ModelName)
                .ThenBy(vm => vm.Year)
                .ToListAsync();
        }

        public async Task<VehicleModel?> GetWithPolicyAsync(Guid modelId)
        {
            return await _dbSet
                .Include(vm => vm.Policy)
                .FirstOrDefaultAsync(vm => vm.Id == modelId);
        }

        // Overide for filtering
        public override async Task<IEnumerable<VehicleModel>> GetAllAsync()
        {
            return await _dbSet
                .OrderBy(vm => vm.Brand)
                .ThenBy(vm => vm.ModelName)
                .ThenBy(vm => vm.Year)
                .ToListAsync();
        }
    }
}
