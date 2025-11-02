using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class PolicyCoveragePartsControllerTests
    {
        private Mock<IPolicyCoveragePartService> _mockService;
        private PolicyCoveragePartsController _sut;

        [SetUp]
        public void SetUp()
        {
            _mockService = new Mock<IPolicyCoveragePartService>();
            _sut = new PolicyCoveragePartsController(_mockService.Object);
        }

        [Test]
        public async Task GetById_Success()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto<PolicyCoveragePartResponse>
            {
                IsSuccess = true,
                Data = new PolicyCoveragePartResponse { Id = id }
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
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = Guid.NewGuid(),
                PartCategoryId = Guid.NewGuid()
            };
            var response = new BaseResponseDto<PolicyCoveragePartResponse>
            {
                IsSuccess = true,
                Data = new PolicyCoveragePartResponse { Id = Guid.NewGuid() }
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
        public async Task Create_PolicyNotFound()
        {
            // Arrange
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = Guid.NewGuid(),
                PartCategoryId = Guid.NewGuid()
            };
            var response = new BaseResponseDto<PolicyCoveragePartResponse>
            {
                IsSuccess = false,
                ErrorCode = "POLICY_NOT_FOUND"
            };
            _mockService.Setup(x => x.CreateAsync(request)).ReturnsAsync(response);

            // Act
            var result = await _sut.Create(request);

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            var badRequest = result as BadRequestObjectResult;
            badRequest.Value.Should().Be(response);
            _mockService.Verify(x => x.CreateAsync(request), Times.Once);
        }

        [Test]
        public async Task Delete_PolicyNotEditable_ReturnsBadRequest()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "POLICY_NOT_EDITABLE"
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
        public async Task GetAll_WithPolicyIdFilter()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var response = new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
            {
                IsSuccess = true,
                Data = new List<PolicyCoveragePartResponse>()
            };
            _mockService.Setup(x => x.GetByPolicyIdAsync(policyId)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(policyId, null);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByPolicyIdAsync(policyId), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithCategoryIdFilter()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            var response = new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
            {
                IsSuccess = true,
                Data = new List<PolicyCoveragePartResponse>()
            };
            _mockService.Setup(x => x.GetByPartCategoryIdAsync(categoryId)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(null, categoryId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByPartCategoryIdAsync(categoryId), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithBothFilters_ReturnsSingleItem()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var singleResponse = new BaseResponseDto<PolicyCoveragePartResponse>
            {
                IsSuccess = true,
                Data = new PolicyCoveragePartResponse { Id = Guid.NewGuid() }
            };
            _mockService.Setup(x => x.GetByPolicyAndCategoryAsync(policyId, categoryId))
                .ReturnsAsync(singleResponse);

            // Act
            var result = await _sut.GetAll(policyId, categoryId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            var response = okResult.Value as BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>;
            response.Should().NotBeNull();
            response.IsSuccess.Should().BeTrue();
            response.Data.Should().HaveCount(1);
            response.Data.First().Should().Be(singleResponse.Data);
            _mockService.Verify(x => x.GetByPolicyAndCategoryAsync(policyId, categoryId), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_NoFilters()
        {
            // Arrange
            var response = new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
            {
                IsSuccess = true,
                Data = new List<PolicyCoveragePartResponse>()
            };
            _mockService.Setup(x => x.GetAllAsync()).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(null, null);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetAllAsync(), Times.Once);
        }
    }
}
