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
    public class PolicyCoveragePartConfiguration : IEntityTypeConfiguration<PolicyCoveragePart>
    {
        public void Configure(EntityTypeBuilder<PolicyCoveragePart> builder)
        {
            builder.ToTable("policy_coverage_parts");

            // Primary Key
            builder.HasKey(pcp => pcp.Id);
            builder.Property(pcp => pcp.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(pcp => pcp.PolicyId)
                .HasColumnName("policy_id")
                .IsRequired();

            builder.Property(pcp => pcp.PartCategoryId)
                .HasColumnName("part_category_id")
                .IsRequired();

            builder.Property(pcp => pcp.CoverageConditions)
                .HasColumnName("coverage_conditions")
                .HasColumnType("text")
                .IsRequired(false);

            builder.Property(pcp => pcp.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(pcp => pcp.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            // Indexes
            builder.HasIndex(pcp => pcp.PolicyId)
                .HasDatabaseName("ix_policy_coverage_parts_policy_id");

            builder.HasIndex(pcp => pcp.PartCategoryId)
                .HasDatabaseName("ix_policy_coverage_parts_part_category_id");

            // Unique constraint: one policy can only cover each category once
            builder.HasIndex(pcp => new { pcp.PolicyId, pcp.PartCategoryId })
                .IsUnique()
                .HasDatabaseName("ix_policy_coverage_parts_policy_category_unique");

            // Relationships
            builder.HasOne(pcp => pcp.Policy)
                .WithMany(wp => wp.CoverageParts)
                .HasForeignKey(pcp => pcp.PolicyId)
                .OnDelete(DeleteBehavior.Cascade)
                .IsRequired();

            builder.HasOne(pcp => pcp.PartCategory)
                .WithMany(pc => pc.PolicyCoverageParts)
                .HasForeignKey(pcp => pcp.PartCategoryId)
                .OnDelete(DeleteBehavior.Restrict)
                .IsRequired();
        }
    }
}
