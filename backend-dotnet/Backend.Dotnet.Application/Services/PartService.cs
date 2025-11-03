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
using static Backend.Dotnet.Application.DTOs.PartDto;

namespace Backend.Dotnet.Application.Services
{
    public class PartService : IPartService
    {
        private readonly IUnitOfWork _unitOfWork;

        public PartService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<PartResponse>> CreateAsync(CreatePartRequest request)
        {
            try
            {
                var exists = await _unitOfWork.Parts.SerialNumberExistsAsync(request.SerialNumber);
                if (exists)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with serial number '{request.SerialNumber}' already exists",
                        ErrorCode = "DUPLICATE_SERIAL_NUMBER"
                    };
                }

                // Verify category exists
                var category = await _unitOfWork.PartCategories.GetByIdAsync(request.CategoryId);
                if (category == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{request.CategoryId}' not found",
                        ErrorCode = "CATEGORY_NOT_FOUND"
                    };
                }

                var part = request.ToEntity();
                await _unitOfWork.Parts.AddAsync(part);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part created successfully",
                    Data = part.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetByIdAsync(id);
                if (part == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part retrieved successfully",
                    Data = part.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }


        public async Task<BaseResponseDto<PartResponse>> ReserveByOfficeIdAndCategoryIdAsync(Guid officeId, Guid categoryId)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetByOfficeIdAndCategoryId(officeId, categoryId);
                if (part == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with office id '{officeId}' and category id '{categoryId}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                part.ChangeStatus(PartStatus.Reserved);
                _unitOfWork.Parts.Update(part);

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part reserved successfully",
                    Data = part.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
        
        public async Task<BaseResponseDto<IEnumerable<PartResponse>>> GetAllAsync()
        {
            try
            {
                var parts = await _unitOfWork.Parts.GetAllWithDetailsAsync();
                var response = parts.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = true,
                    Message = "Parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartWithDetailsResponse>> GetWithDetailsAsync(Guid id)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetWithDetailsAsync(id);
                if (part == null)
                {
                    return new BaseResponseDto<PartWithDetailsResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PartWithDetailsResponse>
                {
                    IsSuccess = true,
                    Message = "Part with details retrieved successfully",
                    Data = part.ToWithDetailsResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartWithDetailsResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartResponse>> GetBySerialNumberAsync(string serialNumber)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetBySerialNumberAsync(serialNumber);
                if (part == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with serial number '{serialNumber}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part retrieved successfully",
                    Data = part.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartResponse>>> GetByCategoryIdAsync(Guid categoryId)
        {
            try
            {
                var category = await _unitOfWork.PartCategories.GetByIdAsync(categoryId);
                if (category == null)
                {
                    return new BaseResponseDto<IEnumerable<PartResponse>>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{categoryId}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var parts = await _unitOfWork.Parts.GetByCategoryIdAsync(categoryId);
                var response = parts.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = true,
                    Message = "Parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartResponse>>> GetByOfficeLocationIdAsync(Guid officeLocationId)
        {
            try
            {
                var parts = await _unitOfWork.Parts.GetByOfficeLocationIdAsync(officeLocationId);
                var response = parts.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = true,
                    Message = "Parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartResponse>>> GetByStatusAsync(string status)
        {
            try
            {
                if (!Enum.TryParse<PartStatus>(status, true, out var statusEnum))
                {
                    return new BaseResponseDto<IEnumerable<PartResponse>>
                    {
                        IsSuccess = false,
                        Message = "Invalid status value. Must be Available, Reserved, Installed, Defective, Obsolete, or Archived",
                        ErrorCode = "INVALID_STATUS"
                    };
                }

                var parts = await _unitOfWork.Parts.GetByStatusAsync(statusEnum);
                var response = parts.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = true,
                    Message = "Parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PartResponse>>> SearchAsync(string searchTerm)
        {
            try
            {
                var parts = await _unitOfWork.Parts.SearchAsync(searchTerm);
                var response = parts.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = true,
                    Message = "Search completed successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while searching parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartResponse>> UpdateAsync(Guid id, UpdatePartRequest request)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetByIdAsync(id);
                if (part == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                request.ApplyToEntity(part);
                _unitOfWork.Parts.Update(part);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part updated successfully",
                    Data = part.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartResponse>> ChangeCategoryAsync(
            Guid id, ChangePartCategoryRequest request)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetByIdAsync(id);
                if (part == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var category = await _unitOfWork.PartCategories.GetByIdAsync(request.CategoryId);
                if (category == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{request.CategoryId}' not found",
                        ErrorCode = "CATEGORY_NOT_FOUND"
                    };
                }

                part.UpdateCategory(request.CategoryId);
                _unitOfWork.Parts.Update(part);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part category changed successfully",
                    Data = part.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while changing part category",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PartResponse>> ChangeStatusAsync(
            Guid id, PartChangeStatusRequest request)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetByIdAsync(id);
                if (part == null)
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                if (!Enum.TryParse<PartStatus>(request.Status, true, out var statusEnum))
                {
                    return new BaseResponseDto<PartResponse>
                    {
                        IsSuccess = false,
                        Message = "Invalid status value. Must be Available, Reserved, Installed, Defective, Obsolete, or Archived",
                        ErrorCode = "INVALID_STATUS"
                    };
                }

                part.ChangeStatus(statusEnum);
                _unitOfWork.Parts.Update(part);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = true,
                    Message = "Part status changed successfully",
                    Data = part.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while changing part status",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            try
            {
                var part = await _unitOfWork.Parts.GetByIdAsync(id);
                if (part == null)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Check can be deleted
                if (part.Status == PartStatus.Reserved || part.Status == PartStatus.Installed)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Cannot delete part with status {part.Status}",
                        ErrorCode = "PART_IN_USE"
                    };
                }

                _unitOfWork.Parts.Remove(part);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto
                {
                    IsSuccess = true,
                    Message = "Part deleted successfully"
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
