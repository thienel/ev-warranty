using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class PartsControllerTests
    {
        private Mock<IPartService> _mockService;
        private PartsController _sut;

        [SetUp]
        public void SetUp()
        {
            _mockService = new Mock<IPartService>();
            _sut = new PartsController(_mockService.Object);
        }

        [Test]
        public async Task GetById_Success()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto<PartResponse>
            {
                IsSuccess = true,
                Data = new PartResponse { Id = id, SerialNumber = "SN123" }
            };
            _mockService.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetById(id);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByIdAsync(id), Times.Once);
        }

        [Test]
        public async Task Create_ValidData()
        {
            // Arrange
            var request = new CreatePartRequest
            {
                SerialNumber = "SN123",
                CategoryId = Guid.NewGuid()
            };
            var response = new BaseResponseDto<PartResponse>
            {
                IsSuccess = true,
                Data = new PartResponse { Id = Guid.NewGuid(), SerialNumber = "SN123" }
            };
            _mockService.Setup(x => x.CreateAsync(request)).ReturnsAsync(response);

            // Act
            var result = await _sut.Create(request);

            // Assert
            result.Should().BeOfType<CreatedAtActionResult>();
            var createdResult = result as CreatedAtActionResult;
            createdResult.ActionName.Should().Be(nameof(_sut.GetById));
            createdResult.Value.Should().Be(response);
            _mockService.Verify(x => x.CreateAsync(request), Times.Once);
        }

        [Test]
        public async Task ChangeStatus_Success()
        {
            // Arrange
            var id = Guid.NewGuid();
            var request = new PartChangeStatusRequest { Status = "Available" };
            var response = new BaseResponseDto<PartResponse>
            {
                IsSuccess = true,
                Data = new PartResponse { Id = id }
            };
            _mockService.Setup(x => x.ChangeStatusAsync(id, request)).ReturnsAsync(response);

            // Act
            var result = await _sut.ChangeStatus(id, request);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.ChangeStatusAsync(id, request), Times.Once);
        }

        [Test]
        public async Task ChangeStatus_InvalidStatus_ReturnsBadRequest()
        {
            // Arrange
            var id = Guid.NewGuid();
            var request = new PartChangeStatusRequest { Status = "Invalid" };
            var response = new BaseResponseDto<PartResponse>
            {
                IsSuccess = false,
                ErrorCode = "INVALID_STATUS"
            };
            _mockService.Setup(x => x.ChangeStatusAsync(id, request)).ReturnsAsync(response);

            // Act
            var result = await _sut.ChangeStatus(id, request);

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            var badRequest = result as BadRequestObjectResult;
            badRequest.Value.Should().Be(response);
            _mockService.Verify(x => x.ChangeStatusAsync(id, request), Times.Once);
        }

        [Test]
        public async Task ChangeCategory_Success()
        {
            // Arrange
            var id = Guid.NewGuid();
            var request = new ChangePartCategoryRequest { CategoryId = Guid.NewGuid() };
            var response = new BaseResponseDto<PartResponse>
            {
                IsSuccess = true,
                Data = new PartResponse { Id = id }
            };
            _mockService.Setup(x => x.ChangeCategoryAsync(id, request)).ReturnsAsync(response);

            // Act
            var result = await _sut.ChangeCategory(id, request);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.ChangeCategoryAsync(id, request), Times.Once);
        }

        [Test]
        public async Task Delete_ReservedPart_ReturnsBadRequest()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "PART_IN_USE"
            };
            _mockService.Setup(x => x.DeleteAsync(id)).ReturnsAsync(response);

            // Act
            var result = await _sut.Delete(id);

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            var badRequest = result as BadRequestObjectResult;
            badRequest.Value.Should().Be(response);
            _mockService.Verify(x => x.DeleteAsync(id), Times.Once);
        }

        [Test]
        public async Task GetAll_WithSerialNumberFilter()
        {
            // Arrange
            var serialNumber = "SN123";
            var response = new BaseResponseDto<PartResponse>
            {
                IsSuccess = true,
                Data = new PartResponse { SerialNumber = serialNumber }
            };
            _mockService.Setup(x => x.GetBySerialNumberAsync(serialNumber)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(serialNumber, null, null, null, null);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetBySerialNumberAsync(serialNumber), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithStatusFilter()
        {
            // Arrange
            var status = "Available";
            var response = new BaseResponseDto<IEnumerable<PartResponse>>
            {
                IsSuccess = true,
                Data = new List<PartResponse>()
            };
            _mockService.Setup(x => x.GetByStatusAsync(status)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(null, status, null, null, null);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByStatusAsync(status), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithCategoryIdFilter()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            var response = new BaseResponseDto<IEnumerable<PartResponse>>
            {
                IsSuccess = true,
                Data = new List<PartResponse>()
            };
            _mockService.Setup(x => x.GetByCategoryIdAsync(categoryId)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(null, null, categoryId, null, null);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByCategoryIdAsync(categoryId), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithSearch()
        {
            // Arrange
            var search = "test";
            var response = new BaseResponseDto<IEnumerable<PartResponse>>
            {
                IsSuccess = true,
                Data = new List<PartResponse>()
            };
            _mockService.Setup(x => x.SearchAsync(search)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(null, null, null, null, search);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.SearchAsync(search), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }
    }
}
