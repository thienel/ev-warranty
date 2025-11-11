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
    public class WorkOrderConfiguration : IEntityTypeConfiguration<WorkOrder>
    {
        public void Configure(EntityTypeBuilder<WorkOrder> builder) {
            builder.ToTable("work_orders");

            builder.HasKey(x => x.Id);
            builder.Property(x => x.Id)
                .HasColumnName("id")
                .IsRequired();

            builder.Property(x => x.ClaimId)
                .HasColumnName("claim_id")
                .IsRequired();

            builder.Property(x => x.AssignedTechnicianId)
                .HasColumnName("assigned_technician_id")
                .IsRequired();

            builder.Property(x => x.Status)
                .HasColumnName("status")
                .HasMaxLength(12)
                .HasConversion<string>()
                .IsRequired()
                .HasDefaultValue(WorkOrderStatus.Pending);

            builder.Property(x => x.ScheduledDate)
                .HasColumnName("scheduled_date")
                .HasColumnType("date")
                .IsRequired();

            builder.Property(x => x.CompletedDate)
                .HasColumnName("completed_date")
                .HasColumnType("date");

            builder.Property(x => x.Note)
                .HasColumnName("note")
                .HasColumnType("nvarchar(max)")
                .IsRequired(false);

            builder.Property(x => x.CreatedAt)
                .HasColumnName("created_at")
                .HasColumnType("datetime2")
                .IsRequired();

            builder.Property(x => x.UpdatedAt)
                .HasColumnName("updated_at")
                .HasColumnType("datetime2")
                .IsRequired(false);

            // Unique constraint on Claim
            builder.HasIndex(x => x.ClaimId)
                .IsUnique()
                .HasDatabaseName("ix_work_order_claim_id");

            builder.HasIndex(x => new { x.AssignedTechnicianId, x.Status })
                .HasDatabaseName("ix_work_order_technician_status");

            builder.HasIndex(x => x.Status)
                .HasDatabaseName("ix_work_order_status");

            builder.HasIndex(x => x.ScheduledDate)
                .HasDatabaseName("ix_work_order_scheduled_date");
        }
    }
}
