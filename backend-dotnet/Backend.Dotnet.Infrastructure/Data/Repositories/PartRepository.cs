using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Infrastructure.Data.Repositories
{
    public class PartRepository : BaseRepository<Part>, IPartRepository
    {
        public PartRepository(DbContext context) : base(context) { }

        public async Task<Part?> GetBySerialNumberAsync(string serialNumber)
        {
            return await _dbSet
                .Include(p => p.Category)
                .FirstOrDefaultAsync(p => p.SerialNumber.ToLower() == serialNumber.ToLower());
        }

        public async Task<bool> SerialNumberExistsAsync(string serialNumber, Guid? excludePartId = null)
        {
            var query = _dbSet.Where(p => p.SerialNumber.ToLower() == serialNumber.ToLower());

            if (excludePartId.HasValue)
                query = query.Where(p => p.Id != excludePartId.Value);

            return await query.AnyAsync();
        }

        public async Task<IEnumerable<Part>> GetByCategoryIdAsync(Guid categoryId, PartStatus? status = null)
        {
            IQueryable<Part> query = _dbSet
                .Include(p => p.Category)
                .Where(p => p.CategoryId == categoryId);

            if (status.HasValue)
                query = query.Where(p => p.Status == status.Value);

            return await query
                .OrderByDescending(p => p.CreatedAt)
                .ToListAsync();
        }

        public async Task<IEnumerable<Part>> GetByOfficeLocationIdAsync(Guid officeLocationId, PartStatus? status = null)
        {
            IQueryable<Part> query = _dbSet
                .Where(p => p.OfficeLocationId == officeLocationId);

            if (status.HasValue)
                query = query.Where(p => p.Status == status.Value);

            return await query
                .OrderByDescending(p => p.CreatedAt)
                .ToListAsync();
        }

        public async Task<IEnumerable<Part>> GetByStatusAsync(PartStatus status)
        {
            return await _dbSet
                .Include(p => p.Category)
                .Where(p => p.Status == status)
                .OrderByDescending(p => p.CreatedAt)
                .ToListAsync();
        }

        public async Task<Part?> GetWithDetailsAsync(Guid partId)
        {
            return await _dbSet
                .Include(p => p.Category)
                .FirstOrDefaultAsync(p => p.Id == partId);
        }

        public async Task<IEnumerable<Part>> GetAllWithDetailsAsync()
        {
            return await _dbSet
                .Include(p => p.Category)
                .OrderByDescending(p => p.CreatedAt)
                .ToListAsync();
        }

        public async Task<IEnumerable<Part>> SearchAsync(string searchTerm)
        {
            if (string.IsNullOrWhiteSpace(searchTerm))
            {
                return await _dbSet
                    .Include(p => p.Category)
                    .OrderByDescending(p => p.CreatedAt)
                    .ToListAsync();
            }

            var lower = searchTerm.ToLower();

            return await _dbSet
                .Include(p => p.Category)
                .Where(p =>
                    p.PartName.ToLower().Contains(lower) ||
                    p.SerialNumber.ToLower().Contains(lower) ||
                    p.Category.CategoryName.ToLower().Contains(lower))
                .OrderByDescending(p => p.CreatedAt)
                .ToListAsync();
        }

        public async Task<int> GetAvailableCountByCategoryAsync(Guid categoryId)
        {
            return await _dbSet
                .CountAsync(p => p.CategoryId == categoryId && p.Status == PartStatus.Available);
        }

        public async Task<int> GetAvailableCountByLocationAsync(Guid officeLocationId)
        {
            return await _dbSet
                .CountAsync(p => p.OfficeLocationId == officeLocationId && p.Status == PartStatus.Available);
        }
    }
}