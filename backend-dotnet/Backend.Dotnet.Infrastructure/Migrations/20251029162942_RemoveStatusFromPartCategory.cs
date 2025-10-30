using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Dotnet.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class RemoveStatusFromPartCategory : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "ix_part_categories_status",
                table: "part_categories");

            migrationBuilder.DropColumn(
                name: "status",
                table: "part_categories");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<string>(
                name: "status",
                table: "part_categories",
                type: "varchar(20)",
                nullable: false,
                defaultValue: "Active");

            migrationBuilder.CreateIndex(
                name: "ix_part_categories_status",
                table: "part_categories",
                column: "status");
        }
    }
}
