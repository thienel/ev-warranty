namespace CustomerVehicleService.Application.Interfaces.Data
{
    public interface IUnitOfWork
    {
        ICustomerRepository Customers { get; }
        IVehicleRepository Vehicles { get; }
        IVehicleModelRepository VehicleModels { get; }

        Task<int> SaveChangesAsync();
        int SaveChanges();
        Task BeginTransactionAsync();
        Task CommitTransactionAsync();
        Task RollbackTransactionAsync();
    }
}
