using CustomerVehicleService.Application.Interfaces.Data;
using CustomerVehicleService.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace CustomerVehicleService.Infrastructure.Data.Repositories
{
    public class VehicleRepository : BaseRepository<Vehicle>, IVehicleRepository
    {
        public VehicleRepository(DbContext context) : base(context) { }

        public async Task<Vehicle?> GetByVinAsync(string vin)
        {
            return await _dbSet
                .FirstOrDefaultAsync(v => v.Vin.ToLower() == vin.ToLower());
        }

        public async Task<bool> VinExistsAsync(string vin, Guid? excludeVehicleId = null)
        {
            var query = _dbSet.Where(v => v.Vin.ToLower() == vin.ToLower());

            if (excludeVehicleId.HasValue)
            {
                query = query.Where(v => v.Id != excludeVehicleId.Value);
            }

            return await query.AnyAsync();
        }

        public async Task<Vehicle?> GetByLicensePlateAsync(string licensePlate)
        {
            return await _dbSet
                .FirstOrDefaultAsync(v => v.LicensePlate.ToLower() == licensePlate.ToLower());
        }

        public async Task<IEnumerable<Vehicle>> GetByCustomerIdAsync(Guid customerId)
        {
            return await _dbSet
                .Include(v => v.Model)
                .Where(v => v.CustomerId == customerId)
                .OrderByDescending(v => v.CreatedAt)
                .ToListAsync();
        }

        public async Task<IEnumerable<Vehicle>> GetByModelIdAsync(Guid modelId)
        {
            return await _dbSet
                .Include(v => v.Customer)
                .Where(v => v.ModelId == modelId)
                .OrderByDescending(v => v.CreatedAt)
                .ToListAsync();
        }

        public async Task<Vehicle?> GetWithDetailsAsync(Guid vehicleId)
        {
            return await _dbSet
                .Include(v => v.Customer)
                .Include(v => v.Model)
                .FirstOrDefaultAsync(v => v.Id == vehicleId);
        }

        public async Task<IEnumerable<Vehicle>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
                return await _dbSet.Include(v => v.Model).ToListAsync();

            var term = searchTerm.ToLower();
            return await _dbSet
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
    }
}
