using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IPartCategoryRepository : IRepository<PartCategory>
    {
        Task<PartCategory?> GetByCategoryNameAsync(string categoryName);
        Task<bool> CategoryNameExistsAsync(string categoryName, Guid? excludeCategoryId = null);

        Task<IEnumerable<PartCategory>> GetByStatusAsync(PartCategoryStatus status);

        Task<IEnumerable<PartCategory>> GetByParentIdAsync(Guid parentCategoryId); // Child list
        Task<PartCategory?> GetWithPartsAsync(Guid categoryId); // Get details of belonging part
        Task<PartCategory?> GetWithHierarchyAsync(Guid categoryId); // Get details of category
        Task<IEnumerable<PartCategory>> GetFullHierarchyAsync(); // Get all category

        Task<int> GetActivePartCountAsync(Guid categoryId);
        Task<int> GetChildCategoryCountAsync(Guid categoryId);
    }
}
