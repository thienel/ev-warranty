package service_test

import (
	"context"
	"ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"ev-warranty-go/internal/application/service"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/mocks"
)

var _ = Describe("ClaimItemService", func() {
	var (
		mockClaimRepo *mocks.ClaimRepository
		mockItemRepo  *mocks.ClaimItemRepository
		mockUserRepo  *mocks.UserRepository
		mockClient    *mocks.Client
		mockTx        *mocks.Tx
		itemService   service.ClaimItemService
		ctx           context.Context
	)

	BeforeEach(func() {
		mockClaimRepo = mocks.NewClaimRepository(GinkgoT())
		mockItemRepo = mocks.NewClaimItemRepository(GinkgoT())
		mockUserRepo = mocks.NewUserRepository(GinkgoT())
		mockClient = mocks.NewClient(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		itemService = service.NewClaimItemService(mockClaimRepo, mockItemRepo, mockUserRepo, mockClient)
		ctx = context.Background()
	})

	Describe("GetByID", func() {
		var itemID uuid.UUID

		BeforeEach(func() {
			itemID = uuid.New()
		})

		Context("when item is found", func() {
			It("should return the item", func() {
				expectedItem := &entity.ClaimItem{
					ID:               itemID,
					ClaimID:          uuid.New(),
					IssueDescription: "Test issue",
					Status:           entity.ClaimItemStatusPending,
					Type:             entity.ClaimItemTypeRepair,
					Cost:             100.0,
				}

				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(expectedItem, nil).Once()

				item, err := itemService.GetByID(ctx, itemID)

				Expect(err).NotTo(HaveOccurred())
				Expect(item).NotTo(BeNil())
				Expect(item.ID).To(Equal(expectedItem.ID))
				Expect(item.IssueDescription).To(Equal(expectedItem.IssueDescription))
			})
		})

		Context("when item is not found", func() {
			It("should return ClaimItemNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(nil, notFoundErr).Once()

				item, err := itemService.GetByID(ctx, itemID)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				dbErr := apperror.ErrDBOperation
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(nil, dbErr).Once()

				item, err := itemService.GetByID(ctx, itemID)

				Expect(err).To(HaveOccurred())
				Expect(item).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("GetByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when items are found", func() {
			It("should return all items for the claim", func() {
				expectedItems := []*entity.ClaimItem{
					{
						ID:      uuid.New(),
						ClaimID: claimID,
						Status:  entity.ClaimItemStatusPending,
						Type:    entity.ClaimItemTypeRepair,
					},
					{
						ID:      uuid.New(),
						ClaimID: claimID,
						Status:  entity.ClaimItemStatusApproved,
						Type:    entity.ClaimItemTypeReplacement,
					},
				}

				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(expectedItems, nil).Once()

				items, err := itemService.GetByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).NotTo(BeNil())
				Expect(items).To(HaveLen(2))
				Expect(items[0].ClaimID).To(Equal(claimID))
				Expect(items[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no items are found", func() {
			It("should return an empty slice", func() {
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entity.ClaimItem{}, nil).Once()

				items, err := itemService.GetByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).NotTo(BeNil())
				Expect(items).To(BeEmpty())
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				dbErr := apperror.ErrDBOperation
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				items, err := itemService.GetByClaimID(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(items).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Create", func() {
		var (
			claimID uuid.UUID
			cmd     *service.CreateClaimItemCommand
		)

		BeforeEach(func() {
			claimID = uuid.New()
			cmd = &service.CreateClaimItemCommand{
				PartCategoryID:   uuid.New(),
				FaultyPartSerial: "Part serial",
				IssueDescription: "Part is damaged",
				Status:           entity.ClaimItemStatusPending,
				Type:             entity.ClaimItemTypeRepair,
			}
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is created successfully", func() {
			It("should return created item", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().Create(mockTx, mock.MatchedBy(func(i *entity.ClaimItem) bool {
					return i.ClaimID == claimID &&
						i.IssueDescription == cmd.IssueDescription &&
						i.Status == cmd.Status &&
						i.Type == cmd.Type
				})).Return(nil).Once()

				item, err := itemService.Create(mockTx, claimID, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(item).NotTo(BeNil())
				Expect(item.ClaimID).To(Equal(claimID))
				Expect(item.IssueDescription).To(Equal(cmd.IssueDescription))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				item, err := itemService.Create(mockTx, claimID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(item).To(BeNil())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when status is invalid", func() {
			BeforeEach(func() {
				cmd.Status = "INVALID_STATUS"
			})

			It("should return InvalidClaimItemStatus error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				item, err := itemService.Create(mockTx, claimID, cmd)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when type is invalid", func() {
			BeforeEach(func() {
				cmd.Type = "INVALID_TYPE"
			})

			It("should return InvalidClaimItemType error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				item, err := itemService.Create(mockTx, claimID, cmd)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when repository create fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entity.ClaimItem")).Return(dbErr).Once()

				item, err := itemService.Create(mockTx, claimID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(item).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			claimID uuid.UUID
			itemID  uuid.UUID
			cmd     *service.UpdateClaimItemCommand
		)

		BeforeEach(func() {
			claimID = uuid.New()
			itemID = uuid.New()
			cmd = &service.UpdateClaimItemCommand{
				IssueDescription: "Updated description",
				Type:             entity.ClaimItemTypeReplacement,
			}
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is updated successfully", func() {
			It("should return nil error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimItemStatusPending,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.MatchedBy(func(i *entity.ClaimItem) bool {
					return i.ID == itemID &&
						i.IssueDescription == cmd.IssueDescription &&
						i.Type == cmd.Type
				})).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(200.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(nil).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not draft or request_info", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusSubmitted,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when item status is approve", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimStatusApproved,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when item is not found", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				notFoundErr := apperror.ErrNotFoundError

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(nil, notFoundErr).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when item status is REJECTED", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimItemStatusRejected,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when type is invalid", func() {
			BeforeEach(func() {
				cmd.Type = "INVALID_TYPE"
			})

			It("should return InvalidClaimItemType error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimItemStatusPending,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when item repository update fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimItemStatusPending,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.ClaimItem")).Return(dbErr).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimItemStatusPending,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.ClaimItem")).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				item := &entity.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entity.ClaimItemStatusPending,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.ClaimItem")).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(200.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(dbErr).Once()

				err := itemService.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("HardDelete", func() {
		var (
			claimID uuid.UUID
			itemID  uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			itemID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is deleted successfully", func() {
			It("should return nil error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(100.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(nil).Once()

				err := itemService.HardDelete(mockTx, claimID, itemID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not draft", func() {
			It("should return NotAllowDeleteClaim error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusSubmitted,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := itemService.HardDelete(mockTx, claimID, itemID)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := itemService.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when item repository delete fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(dbErr).Once()

				err := itemService.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := itemService.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(100.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(dbErr).Once()

				err := itemService.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Approve", func() {
		var (
			claimID uuid.UUID
			itemID  uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			itemID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is approved successfully", func() {
			It("should return nil error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusApproved).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(150.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(nil).Once()

				err := itemService.Approve(mockTx, claimID, itemID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not reviewing", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := itemService.Approve(mockTx, claimID, itemID)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := itemService.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when UpdateStatus fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusApproved).Return(dbErr).Once()

				err := itemService.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusApproved).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := itemService.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusApproved).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(150.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(dbErr).Once()

				err := itemService.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Reject", func() {
		var (
			claimID uuid.UUID
			itemID  uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			itemID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is rejected successfully", func() {
			It("should return nil error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusRejected).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(50.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(nil).Once()

				err := itemService.Reject(mockTx, claimID, itemID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not reviewing", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := itemService.Reject(mockTx, claimID, itemID)

				ExpectAppError(err, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := itemService.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when UpdateStatus fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusRejected).Return(dbErr).Once()

				err := itemService.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusRejected).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := itemService.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entity.Claim{
					ID:     claimID,
					Status: entity.ClaimStatusReviewing,
				}
				dbErr := apperror.ErrDBOperation

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entity.ClaimStatusRejected).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(50.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entity.Claim")).Return(dbErr).Once()

				err := itemService.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})
})
