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
    public class CustomerConfiguration : IEntityTypeConfiguration<Customer>
    {
        public void Configure(EntityTypeBuilder<Customer> builder)
        {
            // Table name
            builder.ToTable("customers");

            // Primary key
            builder.HasKey(c => c.Id);
            builder.Property(c => c.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(c => c.FirstName)
                .HasColumnName("first_name")
                .HasColumnType("varchar(100)")
                .IsRequired();

            builder.Property(c => c.LastName)
                .HasColumnName("last_name")
                .HasColumnType("varchar(100)")
                .IsRequired();

            builder.Property(c => c.PhoneNumber)
                .HasColumnName("phone_number")
                .HasColumnType("varchar(20)")
                .IsRequired();

            builder.Property(c => c.Email)
                .HasColumnName("email")
                .HasColumnType("varchar(255)")
                .IsRequired(false);

            builder.Property(c => c.Address)
                .HasColumnName("address")
                .HasColumnType("text")
                .IsRequired(false);

            // Timestamps
            builder.Property(c => c.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("timestamp")
                .IsRequired();

            builder.Property(c => c.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("timestamp")
                .IsRequired(false);

            builder.Property(c => c.DeletedAt)
                .HasColumnName("deleted_at")
                .HasColumnType("timestamp")
                .IsRequired(false);

            // Computed column (not mapped to database)
            builder.Ignore(c => c.FullName);
            builder.Ignore(c => c.IsDeleted);

            // Indexes
            builder.HasIndex(c => c.Email)
                .IsUnique()
                .HasFilter("deleted_at IS NULL"); // Unique only for non-deleted records

            builder.HasIndex(c => c.DeletedAt);

            // Relationships
            builder.HasMany(c => c.Vehicles)
                .WithOne(v => v.Customer)
                .HasForeignKey(v => v.CustomerId)
                .OnDelete(DeleteBehavior.Restrict);

            // Global query filter for soft delete
            builder.HasQueryFilter(c => c.DeletedAt == null);
        }
    }
}
