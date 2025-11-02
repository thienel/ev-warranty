using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Services;
using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.WarrantyPolicyDto;

namespace Backend.Dotnet.Tests.UnitTests.Services
{
    [TestFixture]
    public class WarrantyPolicyServiceTests
    {
        private Mock<IUnitOfWork> _unitOfWork;
        private Mock<IWarrantyPolicyRepository> _warrantyRepo;
        private WarrantyPolicyService _sut;

        [SetUp]
        public void Setup()
        {
            _unitOfWork = new Mock<IUnitOfWork>();
            _warrantyRepo = new Mock<IWarrantyPolicyRepository>();
            _unitOfWork.Setup(x => x.WarrantyPolicies).Returns(_warrantyRepo.Object);
            _sut = new WarrantyPolicyService(_unitOfWork.Object);
        }
        [Test]
        public async Task CreateAsync_WithDuplicateName_ReturnsDuplicateError()
        {
            // Arrange
            var request = new CreateWarrantyPolicyRequest
            {
                PolicyName = "Existing Policy",
                WarrantyDurationMonths = 24,
                TermsAndConditions = "Terms"
            };
            _warrantyRepo.Setup(x => x.PolicyNameExistsAsync(request.PolicyName, It.IsAny<Guid?>())).ReturnsAsync(true);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.Message.Should().Be("Policy name already exists");
            result.ErrorCode.Should().Be("DUPLICATE_POLICY_NAME");
            _warrantyRepo.Verify(x => x.AddAsync(It.IsAny<WarrantyPolicy>()), Times.Never);
        }

        [Test]
        public async Task CreateAsync_WithValidRequest_ReturnsSuccess()
        {
            // Arrange
            var request = new CreateWarrantyPolicyRequest
            {
                PolicyName = "Standard Warranty",
                WarrantyDurationMonths = 36,
                KilometerLimit = 100000,
                TermsAndConditions = "Standard terms"
            };
            _warrantyRepo.Setup(x => x.PolicyNameExistsAsync(request.PolicyName, It.IsAny<Guid?>())).ReturnsAsync(false);
            _warrantyRepo.Setup(x => x.AddAsync(It.IsAny<WarrantyPolicy>())).Returns(Task.CompletedTask);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.CreateAsync(request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Message.Should().Be("Warranty policy created successfully");
            result.Data.Should().NotBeNull();
            result.Data!.PolicyName.Should().Be(request.PolicyName);
            _warrantyRepo.Verify(x => x.AddAsync(It.IsAny<WarrantyPolicy>()), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task UpdateAsync_WithValidRequest_ReturnsSuccess()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var policy = new WarrantyPolicy("Old Name", 24, 50000, "Old Terms");
            var request = new UpdateWarrantyPolicyRequest
            {
                PolicyName = "New Name",
                WarrantyDurationMonths = 36,
                KilometerLimit = 100000,
                TermsAndConditions = "New Terms"
            };
            _warrantyRepo.Setup(x => x.GetByIdAsync(policyId)).ReturnsAsync(policy);
            _warrantyRepo.Setup(x => x.PolicyNameExistsAsync(request.PolicyName, policyId)).ReturnsAsync(false);
            _unitOfWork.Setup(x => x.SaveChangesAsync()).ReturnsAsync(1);

            // Act
            var result = await _sut.UpdateAsync(policyId, request);

            // Assert
            result.IsSuccess.Should().BeTrue();
            result.Message.Should().Be("Warranty policy updated successfully");
            _warrantyRepo.Verify(x => x.Update(policy), Times.Once);
            _unitOfWork.Verify(x => x.SaveChangesAsync(), Times.Once);
        }

        [Test]
        public async Task UpdateAsync_WithDuplicateName_ReturnsDuplicateError()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var policy = new WarrantyPolicy("Old Name", 24, 50000, "Terms");
            var request = new UpdateWarrantyPolicyRequest
            {
                PolicyName = "Duplicate Name",
                WarrantyDurationMonths = 36,
                TermsAndConditions = "Terms"
            };
            _warrantyRepo.Setup(x => x.GetByIdAsync(policyId)).ReturnsAsync(policy);
            _warrantyRepo.Setup(x => x.PolicyNameExistsAsync(request.PolicyName, policyId)).ReturnsAsync(true);

            // Act
            var result = await _sut.UpdateAsync(policyId, request);

            // Assert
            result.IsSuccess.Should().BeFalse();
            result.ErrorCode.Should().Be("DUPLICATE_POLICY_NAME");
            _warrantyRepo.Verify(x => x.Update(It.IsAny<WarrantyPolicy>()), Times.Never);
        }
    }
}
