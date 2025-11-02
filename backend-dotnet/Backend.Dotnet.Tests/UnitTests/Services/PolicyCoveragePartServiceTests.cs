using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class PolicyCoveragePartServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<IPolicyCoveragePartRepository> _coverageRepo;
        private Mock<IWarrantyPolicyRepository> _policyRepo;
        private Mock<IPartCategoryRepository> _categoryRepo;
        private PolicyCoveragePartService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _coverageRepo = new Mock<IPolicyCoveragePartRepository>();
            _policyRepo = new Mock<IWarrantyPolicyRepository>();
            _categoryRepo = new Mock<IPartCategoryRepository>();
            _unitOfWork.Setup(x => x.PolicyCoverageParts).Returns(_coverageRepo.Object);
            _unitOfWork.Setup(x => x.WarrantyPolicies).Returns(_policyRepo.Object);
            _unitOfWork.Setup(x => x.PartCategories).Returns(_categoryRepo.Object);
            _sut = new PolicyCoveragePartService(_unitOfWork.Object);
        }
        /*
        [Test]
        public async Task CreateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = policyId,
                PartCategoryId = categoryId,
                CoverageConditions = "Test conditions"
            };

            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");
            var category = new PartCategory("Test Category", "Description");
            var coverage = new PolicyCoveragePart(policyId, categoryId, "Test conditions");

            _policyRepo.Setup(x => x.GetByIdAsync(policyId)).ReturnsAsync(policy);
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId)).ReturnsAsync(category);
            _coverageRepo.Setup(x => x.ExistsByPolicyAndCategoryAsync(policyId, categoryId, null)).ReturnsAsync(false);
            _coverageRepo.Setup(x => x.GetWithDetailsAsync(It.IsAny<Guid>())).ReturnsAsync(coverage);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            _coverageRepo.Verify(x => x.AddAsync(It.IsAny<PolicyCoveragePart>()), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }
        */
        [Test]
        public async Task CreateAsync_PolicyNotFound_ReturnsError()
        {
            // Arrange
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = Guid.NewGuid(),
                PartCategoryId = Guid.NewGuid()
            };

            _policyRepo.Setup(x => x.GetByIdAsync(request.PolicyId)).ReturnsAsync((WarrantyPolicy)null);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("POLICY_NOT_FOUND");
            _coverageRepo.Verify(x => x.AddAsync(It.IsAny<PolicyCoveragePart>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task CreateAsync_CategoryNotFound_ReturnsError()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = policyId,
                PartCategoryId = Guid.NewGuid()
            };

            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");

            _policyRepo.Setup(x => x.GetByIdAsync(policyId)).ReturnsAsync(policy);
            _categoryRepo.Setup(x => x.GetByIdAsync(request.PartCategoryId)).ReturnsAsync((PartCategory)null);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("CATEGORY_NOT_FOUND");
            _coverageRepo.Verify(x => x.AddAsync(It.IsAny<PolicyCoveragePart>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }

        [Test]
        public async Task CreateAsync_DuplicateCoverage_ReturnsError()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = policyId,
                PartCategoryId = categoryId
            };

            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");
            var category = new PartCategory("Test Category", "Description");

            _policyRepo.Setup(x => x.GetByIdAsync(policyId)).ReturnsAsync(policy);
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId)).ReturnsAsync(category);
            _coverageRepo.Setup(x => x.ExistsByPolicyAndCategoryAsync(policyId, categoryId, null)).ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_COVERAGE");
            _coverageRepo.Verify(x => x.AddAsync(It.IsAny<PolicyCoveragePart>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }
        /*
        [Test]
        public async Task CreateAsync_BusinessRuleViolation_ReturnsError()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var request = new CreatePolicyCoveragePartRequest
            {
                PolicyId = policyId,
                PartCategoryId = categoryId
            };

            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");
            policy.Activate();
            var category = new PartCategory("Test Category", "Description");

            _policyRepo.Setup(x => x.GetByIdAsync(policyId)).ReturnsAsync(policy);
            _categoryRepo.Setup(x => x.GetByIdAsync(categoryId)).ReturnsAsync(category);
            _coverageRepo.Setup(x => x.ExistsByPolicyAndCategoryAsync(policyId, categoryId, null)).ReturnsAsync(false);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().NotBeNullOrEmpty();
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }
        */
        /*
        [Test]
        public async Task DeleteAsync_PolicyNotEditable_ReturnsError()
        {
            // Arrange
            var id = Guid.NewGuid();
            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");
            policy.Activate();
            var coverage = new PolicyCoveragePart(policy.Id, Guid.NewGuid());

            _coverageRepo.Setup(x => x.GetWithDetailsAsync(id)).ReturnsAsync(coverage);

            // Act
            var result = await _sut.DeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("POLICY_NOT_EDITABLE");
            _coverageRepo.Verify(x => x.Remove(It.IsAny<PolicyCoveragePart>()), Times.Never);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Never);
        }
        */
        [Test]
        public async Task DeleteAsync_DraftPolicy_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");
            var coverage = new PolicyCoveragePart(policy.Id, Guid.NewGuid());

            _coverageRepo.Setup(x => x.GetWithDetailsAsync(id)).ReturnsAsync(coverage);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.DeleteAsync(id);

            // Assert
            result.IsSuccess.Should().BeTrue();
            _coverageRepo.Verify(x => x.Remove(coverage), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task UpdateAsync_ValidData_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var policy = new WarrantyPolicy("Test Policy", 12, 10000, "Terms");
            var coverage = new PolicyCoveragePart(policy.Id, Guid.NewGuid(), "Old conditions");
            var request = new UpdatePolicyCoveragePartRequest
            {
                CoverageConditions = "New conditions"
            };

            _coverageRepo.Setup(x => x.GetWithDetailsAsync(id)).ReturnsAsync(coverage);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.UpdateAsync(id, request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            _coverageRepo.Verify(x => x.Update(coverage), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task GetByPolicyIdAsync_NoCoverage_ReturnsError()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            _coverageRepo.Setup(x => x.GetByPolicyIdAsync(policyId)).ReturnsAsync(new List<PolicyCoveragePart>());

            // Act
            var result = await _sut.GetByPolicyIdAsync(policyId);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("NOT_FOUND");
        }

        [Test]
        public async Task GetByPolicyAndCategoryAsync_Found_ReturnsSuccess()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var coverage = new PolicyCoveragePart(policyId, categoryId);

            _coverageRepo.Setup(x => x.GetByPolicyAndCategoryAsync(policyId, categoryId)).ReturnsAsync(coverage);

            // Act
            var result = await _sut.GetByPolicyAndCategoryAsync(policyId, categoryId);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.PolicyId.Should().Be(policyId);
            result.Data.PartCategoryId.Should().Be(categoryId);
        }

        [Test]
        public async Task GetByPartCategoryIdAsync_Found_ReturnsSuccess()
        {
            // Arrange
            var categoryId = Guid.NewGuid();
            var coverages = new List<PolicyCoveragePart>
            {
                new PolicyCoveragePart(Guid.NewGuid(), categoryId)
            };

            _coverageRepo.Setup(x => x.FindAsync(It.IsAny<System.Linq.Expressions.Expression<System.Func<PolicyCoveragePart, bool>>>()))
                .ReturnsAsync(coverages);

            // Act
            var result = await _sut.GetByPartCategoryIdAsync(categoryId);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Should().HaveCount(1);
        }

        [Test]
        public async Task GetWithDetailsAsync_Found_ReturnsSuccess()
        {
            // Arrange
            var id = Guid.NewGuid();
            var coverage = new PolicyCoveragePart(Guid.NewGuid(), Guid.NewGuid());

            _coverageRepo.Setup(x => x.GetWithDetailsAsync(id)).ReturnsAsync(coverage);

            // Act
            var result = await _sut.GetWithDetailsAsync(id);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Data.Should().NotBeNull();
            result.Data.Id.Should().Be(coverage.Id);
        }
    }
}
