namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IUnitOfWork
    {
        ICustomerRepository Customers { get; }
        IVehicleRepository Vehicles { get; }
        IVehicleModelRepository VehicleModels { get; }
        IWarrantyPolicyRepository WarrantyPolicies { get; }
        IPartCategoryRepository PartCategories { get; }
        IPartRepository Parts { get; }
        IPolicyCoveragePartRepository PolicyCoverageParts { get; }
        IWorkOrderRepository WorkOrderRepository { get; }

        Task<int> SaveChangesAsync();
        int SaveChanges();
        Task BeginTransactionAsync();
        Task CommitTransactionAsync();
        Task RollbackTransactionAsync();
    }
}
