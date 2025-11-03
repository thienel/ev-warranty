package services_test

import (
	"context"
	"errors"
	apperrors2 "ev-warranty-go/pkg/apperror"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/mocks"
)

var _ = Describe("ClaimItemService", func() {
	var (
		mockClaimRepo *mocks.ClaimRepository
		mockItemRepo  *mocks.ClaimItemRepository
		mockTx        *mocks.Tx
		service       services.ClaimItemService
		ctx           context.Context
	)

	BeforeEach(func() {
		mockClaimRepo = mocks.NewClaimRepository(GinkgoT())
		mockItemRepo = mocks.NewClaimItemRepository(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		service = services.NewClaimItemService(mockClaimRepo, mockItemRepo)
		ctx = context.Background()
	})

	Describe("GetByID", func() {
		var itemID uuid.UUID

		BeforeEach(func() {
			itemID = uuid.New()
		})

		Context("when item is found", func() {
			It("should return the item", func() {
				expectedItem := &entities.ClaimItem{
					ID:               itemID,
					ClaimID:          uuid.New(),
					IssueDescription: "Test issue",
					Status:           entities.ClaimItemStatusPending,
					Type:             entities.ClaimItemTypeRepair,
					Cost:             100.0,
				}

				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(expectedItem, nil).Once()

				item, err := service.GetByID(ctx, itemID)

				Expect(err).NotTo(HaveOccurred())
				Expect(item).NotTo(BeNil())
				Expect(item.ID).To(Equal(expectedItem.ID))
				Expect(item.IssueDescription).To(Equal(expectedItem.IssueDescription))
			})
		})

		Context("when item is not found", func() {
			It("should return ClaimItemNotFound error", func() {
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimItemNotFound, errors.New("item not found"))
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(nil, notFoundErr).Once()

				item, err := service.GetByID(ctx, itemID)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperrors2.ErrorCodeClaimItemNotFound)
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(nil, dbErr).Once()

				item, err := service.GetByID(ctx, itemID)

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
				expectedItems := []*entities.ClaimItem{
					{
						ID:      uuid.New(),
						ClaimID: claimID,
						Status:  entities.ClaimItemStatusPending,
						Type:    entities.ClaimItemTypeRepair,
					},
					{
						ID:      uuid.New(),
						ClaimID: claimID,
						Status:  entities.ClaimItemStatusApproved,
						Type:    entities.ClaimItemTypeReplacement,
					},
				}

				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(expectedItems, nil).Once()

				items, err := service.GetByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).NotTo(BeNil())
				Expect(items).To(HaveLen(2))
				Expect(items[0].ClaimID).To(Equal(claimID))
				Expect(items[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no items are found", func() {
			It("should return an empty slice", func() {
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entities.ClaimItem{}, nil).Once()

				items, err := service.GetByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).NotTo(BeNil())
				Expect(items).To(BeEmpty())
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))
				mockItemRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				items, err := service.GetByClaimID(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(items).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Create", func() {
		var (
			claimID uuid.UUID
			cmd     *services.CreateClaimItemCommand
		)

		BeforeEach(func() {
			claimID = uuid.New()
			cmd = &services.CreateClaimItemCommand{
				PartCategoryID:   uuid.New(),
				FaultyPartID:     uuid.New(),
				IssueDescription: "Part is damaged",
				Status:           entities.ClaimItemStatusPending,
				Type:             entities.ClaimItemTypeRepair,
				Cost:             150.0,
			}
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is created successfully", func() {
			It("should return created item", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().Create(mockTx, mock.MatchedBy(func(i *entities.ClaimItem) bool {
					return i.ClaimID == claimID &&
						i.IssueDescription == cmd.IssueDescription &&
						i.Status == cmd.Status &&
						i.Type == cmd.Type &&
						i.Cost == cmd.Cost
				})).Return(nil).Once()

				item, err := service.Create(mockTx, claimID, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(item).NotTo(BeNil())
				Expect(item.ClaimID).To(Equal(claimID))
				Expect(item.IssueDescription).To(Equal(cmd.IssueDescription))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				item, err := service.Create(mockTx, claimID, cmd)

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
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				item, err := service.Create(mockTx, claimID, cmd)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperrors2.ErrorCodeInvalidClaimItemStatus)
			})
		})

		Context("when type is invalid", func() {
			BeforeEach(func() {
				cmd.Type = "INVALID_TYPE"
			})

			It("should return InvalidClaimItemType error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				item, err := service.Create(mockTx, claimID, cmd)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperrors2.ErrorCodeInvalidClaimItemType)
			})
		})

		Context("when repository create fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimItem")).Return(dbErr).Once()

				item, err := service.Create(mockTx, claimID, cmd)

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
			cmd     *services.UpdateClaimItemCommand
		)

		BeforeEach(func() {
			claimID = uuid.New()
			itemID = uuid.New()
			cmd = &services.UpdateClaimItemCommand{
				IssueDescription: "Updated description",
				Type:             entities.ClaimItemTypeReplacement,
				Cost:             200.0,
			}
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when item is updated successfully", func() {
			It("should return nil error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusPending,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.MatchedBy(func(i *entities.ClaimItem) bool {
					return i.ID == itemID &&
						i.IssueDescription == cmd.IssueDescription &&
						i.Type == cmd.Type &&
						i.Cost == cmd.Cost
				})).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(200.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

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

				err := service.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperrors2.ErrorCodeClaimStatusNotAllowedUpdate)
			})
		})

		Context("when item status is approve", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimStatusApproved,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperrors2.ErrorCodeClaimStatusNotAllowedUpdate)
			})
		})

		Context("when claim status is REQUEST_INFO", func() {
			It("should allow update", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusRequestInfo,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusPending,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.ClaimItem")).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(200.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when item is not found", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimItemNotFound, errors.New("item not found"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(nil, notFoundErr).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when item status is REJECTED", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusRejected,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperrors2.ErrorCodeClaimStatusNotAllowedUpdate)
			})
		})

		Context("when type is invalid", func() {
			BeforeEach(func() {
				cmd.Type = "INVALID_TYPE"
			})

			It("should return InvalidClaimItemType error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusPending,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				ExpectAppError(err, apperrors2.ErrorCodeInvalidClaimItemType)
			})
		})

		Context("when item repository update fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusPending,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.ClaimItem")).Return(dbErr).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusPending,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("sum error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.ClaimItem")).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				item := &entities.ClaimItem{
					ID:      itemID,
					ClaimID: claimID,
					Status:  entities.ClaimItemStatusPending,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().FindByID(ctx, itemID).Return(item, nil).Once()
				mockItemRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.ClaimItem")).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(200.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(dbErr).Once()

				err := service.Update(mockTx, claimID, itemID, cmd)

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
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(100.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()

				err := service.HardDelete(mockTx, claimID, itemID)

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

				err := service.HardDelete(mockTx, claimID, itemID)

				ExpectAppError(err, apperrors2.ErrorCodeClaimStatusNotAllowedDelete)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when item repository delete fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(dbErr).Once()

				err := service.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("sum error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := service.HardDelete(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().HardDelete(mockTx, itemID).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(100.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(dbErr).Once()

				err := service.HardDelete(mockTx, claimID, itemID)

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
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusApproved).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(150.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()

				err := service.Approve(mockTx, claimID, itemID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not reviewing", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.Approve(mockTx, claimID, itemID)

				ExpectAppError(err, apperrors2.ErrorCodeClaimStatusNotAllowedUpdate)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when UpdateStatus fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusApproved).Return(dbErr).Once()

				err := service.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("sum error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusApproved).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := service.Approve(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusApproved).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(150.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(dbErr).Once()

				err := service.Approve(mockTx, claimID, itemID)

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
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusRejected).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(50.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(nil).Once()

				err := service.Reject(mockTx, claimID, itemID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not reviewing", func() {
			It("should return NotAllowUpdateClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.Reject(mockTx, claimID, itemID)

				ExpectAppError(err, apperrors2.ErrorCodeClaimStatusNotAllowedUpdate)
			})
		})

		Context("when claim is not found", func() {
			It("should return error", func() {
				notFoundErr := apperrors2.New(404, apperrors2.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when UpdateStatus fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusRejected).Return(dbErr).Once()

				err := service.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when SumCostByClaimID fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("sum error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusRejected).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(0.0, dbErr).Once()

				err := service.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when claim repository update fails", func() {
			It("should return error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusReviewing,
				}
				dbErr := apperrors2.New(500, apperrors2.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockItemRepo.EXPECT().UpdateStatus(mockTx, itemID, entities.ClaimStatusRejected).Return(nil).Once()
				mockItemRepo.EXPECT().SumCostByClaimID(mockTx, claimID).Return(50.0, nil).Once()
				mockClaimRepo.EXPECT().Update(mockTx, mock.AnythingOfType("*entities.Claim")).Return(dbErr).Once()

				err := service.Reject(mockTx, claimID, itemID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})
})
