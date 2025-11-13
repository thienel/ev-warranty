using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Infrastructure.Data.Configurations
{
    public class PartConfiguration : IEntityTypeConfiguration<Part>
    {
        public void Configure(EntityTypeBuilder<Part> builder)
        {
            builder.ToTable("parts");

            // Primary Key
            builder.HasKey(p => p.Id);
            builder.Property(p => p.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(p => p.SerialNumber)
                .HasColumnName("serial_number")
                .HasMaxLength(255)
                .IsRequired();

            builder.Property(p => p.PartName)
                .HasColumnName("part_name")
                .HasMaxLength(255)
                .IsRequired();

            builder.Property(p => p.UnitPrice)
                .HasColumnName("unit_price")
                .HasColumnType("decimal(18,2)")
                .IsRequired();

            builder.Property(p => p.CategoryId)
                .HasColumnName("category_id")
                .IsRequired();

            builder.Property(p => p.OfficeLocationId)
                .HasColumnName("office_location_id")
                .IsRequired(false);

            builder.Property(p => p.Status)
                .HasColumnName("status")
                .HasMaxLength(50)
                .IsRequired()
                .HasConversion<string>()
                .HasDefaultValue(PartStatus.Available);

            builder.Property(p => p.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(p => p.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            // Indexes
            builder.HasIndex(p => p.SerialNumber)
                .IsUnique()
                .HasDatabaseName("ix_parts_serial_number");

            builder.HasIndex(p => p.CategoryId)
                .HasDatabaseName("ix_parts_category_id");

            builder.HasIndex(p => p.OfficeLocationId)
                .HasDatabaseName("ix_parts_office_location_id");

            builder.HasIndex(p => p.Status)
                .HasDatabaseName("ix_parts_status");

            // Index
            builder.HasIndex(p => new { p.CategoryId, p.Status })
                .HasDatabaseName("ix_parts_category_status");

            // Relationships
            builder.HasOne(p => p.Category)
                .WithMany(pc => pc.Parts)
                .HasForeignKey(p => p.CategoryId)
                .OnDelete(DeleteBehavior.Restrict)
                .IsRequired();
        }
    }
}
