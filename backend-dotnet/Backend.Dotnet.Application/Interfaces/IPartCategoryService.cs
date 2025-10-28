using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IPartCategoryService
    {
        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> CreateAsync(
            PartCategoryDto.CreatePartCategoryRequest request);

        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PartCategoryDto.PartCategoryResponse>>> GetAllAsync();
        Task<BaseResponseDto<PartCategoryDto.PartCategoryWithHierarchyResponse>> GetWithHierarchyAsync(Guid id);

        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> GetByCategoryNameAsync(string categoryName);

        Task<BaseResponseDto<IEnumerable<PartCategoryDto.PartCategoryResponse>>> GetByStatusAsync(string status);
        Task<BaseResponseDto<IEnumerable<PartCategoryDto.PartCategoryResponse>>> GetActiveCategoriesAsync();

        Task<BaseResponseDto<IEnumerable<PartCategoryDto.PartCategoryResponse>>> GetByParentIdAsync(Guid parentId);
        Task<BaseResponseDto<IEnumerable<PartCategoryDto.PartCategoryResponse>>> GetRootCategoriesAsync();
        Task<BaseResponseDto<IEnumerable<PartCategoryDto.PartCategoryWithHierarchyResponse>>> GetFullHierarchyAsync();

        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> UpdateAsync(
            Guid id, PartCategoryDto.UpdatePartCategoryRequest request);
        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> ChangeParentAsync(
            Guid id, PartCategoryDto.ChangeParentCategoryRequest request);

        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> ChangeStatusAsync(
            Guid id, PartCategoryDto.ChangeStatusRequest request);
        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> ActivateAsync(Guid id);
        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> MakeReadOnlyAsync(Guid id);
        Task<BaseResponseDto<PartCategoryDto.PartCategoryResponse>> ArchiveAsync(Guid id);

        // Soft delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
