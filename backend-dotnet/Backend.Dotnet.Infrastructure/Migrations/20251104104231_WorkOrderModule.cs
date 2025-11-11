using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Dotnet.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class WorkOrderModule : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "work_orders",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    claim_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    assigned_technician_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    status = table.Column<string>(type: "nvarchar(12)", maxLength: 12, nullable: false, defaultValue: "Pending"),
                    scheduled_date = table.Column<DateTime>(type: "date", nullable: false),
                    completed_date = table.Column<DateTime>(type: "date", nullable: true),
                    note = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    created_at = table.Column<DateTime>(type: "datetime2", nullable: false),
                    updated_at = table.Column<DateTime>(type: "datetime2", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_work_orders", x => x.id);
                });

            migrationBuilder.CreateIndex(
                name: "ix_work_order_claim_id",
                table: "work_orders",
                column: "claim_id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_work_order_scheduled_date",
                table: "work_orders",
                column: "scheduled_date");

            migrationBuilder.CreateIndex(
                name: "ix_work_order_status",
                table: "work_orders",
                column: "status");

            migrationBuilder.CreateIndex(
                name: "ix_work_order_technician_status",
                table: "work_orders",
                columns: new[] { "assigned_technician_id", "status" });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "work_orders");
        }
    }
}
