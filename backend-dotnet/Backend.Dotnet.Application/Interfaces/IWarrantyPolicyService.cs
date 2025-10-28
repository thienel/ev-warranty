using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IWarrantyPolicyService
    {
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> CreateAsync(
            WarrantyPolicyDto.CreateWarrantyPolicyRequest request);
        
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<WarrantyPolicyDto.WarrantyPolicyResponse>>> GetAllAsync();
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyWithDetailsResponse>> GetWithDetailsAsync(Guid id);

        Task<BaseResponseDto<IEnumerable<WarrantyPolicyDto.WarrantyPolicyResponse>>> GetByModelIdAsync(Guid modelId);
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> GetActiveByModelIdAsync(Guid modelId);

        Task<BaseResponseDto<IEnumerable<WarrantyPolicyDto.WarrantyPolicyResponse>>> GetByStatusAsync(string status);
        Task<BaseResponseDto<IEnumerable<WarrantyPolicyDto.WarrantyPolicyResponse>>> GetActivePoliciesAsync();
        Task<BaseResponseDto<IEnumerable<WarrantyPolicyDto.WarrantyPolicyResponse>>> GetDraftPoliciesAsync();

        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> GetByPolicyNameAsync(string policyName);

        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> UpdateAsync(
            Guid id, WarrantyPolicyDto.UpdateWarrantyPolicyRequest request);

        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> ChangeStatusAsync(
            Guid id, WarrantyPolicyDto.ChangeStatusRequest request);
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> ActivateAsync(Guid id);
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> ExpireAsync(Guid id);
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> SupersedeAsync(Guid id);
        Task<BaseResponseDto<WarrantyPolicyDto.WarrantyPolicyResponse>> ArchiveAsync(Guid id);

        // Soft delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }

}
