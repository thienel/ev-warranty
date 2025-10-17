using Backend.Dotnet.Domain.Entities;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IVehicleModelRepository : IRepository<VehicleModel>
    {
        // VehicleModel-specific queries
        Task<VehicleModel?> GetByBrandModelYearAsync(string brand, string modelName, int year);
        Task<bool> ExistsByBrandModelYearAsync(string brand, string modelName, int year, Guid? excludeModelId = null);
        Task<IEnumerable<VehicleModel>> GetByBrandAsync(string brand);
        Task<IEnumerable<VehicleModel>> SearchAsync(string searchTerm);
        Task<IEnumerable<string>> GetAllBrandsAsync();

        // Hard delete assurance check
        Task<bool> HasActiveVehiclesAsync(Guid modelId);

        Task<int> GetActiveVehicleCountAsync(Guid modelId);
        // thong nhat su dung truy van xe active/inactive - non-deleted/deleted
        // VehicleRepo: GetByModelIdAsync GetDeletedVehicleAsync GetByIdIncludingDeletedAsync
        //Task<IEnumerable<Vehicle>> GetActiveVehicleAsync(Guid modelId);
        //Task<IEnumerable<VehicleModel>> GetInactiveVehicleAsync(Guid modelId); // relate to restore procudure

        // *u: count feat later due to needed
    }
}
