using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Backend.Dotnet.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class RemoveConstaintPolicyCoveragePart : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "ix_policy_coverage_parts_policy_category_unique",
                table: "policy_coverage_parts");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateIndex(
                name: "ix_policy_coverage_parts_policy_category_unique",
                table: "policy_coverage_parts",
                columns: new[] { "policy_id", "part_category_id" },
                unique: true);
        }
    }
}
