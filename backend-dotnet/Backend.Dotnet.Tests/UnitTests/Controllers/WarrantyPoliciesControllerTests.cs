using Backend.Dotnet.API.Controllers;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Domain.Entities;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;
using static Backend.Dotnet.Application.DTOs.WarrantyPolicyDto;

namespace Backend.Dotnet.Tests.UnitTests.Controllers
{
    [TestFixture]
    public class WarrantyPoliciesControllerTests
    {
        private Mock<IWarrantyPolicyService> _warrantyPolicyService;
        private WarrantyPoliciesController _controller;

        [SetUp]
        public void Setup()
        {
            _warrantyPolicyService = new Mock<IWarrantyPolicyService>();
            _controller = new WarrantyPoliciesController(_warrantyPolicyService.Object);
        }
        
        [Test]
        public async Task GetById_ExistingPolicy_ReturnsOk()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = true,
                Data = new WarrantyPolicyResponse
                {
                    Id = policyId,
                    PolicyName = "Premium EV Warranty"
                }
            };

            _warrantyPolicyService.Setup(x => x.GetByIdAsync(policyId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(policyId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            okResult.Value.Should().BeEquivalentTo(serviceResponse);
        }

        [Test]
        public async Task GetById_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _warrantyPolicyService.Setup(x => x.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetById(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetWithDetails_ExistingPolicy_ReturnsOk()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<WarrantyPolicyWithDetailsResponse>
            {
                IsSuccess = true,
                Data = new WarrantyPolicyWithDetailsResponse
                {
                    Id = policyId,
                    PolicyName = "Premium",
                    CoveredParts = new List<PolicyCoveragePartResponse>()
                }
            };

            _warrantyPolicyService.Setup(x => x.GetWithDetailsAsync(policyId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetWithDetails(policyId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task GetWithDetails_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyWithDetailsResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _warrantyPolicyService.Setup(x => x.GetWithDetailsAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetWithDetails(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_WithStatusFilter_CallsGetByStatusAsync()
        {
            // Arrange
            var status = "Active";
            var serviceResponse = new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
            {
                IsSuccess = true,
                Data = new List<WarrantyPolicyResponse>()
            };

            _warrantyPolicyService.Setup(x => x.GetByStatusAsync(status))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(status: status);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _warrantyPolicyService.Verify(x => x.GetByStatusAsync(status), Times.Once);
            _warrantyPolicyService.Verify(x => x.GetAllAsync(), Times.Never);
        }

        [Test]
        public async Task GetAll_WithInvalidStatus_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
            {
                IsSuccess = false,
                ErrorCode = "INVALID_STATUS"
            };

            _warrantyPolicyService.Setup(x => x.GetByStatusAsync(It.IsAny<string>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(status: "InvalidStatus");

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task GetAll_WithPolicyNameFilter_CallsGetByPolicyNameAsync()
        {
            // Arrange
            var policyName = "Premium Warranty";
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = true,
                Data = new WarrantyPolicyResponse { PolicyName = policyName }
            };

            _warrantyPolicyService.Setup(x => x.GetByPolicyNameAsync(policyName))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(policyName: policyName);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _warrantyPolicyService.Verify(x => x.GetByPolicyNameAsync(policyName), Times.Once);
            _warrantyPolicyService.Verify(x => x.GetByStatusAsync(It.IsAny<string>()), Times.Never);
        }

        [Test]
        public async Task GetAll_WithPolicyNameFilter_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _warrantyPolicyService.Setup(x => x.GetByPolicyNameAsync(It.IsAny<string>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll(policyName: "Unknown");

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task GetAll_NoFilters_CallsGetAllAsync()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
            {
                IsSuccess = true,
                Data = new List<WarrantyPolicyResponse>()
            };

            _warrantyPolicyService.Setup(x => x.GetAllAsync())
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll();

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            _warrantyPolicyService.Verify(x => x.GetAllAsync(), Times.Once);
        }

        [Test]
        public async Task GetAll_NoFilters_Failed_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
            {
                IsSuccess = false,
                ErrorCode = "INTERNAL_ERROR"
            };

            _warrantyPolicyService.Setup(x => x.GetAllAsync())
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.GetAll();

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Create_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("PolicyName", "Required");

            // Act
            var result = await _controller.Create(new CreateWarrantyPolicyRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _warrantyPolicyService.Verify(x => x.CreateAsync(It.IsAny<CreateWarrantyPolicyRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Create_ValidData_ReturnsCreatedAtAction()
        {
            // Arrange
            var request = new CreateWarrantyPolicyRequest
            {
                PolicyName = "Premium Warranty",
                WarrantyDurationMonths = 36,
                KilometerLimit = 50000
            };
            var policyId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = true,
                Data = new WarrantyPolicyResponse
                {
                    Id = policyId,
                    PolicyName = "Premium Warranty"
                }
            };

            _warrantyPolicyService.Setup(x => x.CreateAsync(request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(request);

            // Assert
            result.Should().BeOfType<CreatedAtActionResult>();
            var createdResult = result as CreatedAtActionResult;
            createdResult.ActionName.Should().Be(nameof(_controller.GetById));
            createdResult.RouteValues["id"].Should().Be(policyId);
        }

        [Test]
        public async Task Create_DuplicatePolicyName_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_POLICY_NAME"
            };

            _warrantyPolicyService.Setup(x => x.CreateAsync(It.IsAny<CreateWarrantyPolicyRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Create(new CreateWarrantyPolicyRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Update_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("DurationMonths", "Invalid");

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateWarrantyPolicyRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _warrantyPolicyService.Verify(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateWarrantyPolicyRequest>()),
                Times.Never);
        }

        [Test]
        public async Task Update_ValidData_ReturnsOk()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var request = new UpdateWarrantyPolicyRequest();
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = true,
                Data = new WarrantyPolicyResponse { Id = policyId }
            };

            _warrantyPolicyService.Setup(x => x.UpdateAsync(policyId, request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(policyId, request);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Update_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _warrantyPolicyService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateWarrantyPolicyRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateWarrantyPolicyRequest());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Update_BusinessError_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "DUPLICATE_POLICY_NAME"
            };

            _warrantyPolicyService.Setup(x => x.UpdateAsync(It.IsAny<Guid>(), It.IsAny<UpdateWarrantyPolicyRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Update(Guid.NewGuid(), new UpdateWarrantyPolicyRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task ChangeStatus_InvalidModelState_ReturnsBadRequest()
        {
            // Arrange
            _controller.ModelState.AddModelError("Status", "Required");

            // Act
            var result = await _controller.ChangeStatus(Guid.NewGuid(), new ChangeStatusRequest());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
            _warrantyPolicyService.Verify(x => x.ChangeStatusAsync(It.IsAny<Guid>(), It.IsAny<ChangeStatusRequest>()),
                Times.Never);
        }

        [Test]
        public async Task ChangeStatus_ValidStatus_ReturnsOk()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var request = new ChangeStatusRequest { Status = "Active" };
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = true,
                Data = new WarrantyPolicyResponse
                {
                    Id = policyId,
                    Status = WarrantyPolicyStatus.Active.ToString()
                }
            };

            _warrantyPolicyService.Setup(x => x.ChangeStatusAsync(policyId, request))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.ChangeStatus(policyId, request);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
            var okResult = result as OkObjectResult;
            var response = okResult.Value as BaseResponseDto<WarrantyPolicyResponse>;
            response.Data.Status.Should().Be(WarrantyPolicyStatus.Active.ToString());
        }

        [Test]
        public async Task ChangeStatus_InvalidStatus_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "INVALID_STATUS"
            };

            _warrantyPolicyService.Setup(x => x.ChangeStatusAsync(It.IsAny<Guid>(), It.IsAny<ChangeStatusRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.ChangeStatus(Guid.NewGuid(),
                new ChangeStatusRequest { Status = "InvalidStatus" });

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task ChangeStatus_PolicyNotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _warrantyPolicyService.Setup(x => x.ChangeStatusAsync(It.IsAny<Guid>(), It.IsAny<ChangeStatusRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.ChangeStatus(Guid.NewGuid(), new ChangeStatusRequest());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task ChangeStatus_BusinessRuleViolation_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto<WarrantyPolicyResponse>
            {
                IsSuccess = false,
                ErrorCode = "POLICY_HAS_NO_COVERAGE"
            };

            _warrantyPolicyService.Setup(x => x.ChangeStatusAsync(It.IsAny<Guid>(), It.IsAny<ChangeStatusRequest>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.ChangeStatus(Guid.NewGuid(),
                new ChangeStatusRequest { Status = "Active" });

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }

        [Test]
        public async Task Delete_ExistingPolicy_ReturnsOk()
        {
            // Arrange
            var policyId = Guid.NewGuid();
            var serviceResponse = new BaseResponseDto
            {
                IsSuccess = true,
                Message = "Policy deleted successfully"
            };

            _warrantyPolicyService.Setup(x => x.DeleteAsync(policyId))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Delete(policyId);

            // Assert
            result.Should().BeOfType<OkObjectResult>();
        }

        [Test]
        public async Task Delete_NotFound_ReturnsNotFound()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "NOT_FOUND"
            };

            _warrantyPolicyService.Setup(x => x.DeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Delete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<NotFoundObjectResult>();
        }

        [Test]
        public async Task Delete_PolicyInUse_ReturnsBadRequest()
        {
            // Arrange
            var serviceResponse = new BaseResponseDto
            {
                IsSuccess = false,
                ErrorCode = "POLICY_ASSIGNED_TO_VEHICLES"
            };

            _warrantyPolicyService.Setup(x => x.DeleteAsync(It.IsAny<Guid>()))
                .ReturnsAsync(serviceResponse);

            // Act
            var result = await _controller.Delete(Guid.NewGuid());

            // Assert
            result.Should().BeOfType<BadRequestObjectResult>();
        }
    }
}
