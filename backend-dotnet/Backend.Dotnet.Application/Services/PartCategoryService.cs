using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartCategoryDto;

namespace Backend.Dotnet.Application.Services
{
    public class PartCategoryService : IPartCategoryService
    {
        private readonly IUnitOfWork _unitOfWork;

        public PartCategoryService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<PartCategoryResponse>> CreateAsync(
            CreatePartCategoryRequest request)
        {
            try
            {
                var exists = await _unitOfWork.PartCategories.CategoryNameExistsAsync(request.CategoryName);
                if (exists)
                {
                    return new BaseResponseDto<PartCategoryResponse>
                    {
                        IsSuccess = false,
                        Message = $"Category with name '{request.CategoryName}' already exists",
                        ErrorCode = "DUPLICATE_CATEGORY_NAME"
                    };
                }

                if (request.ParentCategoryId.HasValue)
                {
                    var parent = await _unitOfWork.PartCategories.GetByIdAsync(request.ParentCategoryId.Value);
                    if (parent == null)
                    {
                        return new BaseResponseDto<PartCategoryResponse>
                        {
                            IsSuccess = false,
                            Message = $"Parent category with ID '{request.ParentCategoryId}' not found",
                            ErrorCode = "PARENT_CATEGORY_NOT_FOUND"
                        };
                    }
                }

                var category = request.ToEntity();
                await _unitOfWork.PartCategories.AddAsync(category);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = true,
                    Message = "Part category created successfully",
                    Data = category.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartCategoryResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var category = await _unitOfWork.PartCategories.GetByIdAsync(id);
                if (category == null)
                {
                    return new BaseResponseDto<PartCategoryResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = true,
                    Message = "Part category retrieved successfully",
                    Data = category.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartCategoryResponse>>> GetAllAsync()
        {
            try
            {
                var categories = await _unitOfWork.PartCategories.GetAllAsync();
                var response = categories.Select(c => c.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartCategoryResponse>>
                {
                    IsSuccess = true,
                    Message = "Part categories retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartCategoryResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part categories",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartCategoryResponse>> GetByCategoryNameAsync(string categoryName)
        {
            try
            {
                var category = await _unitOfWork.PartCategories.GetByCategoryNameAsync(categoryName);
                if (category == null)
                {
                    return new BaseResponseDto<PartCategoryResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with name '{categoryName}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = true,
                    Message = "Part category retrieved successfully",
                    Data = category.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartCategoryResponse>>> GetByParentIdAsync(Guid parentId)
        {
            try
            {
                var parent = await _unitOfWork.PartCategories.GetByIdAsync(parentId);
                if (parent == null)
                {
                    return new BaseResponseDto<IEnumerable<PartCategoryResponse>>
                    {
                        IsSuccess = false,
                        Message = $"Parent category with ID '{parentId}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var categories = await _unitOfWork.PartCategories.GetByParentIdAsync(parentId);
                var response = categories.Select(c => c.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartCategoryResponse>>
                {
                    IsSuccess = true,
                    Message = "Child categories retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartCategoryResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part categories",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartCategoryWithHierarchyResponse>> GetWithHierarchyAsync(Guid id)
        {
            try
            {
                var category = await _unitOfWork.PartCategories.GetWithHierarchyAsync(id);
                if (category == null)
                {
                    return new BaseResponseDto<PartCategoryWithHierarchyResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PartCategoryWithHierarchyResponse>
                {
                    IsSuccess = true,
                    Message = "Part category with hierarchy retrieved successfully",
                    Data = category.ToWithHierarchyResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartCategoryWithHierarchyResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartCategoryWithHierarchyResponse>>> GetFullHierarchyAsync()
        {
            try
            {
                var categories = await _unitOfWork.PartCategories.GetFullHierarchyAsync();
                var rootCategories = categories.Where(c => !c.ParentCategoryId.HasValue);
                var response = rootCategories.Select(c => c.ToWithHierarchyResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartCategoryWithHierarchyResponse>>
                {
                    IsSuccess = true,
                    Message = "Full hierarchy retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartCategoryWithHierarchyResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving hierarchy",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartCategoryResponse>> UpdateAsync(
            Guid id, UpdatePartCategoryRequest request)
        {
            try
            {
                var category = await _unitOfWork.PartCategories.GetByIdAsync(id);
                if (category == null)
                {
                    return new BaseResponseDto<PartCategoryResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var exists = await _unitOfWork.PartCategories.CategoryNameExistsAsync(request.CategoryName, id);
                if (exists)
                {
                    return new BaseResponseDto<PartCategoryResponse>
                    {
                        IsSuccess = false,
                        Message = $"Category with name '{request.CategoryName}' already exists",
                        ErrorCode = "DUPLICATE_CATEGORY_NAME"
                    };
                }

                request.ApplyToEntity(category);
                _unitOfWork.PartCategories.Update(category);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = true,
                    Message = "Part category updated successfully",
                    Data = category.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartCategoryResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            try
            {
                var category = await _unitOfWork.PartCategories.GetByIdAsync(id);
                if (category == null)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var activePartCount = await _unitOfWork.PartCategories.GetActivePartCountAsync(id);
                if (activePartCount > 0)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Cannot delete category. {activePartCount} active part(s) are using this category.",
                        ErrorCode = "CATEGORY_HAS_ACTIVE_PARTS"
                    };
                }

                var childCount = await _unitOfWork.PartCategories.GetChildCategoryCountAsync(id);
                if (childCount > 0)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Cannot delete category. {childCount} child categor(ies) exist.",
                        ErrorCode = "CATEGORY_HAS_CHILDREN"
                    };
                }

                _unitOfWork.PartCategories.Remove(category);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto
                {
                    IsSuccess = true,
                    Message = "Part category deleted successfully"
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
