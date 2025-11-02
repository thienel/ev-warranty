using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.VehicleDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class VehicleServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<ICustomerRepository> _customerRepo;
        private Mock<IVehicleModelRepository> _modelRepo;
        private Mock<IVehicleRepository> _vehicleRepo;
        private VehicleService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _customerRepo = new Mock<ICustomerRepository>();
            _modelRepo = new Mock<IVehicleModelRepository>();
            _vehicleRepo = new Mock<IVehicleRepository>();
            _unitOfWork.Setup(x => x.Vehicles).Returns(_vehicleRepo.Object);
            _unitOfWork.Setup(x => x.Customers).Returns(_customerRepo.Object);
            _unitOfWork.Setup(x => x.VehicleModels).Returns(_modelRepo.Object);
            _sut = new VehicleService(_unitOfWork.Object);
        }

        [Test]
        public async Task CreateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var customerId = Guid.NewGuid();
            var modelId = Guid.NewGuid();
            var request = new CreateVehicleRequest
            {
                Vin = "1HGBH41JXMN109186",
                LicensePlate = "ABC123",
                CustomerId = customerId,
                ModelId = modelId,
                PurchaseDate = DateTime.UtcNow
            };

            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");
            var model = new VehicleModel("Tesla", "Model 3", 2023);

            _vehicleRepo.Setup(x => x.VinExistsAsync(request.Vin, null)).ReturnsAsync(false);
            _customerRepo.Setup(x => x.GetByIdAsync(customerId)).ReturnsAsync(customer);
            _modelRepo.Setup(x => x.GetByIdAsync(modelId)).ReturnsAsync(model);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Vin.Should().Be("1HGBH41JXMN109186");
            _vehicleRepo.Verify(x => x.AddAsync(It.IsAny<Vehicle>()), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task CreateAsync_DuplicateVIN_ReturnsError()
        {
            // Arrange
            var request = new CreateVehicleRequest
            {
                Vin = "1HGBH41JXMN109186",
                CustomerId = Guid.NewGuid(),
                ModelId = Guid.NewGuid()
            };

            _vehicleRepo.Setup(x => x.VinExistsAsync(request.Vin, null)).ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_VIN");
            _vehicleRepo.Verify(x => x.AddAsync(It.IsAny<Vehicle>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task CreateAsync_CustomerNotFound_ReturnsError()
        {
            // Arrange
            var customerId = Guid.NewGuid();
            var request = new CreateVehicleRequest
            {
                Vin = "1HGBH41JXMN109186",
                CustomerId = customerId,
                ModelId = Guid.NewGuid()
            };

            _vehicleRepo.Setup(x => x.VinExistsAsync(request.Vin, null)).ReturnsAsync(false);
            _customerRepo.Setup(x => x.GetByIdAsync(customerId)).ReturnsAsync((Customer)null);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CUSTOMER_NOT_FOUND");
            _vehicleRepo.Verify(x => x.AddAsync(It.IsAny<Vehicle>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task CreateAsync_ModelNotFound_ReturnsError()
        {
            // Arrange
            var customerId = Guid.NewGuid();
            var modelId = Guid.NewGuid();
            var request = new CreateVehicleRequest
            {
                Vin = "1HGBH41JXMN109186",
                CustomerId = customerId,
                ModelId = modelId
            };

            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");

            _vehicleRepo.Setup(x => x.VinExistsAsync(request.Vin, null)).ReturnsAsync(false);
            _customerRepo.Setup(x => x.GetByIdAsync(customerId)).ReturnsAsync(customer);
            _modelRepo.Setup(x => x.GetByIdAsync(modelId)).ReturnsAsync((VehicleModel)null);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("MODEL_NOT_FOUND");
            _vehicleRepo.Verify(x => x.AddAsync(It.IsAny<Vehicle>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task CreateAsync_InvalidVIN_ReturnsBusinessRuleError()
        {
            // Arrange
            var customerId = Guid.NewGuid();
            var modelId = Guid.NewGuid();
            var request = new CreateVehicleRequest
            {
                Vin = "INVALID_VIN",
                CustomerId = customerId,
                ModelId = modelId
            };

            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");
            var model = new VehicleModel("Tesla", "Model 3", 2023);

            _vehicleRepo.Setup(x => x.VinExistsAsync(request.Vin, null)).ReturnsAsync(false);
            _customerRepo.Setup(x => x.GetByIdAsync(customerId)).ReturnsAsync(customer);
            _modelRepo.Setup(x => x.GetByIdAsync(modelId)).ReturnsAsync(model);
            _vehicleRepo.Setup(x => x.AddAsync(It.IsAny<Vehicle>()))
                .ThrowsAsync(new BusinessRuleViolationException("INVALID_VIN"));

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("VIOLATION_000");
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task TransferOwnership_ValidData_UpdatesCustomerId()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var oldCustomerId = Guid.NewGuid();
            var newCustomerId = Guid.NewGuid();
            var vehicle = new Vehicle("1HGBH41JXMN109186", oldCustomerId, Guid.NewGuid(), "ABC123", DateTime.UtcNow);
            var newCustomer = new Customer("Jane", "Smith", "jane@test.com", "789012", "Address");
            var command = new TransferVehicleCommand { NewCustomerId = newCustomerId };

            _vehicleRepo.Setup(x => x.GetByIdAsync(vehicleId)).ReturnsAsync(vehicle);
            _customerRepo.Setup(x => x.GetByIdAsync(newCustomerId)).ReturnsAsync(newCustomer);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.TransferOwnershipAsync(vehicleId, command);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.CustomerId.Should().Be(newCustomerId);
            _vehicleRepo.Verify(x => x.Update(vehicle), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task TransferOwnership_NewCustomerNotFound_ReturnsError()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var newCustomerId = Guid.NewGuid();
            var vehicle = new Vehicle("1HGBH41JXMN109186", Guid.NewGuid(), Guid.NewGuid(), "ABC123", DateTime.UtcNow);
            var command = new TransferVehicleCommand { NewCustomerId = newCustomerId };

            _vehicleRepo.Setup(x => x.GetByIdAsync(vehicleId)).ReturnsAsync(vehicle);
            _customerRepo.Setup(x => x.GetByIdAsync(newCustomerId)).ReturnsAsync((Customer)null);

            // Act
            var result = await _sut.TransferOwnershipAsync(vehicleId, command);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CUSTOMER_NOT_FOUND");
            _vehicleRepo.Verify(x => x.Update(It.IsAny<Vehicle>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task TransferOwnership_VehicleNotFound_ReturnsError()
        {
            // Arrange
            var vehicleId = Guid.NewGuid();
            var command = new TransferVehicleCommand { NewCustomerId = Guid.NewGuid() };

            _vehicleRepo.Setup(x => x.GetByIdAsync(vehicleId)).ReturnsAsync((Vehicle)null);

            // Act
            var result = await _sut.TransferOwnershipAsync(vehicleId, command);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
            _customerRepo.Verify(x => x.GetByIdAsync(It.IsAny<Guid>()), Times.Never);
            _vehicleRepo.Verify(x => x.Update(It.IsAny<Vehicle>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

    }
}
