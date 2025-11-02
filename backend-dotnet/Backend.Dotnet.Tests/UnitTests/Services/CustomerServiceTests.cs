using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.CustomerDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class CustomerServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<ICustomerRepository> _customerRepo;
        private CustomerService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _customerRepo = new Mock<ICustomerRepository>();
            _unitOfWork.Setup(x => x.Customers).Returns(_customerRepo.Object);
            _sut = new CustomerService(_unitOfWork.Object);
        }

        // Create
        [Test]
        public async Task CreateAsync_ValidData_ReturnsSuccess()
        {
            //Arrange
            var request = new CreateCustomerRequest
            {
                FirstName = "Mark",
                LastName = "Le",
                Email = "mark.lev1369@gmail.com",
                PhoneNumber = "1234567890",
                Address = ""
            };

            _customerRepo.Setup(x => x.EmailExistsAsync(It.IsAny<string>(), null))
                .ReturnsAsync(false);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            //Act
            var result = await _sut.CreateAsync(request);

            //Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Email.Should().Be("mark.lev1369@gmail.com");
            _customerRepo.Verify(x => x.AddAsync(It.IsAny<Customer>()), Times.Once());
        }

        [Test]
        public async Task CreateAsync_DuplicateEmail_ReturnsError()
        {
            // Arrange
            var request = new CreateCustomerRequest { Email = "exist@test.com" };
            _customerRepo.Setup(x => x.EmailExistsAsync("exist@test.com", null))
                .ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_EMAIL");
            _customerRepo.Verify(x => x.AddAsync(It.IsAny<Customer>()), Times.Never);
        }

        [Test]
        public async Task CreateAsync_BusinessRuleViolation_ReturnsError()
        {
            // Arrange
            var request = new CreateCustomerRequest { FirstName = "" }; // Invalid

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("VIOLATION_000");
        }

        [Test]
        public async Task CreateAsync_UnexpectedException_ReturnsInternalError()
        {
            // Arrange
            var request = new CreateCustomerRequest
            {
                FirstName = "Mark",
                LastName = "Le"
            };
            _customerRepo.Setup(x => x.EmailExistsAsync(It.IsAny<string>(), null))
                .ThrowsAsync(new Exception("DB Error"));

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("INTERNAL_ERROR");

        }

        [Test]
        public async Task UpdateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var existingCustomer = new Customer("John", "Doe", "john@test.com", "123456", "Address");
            var request = new UpdateCustomerRequest
            {
                FirstName = "Jane",
                LastName = "Smith",
                Email = "jane@test.com",
                PhoneNumber = "789012",
                Address = "New Address"
            };

            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(existingCustomer);
            _customerRepo.Setup(x => x.EmailExistsAsync(request.Email, id)).ReturnsAsync(false);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.FirstName.Should().Be("Jane");
            result.Data.LastName.Should().Be("Smith");
            _customerRepo.Verify(x => x.Update(existingCustomer), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task UpdateAsync_NotFound_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var request = new UpdateCustomerRequest
            {
                FirstName = "Jane",
                LastName = "Smith",
                Email = "jane@test.com"
            };

            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync((Customer)null);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
            _customerRepo.Verify(x => x.Update(It.IsAny<Customer>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task UpdateAsync_DuplicateEmailForOtherCustomer_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var existingCustomer = new Customer("John", "Doe", "john@test.com", "123456", "Address");
            var request = new UpdateCustomerRequest
            {
                FirstName = "Jane",
                LastName = "Smith",
                Email = "existing@test.com"
            };

            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(existingCustomer);
            _customerRepo.Setup(x => x.EmailExistsAsync(request.Email, id)).ReturnsAsync(true);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_EMAIL");
            _customerRepo.Verify(x => x.Update(It.IsAny<Customer>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task SoftDeleteAsync_ExistingCustomer_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");

            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(customer);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.SoftDeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.IsDeleted.Should().BeTrue();
            _customerRepo.Verify(x => x.Update(customer), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task SoftDeleteAsync_AlreadyDeleted_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");
            customer.Delete();

            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(customer);

            // Act
            var result = await _sut.SoftDeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("VIOLATION_000");
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task RestoreAsync_DeletedCustomer_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");
            customer.Delete();

            _customerRepo.Setup(x => x.GetByIdIncludingDeletedAsync(id)).ReturnsAsync(customer);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.RestoreAsync(id);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.IsDeleted.Should().BeFalse();
            _customerRepo.Verify(x => x.Update(customer), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task RestoreAsync_NotDeleted_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");

            _customerRepo.Setup(x => x.GetByIdIncludingDeletedAsync(id)).ReturnsAsync(customer);

            // Act
            var result = await _sut.RestoreAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("VIOLATION_000");
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task GetByIdAsync_ExistingId_ReturnsCustomer()
        {
            // Arrange
            var id = Guid.NewGuid();
            var customer = new Customer("John", "Doe", "john@test.com", "123456", "Address");

            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(customer);

            // Act
            var result = await _sut.GetByIdAsync(id);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Id.Should().Be(customer.Id);
            result.Data.Email.Should().Be("john@test.com");
            _customerRepo.Verify(x => x.GetByIdAsync(id), Times.Once);
        }

        [Test]
        public async Task GetByIdAsync_NonExistingId_ReturnsNotFound()
        {
            // Arrange
            var id = Guid.NewGuid();
            _customerRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync((Customer)null);

            // Act
            var result = await _sut.GetByIdAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
            result.Data.Should().BeNull();
            _customerRepo.Verify(x => x.GetByIdAsync(id), Times.Once);
        }

        [Test]
        public async Task GetAllAsync_ReturnsAllCustomers()
        {
            // Arrange
            var customers = new List<Customer>
            {
                new Customer("John", "Doe", "john@test.com", "123456", "Address1"),
                new Customer("Jane", "Smith", "jane@test.com", "789012", "Address2")
            };

            _customerRepo.Setup(x => x.GetAllAsync()).ReturnsAsync(customers);

            // Act
            var result = await _sut.GetAllAsync();

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Should().HaveCount(2);
            result.Data.First().Email.Should().Be("john@test.com");
            _customerRepo.Verify(x => x.GetAllAsync(), Times.Once);
        }

        /*
         * Template for refer
        [Test]
        public async Task Async__ReturnsSuccessError()
        {
            //Arrange

            //Act

            //Assert

        }
        */
    }
}
