using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Backend.Dotnet.Infrastructure.Data.Configurations
{
    public class VehicleConfiguration : IEntityTypeConfiguration<Vehicle>
    {
        public void Configure(EntityTypeBuilder<Vehicle> builder)
        {
            // Table name
            builder.ToTable("vehicles");

            // Primary key
            builder.HasKey(v => v.Id);
            builder.Property(v => v.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(v => v.Vin)
                .HasColumnName("vin")
                .HasColumnType("varchar(17)")
                .IsRequired();

            builder.Property(v => v.LicensePlate)
                .HasColumnName("license_plate")
                .HasColumnType("varchar(20)")
                .IsRequired(false);

            builder.Property(v => v.CustomerId)
                .HasColumnName("customer_id")
                .IsRequired();

            builder.Property(v => v.ModelId)
                .HasColumnName("model_id")
                .IsRequired();

            builder.Property(v => v.PurchaseDate)
                .HasColumnName("purchase_date")
                .HasColumnType("date")
                .IsRequired(false);

            // datetime2s
            builder.Property(v => v.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(v => v.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            builder.Property(v => v.DeletedAt)
                .HasColumnName("deleted_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            builder.Ignore(v => v.IsDeleted);

            // Indexes
            builder.HasIndex(v => v.Vin)
                .IsUnique()
                .HasDatabaseName("ix_vehicles_vin");

            builder.HasIndex(v => v.CustomerId)
                .HasDatabaseName("ix_vehicles_customer_id");

            builder.HasIndex(v => v.ModelId)
                .HasDatabaseName("ix_vehicles_model_id");

            builder.HasIndex(v => v.LicensePlate)
                .HasDatabaseName("ix_vehicles_license_plate");

            builder.HasIndex(v => v.DeletedAt);

            // Relationships
            builder.HasOne(v => v.Customer)
                .WithMany(c => c.Vehicles)
                .HasForeignKey(v => v.CustomerId)
                .OnDelete(DeleteBehavior.Restrict);

            builder.HasOne(v => v.Model)
                .WithMany(vm => vm.Vehicles)
                .HasForeignKey(v => v.ModelId)
                .OnDelete(DeleteBehavior.Restrict);

            builder.HasQueryFilter(v => v.DeletedAt == null);
        }
    }
}
