using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class PartServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<IPartRepository> _partRepo;
        private Mock<IPartCategoryRepository> _categoryRepo;
        private PartService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _partRepo = new Mock<IPartRepository>();
            _categoryRepo = new Mock<IPartCategoryRepository>();

            _unitOfWork.Setup(x => x.Parts).Returns(_partRepo.Object);
            _unitOfWork.Setup(x => x.PartCategories).Returns(_categoryRepo.Object);

            _sut = new PartService(_unitOfWork.Object);
        }
        /*
        [Test]
        public async Task CreateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var request = new CreatePartRequest
            {
                SerialNumber = "HEATC3WAYV1SVC000",
                PartName = "Heat Gird",
                UnitPrice = 80000,
                CategoryId = Guid.NewGuid(),
                OfficeLocationId = Guid.NewGuid()
            };

            _partRepo.Setup(x => x.SerialNumberExistsAsync(It.IsAny<string>(), null))
                .ReturnsAsync(false);
            _categoryRepo.Setup(x => x.GetByIdAsync(Guid.NewGuid()))
                .ReturnsAsync(new PartCategory("Cooling"));

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.SerialNumber.Should().Be("BAT-001");
            _partRepo.Verify(x => x.AddAsync(It.IsAny<Part>()), Times.Once);
        }
        */
        [Test]
        public async Task CreateAsync_DuplicateSerialNumber_ReturnsError()
        {
            // Arrange
            var request = new CreatePartRequest { SerialNumber = "EXIST-001" };
            _partRepo.Setup(x => x.SerialNumberExistsAsync("EXIST-001", null))
                .ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_SERIAL_NUMBER");
        }

        [Test]
        public async Task CreateAsync_CategoryNotFound_ReturnsError()
        {
            // Arrange
            var request = new CreatePartRequest
            {
                SerialNumber = "NEW-001",
                CategoryId = Guid.NewGuid()
            };

            _partRepo.Setup(x => x.SerialNumberExistsAsync(It.IsAny<string>(), null))
                .ReturnsAsync(false);
            _categoryRepo.Setup(x => x.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync((PartCategory)null);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CATEGORY_NOT_FOUND");
        }

        [Test]
        public async Task CreateAsync_InvalidData_ReturnsBusinessRuleError() { /* TODO */ }

        [Test]
        public async Task ChangeStatusAsync_ValidStatus_UpdatesSuccessfully()
        {
            // Arrange
            var partId = Guid.NewGuid();
            var part = new Part("VALVE3WAYV1SVC001", "Test Part", 300000, Guid.NewGuid(), null);

            _partRepo.Setup(x => x.GetByIdAsync(partId)).ReturnsAsync(part);

            var request = new PartChangeStatusRequest { Status = "Reserved" };

            // Act
            var result = await _sut.ChangeStatusAsync(partId, request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Status.Should().Be(PartStatus.Reserved.ToString());
        }

        [Test]
        public async Task ChangeStatusAsync_InvalidStatus_ReturnsError()
        {
            // Arrange
            var partId = Guid.NewGuid();
            _partRepo.Setup(x => x.GetByIdAsync(partId))
                .ReturnsAsync(new Part("VALVE3WAYV1SVC111", "Test Change Part", 2500000, Guid.NewGuid(), null));

            var request = new PartChangeStatusRequest { Status = "InvalidStatus" };

            // Act
            var result = await _sut.ChangeStatusAsync(partId, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("INVALID_STATUS");
            result.Message.Should().Contain("Available, Reserved, Installed");
        }

        [Test]
        public async Task ChangeStatusAsync_PartNotFound_ReturnsError() { }

        [Test]
        public async Task DeleteAsync_PartWithReservedStatus_ReturnsError()
        {
            // Arrange
            var partId = Guid.NewGuid();
            var part = new Part("BATERY3WAYV1SVC001", "Test Battery", 10000, Guid.NewGuid(), null);
            part.ChangeStatus(PartStatus.Reserved);

            _partRepo.Setup(x => x.GetByIdAsync(partId)).ReturnsAsync(part);

            // Act
            var result = await _sut.DeleteAsync(partId);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("PART_IN_USE");
            result.Message.Should().Contain("Reserved");
            _partRepo.Verify(x => x.Remove(It.IsAny<Part>()), Times.Never);
        }

        [Test]
        public async Task DeleteAsync_PartWithInstalledStatus_ReturnsError()
        {
            // Arrange
            var part = new Part("NGINE3WAYV1SVC001", "Test Engine", 59400000, Guid.NewGuid(), null);
            part.ChangeStatus(PartStatus.Installed);

            _partRepo.Setup(x => x.GetByIdAsync(It.IsAny<Guid>())).ReturnsAsync(part);

            // Act
            var result = await _sut.DeleteAsync(Guid.NewGuid());

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("PART_IN_USE");
        }

        [Test]
        public async Task DeleteAsync_AvailablePart_DeletesSuccessfully()
        {
            // Arrange
            var partId = Guid.NewGuid();
            var part = new Part("FDOOR3WAYV1SVC001", "Test Front Door", 300000, Guid.NewGuid(), null);

            _partRepo.Setup(x => x.GetByIdAsync(partId)).ReturnsAsync(part);

            // Act
            var result = await _sut.DeleteAsync(partId);

            // Assert
            result.IsSuccess.Should().BeTrue();
            _partRepo.Verify(x => x.Remove(part), Times.Once);
        }

        [Test]
        public async Task ChangeCategoryAsync_ValidCategory_UpdatesSuccessfully()
        {
            // Arrange
            var partId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var part = new Part("SN123", "PartName", 100m, Guid.NewGuid());
            var category = new PartCategory("NewCategory", "Description");
            var request = new ChangePartCategoryRequest { CategoryId = categoryId };

            _partRepo.Setup(x => x.GetByIdAsync(partId)).ReturnsAsync(part);
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId)).ReturnsAsync(category);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.ChangeCategoryAsync(partId, request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.CategoryId.Should().Be(categoryId);
            _partRepo.Verify(x => x.Update(part), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task ChangeCategoryAsync_CategoryNotFound_ReturnsError()
        {
            // Arrange
            var partId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var part = new Part("SN123", "PartName", 100m, Guid.NewGuid());
            var request = new ChangePartCategoryRequest { CategoryId = categoryId };

            _partRepo.Setup(x => x.GetByIdAsync(partId)).ReturnsAsync(part);
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId)).ReturnsAsync((PartCategory)null);

            // Act
            var result = await _sut.ChangeCategoryAsync(partId, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CATEGORY_NOT_FOUND");
            _partRepo.Verify(x => x.Update(It.IsAny<Part>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task GetByStatusAsync_ValidStatus_ReturnsPartsList()
        {
            // Arrange
            var parts = new List<Part>
            {
                new Part("SN001", "Part1", 100m, Guid.NewGuid()),
                new Part("SN002", "Part2", 200m, Guid.NewGuid())
            };

            _partRepo.Setup(x => x.GetByStatusAsync(PartStatus.Available)).ReturnsAsync(parts);

            // Act
            var result = await _sut.GetByStatusAsync("Available");

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Should().HaveCount(2);
            _partRepo.Verify(x => x.GetByStatusAsync(PartStatus.Available), Times.Once);
        }

        [Test]
        public async Task GetByStatusAsync_InvalidStatus_ReturnsError()
        {
            // Arrange
            var invalidStatus = "InvalidStatus";

            // Act
            var result = await _sut.GetByStatusAsync(invalidStatus);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("INVALID_STATUS");
            result.Message.Should().Contain("Invalid status value");
            _partRepo.Verify(x => x.GetByStatusAsync(It.IsAny<PartStatus>()), Times.Never);
        }
    }
}
