using CustomerVehicleService.Application.Interfaces;
using CustomerVehicleService.Application.Interfaces.Data;
using CustomerVehicleService.Application.Services;
using CustomerVehicleService.Infrastructure.Data.Context;
using CustomerVehicleService.Infrastructure.Data.Repositories;
using CustomerVehicleService.Infrastructure.Data.UnitOfWork;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Options;


namespace CustomerVehicleService.Infrastructure
{
    public static class DependencyInjection
    {
        public static IServiceCollection AddInfrastructure(
            this IServiceCollection services,
            IConfiguration configuration)
        {
            services.AddDbContext<AppDbContext>(options => options.UseSqlServer(configuration.GetConnectionString("SqlServer")));

            //services.AddDbContext<AppDbContext>(options =>
            //{
            //    var connectionString = configuration.GetConnectionString("SqlServer")
            //        ?? throw new InvalidOperationException("Connection string 'SqlServer' not found.");

            //});

            services.AddScoped<IUnitOfWork, UnitOfWork>();

            //services.AddScoped<ICustomerRepository, CustomerRepository>();
            //services.AddScoped<IVehicleRepository, VehicleRepository>();
            //services.AddScoped<IVehicleModelRepository, VehicleModelRepository>();
            //services.AddScoped(typeof(IRepository<>), typeof(BaseRepository<>));

            services.AddScoped<ICustomerService, CustomerService>();
            services.AddScoped<IVehicleService, VehicleService>();
            services.AddScoped<IVehicleModelService, VehicleModelService>();

            return services;
        }
    }
}
