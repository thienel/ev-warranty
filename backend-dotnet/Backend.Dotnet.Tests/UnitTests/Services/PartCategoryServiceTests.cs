using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartCategoryDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class PartCategoryServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<IPartCategoryRepository> _categoryRepo;
        private PartCategoryService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _categoryRepo = new Mock<IPartCategoryRepository>();
            _unitOfWork.Setup(x => x.PartCategories).Returns(_categoryRepo.Object);
            _sut = new PartCategoryService(_unitOfWork.Object);
        }

        [Test]
        public async Task CreateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var request = new CreatePartCategoryRequest
            {
                CategoryName = "Battery",
                Description = "Battery parts"
            };

            _categoryRepo.Setup(x => x.CategoryNameExistsAsync(It.IsAny<string>(), null))
                .ReturnsAsync(false);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.CategoryName.Should().Be("Battery");
            _categoryRepo.Verify(x => x.AddAsync(It.IsAny<PartCategory>()), Times.Once);
        }

        [Test]
        public async Task CreateAsync_DuplicateCategoryName_ReturnsError()
        {
            // Arrange
            var request = new CreatePartCategoryRequest { CategoryName = "Existing" };
            _categoryRepo.Setup(x => x.CategoryNameExistsAsync("Existing", null))
                .ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_CATEGORY_NAME");
            _categoryRepo.Verify(x => x.AddAsync(It.IsAny<PartCategory>()), Times.Never);
        }

        [Test]
        public async Task CreateAsync_WithParentCategory_ValidatesParentExists()
        {
            // Arrange
            var parentId = Guid.NewGuid();
            var request = new CreatePartCategoryRequest
            {
                CategoryName = "SubCategory",
                ParentCategoryId = parentId
            };

            _categoryRepo.Setup(x => x.CategoryNameExistsAsync(It.IsAny<string>(), null))
                .ReturnsAsync(false);
            _categoryRepo.Setup(x => x.GetByIdAsync(parentId))
                .ReturnsAsync(new PartCategory("ParentCat"));

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            _categoryRepo.Verify(x => x.GetByIdAsync(parentId), Times.Once);
        }

        [Test]
        public async Task CreateAsync_ParentCategoryNotFound_ReturnsError()
        {
            // Arrange
            var request = new CreatePartCategoryRequest
            {
                CategoryName = "SubCat",
                ParentCategoryId = Guid.NewGuid()
            };

            _categoryRepo.Setup(x => x.CategoryNameExistsAsync(It.IsAny<string>(), null))
                .ReturnsAsync(false);
            _categoryRepo.Setup(x => x.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync((PartCategory)null);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("PARENT_CATEGORY_NOT_FOUND");
        }

        [Test]
        public async Task DeleteAsync_CategoryWithActiveParts_ReturnsError()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId))
                .ReturnsAsync(new PartCategory("Battery"));
            _categoryRepo.Setup(x => x.GetActivePartCountAsync(categoryId))
                .ReturnsAsync(5);

            // Act
            var result = await _sut.DeleteAsync(categoryId);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CATEGORY_HAS_ACTIVE_PARTS");
            result.Message.Should().Contain("5 active part(s)");
            _categoryRepo.Verify(x => x.Remove(It.IsAny<PartCategory>()), Times.Never);
        }

        [Test]
        public async Task DeleteAsync_CategoryWithChildren_ReturnsError()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId))
                .ReturnsAsync(new PartCategory("Parent"));
            _categoryRepo.Setup(x => x.GetActivePartCountAsync(categoryId))
                .ReturnsAsync(0);
            _categoryRepo.Setup(x => x.GetChildCategoryCountAsync(categoryId))
                .ReturnsAsync(3);

            // Act
            var result = await _sut.DeleteAsync(categoryId);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CATEGORY_HAS_CHILDREN");
            result.Message.Should().Contain("3 child categor(ies)");
        }

        [Test]
        public async Task DeleteAsync_LeafCategoryWithNoParts_DeletesSuccessfully()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            var category = new PartCategory("LeafCategory");

            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId))
                .ReturnsAsync(category);
            _categoryRepo.Setup(x => x.GetActivePartCountAsync(categoryId))
                .ReturnsAsync(0);
            _categoryRepo.Setup(x => x.GetChildCategoryCountAsync(categoryId))
                .ReturnsAsync(0);

            // Act
            var result = await _sut.DeleteAsync(categoryId);

            // Assert
            result.IsSuccess.Should().BeTrue();
            _categoryRepo.Verify(x => x.Remove(category), Times.Once);
        }

        [Test]
        public async Task DeleteAsync_CategoryNotFound_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            _categoryRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync((PartCategory)null);

            // Act
            var result = await _sut.DeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
            _categoryRepo.Verify(x => x.Remove(It.IsAny<PartCategory>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task GetWithHierarchyAsync_ExistingCategory_ReturnsHierarchy()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            var category = new PartCategory("Parent");
            // Mock Include relationships

            _categoryRepo.Setup(x => x.GetWithHierarchyAsync(categoryId))
                .ReturnsAsync(category);

            // Act
            var result = await _sut.GetWithHierarchyAsync(categoryId);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
        }

        [Test]
        public async Task GetFullHierarchyAsync_ReturnsOnlyRootCategories()
        {
            // Arrange
            var rootCategory1 = new PartCategory("Root1", "Description1");
            var rootCategory2 = new PartCategory("Root2", "Description2");
            var childCategory = new PartCategory("Child1", "Description3", rootCategory1.Id);

            var allCategories = new List<PartCategory> { rootCategory1, rootCategory2, childCategory };
            _categoryRepo.Setup(x => x.GetFullHierarchyAsync()).ReturnsAsync(allCategories);

            // Act
            var result = await _sut.GetFullHierarchyAsync();

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Should().HaveCount(2);
            result.Data.Should().OnlyContain(c => c.ParentCategoryId == null);
            _categoryRepo.Verify(x => x.GetFullHierarchyAsync(), Times.Once);
        }

        [Test]
        public async Task UpdateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var category = new PartCategory("OldName", "OldDesc");
            var request = new UpdatePartCategoryRequest
            {
                CategoryName = "NewName",
                Description = "NewDesc"
            };

            _categoryRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(category);
            _categoryRepo.Setup(x => x.CategoryNameExistsAsync(request.CategoryName, id)).ReturnsAsync(false);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.CategoryName.Should().Be("NewName");
            _categoryRepo.Verify(x => x.Update(category), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task UpdateAsync_DuplicateName_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var category = new PartCategory("OldName", "OldDesc");
            var request = new UpdatePartCategoryRequest
            {
                CategoryName = "ExistingName",
                Description = "NewDesc"
            };

            _categoryRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync(category);
            _categoryRepo.Setup(x => x.CategoryNameExistsAsync(request.CategoryName, id)).ReturnsAsync(true);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_CATEGORY_NAME");
            _categoryRepo.Verify(x => x.Update(It.IsAny<PartCategory>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task UpdateAsync_NotFound_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var request = new UpdatePartCategoryRequest
            {
                CategoryName = "NewName",
                Description = "NewDesc"
            };

            _categoryRepo.Setup(x => x.GetByIdAsync(id)).ReturnsAsync((PartCategory)null);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
            _categoryRepo.Verify(x => x.Update(It.IsAny<PartCategory>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }
    }
}