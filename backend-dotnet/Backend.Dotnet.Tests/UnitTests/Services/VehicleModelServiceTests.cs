using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.VehicleModelDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class VehicleModelServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<IVehicleModelRepository> _modelRepo;
        private VehicleModelService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _modelRepo = new Mock<IVehicleModelRepository>();
            _unitOfWork.Setup(x => x.VehicleModels).Returns(_modelRepo.Object);
            _sut = new VehicleModelService(_unitOfWork.Object);
        }

        [Test]
        public async Task CreateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var request = new CreateVehicleModelRequest
            {
                Brand = "Tesla",
                ModelName = "Model 3",
                Year = 2023
            };

            _modelRepo.Setup(x => x.ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year, null))
                .ReturnsAsync(false);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Brand.Should().Be("Tesla");
            result.Data.ModelName.Should().Be("Model 3");
            result.Data.Year.Should().Be(2023);
            _modelRepo.Verify(x => x.AddAsync(It.IsAny<VehicleModel>()), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task CreateAsync_DuplicateBrandModelYear_ReturnsError()
        {
            // Arrange
            var request = new CreateVehicleModelRequest
            {
                Brand = "Tesla",
                ModelName = "Model 3",
                Year = 2023
            };

            _modelRepo.Setup(x => x.ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year, null))
                .ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_MODEL");
            _modelRepo.Verify(x => x.AddAsync(It.IsAny<VehicleModel>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task CreateAsync_InvalidYear_ReturnsBusinessRuleError()
        {
            // Arrange
            var request = new CreateVehicleModelRequest
            {
                Brand = "Tesla",
                ModelName = "Model 3",
                Year = 1999
            };

            _modelRepo.Setup(x => x.ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year, null))
                .ReturnsAsync(false);
            _modelRepo.Setup(x => x.AddAsync(It.IsAny<VehicleModel>()))
                .ThrowsAsync(new BusinessRuleViolationException("INVALID_YEAR"));

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("VIOLATION_000");
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task DeleteAsync_NoActiveVehicles_DeletesSuccessfully()
        {
            // Arrange
            var modelId = Guid.NewGuid();
            var model = new VehicleModel("Tesla", "Model 3", 2023);

            _modelRepo.Setup(x => x.GetByIdAsync(modelId)).ReturnsAsync(model);
            _modelRepo.Setup(x => x.HasActiveVehiclesAsync(modelId)).ReturnsAsync(false);

            // Act
            var result = await _sut.DeleteAsync(modelId);

            // Assert
            result.IsSuccess.Should().BeTrue();
            _modelRepo.Verify(x => x.Remove(model), Times.Once);
        }

        [Test]
        public async Task DeleteAsync_HasActiveVehicles_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var model = new VehicleModel("Tesla", "Model 3", 2023);

            _modelRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(model);
            _modelRepo.Setup(x => x.HasActiveVehiclesAsync(id)).ReturnsAsync(true);
            _modelRepo.Setup(x => x.GetActiveVehicleCountAsync(id)).ReturnsAsync(5);

            // Act
            var result = await _sut.DeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("MODEL_IN_USE");
            result.Message.Should().Contain("5 active vehicle(s)");
            _modelRepo.Verify(x => x.Remove(It.IsAny<VehicleModel>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task DeleteAsync_ModelNotFound_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            _modelRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync((VehicleModel)null);

            // Act
            var result = await _sut.DeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
            _modelRepo.Verify(x => x.Remove(It.IsAny<VehicleModel>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }
    }
}
