using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartCategoryDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class PartCategoriesControllerTests
    {
        private Mock<IPartCategoryService> _mockService;
        private PartCategoriesController _sut;

        [SetUp]
        public void SetUp()
        {
            _mockService = new Mock<IPartCategoryService>();
            _sut = new PartCategoriesController(_mockService.Object);
        }

        [Test]
        public async Task GetById_Success()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto<PartCategoryResponse>
            {
                IsSuccess = true,
                Data = new PartCategoryResponse { Id = id, CategoryName = "Test" }
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
        public async Task Create_WithParent_Success()
        {
            // Arrange
            var parentId = Guid.NewGuid();
            var request = new CreatePartCategoryRequest
            {
                CategoryName = "Child",
                ParentCategoryId = parentId
            };
            var response = new BaseResponseDto<PartCategoryResponse>
            {
                IsSuccess = true,
                Data = new PartCategoryResponse { Id = Guid.NewGuid(), CategoryName = "Child" }
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
        public async Task Create_ParentNotFound()
        {
            // Arrange
            var request = new CreatePartCategoryRequest
            {
                CategoryName = "Test",
                ParentCategoryId = Guid.NewGuid()
            };
            var response = new BaseResponseDto<PartCategoryResponse>
            {
                IsSuccess = false,
                ErrorCode = "PARENT_CATEGORY_NOT_FOUND"
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
        public async Task Delete_WithChildren_ReturnsBadRequest()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "CATEGORY_HAS_CHILDREN"
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
        public async Task GetWithHierarchy_Success()
        {
            // Arrange
            var id = Guid.NewGuid();
            var response = new BaseResponseDto<PartCategoryWithHierarchyResponse>
            {
                IsSuccess = true,
                Data = new PartCategoryWithHierarchyResponse { Id = id }
            };
            _mockService.Setup(x => x.GetWithHierarchyAsync(id)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetWithHierarchy(id);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetWithHierarchyAsync(id), Times.Once);
        }

        [Test]
        public async Task GetFullHierarchy_Success()
        {
            // Arrange
            var response = new BaseResponseDto<IEnumerable<PartCategoryWithHierarchyResponse>>
            {
                IsSuccess = true,
                Data = new List<PartCategoryWithHierarchyResponse>()
            };
            _mockService.Setup(x => x.GetFullHierarchyAsync()).ReturnsAsync(response);

            // Act
            var result = await _sut.GetFullHierarchy();

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetFullHierarchyAsync(), Times.Once);
        }

        [Test]
        public async Task GetAll_WithParentIdFilter()
        {
            // Arrange
            var parentId = Guid.NewGuid();
            var response = new BaseResponseDto<IEnumerable<PartCategoryResponse>>
            {
                IsSuccess = true,
                Data = new List<PartCategoryResponse>()
            };
            _mockService.Setup(x => x.GetByParentIdAsync(parentId)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(null, parentId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByParentIdAsync(parentId), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithNameFilter()
        {
            // Arrange
            var name = "Test";
            var response = new BaseResponseDto<PartCategoryResponse>
            {
                IsSuccess = true,
                Data = new PartCategoryResponse { CategoryName = name }
            };
            _mockService.Setup(x => x.GetByCategoryNameAsync(name)).ReturnsAsync(response);

            // Act
            var result = await _sut.GetAll(name, null);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().Be(response);
            _mockService.Verify(x => x.GetByCategoryNameAsync(name), Times.Once);
            _mockService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_NoFilters()
        {
            // Arrange
            var response = new BaseResponseDto<IEnumerable<PartCategoryResponse>>
            {
                IsSuccess = true,
                Data = new List<PartCategoryResponse>()
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
