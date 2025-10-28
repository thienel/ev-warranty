using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Dotnet.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class UpdateInventoryPolicyModule : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_warranty_policies_vehicle_models_model_id",
                table: "warranty_policies");

            migrationBuilder.DropIndex(
                name: "ix_warranty_policies_model_id",
                table: "warranty_policies");

            migrationBuilder.DropIndex(
                name: "ix_warranty_policies_model_status",
                table: "warranty_policies");

            migrationBuilder.DropColumn(
                name: "model_id",
                table: "warranty_policies");

            migrationBuilder.AddColumn<Guid>(
                name: "policy_id",
                table: "vehicle_models",
                type: "uniqueidentifier",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "ix_vehicle_models_policy_id",
                table: "vehicle_models",
                column: "policy_id");

            migrationBuilder.AddForeignKey(
                name: "FK_vehicle_models_warranty_policies_policy_id",
                table: "vehicle_models",
                column: "policy_id",
                principalTable: "warranty_policies",
                principalColumn: "id",
                onDelete: ReferentialAction.SetNull);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_vehicle_models_warranty_policies_policy_id",
                table: "vehicle_models");

            migrationBuilder.DropIndex(
                name: "ix_vehicle_models_policy_id",
                table: "vehicle_models");

            migrationBuilder.DropColumn(
                name: "policy_id",
                table: "vehicle_models");

            migrationBuilder.AddColumn<Guid>(
                name: "model_id",
                table: "warranty_policies",
                type: "uniqueidentifier",
                nullable: false,
                defaultValue: new Guid("00000000-0000-0000-0000-000000000000"));

            migrationBuilder.CreateIndex(
                name: "ix_warranty_policies_model_id",
                table: "warranty_policies",
                column: "model_id");

            migrationBuilder.CreateIndex(
                name: "ix_warranty_policies_model_status",
                table: "warranty_policies",
                columns: new[] { "model_id", "status" });

            migrationBuilder.AddForeignKey(
                name: "FK_warranty_policies_vehicle_models_model_id",
                table: "warranty_policies",
                column: "model_id",
                principalTable: "vehicle_models",
                principalColumn: "id",
                onDelete: ReferentialAction.Restrict);
        }
    }
}
