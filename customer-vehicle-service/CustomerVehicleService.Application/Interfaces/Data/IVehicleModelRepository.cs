using CustomerVehicleService.Domain.Entities;

namespace CustomerVehicleService.Application.Interfaces.Data
{
    public interface IVehicleModelRepository : IRepository<VehicleModel>
    {
        // VehicleModel-specific queries
        Task<VehicleModel?> GetByBrandModelYearAsync(string brand, string modelName, int year);
        Task<bool> ExistsByBrandModelYearAsync(string brand, string modelName, int year, Guid? excludeModelId = null);
        Task<IEnumerable<VehicleModel>> GetByBrandAsync(string brand);
        Task<IEnumerable<VehicleModel>> SearchAsync(string searchTerm);
        Task<IEnumerable<string>> GetAllBrandsAsync();
    }
}
