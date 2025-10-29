using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Dotnet.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class InventoryPolicyModule : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "part_categories",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    category_name = table.Column<string>(type: "varchar(255)", nullable: false),
                    description = table.Column<string>(type: "text", nullable: true),
                    parent_category_id = table.Column<Guid>(type: "uniqueidentifier", nullable: true),
                    status = table.Column<string>(type: "varchar(20)", nullable: false, defaultValue: "Active"),
                    created_at = table.Column<DateTime>(type: "datetime2", nullable: false),
                    updated_at = table.Column<DateTime>(type: "datetime2", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_part_categories", x => x.id);
                    table.ForeignKey(
                        name: "FK_part_categories_part_categories_parent_category_id",
                        column: x => x.parent_category_id,
                        principalTable: "part_categories",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateTable(
                name: "warranty_policies",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    policy_name = table.Column<string>(type: "varchar(255)", nullable: false),
                    model_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    warranty_duration_months = table.Column<int>(type: "int", nullable: false),
                    kilometer_limit = table.Column<int>(type: "int", nullable: true),
                    terms_and_conditions = table.Column<string>(type: "text", nullable: false),
                    status = table.Column<string>(type: "varchar(50)", nullable: false, defaultValue: "Draft"),
                    created_at = table.Column<DateTime>(type: "datetime2", nullable: false),
                    updated_at = table.Column<DateTime>(type: "datetime2", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_warranty_policies", x => x.id);
                    table.ForeignKey(
                        name: "FK_warranty_policies_vehicle_models_model_id",
                        column: x => x.model_id,
                        principalTable: "vehicle_models",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateTable(
                name: "parts",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    serial_number = table.Column<string>(type: "varchar(255)", nullable: false),
                    part_name = table.Column<string>(type: "varchar(255)", nullable: false),
                    unit_price = table.Column<decimal>(type: "decimal(10,2)", nullable: false),
                    category_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    office_location_id = table.Column<Guid>(type: "uniqueidentifier", nullable: true),
                    status = table.Column<string>(type: "varchar(50)", nullable: false, defaultValue: "Available"),
                    created_at = table.Column<DateTime>(type: "datetime2", nullable: false),
                    updated_at = table.Column<DateTime>(type: "datetime2", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_parts", x => x.id);
                    table.ForeignKey(
                        name: "FK_parts_part_categories_category_id",
                        column: x => x.category_id,
                        principalTable: "part_categories",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateTable(
                name: "policy_coverage_parts",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    policy_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    part_category_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    coverage_conditions = table.Column<string>(type: "text", nullable: true),
                    created_at = table.Column<DateTime>(type: "datetime2", nullable: false),
                    updated_at = table.Column<DateTime>(type: "datetime2", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_policy_coverage_parts", x => x.id);
                    table.ForeignKey(
                        name: "FK_policy_coverage_parts_part_categories_part_category_id",
                        column: x => x.part_category_id,
                        principalTable: "part_categories",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_policy_coverage_parts_warranty_policies_policy_id",
                        column: x => x.policy_id,
                        principalTable: "warranty_policies",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "ix_part_categories_category_name",
                table: "part_categories",
                column: "category_name",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_part_categories_parent_category_id",
                table: "part_categories",
                column: "parent_category_id");

            migrationBuilder.CreateIndex(
                name: "ix_part_categories_status",
                table: "part_categories",
                column: "status");

            migrationBuilder.CreateIndex(
                name: "ix_parts_category_id",
                table: "parts",
                column: "category_id");

            migrationBuilder.CreateIndex(
                name: "ix_parts_category_status",
                table: "parts",
                columns: new[] { "category_id", "status" });

            migrationBuilder.CreateIndex(
                name: "ix_parts_office_location_id",
                table: "parts",
                column: "office_location_id");

            migrationBuilder.CreateIndex(
                name: "ix_parts_serial_number",
                table: "parts",
                column: "serial_number",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_parts_status",
                table: "parts",
                column: "status");

            migrationBuilder.CreateIndex(
                name: "ix_policy_coverage_parts_part_category_id",
                table: "policy_coverage_parts",
                column: "part_category_id");

            migrationBuilder.CreateIndex(
                name: "ix_policy_coverage_parts_policy_category_unique",
                table: "policy_coverage_parts",
                columns: new[] { "policy_id", "part_category_id" },
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_policy_coverage_parts_policy_id",
                table: "policy_coverage_parts",
                column: "policy_id");

            migrationBuilder.CreateIndex(
                name: "ix_warranty_policies_model_id",
                table: "warranty_policies",
                column: "model_id");

            migrationBuilder.CreateIndex(
                name: "ix_warranty_policies_model_status",
                table: "warranty_policies",
                columns: new[] { "model_id", "status" });

            migrationBuilder.CreateIndex(
                name: "ix_warranty_policies_policy_name",
                table: "warranty_policies",
                column: "policy_name");

            migrationBuilder.CreateIndex(
                name: "ix_warranty_policies_status",
                table: "warranty_policies",
                column: "status");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "parts");

            migrationBuilder.DropTable(
                name: "policy_coverage_parts");

            migrationBuilder.DropTable(
                name: "part_categories");

            migrationBuilder.DropTable(
                name: "warranty_policies");
        }
    }
}
