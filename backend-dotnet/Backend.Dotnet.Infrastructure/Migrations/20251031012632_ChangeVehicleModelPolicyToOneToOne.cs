using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Dotnet.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class ChangeVehicleModelPolicyToOneToOne : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "ix_vehicle_models_policy_id",
                table: "vehicle_models");

            migrationBuilder.AddColumn<Guid>(
                name: "AssignedModelId",
                table: "warranty_policies",
                type: "uniqueidentifier",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_warranty_policies_AssignedModelId",
                table: "warranty_policies",
                column: "AssignedModelId");

            migrationBuilder.CreateIndex(
                name: "ix_vehicle_models_policy_id",
                table: "vehicle_models",
                column: "policy_id",
                unique: true,
                filter: "policy_id IS NOT NULL");

            migrationBuilder.AddForeignKey(
                name: "FK_warranty_policies_vehicle_models_AssignedModelId",
                table: "warranty_policies",
                column: "AssignedModelId",
                principalTable: "vehicle_models",
                principalColumn: "id");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_warranty_policies_vehicle_models_AssignedModelId",
                table: "warranty_policies");

            migrationBuilder.DropIndex(
                name: "IX_warranty_policies_AssignedModelId",
                table: "warranty_policies");

            migrationBuilder.DropIndex(
                name: "ix_vehicle_models_policy_id",
                table: "vehicle_models");

            migrationBuilder.DropColumn(
                name: "AssignedModelId",
                table: "warranty_policies");

            migrationBuilder.CreateIndex(
                name: "ix_vehicle_models_policy_id",
                table: "vehicle_models",
                column: "policy_id");
        }
    }
}
