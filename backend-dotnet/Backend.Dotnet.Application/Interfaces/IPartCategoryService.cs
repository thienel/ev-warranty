using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartCategoryDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IPartCategoryService
    {
        Task<BaseResponseDto<PartCategoryResponse>> CreateAsync(
            CreatePartCategoryRequest request);

        Task<BaseResponseDto<PartCategoryResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PartCategoryResponse>>> GetAllAsync();

        Task<BaseResponseDto<PartCategoryResponse>> GetByCategoryNameAsync(string categoryName);

        Task<BaseResponseDto<IEnumerable<PartCategoryResponse>>> GetByParentIdAsync(Guid parentId);
        Task<BaseResponseDto<PartCategoryWithHierarchyResponse>> GetWithHierarchyAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PartCategoryWithHierarchyResponse>>> GetFullHierarchyAsync();

        Task<BaseResponseDto<PartCategoryResponse>> UpdateAsync(
            Guid id, UpdatePartCategoryRequest request);

        // Soft delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
