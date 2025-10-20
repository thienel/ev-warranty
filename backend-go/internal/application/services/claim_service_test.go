package services_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/mocks"
)

var _ = Describe("ClaimService", func() {
	var (
		mockLogger     *mocks.Logger
		mockClaimRepo  *mocks.ClaimRepository
		mockItemRepo   *mocks.ClaimItemRepository
		mockAttachRepo *mocks.ClaimAttachmentRepository
		mockHistRepo   *mocks.ClaimHistoryRepository
		mockCloudServ  *mocks.CloudinaryService
		mockTx         *mocks.Tx
		service        services.ClaimService
		ctx            context.Context
	)

	BeforeEach(func() {
		mockLogger = mocks.NewLogger(GinkgoT())
		mockClaimRepo = mocks.NewClaimRepository(GinkgoT())
		mockItemRepo = mocks.NewClaimItemRepository(GinkgoT())
		mockAttachRepo = mocks.NewClaimAttachmentRepository(GinkgoT())
		mockHistRepo = mocks.NewClaimHistoryRepository(GinkgoT())
		mockCloudServ = mocks.NewCloudinaryService(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		service = services.NewClaimService(mockLogger, mockClaimRepo, mockItemRepo, mockAttachRepo, mockHistRepo, mockCloudServ)
		ctx = context.Background()
	})

	Describe("GetByID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim is found", func() {
			It("should return the claim", func() {
				expectedClaim := &entities.Claim{
					ID:          claimID,
					VehicleID:   uuid.New(),
					CustomerID:  uuid.New(),
					Description: "Test claim",
					Status:      entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(expectedClaim, nil).Once()

				claim, err := service.GetByID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claim).NotTo(BeNil())
				Expect(claim.ID).To(Equal(expectedClaim.ID))
				Expect(claim.Description).To(Equal(expectedClaim.Description))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				claim, err := service.GetByID(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(claim).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimNotFound))
			})
		})
	})

	Describe("GetAll", func() {
		Context("when claims are found", func() {
			It("should return paginated claims", func() {
				filters := services.ClaimFilters{}
				pagination := services.Pagination{
					Page:     1,
					PageSize: 10,
					SortBy:   "created_at",
					SortDir:  "desc",
				}

				expectedClaims := []*entities.Claim{
					{ID: uuid.New(), Status: entities.ClaimStatusDraft},
					{ID: uuid.New(), Status: entities.ClaimStatusSubmitted},
				}

				mockClaimRepo.EXPECT().FindAll(ctx, mock.AnythingOfType("repositories.ClaimFilters"),
					mock.AnythingOfType("repositories.Pagination")).Return(expectedClaims, int64(2), nil).Once()

				result, err := service.GetAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Claims).To(HaveLen(2))
				Expect(result.Total).To(Equal(int64(2)))
				Expect(result.Page).To(Equal(1))
				Expect(result.PageSize).To(Equal(10))
				Expect(result.TotalPages).To(Equal(1))
			})
		})

		Context("when no claims are found", func() {
			It("should return empty result", func() {
				filters := services.ClaimFilters{}
				pagination := services.Pagination{Page: 1, PageSize: 10}

				mockClaimRepo.EXPECT().FindAll(ctx, mock.AnythingOfType("repositories.ClaimFilters"),
					mock.AnythingOfType("repositories.Pagination")).Return([]*entities.Claim{}, int64(0), nil).Once()

				result, err := service.GetAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Claims).To(BeEmpty())
				Expect(result.Total).To(Equal(int64(0)))
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				filters := services.ClaimFilters{}
				pagination := services.Pagination{Page: 1, PageSize: 10}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindAll(ctx, mock.AnythingOfType("repositories.ClaimFilters"),
					mock.AnythingOfType("repositories.Pagination")).Return(nil, int64(0), dbErr).Once()

				result, err := service.GetAll(ctx, filters, pagination)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when PageSize is 0", func() {
			It("should return totalPages as 0", func() {
				filters := services.ClaimFilters{}
				pagination := services.Pagination{Page: 1, PageSize: 0}

				mockClaimRepo.EXPECT().FindAll(ctx, mock.AnythingOfType("repositories.ClaimFilters"),
					mock.AnythingOfType("repositories.Pagination")).Return([]*entities.Claim{}, int64(5), nil).Once()

				result, err := service.GetAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.TotalPages).To(Equal(0))
			})
		})
	})

	Describe("Create", func() {
		var cmd *services.CreateClaimCommand

		BeforeEach(func() {
			cmd = &services.CreateClaimCommand{
				VehicleID:   uuid.New(),
				CustomerID:  uuid.New(),
				CreatorID:   uuid.New(),
				Description: "Test claim",
			}
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when claim is created successfully", func() {
			It("should create claim and history", func() {
				mockClaimRepo.EXPECT().Create(mockTx, mock.MatchedBy(func(c *entities.Claim) bool {
					return c.VehicleID == cmd.VehicleID &&
						c.CustomerID == cmd.CustomerID &&
						c.Description == cmd.Description &&
						c.Status == entities.ClaimStatusDraft
				})).Return(nil).Once()

				mockHistRepo.EXPECT().Create(mockTx, mock.MatchedBy(func(h *entities.ClaimHistory) bool {
					return h.Status == entities.ClaimStatusDraft &&
						h.ChangedBy == cmd.CreatorID
				})).Return(nil).Once()

				claim, err := service.Create(mockTx, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(claim).NotTo(BeNil())
				Expect(claim.VehicleID).To(Equal(cmd.VehicleID))
				Expect(claim.CustomerID).To(Equal(cmd.CustomerID))
				Expect(claim.Status).To(Equal(entities.ClaimStatusDraft))
			})
		})

		Context("when claim repository fails", func() {
			It("should return error", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockClaimRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.Claim")).Return(dbErr).Once()

				claim, err := service.Create(mockTx, cmd)

				Expect(err).To(HaveOccurred())
				Expect(claim).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when history creation fails", func() {
			It("should return error", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockClaimRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(dbErr).Once()

				claim, err := service.Create(mockTx, cmd)

				Expect(err).To(HaveOccurred())
				Expect(claim).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			claimID uuid.UUID
			cmd     *services.UpdateClaimCommand
		)

		BeforeEach(func() {
			claimID = uuid.New()
			cmd = &services.UpdateClaimCommand{
				Description: "Updated description",
			}
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when claim is updated successfully", func() {
			It("should update claim description", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.MatchedBy(func(c *entities.Claim) bool {
					return c.ID == claimID && c.Description == cmd.Description
				})).Return(nil).Once()

				err := service.Update(mockTx, claimID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not draft or request_info", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusSubmitted,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.Update(mockTx, claimID, cmd)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimStatusNotAllowedUpdate))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.Update(mockTx, claimID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when claim status is REQUEST_INFO", func() {
			It("should allow update", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusRequestInfo,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()

				err := service.Update(mockTx, claimID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when repository update fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(dbErr).Once()

				err := service.Update(mockTx, claimID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("HardDelete", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when claim is deleted successfully", func() {
			It("should delete claim and cloud files", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				attachments := []*entities.ClaimAttachment{
					{ID: uuid.New(), URL: "https://example.com/file1.jpg"},
					{ID: uuid.New(), URL: "https://example.com/file2.jpg"},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(attachments, nil).Once()
				mockClaimRepo.EXPECT().HardDelete(mockTx, claimID).Return(nil).Once()
				mockCloudServ.EXPECT().DeleteFileByURL(mock.Anything, mock.Anything).Return(nil).Maybe()
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

				err := service.HardDelete(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not draft", func() {
			It("should return NotAllowDeleteClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusSubmitted,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.HardDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimStatusNotAllowedDelete))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.HardDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when finding attachments fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				err := service.HardDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when hard delete fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entities.ClaimAttachment{}, nil).Once()
				mockClaimRepo.EXPECT().HardDelete(mockTx, claimID).Return(dbErr).Once()

				err := service.HardDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("SoftDelete", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when claim is soft deleted successfully", func() {
			It("should soft delete claim and related records", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusCancelled,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().SoftDelete(mockTx, claimID).Return(nil).Once()
				mockItemRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(nil).Once()
				mockAttachRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(nil).Once()
				mockHistRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(nil).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not cancelled", func() {
			It("should return NotAllowDeleteClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimStatusNotAllowedDelete))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when soft deleting claim fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusCancelled,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().SoftDelete(mockTx, claimID).Return(dbErr).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when soft deleting items fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusCancelled,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().SoftDelete(mockTx, claimID).Return(nil).Once()
				mockItemRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(dbErr).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when soft deleting attachments fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusCancelled,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().SoftDelete(mockTx, claimID).Return(nil).Once()
				mockItemRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(nil).Once()
				mockAttachRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(dbErr).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when soft deleting history fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusCancelled,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().SoftDelete(mockTx, claimID).Return(nil).Once()
				mockItemRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(nil).Once()
				mockAttachRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(nil).Once()
				mockHistRepo.EXPECT().SoftDeleteByClaimID(mockTx, claimID).Return(dbErr).Once()

				err := service.SoftDelete(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("UpdateStatus", func() {
		var (
			claimID   uuid.UUID
			changedBy uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			changedBy = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when status is updated successfully", func() {
			It("should update status and create history", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.MatchedBy(func(h *entities.ClaimHistory) bool {
					return h.ClaimID == claimID &&
						h.Status == entities.ClaimStatusSubmitted &&
						h.ChangedBy == changedBy
				})).Return(nil).Once()

				err := service.UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted, changedBy)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when status is invalid", func() {
			It("should return InvalidClaimStatus error", func() {
				err := service.UpdateStatus(mockTx, claimID, "INVALID_STATUS", changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeInvalidClaimStatus))
			})
		})

		Context("when status transition is invalid", func() {
			It("should return InvalidClaimAction error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusApproved,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.UpdateStatus(mockTx, claimID, entities.ClaimStatusDraft, changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeInvalidClaimAction))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when updating status fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted).Return(dbErr).Once()

				err := service.UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when creating history fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(dbErr).Once()

				err := service.UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Submit", func() {
		var (
			claimID   uuid.UUID
			changedBy uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			changedBy = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when claim is submitted successfully", func() {
			It("should update status to submitted", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), ClaimID: claimID},
				}
				attachments := []*entities.ClaimAttachment{
					{ID: uuid.New(), ClaimID: claimID},
					{ID: uuid.New(), ClaimID: claimID},
					{ID: uuid.New(), ClaimID: claimID},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(attachments, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(nil).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim has insufficient items", func() {
			It("should return MissingInformationClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				attachments := []*entities.ClaimAttachment{
					{ID: uuid.New(), ClaimID: claimID},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entities.ClaimItem{}, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(attachments, nil).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimMissingInformation))
			})
		})

		Context("when claim has insufficient attachments", func() {
			It("should return MissingInformationClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), ClaimID: claimID},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entities.ClaimAttachment{}, nil).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimMissingInformation))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when status transition is invalid", func() {
			It("should return InvalidClaimAction error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusApproved,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeInvalidClaimAction))
			})
		})

		Context("when finding items fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when finding attachments fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				items := []*entities.ClaimItem{{ID: uuid.New()}}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when updating status fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				items := []*entities.ClaimItem{{ID: uuid.New()}}
				attachments := []*entities.ClaimAttachment{
					{ID: uuid.New()}, {ID: uuid.New()}, {ID: uuid.New()},
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(attachments, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted).Return(dbErr).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when creating history fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				items := []*entities.ClaimItem{{ID: uuid.New()}}
				attachments := []*entities.ClaimAttachment{
					{ID: uuid.New()}, {ID: uuid.New()}, {ID: uuid.New()},
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(attachments, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusSubmitted).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(dbErr).Once()

				err := service.Submit(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Complete", func() {
		var (
			claimID   uuid.UUID
			changedBy uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			changedBy = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when all items are approved", func() {
			It("should set status to approved", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), Status: entities.ClaimItemStatusApproved},
					{ID: uuid.New(), Status: entities.ClaimItemStatusApproved},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusPartiallyApproved).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(nil).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when some items are rejected", func() {
			It("should set status to partially approved", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), Status: entities.ClaimItemStatusApproved},
					{ID: uuid.New(), Status: entities.ClaimItemStatusRejected},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusPartiallyApproved).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(nil).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when items have pending status", func() {
			It("should return InvalidClaimAction error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), Status: entities.ClaimItemStatusPending},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeInvalidClaimAction))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when finding items fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when status transition is invalid", func() {
			It("should return InvalidClaimAction error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), Status: entities.ClaimItemStatusApproved},
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeInvalidClaimAction))
			})
		})

		Context("when updating status fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), Status: entities.ClaimItemStatusApproved},
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusPartiallyApproved).Return(dbErr).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when creating history fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				items := []*entities.ClaimItem{
					{ID: uuid.New(), Status: entities.ClaimItemStatusApproved},
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(items, nil).Once()
				mockClaimRepo.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusPartiallyApproved).Return(nil).Once()
				mockHistRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimHistory")).Return(dbErr).Once()

				err := service.Complete(mockTx, claimID, changedBy)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("GetHistory", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when history is found", func() {
			It("should return claim histories", func() {
				expectedHistories := []*entities.ClaimHistory{
					{
						ID:        uuid.New(),
						ClaimID:   claimID,
						Status:    entities.ClaimStatusDraft,
						ChangedBy: uuid.New(),
					},
					{
						ID:        uuid.New(),
						ClaimID:   claimID,
						Status:    entities.ClaimStatusSubmitted,
						ChangedBy: uuid.New(),
					},
				}

				mockHistRepo.EXPECT().FindByClaimID(ctx, claimID).Return(expectedHistories, nil).Once()

				histories, err := service.GetHistory(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).NotTo(BeNil())
				Expect(histories).To(HaveLen(2))
			})
		})

		Context("when no history is found", func() {
			It("should return an empty slice", func() {
				mockHistRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entities.ClaimHistory{}, nil).Once()

				histories, err := service.GetHistory(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).NotTo(BeNil())
				Expect(histories).To(BeEmpty())
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockHistRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				histories, err := service.GetHistory(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(histories).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})
})
