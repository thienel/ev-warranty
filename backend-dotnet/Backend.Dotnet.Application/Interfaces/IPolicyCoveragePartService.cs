using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IPolicyCoveragePartService
    {
        Task<BaseResponseDto<PolicyCoveragePartDto.PolicyCoveragePartResponse>> CreateAsync(
            PolicyCoveragePartDto.CreatePolicyCoveragePartRequest request);

        Task<BaseResponseDto<PolicyCoveragePartDto.PolicyCoveragePartResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PolicyCoveragePartDto.PolicyCoveragePartResponse>>> GetAllAsync();
        Task<BaseResponseDto<PolicyCoveragePartDto.PolicyCoveragePartDetailResponse>> GetWithDetailsAsync(Guid id);

        Task<BaseResponseDto<IEnumerable<PolicyCoveragePartDto.PolicyCoveragePartResponse>>> GetByPolicyIdAsync(Guid policyId);

        Task<BaseResponseDto<IEnumerable<PolicyCoveragePartDto.PolicyCoveragePartResponse>>> GetByPartCategoryIdAsync(Guid partCategoryId);

        Task<BaseResponseDto<PolicyCoveragePartDto.PolicyCoveragePartResponse>> GetByPolicyAndCategoryAsync(
            Guid policyId, Guid partCategoryId);
        //Task<BaseResponseDto<bool>> IsCategoryCoverageAsync(Guid policyId, Guid partCategoryId);

        Task<BaseResponseDto<PolicyCoveragePartDto.PolicyCoveragePartResponse>> UpdateAsync(
            Guid id, PolicyCoveragePartDto.UpdatePolicyCoveragePartRequest request);

        // Hard delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
        Task<BaseResponseDto> RemoveByPolicyIdAsync(Guid policyId);
        Task<BaseResponseDto> RemoveCategoryFromPolicyAsync(Guid policyId, Guid partCategoryId);
    }
}
