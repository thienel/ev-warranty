using CustomerVehicleService.Domain.Abstractions;
using CustomerVehicleService.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Infrastructure.Data.Context
{
    public class CustomerVehicleDbContext : DbContext
    {
        public DbSet<Customer> Customers { get; set; }
        public DbSet<Vehicle> Vehicles { get; set; }
        public DbSet<VehicleModel> VehicleModels { get; set; }

        public CustomerVehicleDbContext(DbContextOptions<CustomerVehicleDbContext> options)
            : base(options)
        {
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);

            modelBuilder.ApplyConfigurationsFromAssembly(typeof(CustomerVehicleDbContext).Assembly);
        }

        public override async Task<int> SaveChangesAsync(CancellationToken cancellationToken = default)
        {
            UpdateTimestamps();
            return await base.SaveChangesAsync(cancellationToken);
        }

        public override int SaveChanges()
        {
            UpdateTimestamps();
            return base.SaveChanges();
        }

        private void UpdateTimestamps()
        {
            var entries = ChangeTracker.Entries<BaseEntity>();
            var now = DateTime.UtcNow;

            foreach (var entry in entries)
            {
                switch (entry.State)
                {
                    case EntityState.Added:
                        // Use reflection to set protected properties
                        SetPropertyValue(entry.Entity, nameof(BaseEntity.CreatedAt), now);
                        SetPropertyValue(entry.Entity, nameof(BaseEntity.UpdatedAt), now);
                        break;

                    case EntityState.Modified:
                        SetPropertyValue(entry.Entity, nameof(BaseEntity.UpdatedAt), now);
                        // Prevent overwriting CreatedAt
                        entry.Property(nameof(BaseEntity.CreatedAt)).IsModified = false;
                        break;
                }
            }
        }

        private void SetPropertyValue(object obj, string propertyName, object value)
        {
            var property = obj.GetType()
                .GetProperty(propertyName, BindingFlags.Public | BindingFlags.Instance);

            if (property != null && property.CanWrite)
            {
                property.SetValue(obj, value);
            }
        }
    }
}
