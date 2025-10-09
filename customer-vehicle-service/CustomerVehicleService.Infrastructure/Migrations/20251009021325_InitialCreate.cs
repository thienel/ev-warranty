using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace CustomerVehicleService.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "customers",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    first_name = table.Column<string>(type: "varchar(100)", nullable: false),
                    last_name = table.Column<string>(type: "varchar(100)", nullable: false),
                    phone_number = table.Column<string>(type: "varchar(20)", nullable: false),
                    email = table.Column<string>(type: "varchar(255)", nullable: true),
                    address = table.Column<string>(type: "text", nullable: true),
                    deleted_at = table.Column<byte[]>(type: "timestamp", nullable: true),
                    created_at = table.Column<byte[]>(type: "timestamp", nullable: false),
                    updated_at = table.Column<byte[]>(type: "timestamp", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_customers", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "vehicle_models",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    brand = table.Column<string>(type: "varchar(100)", nullable: false),
                    model_name = table.Column<string>(type: "varchar(100)", nullable: false),
                    year = table.Column<int>(type: "integer", nullable: false),
                    created_at = table.Column<byte[]>(type: "timestamp", nullable: false),
                    updated_at = table.Column<byte[]>(type: "timestamp", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_vehicle_models", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "vehicles",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    vin = table.Column<string>(type: "varchar(17)", nullable: false),
                    license_plate = table.Column<string>(type: "varchar(20)", nullable: true),
                    customer_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    model_id = table.Column<Guid>(type: "uniqueidentifier", nullable: false),
                    purchase_date = table.Column<DateTime>(type: "date", nullable: true),
                    created_at = table.Column<byte[]>(type: "timestamp", nullable: false),
                    updated_at = table.Column<byte[]>(type: "timestamp", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_vehicles", x => x.id);
                    table.ForeignKey(
                        name: "FK_vehicles_customers_customer_id",
                        column: x => x.customer_id,
                        principalTable: "customers",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_vehicles_vehicle_models_model_id",
                        column: x => x.model_id,
                        principalTable: "vehicle_models",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateIndex(
                name: "IX_customers_deleted_at",
                table: "customers",
                column: "deleted_at");

            migrationBuilder.CreateIndex(
                name: "IX_customers_email",
                table: "customers",
                column: "email",
                unique: true,
                filter: "deleted_at IS NULL");

            migrationBuilder.CreateIndex(
                name: "ix_vehicle_models_brand_model_year",
                table: "vehicle_models",
                columns: new[] { "brand", "model_name", "year" },
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_vehicles_customer_id",
                table: "vehicles",
                column: "customer_id");

            migrationBuilder.CreateIndex(
                name: "ix_vehicles_license_plate",
                table: "vehicles",
                column: "license_plate");

            migrationBuilder.CreateIndex(
                name: "ix_vehicles_model_id",
                table: "vehicles",
                column: "model_id");

            migrationBuilder.CreateIndex(
                name: "ix_vehicles_vin",
                table: "vehicles",
                column: "vin",
                unique: true);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "vehicles");

            migrationBuilder.DropTable(
                name: "customers");

            migrationBuilder.DropTable(
                name: "vehicle_models");
        }
    }
}
