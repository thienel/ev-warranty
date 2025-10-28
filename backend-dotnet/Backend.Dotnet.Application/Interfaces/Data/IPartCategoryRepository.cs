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

        Task<IEnumerable<PartCategory>> GetByParentIdAsync(Guid parentCategoryId);
        Task<IEnumerable<PartCategory>> GetRootCategoriesAsync();
        Task<PartCategory?> GetWithHierarchyAsync(Guid categoryId);
        Task<IEnumerable<PartCategory>> GetFullHierarchyAsync();

        Task<PartCategory?> GetWithParentAsync(Guid categoryId);
        Task<PartCategory?> GetWithChildrenAsync(Guid categoryId);
        Task<PartCategory?> GetWithPartsAsync(Guid categoryId);

        Task<bool> HasActivePartsAsync(Guid categoryId);
        Task<bool> HasActiveChildrenAsync(Guid categoryId);
        Task<int> GetActivePartCountAsync(Guid categoryId);
        Task<int> GetChildCategoryCountAsync(Guid categoryId);
        Task<bool> CanBeUsedForNewPartsAsync(Guid categoryId);
        Task<bool> IsDescendantOfAsync(Guid categoryId, Guid potentialAncestorId);
    }
}
