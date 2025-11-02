using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.VehicleDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class VehiclesControllerTests
    {
        private Mock<IVehicleService> _vehicleService;
        private VehiclesController _controller;

        [SetUp]
        public void Setup()
        {
            _vehicleService = new Mock<IVehicleService>();
            _controller = new VehiclesController(_vehicleService.Object);
        }

        [Test]
        public async Task GetById_ExistingVehicle_ReturnsOk()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse
                {
                    Id = vehicleId,
                    Vin = "1HGBH41JXMN109186"
                }
            };

            _vehicleService.Setup(x => x.GetByIdAsync(vehicleId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(vehicleId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().BeEquivalentTo(serviceResponse);
        }

        [Test]
        public async Task GetById_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleService.Setup(x => x.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_WithVinFilter_CallsGetByVinAsync()
        {
            // Arrange
            var vin = "1HGBH41JXMN109186";
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse { Vin = vin }
            };

            _vehicleService.Setup(x => x.GetByVinAsync(vin))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(vin: vin);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleService.Verify(x => x.GetByVinAsync(vin), Times.Once);
            _vehicleService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithVinFilter_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleService.Setup(x => x.GetByVinAsync(It.IsAny<string>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(vin: "UNKNOWN");

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_WithLicensePlateFilter_CallsGetByLicensePlateAsync()
        {
            // Arrange
            var licensePlate = "ABC-1234";
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse { LicensePlate = licensePlate }
            };

            _vehicleService.Setup(x => x.GetByLicensePlateAsync(licensePlate))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(licensePlate: licensePlate);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleService.Verify(x => x.GetByLicensePlateAsync(licensePlate), Times.Once);
            _vehicleService.Verify(x => x.GetByVinAsync(It.IsAny<string>()), Times.Never);
        }

        [Test]
        public async Task GetAll_WithCustomerIdFilter_CallsGetByCustomerIdAsync()
        {
            // Arrange
            var customerId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleResponse> { new VehicleResponse() }
            };

            _vehicleService.Setup(x => x.GetByCustomerIdAsync(customerId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(customerId: customerId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleService.Verify(x => x.GetByCustomerIdAsync(customerId), Times.Once);
        }

        [Test]
        public async Task GetAll_WithCustomerIdFilter_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleResponse>>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleService.Setup(x => x.GetByCustomerIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(customerId: Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_WithModelIdFilter_CallsGetByModelIdAsync()
        {
            // Arrange
            var modelId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleResponse>()
            };

            _vehicleService.Setup(x => x.GetByModelIdAsync(modelId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(modelId: modelId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleService.Verify(x => x.GetByModelIdAsync(modelId), Times.Once);
        }

        [Test]
        public async Task GetAll_NoFilters_CallsGetAllAsync()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleResponse>()
            };

            _vehicleService.Setup(x => x.GetAllAsync())
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll();

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleService.Verify(x => x.GetAllAsync(), Times.Once);
        }

        [Test]
        public async Task GetAll_NoFilters_Failed_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleResponse>>
            {
                IsSuccess = false,
                ErrorCode = "INTERNAL_ERROR"
            };

            _vehicleService.Setup(x => x.GetAllAsync())
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll();

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Create_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("Vin", "Required");
            var request = new CreateVehicleRequest();

            // Act
            var result = await _controller.Create(request);

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _vehicleService.Verify(x => x.CreateAsync(It.IsAny<CreateVehicleRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Create_ValidData_ReturnsCreatedAtAction()
        {
            // Arrange
            var request = new CreateVehicleRequest
            {
                Vin = "1HGBH41JXMN109186",
                CustomerId = Guid.NewGuid(),
                ModelId = Guid.NewGuid()
            };
            var vehicleId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse { Id = vehicleId, Vin = "1HGBH41JXMN109186" }
            };

            _vehicleService.Setup(x => x.CreateAsync(request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(request);

            // Assert
            result.Should().BeOfType<CreatedAtActionResult>();
            var createdResult = result as CreatedAtActionResult;
            createdResult.ActionName.Should().Be(nameof(_controller.GetById));
            createdResult.RouteValues["id"].Should().Be(vehicleId);
            createdResult.Value.Should().BeEquivalentTo(serviceResponse);
        }

        [Test]
        public async Task Create_DuplicateVin_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_VIN"
            };

            _vehicleService.Setup(x => x.CreateAsync(It.IsAny<CreateVehicleRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(new CreateVehicleRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Update_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("Vin", "Invalid");

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateVehicleRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _vehicleService.Verify(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateVehicleRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Update_ValidData_ReturnsOk()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var request = new UpdateVehicleRequest();
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse { Id = vehicleId }
            };

            _vehicleService.Setup(x => x.UpdateAsync(vehicleId, request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(vehicleId, request);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Update_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateVehicleRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateVehicleRequest());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Update_BusinessError_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_VIN"
            };

            _vehicleService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateVehicleRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateVehicleRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task SoftDelete_ExistingVehicle_ReturnsOk()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse { Id = vehicleId }
            };

            _vehicleService.Setup(x => x.SoftDeleteAsync(vehicleId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.SoftDelete(vehicleId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task SoftDelete_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleService.Setup(x => x.SoftDeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.SoftDelete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task SoftDelete_BusinessError_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "VEHICLE_HAS_ACTIVE_SERVICES"
            };

            _vehicleService.Setup(x => x.SoftDeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.SoftDelete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Restore_DeletedVehicle_ReturnsOk()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = true,
                Data = new VehicleResponse { Id = vehicleId }
            };

            _vehicleService.Setup(x => x.RestoreAsync(vehicleId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Restore(vehicleId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Restore_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleService.Setup(x => x.RestoreAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Restore(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Restore_AlreadyActive_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleResponse>
            {
                IsSuccess = false,
                ErrorCode = "VEHICLE_NOT_DELETED"
            };

            _vehicleService.Setup(x => x.RestoreAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Restore(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }
    }
}
