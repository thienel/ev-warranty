using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IPolicyCoveragePartService
    {
        Task<BaseResponseDto<PolicyCoveragePartResponse>> CreateAsync(
            CreatePolicyCoveragePartRequest request);

        Task<BaseResponseDto<PolicyCoveragePartResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>> GetAllAsync();
        Task<BaseResponseDto<PolicyCoveragePartDetailResponse>> GetWithDetailsAsync(Guid id);

        Task<BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>> GetByPolicyIdAsync(Guid policyId);
        Task<BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>> GetByPartCategoryIdAsync(Guid partCategoryId);
        Task<BaseResponseDto<PolicyCoveragePartResponse>> GetByPolicyAndCategoryAsync(
            Guid policyId, Guid partCategoryId);

        Task<BaseResponseDto<PolicyCoveragePartResponse>> UpdateAsync(
            Guid id, UpdatePolicyCoveragePartRequest request);

        // Hard delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
