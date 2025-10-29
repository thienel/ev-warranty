using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.WarrantyPolicyDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IWarrantyPolicyService
    {
        Task<BaseResponseDto<WarrantyPolicyResponse>> CreateAsync(
            CreateWarrantyPolicyRequest request);
        
        Task<BaseResponseDto<WarrantyPolicyResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<WarrantyPolicyWithDetailsResponse>> GetWithDetailsAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>> GetAllAsync();

        Task<BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>> GetByStatusAsync(string status);

        Task<BaseResponseDto<WarrantyPolicyResponse>> GetByPolicyNameAsync(string policyName);

        Task<BaseResponseDto<WarrantyPolicyResponse>> UpdateAsync(
            Guid id, UpdateWarrantyPolicyRequest request);

        Task<BaseResponseDto<WarrantyPolicyResponse>> ChangeStatusAsync(
            Guid id, ChangeStatusRequest request);
        // Soft delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
