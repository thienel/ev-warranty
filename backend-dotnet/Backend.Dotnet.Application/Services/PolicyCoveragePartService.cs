using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Exceptions;
using Microsoft.IdentityModel.Tokens;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;

namespace Backend.Dotnet.Application.Services
{
    public class PolicyCoveragePartService : IPolicyCoveragePartService
    {
        private readonly IUnitOfWork _unitOfWork;

        public PolicyCoveragePartService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<PolicyCoveragePartResponse>> CreateAsync(CreatePolicyCoveragePartRequest request)
        {
            try
            {
                // Check if policy exists
                var policy = await _unitOfWork.WarrantyPolicies.GetByIdAsync(request.PolicyId);
                if (policy == null)
                {
                    return new BaseResponseDto<PolicyCoveragePartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with ID '{request.PolicyId}' not found",
                        ErrorCode = "POLICY_NOT_FOUND"
                    };
                }

                // Check if part category exists
                var partCategory = await _unitOfWork.PartCategories.GetByIdAsync(request.PartCategoryId);
                if (partCategory == null)
                {
                    return new BaseResponseDto<PolicyCoveragePartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{request.PartCategoryId}' not found",
                        ErrorCode = "CATEGORY_NOT_FOUND"
                    };
                }

                // Check if combination already exists
                var exists = await _unitOfWork.PolicyCoverageParts
                    .ExistsByPolicyAndCategoryAsync(request.PolicyId, request.PartCategoryId);

                if (exists)
                {
                    return new BaseResponseDto<PolicyCoveragePartResponse>
                    {
                        IsSuccess = false,
                        Message = "This part category is already covered by the policy",
                        ErrorCode = "DUPLICATE_COVERAGE"
                    };
                }

                var coverage = request.ToEntity();

                // Validate business rules
                coverage.ValidateAgainstPolicy(policy);
                coverage.ValidateAgainstCategory(partCategory);

                await _unitOfWork.PolicyCoverageParts.AddAsync(coverage);
                await _unitOfWork.SaveChangesAsync();

