using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Interfaces.External;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Infrastructure.Data.Context;
using Backend.Dotnet.Infrastructure.Data.UnitOfWork;
using Backend.Dotnet.Infrastructure.External.Clients;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;


namespace Backend.Dotnet.Infrastructure
{
    public static class DependencyInjection
    {
        public static IServiceCollection AddInfrastructure(
            this IServiceCollection services,
            IConfiguration configuration)
        {
            services.AddDbContext<AppDbContext>(options =>
                options.UseSqlServer(configuration.GetConnectionString("SqlServer")));

            services.AddScoped<IUnitOfWork, UnitOfWork>();

            services.AddScoped<ICustomerService, CustomerService>();
            services.AddScoped<IVehicleService, VehicleService>();
            services.AddScoped<IVehicleModelService, VehicleModelService>();
            services.AddScoped<IWarrantyPolicyService, WarrantyPolicyService>();
            services.AddScoped<IPartCategoryService, PartCategoryService>();
            services.AddScoped<IPartService, PartService>();
            services.AddScoped<IPolicyCoveragePartService, PolicyCoveragePartService>();
            services.AddScoped<IWorkOrderService, WorkOrderService>();

            services.AddHttpClient<IExternalServiceClient, ExternalServiceClient>();

            return services;
        }
    }
}
