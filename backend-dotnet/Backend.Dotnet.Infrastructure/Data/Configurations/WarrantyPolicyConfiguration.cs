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
    public class WarrantyPolicyConfiguration : IEntityTypeConfiguration<WarrantyPolicy>
    {
        public void Configure(EntityTypeBuilder<WarrantyPolicy> builder)
        {
            builder.ToTable("warranty_policies");

            // Primary Key
            builder.HasKey(wp => wp.Id);
            builder.Property(wp => wp.Id)
                .HasColumnName("id")
                .IsRequired();

            // Properties
            builder.Property(wp => wp.PolicyName)
                .HasColumnName("policy_name")
                .HasColumnType("varchar(255)")
                .IsRequired();

            builder.Property(wp => wp.ModelId)
                .HasColumnName("model_id")
                .IsRequired();

            builder.Property(wp => wp.WarrantyDurationMonths)
                .HasColumnName("warranty_duration_months")
                .HasColumnType("int")
                .IsRequired();

            builder.Property(wp => wp.KilometerLimit)
                .HasColumnName("kilometer_limit")
                .HasColumnType("int")
                .IsRequired(false);

            builder.Property(wp => wp.TermsAndConditions)
                .HasColumnName("terms_and_conditions")
                .HasColumnType("text")
                .IsRequired();

            builder.Property(wp => wp.Status)
                .HasColumnName("status")
                .HasColumnType("varchar(50)")
                .IsRequired()
                .HasConversion<string>()
                .HasDefaultValue(WarrantyPolicyStatus.Draft);

            builder.Property(wp => wp.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(wp => wp.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            // Indexes
            builder.HasIndex(wp => wp.PolicyName)
                .HasDatabaseName("ix_warranty_policies_policy_name");

            builder.HasIndex(wp => wp.ModelId)
                .HasDatabaseName("ix_warranty_policies_model_id");

            builder.HasIndex(wp => wp.Status)
                .HasDatabaseName("ix_warranty_policies_status");

            // Composite index for common queries (one active policy per model)
            builder.HasIndex(wp => new { wp.ModelId, wp.Status })
                .HasDatabaseName("ix_warranty_policies_model_status");

            // Relationships
            builder.HasOne(wp => wp.Model)
                .WithMany()
                .HasForeignKey(wp => wp.ModelId)
                .OnDelete(DeleteBehavior.Restrict)
                .IsRequired();

            builder.HasMany(wp => wp.CoverageParts)
                .WithOne(pcp => pcp.Policy)
                .HasForeignKey(pcp => pcp.PolicyId)
                .OnDelete(DeleteBehavior.Cascade);
        }
    }
}
