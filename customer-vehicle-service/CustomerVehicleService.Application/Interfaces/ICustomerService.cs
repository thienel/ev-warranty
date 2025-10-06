using CustomerVehicleService.Application.DTOs;
using static CustomerVehicleService.Application.DTOs.CustomerDto;

namespace CustomerVehicleService.Application.Interfaces
{
    public interface ICustomerService
    {
        Task<BaseResponseDto<CustomerResponse>> CreateAsync(CreateCustomerRequest request);

        Task<BaseResponseDto<CustomerResponse>> GetByIdAsync(Guid id);
        //Task<BaseResponseDto<CustomerWithVehiclesResponse>> GetWithVehiclesAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<CustomerResponse>>> GetAllAsync();
        Task<BaseResponseDto<CustomerResponse>> GetByEmailAsync(string email);
        Task<BaseResponseDto<IEnumerable<CustomerResponse>>> SearchAsync(string searchTerm);

        Task<BaseResponseDto<CustomerResponse>> UpdateAsync(Guid id, UpdateCustomerRequest request);
        Task<BaseResponseDto<CustomerResponse>> UpdateEmailAsync(Guid id, string email);
        Task<BaseResponseDto<CustomerResponse>> UpdatePhoneNumberAsync(Guid id, string phoneNumber);
        Task<BaseResponseDto<CustomerResponse>> UpdateAddressAsync(Guid id, string address);
    }
}