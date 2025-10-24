using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class VehicleRepository : BaseRepository<Vehicle>, IVehicleRepository
    {
        public VehicleRepository(DbContext context) : base(context) { }

        public async Task<Vehicle?> GetByVinAsync(string vin)
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null)
                .Include(v => v.Customer)
                .Include(v => v.Model)
                .FirstOrDefaultAsync(v => v.Vin.ToLower() == vin.ToLower());
        }

        public async Task<bool> VinExistsAsync(string vin, Guid? excludeVehicleId = null)
        {
            var query = _dbSet.Where(v => v.DeletedAt == null && v.Vin.ToLower() == vin.ToLower());

            if (excludeVehicleId.HasValue)
            {
                query = query.Where(v => v.Id != excludeVehicleId.Value);
            }

            return await query.AnyAsync();
        }

        public async Task<Vehicle?> GetByLicensePlateAsync(string licensePlate)
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null)
                .Include(v => v.Customer)
                .Include(v => v.Model)
                .FirstOrDefaultAsync(v => v.LicensePlate.ToLower() == licensePlate.ToLower());
        }

        public async Task<IEnumerable<Vehicle>> GetByCustomerIdAsync(Guid customerId)
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null && v.CustomerId == customerId)
                .Include(v => v.Model)
                .OrderByDescending(v => v.CreatedAt)
                .ToListAsync();
        }

        public async Task<IEnumerable<Vehicle>> GetByModelIdAsync(Guid modelId)
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null && v.ModelId == modelId)
                .Include(v => v.Customer)
                .OrderByDescending(v => v.CreatedAt)
                .ToListAsync();
        }

        public async Task<Vehicle?> GetWithDetailsAsync(Guid vehicleId)
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null)
                .Include(v => v.Customer)
                .Include(v => v.Model)
                .FirstOrDefaultAsync(v => v.Id == vehicleId);
        }

        public async Task<IEnumerable<Vehicle>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
                return await _dbSet
                    .Where(v => v.DeletedAt == null)
                    .Include(v => v.Model)
                    .ToListAsync();

            var term = searchTerm.ToLower();
            return await _dbSet
                .Where(v => v.DeletedAt == null)
                .Include(v => v.Model)
                .Include(v => v.Customer)
                .Where(v =>
                    v.Vin.ToLower().Contains(term) ||
                    v.LicensePlate.ToLower().Contains(term) ||
                    v.Model.Brand.ToLower().Contains(term) ||
                    v.Model.ModelName.ToLower().Contains(term))
                .OrderByDescending(v => v.CreatedAt)
                .ToListAsync();
        }

        // Soft Delete
        public async Task<Vehicle?> GetByIdIncludingDeletedAsync(Guid id)
        {
            return await _dbSet
                .IgnoreQueryFilters()
                .FirstOrDefaultAsync(v => v.Id == id);
        }

        public async Task<IEnumerable<Vehicle>> GetDeletedVehicleAsync()
        {
            return await _dbSet
                .IgnoreQueryFilters()
                .Where(v => v.DeletedAt != null)
                .ToListAsync();
        }

        // Override to exclude soft-deleted vehicle by default
        public override async Task<Vehicle?> GetByIdAsync(Guid id)
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null)
                .Include(v => v.Customer)
                .Include(v => v.Model)
                .FirstOrDefaultAsync(v => v.Id == id);
        }

        public override async Task<IEnumerable<Vehicle>> GetAllAsync()
        {
            return await _dbSet
                .Where(v => v.DeletedAt == null)
                .Include(v => v.Customer)
                .Include(v => v.Model)
                .OrderByDescending(v => v.CreatedAt)
                .ToListAsync();
        }
    }
}
