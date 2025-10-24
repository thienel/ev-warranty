using Backend.Dotnet.Domain.Entities;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IVehicleModelRepository : IRepository<VehicleModel>
    {
        Task<VehicleModel?> GetByBrandModelYearAsync(string brand, string modelName, int year);
        Task<bool> ExistsByBrandModelYearAsync(string brand, string modelName, int year, Guid? excludeModelId = null);
        Task<IEnumerable<VehicleModel>> GetByBrandAsync(string brand);
        Task<IEnumerable<VehicleModel>> GetByModelNameAsync(string modelName);
        Task<IEnumerable<VehicleModel>> GetByYearAsync(int year);
        Task<IEnumerable<string>> GetAllBrandsAsync();

        Task<bool> HasActiveVehiclesAsync(Guid modelId);
        Task<int> GetActiveVehicleCountAsync(Guid modelId);
        // *u: count feat later due to needed
    }
}
