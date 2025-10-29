using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Backend.Dotnet.Infrastructure.Data.Configurations
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

            builder.Property(vm => vm.PolicyId)
                .HasColumnName("policy_id")
                .IsRequired(false);

            // DateTime
            builder.Property(vm => vm.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(vm => vm.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            // Indexes - Composite unique index
            builder.HasIndex(vm => new { vm.Brand, vm.ModelName, vm.Year })
                .IsUnique()
                .HasDatabaseName("ix_vehicle_models_brand_model_year");

            builder.HasIndex(vm => vm.PolicyId)
               .HasDatabaseName("ix_vehicle_models_policy_id");

            // Relationships
            builder.HasMany(vm => vm.Vehicles)
                .WithOne(v => v.Model)
                .HasForeignKey(v => v.ModelId)
                .OnDelete(DeleteBehavior.Restrict);

            builder.HasOne(vm => vm.Policy)
               .WithMany(wp => wp.VehicleModels)
               .HasForeignKey(vm => vm.PolicyId)
               .OnDelete(DeleteBehavior.SetNull)
               .IsRequired(false);
        }
    }
}
