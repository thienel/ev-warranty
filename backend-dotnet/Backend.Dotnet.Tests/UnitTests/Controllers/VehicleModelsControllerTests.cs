using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.VehicleModelDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class VehicleModelsControllerTests
    {
        private Mock<IVehicleModelService> _vehicleModelService;
        private VehicleModelsController _controller;

        [SetUp]
        public void Setup()
        {
            _vehicleModelService = new Mock<IVehicleModelService>();
            _controller = new VehicleModelsController(_vehicleModelService.Object);
        }

        [Test]
        public async Task GetById_ExistingModel_ReturnsOk()
        {
            // Arrange
            var modelId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = true,
                Data = new VehicleModelResponse { Id = modelId, Brand = "Tesla" }
            };

            _vehicleModelService.Setup(x => x.GetByIdAsync(modelId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(modelId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().BeEquivalentTo(serviceResponse);
        }

        [Test]
        public async Task GetById_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleModelService.Setup(x => x.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_WithBrandModelYear_CallsGetByBrandModelYearAsync()
        {
            // Arrange
            var brand = "Tesla";
            var model = "Model 3";
            var year = 2023;
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleModelResponse>()
            };

            _vehicleModelService.Setup(x => x.GetByBrandModelYearAsync(brand, model, year))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(brand: brand, model: model, year: year);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleModelService.Verify(x => x.GetByBrandModelYearAsync(brand, model, year), Times.Once);
            _vehicleModelService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithBrandModelYear_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleModelService.Setup(x => x.GetByBrandModelYearAsync(
                It.IsAny<string>(), It.IsAny<string>(), It.IsAny<int>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(brand: "Unknown", model: "X", year: 2020);

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_WithBrandOnly_CallsGetByBrandAsync()
        {
            // Arrange
            var brand = "Tesla";
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleModelResponse>()
            };

            _vehicleModelService.Setup(x => x.GetByBrandAsync(brand))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(brand: brand);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleModelService.Verify(x => x.GetByBrandAsync(brand), Times.Once);
            _vehicleModelService.Verify(x => x.GetByBrandModelYearAsync(
                It.IsAny<string>(), It.IsAny<string>(), It.IsAny<int>()), Times.Never);
        }

        [Test]
        public async Task GetAll_WithModelOnly_CallsGetByModelNameAsync()
        {
            // Arrange
            var modelName = "Model 3";
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleModelResponse>()
            };

            _vehicleModelService.Setup(x => x.GetByModelNameAsync(modelName))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(model: modelName);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleModelService.Verify(x => x.GetByModelNameAsync(modelName), Times.Once);
        }

        [Test]
        public async Task GetAll_WithYearOnly_CallsGetByYearAsync()
        {
            // Arrange
            var year = 2023;
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleModelResponse>()
            };

            _vehicleModelService.Setup(x => x.GetByYearAsync(year))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(year: year);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleModelService.Verify(x => x.GetByYearAsync(year), Times.Once);
        }

        [Test]
        public async Task GetAll_NoFilters_CallsGetAllAsync()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = true,
                Data = new List<VehicleModelResponse>()
            };

            _vehicleModelService.Setup(x => x.GetAllAsync())
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll();

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _vehicleModelService.Verify(x => x.GetAllAsync(), Times.Once);
        }

        [Test]
        public async Task GetAll_NoFilters_Failed_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<VehicleModelResponse>>
            {
                IsSuccess = false,
                ErrorCode = "INTERNAL_ERROR"
            };

            _vehicleModelService.Setup(x => x.GetAllAsync())
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
            _controller.ModelState.AddModelError("Brand", "Required");

            // Act
            var result = await _controller.Create(new CreateVehicleModelRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _vehicleModelService.Verify(x => x.CreateAsync(It.IsAny<CreateVehicleModelRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Create_ValidData_ReturnsCreatedAtAction()
        {
            // Arrange
            var request = new CreateVehicleModelRequest
            {
                Brand = "Tesla",
                ModelName = "Model 3",
                Year = 2023
            };
            var modelId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = true,
                Data = new VehicleModelResponse { Id = modelId, Brand = "Tesla" }
            };

            _vehicleModelService.Setup(x => x.CreateAsync(request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(request);

            // Assert
            result.Should().BeOfType<CreatedAtActionResult>();
            var createdResult = result as CreatedAtActionResult;
            createdResult.ActionName.Should().Be(nameof(_controller.GetById));
            createdResult.RouteValues["id"].Should().Be(modelId);
        }

        [Test]
        public async Task Create_DuplicateModel_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_MODEL"
            };

            _vehicleModelService.Setup(x => x.CreateAsync(It.IsAny<CreateVehicleModelRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(new CreateVehicleModelRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Update_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("Year", "Invalid");

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateVehicleModelRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _vehicleModelService.Verify(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateVehicleModelRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Update_ValidData_ReturnsOk()
        {
            // Arrange
            var modelId = Guid.NewGuid();
            var request = new UpdateVehicleModelRequest();
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = true,
                Data = new VehicleModelResponse { Id = modelId }
            };

            _vehicleModelService.Setup(x => x.UpdateAsync(modelId, request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(modelId, request);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Update_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleModelService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateVehicleModelRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateVehicleModelRequest());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Update_DuplicateModel_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<VehicleModelResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_MODEL"
            };

            _vehicleModelService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateVehicleModelRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateVehicleModelRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Delete_ExistingModel_ReturnsOk()
        {
            // Arrange
            var modelId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto
            {
                IsSuccess = true,
                Message = "Deleted successfully"
            };

            _vehicleModelService.Setup(x => x.DeleteAsync(modelId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Delete(modelId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Delete_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _vehicleModelService.Setup(x => x.DeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Delete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Delete_ModelInUse_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "MODEL_IN_USE"
            };

            _vehicleModelService.Setup(x => x.DeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Delete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }
    }
}
