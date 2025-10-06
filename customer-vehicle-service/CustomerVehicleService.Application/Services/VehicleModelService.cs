using CustomerVehicleService.Application.DTOs;
using CustomerVehicleService.Application.Interfaces;
using CustomerVehicleService.Application.Interfaces.Data;
using CustomerVehicleService.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static CustomerVehicleService.Application.DTOs.VehicleModelDto;

namespace CustomerVehicleService.Application.Services
{
    public class VehicleModelService : IVehicleModelService
    {
        private readonly IUnitOfWork _unitOfWork;

        public VehicleModelService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> CreateAsync(CreateVehicleModelRequest request)
        {
            try
            {
                // Check if combination already exists
                var exists = await _unitOfWork.VehicleModels
                    .ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year);

                if (exists)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model '{request.Brand} {request.ModelName} {request.Year}' already exists",
                        ErrorCode = "DUPLICATE_MODEL"
                    };
                }

                var model = request.ToEntity();
                await _unitOfWork.VehicleModels.AddAsync(model);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model created successfully",
                    Data = model.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(id);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model retrieved successfully",
                    Data = model.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        //public async Task<BaseResponseDto<VehicleModelWithStatsResponse>> GetWithStatsAsync(Guid id) => throw new NotImplementedException();

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetAllAsync()
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.GetAllAsync();
                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicle models retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> GetByBrandModelYearAsync(string brand, string modelName, int year)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByBrandModelYearAsync(brand, modelName, year);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model '{brand} {modelName} {year}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model retrieved successfully",
                    Data = model.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetByBrandAsync(string brand)
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.GetByBrandAsync(brand);
                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicle models retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<string>>> GetAllBrandsAsync()
        {
            try
            {
                var brands = await _unitOfWork.VehicleModels.GetAllBrandsAsync();

                return new BaseResponseDto<IEnumerable<string>>
                {
                    IsSuccess = true,
                    Message = "Brands retrieved successfully",
                    Data = brands
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<string>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving brands",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> SearchAsync(string searchTerm)
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.SearchAsync(searchTerm);
                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Search completed successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while searching vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> UpdateAsync(Guid id, UpdateVehicleModelRequest request)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(id);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Check if combination already exists for another model
                var exists = await _unitOfWork.VehicleModels
                    .ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year, id);

                if (exists)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model '{request.Brand} {request.ModelName} {request.Year}' already exists",
                        ErrorCode = "DUPLICATE_MODEL"
                    };
                }

                request.ApplyToEntity(model);
                _unitOfWork.VehicleModels.Update(model);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model updated successfully",
                    Data = model.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
