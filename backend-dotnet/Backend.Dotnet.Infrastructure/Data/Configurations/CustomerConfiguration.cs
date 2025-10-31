using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Backend.Dotnet.Infrastructure.Data.Configurations
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
                .HasMaxLength(100)
                .IsRequired();

            builder.Property(c => c.LastName)
                .HasColumnName("last_name")
                .HasMaxLength(100)
                .IsRequired();

            builder.Property(c => c.PhoneNumber)
                .HasColumnName("phone_number")
                .HasMaxLength(20)
                .IsRequired();

            builder.Property(c => c.Email)
                .HasColumnName("email")
                .HasMaxLength(255)
                .IsRequired(false);

            builder.Property(c => c.Address)
                .HasColumnName("address")
                .HasColumnType("nvarchar(max)")
                .IsRequired(false);

            // datetime2s
            builder.Property(c => c.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(c => c.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            builder.Property(c => c.DeletedAt)
                .HasColumnName("deleted_at")
                .HasColumnType("datetime2")
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
