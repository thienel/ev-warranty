using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IPartRepository : IRepository<Part>
    {
        Task<Part?> GetBySerialNumberAsync(string serialNumber);
        Task<bool> SerialNumberExistsAsync(string serialNumber, Guid? excludePartId = null);

        Task<IEnumerable<Part>> GetByCategoryIdAsync(Guid categoryId);
        Task<IEnumerable<Part>> GetAvailableByCategoryIdAsync(Guid categoryId);

        Task<IEnumerable<Part>> GetByOfficeLocationIdAsync(Guid officeLocationId);
        Task<IEnumerable<Part>> GetAvailableByOfficeLocationIdAsync(Guid officeLocationId);

        Task<IEnumerable<Part>> GetByStatusAsync(PartStatus status);

        Task<IEnumerable<Part>> GetByPartNameAsync(string partName);

        Task<Part?> GetWithDetailsAsync(Guid partId);
        Task<Part?> GetWithCategoryAsync(Guid partId);
        Task<IEnumerable<Part>> GetAllWithDetailsAsync();

        Task<IEnumerable<Part>> SearchAsync(string searchTerm);
        Task<IEnumerable<Part>> SearchAvailableAsync(string searchTerm);

        Task<bool> CanBeUsedInWorkOrderAsync(Guid partId);
        Task<bool> IsInStockAsync(Guid partId);
        Task<int> GetAvailableCountByCategoryAsync(Guid categoryId);
        Task<int> GetAvailableCountByLocationAsync(Guid officeLocationId);

        Task<decimal> GetAveragePriceByCategoryAsync(Guid categoryId);
        Task<IEnumerable<Part>> GetByPriceRangeAsync(decimal minPrice, decimal maxPrice);
    }
}
