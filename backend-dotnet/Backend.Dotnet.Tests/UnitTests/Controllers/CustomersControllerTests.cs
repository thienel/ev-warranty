using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.CustomerDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class CustomersControllerTests
    {
        private Mock<ICustomerService> _customerService;
        private CustomersController _controller;

        [SetUp]
        public void Setup()
        {
            _customerService = new Mock<ICustomerService>();
            _controller = new CustomersController(_customerService.Object);
        }

        [Test]
        public async Task GetById_ExistingId_ReturnsOk()
        {
            // Arrange
            var customerId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = true,
                Data = new CustomerResponse { Id = customerId, FirstName = "John" }
            };

            _customerService.Setup(x => x.GetByIdAsync(customerId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(customerId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().BeEquivalentTo(serviceResponse);
        }

        [Test]
        public async Task GetById_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _customerService.Setup(x => x.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Create_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("FirstName", "Required");
            var request = new CreateCustomerRequest();

            // Act
            var result = await _controller.Create(request);

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _customerService.Verify(x => x.CreateAsync(It.IsAny<CreateCustomerRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Create_ValidData_ReturnsCreatedAtAction()
        {
            // Arrange
            var request = new CreateCustomerRequest
            {
                FirstName = "John",
                LastName = "Doe"
            };
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = true,
                Data = new CustomerResponse { Id = Guid.NewGuid(), FirstName = "John" }
            };

            _customerService.Setup(x => x.CreateAsync(request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(request);

            // Assert
            result.Should().BeOfType<CreatedAtActionResult>();
            var createdResult = result as CreatedAtActionResult;
            createdResult.ActionName.Should().Be(nameof(_controller.GetById));
            createdResult.RouteValues["id"].Should().Be(serviceResponse.Data.Id);
        }

        [Test]
        public async Task Create_DuplicateEmail_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_EMAIL"
            };

            _customerService.Setup(x => x.CreateAsync(It.IsAny<CreateCustomerRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(new CreateCustomerRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Update_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _customerService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateCustomerRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateCustomerRequest());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Update_DuplicateEmail_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_EMAIL"
            };

            _customerService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateCustomerRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateCustomerRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task GetAll_WithEmailFilter_CallsGetByEmailAsync()
        {
            // Arrange
            var email = "test@example.com";
            var serviceResponse = new BaseResponseDto<IEnumerable<CustomerResponse>>
            {
                IsSuccess = true,
                Data = new List<CustomerResponse>()
            };

            _customerService.Setup(x => x.GetByEmailAsync(email))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(email: email);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _customerService.Verify(x => x.GetByEmailAsync(email), Times.Once);
            _customerService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithPhoneFilter_CallsGetByPhoneAsync()
        {
            // Arrange
            var phone = "0123456789";
            var serviceResponse = new BaseResponseDto<IEnumerable<CustomerResponse>>
            {
                IsSuccess = true,
                Data = new List<CustomerResponse>()
            };

            _customerService.Setup(x => x.GetByPhoneAsync(phone))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(phone: phone);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _customerService.Verify(x => x.GetByPhoneAsync(phone), Times.Once);
        }

        [Test]
        public async Task GetAll_NoFilters_CallsGetAllAsync()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<CustomerResponse>>
            {
                IsSuccess = true,
                Data = new List<CustomerResponse>()
            };

            _customerService.Setup(x => x.GetAllAsync())
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll();

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _customerService.Verify(x => x.GetAllAsync(), Times.Once);
        }

        [Test]
        public async Task SoftDelete_ExistingCustomer_ReturnsOk()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = true,
                Data = new CustomerResponse()
            };

            _customerService.Setup(x => x.SoftDeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.SoftDelete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Restore_DeletedCustomer_ReturnsOk()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<CustomerResponse>
            {
                IsSuccess = true,
                Data = new CustomerResponse()
            };

            _customerService.Setup(x => x.RestoreAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Restore(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }
    }
}