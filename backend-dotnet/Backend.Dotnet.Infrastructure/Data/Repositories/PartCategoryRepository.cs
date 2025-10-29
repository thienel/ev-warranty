using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Infrastructure.Data.Context;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class PartCategoryRepository : BaseRepository<PartCategory>, IPartCategoryRepository
    {
        public PartCategoryRepository(DbContext context) : base(context) { }

        public async Task<PartCategory?> GetByCategoryNameAsync(string categoryName)
        {
            return await _dbSet
                .FirstOrDefaultAsync(pc => pc.CategoryName.ToLower() == categoryName.ToLower());
        }

        public async Task<bool> CategoryNameExistsAsync(string categoryName, Guid? excludeCategoryId = null)
        {
            var query = _dbSet.Where(pc => pc.CategoryName.ToLower() == categoryName.ToLower());

            if (excludeCategoryId.HasValue)
                query = query.Where(pc => pc.Id != excludeCategoryId.Value);

            return await query.AnyAsync();
        }

        public async Task<IEnumerable<PartCategory>> GetByParentIdAsync(Guid parentCategoryId)
        {
            return await _dbSet
                .Where(pc => pc.ParentCategoryId == parentCategoryId)
                .OrderBy(pc => pc.CategoryName)
                .ToListAsync();
        }

        public async Task<PartCategory?> GetWithPartsAsync(Guid categoryId)
        {
            return await _dbSet
                .Include(pc => pc.Parts)
                .FirstOrDefaultAsync(pc => pc.Id == categoryId);
        }

        public async Task<PartCategory?> GetWithHierarchyAsync(Guid categoryId)
        {
            return await _dbSet
                .Include(pc => pc.Parts)
                .Include(pc => pc.ChildCategories)
                    .ThenInclude(child => child.Parts)
                .Include(pc => pc.ChildCategories)
                    .ThenInclude(child => child.ChildCategories)
                .FirstOrDefaultAsync(pc => pc.Id == categoryId);
        }

        public async Task<IEnumerable<PartCategory>> GetFullHierarchyAsync()
        {
            return await _dbSet
                .Include(pc => pc.Parts)
                .Include(pc => pc.ChildCategories)
                    .ThenInclude(child => child.Parts)
                .Include(pc => pc.ChildCategories)
                    .ThenInclude(child => child.ChildCategories)
                .OrderBy(pc => pc.CategoryName)
                .ToListAsync();
        }

        public async Task<int> GetActivePartCountAsync(Guid categoryId)
        {
            return await _context.Set<Part>()
                .CountAsync(p => p.CategoryId == categoryId && p.Status == PartStatus.Available);
        }

        public async Task<int> GetChildCategoryCountAsync(Guid categoryId)
        {
            return await _dbSet
                .CountAsync(pc => pc.ParentCategoryId == categoryId);
        }
    }
}