                var result = await _unitOfWork.PolicyCoverageParts.GetWithDetailsAsync(coverage.Id);

                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = true,
                    Message = "Policy coverage part created successfully",
                    Data = result.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating policy coverage part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PolicyCoveragePartResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var coverage = await _unitOfWork.PolicyCoverageParts.GetWithDetailsAsync(id);
                if (coverage == null)
                {
                    return new BaseResponseDto<PolicyCoveragePartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Policy coverage part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = true,
                    Message = "Policy coverage part retrieved successfully",
                    Data = coverage.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving policy coverage part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>> GetAllAsync()
        {
            try
            {
                var coverages = await _unitOfWork.PolicyCoverageParts.GetAllWithDetailsAsync();
                var response = coverages.Select(c => c.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = true,
                    Message = "Policy coverage parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving policy coverage parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PolicyCoveragePartDetailResponse>> GetWithDetailsAsync(Guid id)
        {
            try
            {
                var coverage = await _unitOfWork.PolicyCoverageParts.GetWithDetailsAsync(id);
                if (coverage == null)
                {
                    return new BaseResponseDto<PolicyCoveragePartDetailResponse>
                    {
                        IsSuccess = false,
                        Message = $"Policy coverage part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PolicyCoveragePartDetailResponse>
                {
                    IsSuccess = true,
                    Message = "Policy coverage part retrieved successfully",
                    Data = coverage.ToDetailResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PolicyCoveragePartDetailResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving policy coverage part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>> GetByPolicyIdAsync(Guid policyId)
        {
            try
            {
                var coverages = await _unitOfWork.PolicyCoverageParts.GetByPolicyIdAsync(policyId);
                if (!coverages.Any())
                {
                    return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                    {
                        IsSuccess = false,
                        Message = $"No coverage parts found for policy ID '{policyId}'",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var response = coverages.Select(c => c.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = true,
                    Message = "Policy coverage parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving policy coverage parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>> GetByPartCategoryIdAsync(Guid partCategoryId)
        {
            try
            {
                var coverages = await _unitOfWork.PolicyCoverageParts
                    .FindAsync(cp => cp.PartCategoryId == partCategoryId);

                if (!coverages.Any())
                {
                    return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                    {
                        IsSuccess = false,
                        Message = $"No coverage parts found for part category ID '{partCategoryId}'",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var response = coverages.Select(c => c.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = true,
                    Message = "Policy coverage parts retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving policy coverage parts",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PolicyCoveragePartResponse>> GetByPolicyAndCategoryAsync(Guid policyId, Guid partCategoryId)
        {
            try
            {
                var coverage = await _unitOfWork.PolicyCoverageParts
                    .GetByPolicyAndCategoryAsync(policyId, partCategoryId);

                if (coverage == null)
                {
                    return new BaseResponseDto<PolicyCoveragePartResponse>
                    {
                        IsSuccess = false,
                        Message = "Coverage not found for the specified policy and part category",
                        ErrorCode = "COVERAGE_NOT_FOUND"
                    };
                }

                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = true,
                    Message = "Policy coverage part retrieved successfully",
                    Data = coverage.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving policy coverage part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<PolicyCoveragePartResponse>> UpdateAsync(Guid id, UpdatePolicyCoveragePartRequest request)
        {
            try
            {
                var coverage = await _unitOfWork.PolicyCoverageParts.GetWithDetailsAsync(id);
                if (coverage == null)
                {
                    return new BaseResponseDto<PolicyCoveragePartResponse>
                    {
                        IsSuccess = false,
                        Message = $"Policy coverage part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                request.ApplyToEntity(coverage);
                _unitOfWork.PolicyCoverageParts.Update(coverage);
                await _unitOfWork.SaveChangesAsync();

                var updated = await _unitOfWork.PolicyCoverageParts.GetWithDetailsAsync(id);

                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = true,
                    Message = "Policy coverage part updated successfully",
                    Data = updated.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<PolicyCoveragePartResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating policy coverage part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            try
            {
                var coverage = await _unitOfWork.PolicyCoverageParts.GetWithDetailsAsync(id);
                if (coverage == null)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Policy coverage part with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Validate if policy is editable
                if (coverage.Policy != null && !coverage.Policy.IsEditable())
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = "Cannot delete coverage from non-draft policy",
                        ErrorCode = "POLICY_NOT_EDITABLE"
                    };
                }

                _unitOfWork.PolicyCoverageParts.Remove(coverage);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto
                {
                    IsSuccess = true,
                    Message = "Policy coverage part deleted successfully",
                    ErrorCode = null
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting policy coverage part",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CoverageDetailsResponse>> GetCoverageDetailsAsync(Guid policyId, Guid partCategoryId)
        {
            try
            {
                var currentCategory = await _unitOfWork.PartCategories.GetByIdAsync(partCategoryId);

                if (currentCategory == null)
                {
                    return new BaseResponseDto<CoverageDetailsResponse>
                    {
                        IsSuccess = false,
                        Message = $"Part category with ID '{partCategoryId}' not found",
                        ErrorCode = "CATEGORY_NOT_FOUND"
                    };
                }

                Guid? currentCategoryId = partCategoryId;
                while (currentCategory != null)
                {
                    var coverage = await _unitOfWork.PolicyCoverageParts
                        .GetByPolicyAndCategoryAsync(policyId, currentCategory.Id);

                    if (coverage != null && !coverage.CoverageConditions.IsNullOrEmpty())
                    {
                        return new BaseResponseDto<CoverageDetailsResponse>
                        {
                            IsSuccess = true,
                            Message = "Coverage details retrieved successfully",
                            //Data = coverage.ToCoverageDetailsResponse(partCategoryId, requestedCategory.CategoryName)
                            Data = coverage.ToCoverageDetailsResponse()
                        };
                    }

                    if (currentCategory.ParentCategoryId.HasValue)
                    {
                        currentCategory = await _unitOfWork.PartCategories.GetByIdAsync(currentCategory.ParentCategoryId.Value);
                    }
                    else
                    {
                        currentCategory = null; // Touch the Root - end loop
                    }
                }

                return new BaseResponseDto<CoverageDetailsResponse>
                {
                    IsSuccess = false,
                    Message = "No coverage found for this part category or its parent categories",
                    ErrorCode = "COVERAGE_NOT_FOUND"
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<CoverageDetailsResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving coverage details",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
