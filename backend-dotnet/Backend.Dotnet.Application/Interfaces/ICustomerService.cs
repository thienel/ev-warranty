using Backend.Dotnet.Application.DTOs;
using static Backend.Dotnet.Application.DTOs.CustomerDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface ICustomerService
    {
        Task<BaseResponseDto<CustomerResponse>> CreateAsync(CreateCustomerRequest request);

        Task<BaseResponseDto<CustomerResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<CustomerResponse>>> GetAllAsync();
        Task<BaseResponseDto<CustomerResponse>> GetByEmailAsync(string email);
        Task<BaseResponseDto<CustomerResponse>> GetByPhoneAsync(string phone);
        Task<BaseResponseDto<CustomerResponse>> GetByNameAsync(string name);

        Task<BaseResponseDto<CustomerResponse>> UpdateAsync(Guid id, UpdateCustomerRequest request);
        
        // Soft delete operations
        Task<BaseResponseDto<CustomerResponse>> SoftDeleteAsync(Guid id);
        Task<BaseResponseDto<CustomerResponse>> RestoreAsync(Guid id);
    }
}