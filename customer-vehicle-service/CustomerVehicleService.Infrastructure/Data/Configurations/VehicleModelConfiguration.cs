using CustomerVehicleService.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Infrastructure.Data.Configurations
{
    public class VehicleModelConfiguration : IEntityTypeConfiguration<VehicleModel>
    {
        public void Configure(EntityTypeBuilder<VehicleModel> builder)
        {
            // Table name
            builder.ToTable("vehicle_models");

            // Primary key
            builder.HasKey(vm => vm.Id);
            builder.Property(vm => vm.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(vm => vm.Brand)
                .HasColumnName("brand")
                .HasColumnType("varchar(100)")
                .IsRequired();

            builder.Property(vm => vm.ModelName)
                .HasColumnName("model_name")
                .HasColumnType("varchar(100)")
                .IsRequired();

            builder.Property(vm => vm.Year)
                .HasColumnName("year")
                .HasColumnType("integer")
                .IsRequired();

            // Timestamps
            builder.Property(vm => vm.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("timestamp")
                .IsRequired();

            builder.Property(vm => vm.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("timestamp")
                .IsRequired(false);

            // Indexes - Composite unique index
            builder.HasIndex(vm => new { vm.Brand, vm.ModelName, vm.Year })
                .IsUnique()
                .HasDatabaseName("ix_vehicle_models_brand_model_year");

            // Relationships
            builder.HasMany(vm => vm.Vehicles)
                .WithOne(v => v.Model)
                .HasForeignKey(v => v.ModelId)
                .OnDelete(DeleteBehavior.Restrict);
        }
    }
}
