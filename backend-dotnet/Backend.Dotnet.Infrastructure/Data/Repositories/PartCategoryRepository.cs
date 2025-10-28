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
        public PartCategoryRepository(AppDbContext context) : base(context)
        {
        }

        public async Task<PartCategory> GetByIdWithChildrenAsync(Guid id)
        {
            return await _dbSet
                .Include(pc => pc.ChildCategories)
                .FirstOrDefaultAsync(pc => pc.Id == id);
        }

        public async Task<PartCategory> GetByIdWithParentAsync(Guid id)
        {
            return await _dbSet
                .Include(pc => pc.ParentCategory)
                .FirstOrDefaultAsync(pc => pc.Id == id);
        }

        public async Task<List<PartCategory>> GetByParentIdAsync(Guid? parentId)
        {
            return await _dbSet
                .Where(pc => pc.ParentCategoryId == parentId)
                .OrderBy(pc => pc.CategoryName)
                .ToListAsync();
        }

        public async Task<List<PartCategory>> GetByStatusAsync(PartCategoryStatus status)
        {
            return await _dbSet
                .Where(pc => pc.Status == status)
                .OrderBy(pc => pc.CategoryName)
                .ToListAsync();
        }

        public async Task<List<PartCategory>> GetActiveHierarchyAsync()
        {
            return await _dbSet
                .Include(pc => pc.ParentCategory)
                .Include(pc => pc.ChildCategories)
                .Where(pc => pc.Status == PartCategoryStatus.Active)
                .OrderBy(pc => pc.CategoryName)
                .ToListAsync();
        }

        public async Task<bool> ExistsByNameAsync(string categoryName)
        {
            return await _dbSet.AnyAsync(pc => pc.CategoryName == categoryName);
        }

        public async Task<bool> HasChildCategoriesAsync(Guid id)
        {
            return await _dbSet.AnyAsync(pc => pc.ParentCategoryId == id);
        }

        public async Task<bool> HasPartsAsync(Guid id)
        {
            //return await _context.Parts.AnyAsync(p => p.CategoryId == id);
            return await _dbSet.AnyAsync(p => p.ParentCategoryId == id); //temp for run test
        }

        public override async Task<IEnumerable<PartCategory>> GetAllAsync()
        {
            return await _dbSet
                .Include(pc => pc.ParentCategory)
                .OrderBy(pc => pc.CategoryName)
                .ToListAsync();
        }
    }
}
