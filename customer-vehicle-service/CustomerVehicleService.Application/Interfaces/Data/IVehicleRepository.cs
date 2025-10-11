using CustomerVehicleService.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Application.Interfaces.Data
{
    public interface IVehicleRepository : IRepository<Vehicle>
    {
        // Vehicle-specific queries
        Task<Vehicle?> GetByVinAsync(string vin);
        Task<bool> VinExistsAsync(string vin, Guid? excludeVehicleId = null);
        Task<Vehicle?> GetByLicensePlateAsync(string licensePlate);
        Task<IEnumerable<Vehicle>> GetByCustomerIdAsync(Guid customerId);
        Task<IEnumerable<Vehicle>> GetByModelIdAsync(Guid modelId);
        Task<Vehicle?> GetWithDetailsAsync(Guid vehicleId);
        Task<IEnumerable<Vehicle>> SearchAsync(string searchTerm);

        // Soft delete operations
        Task<Vehicle?> GetByIdIncludingDeletedAsync(Guid id);
        Task<IEnumerable<Vehicle>> GetDeletedVehicleAsync();
    }
}
