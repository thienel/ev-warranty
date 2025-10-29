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
    public class PartCategoryConfiguration : IEntityTypeConfiguration<PartCategory>
    {
        public void Configure(EntityTypeBuilder<PartCategory> builder)
        {
            builder.ToTable("part_categories");

            // Primary Key
            builder.HasKey(pc => pc.Id);
            builder.Property(pc => pc.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(pc => pc.CategoryName)
                .HasColumnName("category_name")
                .HasColumnType("varchar(255)")
                .IsRequired();

            builder.Property(pc => pc.Description)
                .HasColumnName("description")
                .HasColumnType("text")
                .IsRequired(false);

            builder.Property(pc => pc.ParentCategoryId)
                .HasColumnName("parent_category_id")
                .IsRequired(false);

            builder.Property(pc => pc.Status)
                .HasColumnName("status")
                .HasColumnType("varchar(20)")
                .IsRequired()
                .HasConversion<string>()
                .HasDefaultValue(PartCategoryStatus.Active);

            builder.Property(pc => pc.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(pc => pc.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            // Indexes
            builder.HasIndex(pc => pc.CategoryName)
                .IsUnique()
                .HasDatabaseName("ix_part_categories_category_name");

            builder.HasIndex(pc => pc.ParentCategoryId)
                .HasDatabaseName("ix_part_categories_parent_category_id");

            builder.HasIndex(pc => pc.Status)
                .HasDatabaseName("ix_part_categories_status");

            // Self-referencing relationship
            builder.HasOne(pc => pc.ParentCategory)
                .WithMany(pc => pc.ChildCategories)
                .HasForeignKey(pc => pc.ParentCategoryId)
                .OnDelete(DeleteBehavior.Restrict)
                .IsRequired(false);

            // Relationships
            builder.HasMany(pc => pc.Parts)
                .WithOne(p => p.Category)
                .HasForeignKey(p => p.CategoryId)
                .OnDelete(DeleteBehavior.Restrict);

            builder.HasMany(pc => pc.PolicyCoverageParts)
                .WithOne(pcp => pcp.PartCategory)
                .HasForeignKey(pcp => pcp.PartCategoryId)
                .OnDelete(DeleteBehavior.Restrict);
        }
    }
}
