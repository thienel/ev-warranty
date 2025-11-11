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

var _ = Describe("OfficeService", func() {
	var (
		mockRepo      *mocks.OfficeRepository
		officeService service.OfficeService
		ctx           context.Context
	)

	BeforeEach(func() {
		mockRepo = mocks.NewOfficeRepository(GinkgoT())
		officeService = service.NewOfficeService(mockRepo)
		ctx = context.Background()
	})

	Describe("Create", func() {
		var cmd *service.CreateOfficeCommand

		Context("when office is created successfully with valid office type EVM", func() {
			BeforeEach(func() {
				cmd = &service.CreateOfficeCommand{
					OfficeName: "Test Office",
					OfficeType: entity.OfficeTypeEVM,
					Address:    "123 Test St",
					IsActive:   true,
				}
			})

			It("should return created office", func() {
				mockRepo.EXPECT().Create(ctx, MatchOffice(cmd)).Return(nil).Once()

				office, err := officeService.Create(ctx, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(office).NotTo(BeNil())
				Expect(office.OfficeName).To(Equal(cmd.OfficeName))
				Expect(office.OfficeType).To(Equal(cmd.OfficeType))
				Expect(office.Address).To(Equal(cmd.Address))
				Expect(office.IsActive).To(Equal(cmd.IsActive))
				Expect(office.ID).NotTo(Equal(uuid.Nil))
			})
		})

		Context("when office is created successfully with valid office type SC", func() {
			BeforeEach(func() {
				cmd = &service.CreateOfficeCommand{
					OfficeName: "Service Center",
					OfficeType: entity.OfficeTypeSC,
					Address:    "456 Main St",
					IsActive:   false,
				}
			})

			It("should return created office", func() {
				mockRepo.EXPECT().Create(ctx, MatchOffice(cmd)).Return(nil).Once()

				office, err := officeService.Create(ctx, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(office).NotTo(BeNil())
				Expect(office.OfficeName).To(Equal(cmd.OfficeName))
				Expect(office.OfficeType).To(Equal(cmd.OfficeType))
				Expect(office.Address).To(Equal(cmd.Address))
				Expect(office.IsActive).To(Equal(cmd.IsActive))
			})
		})

		Context("when office type is invalid", func() {
			BeforeEach(func() {
				cmd = &service.CreateOfficeCommand{
					OfficeName: "Test Office",
					OfficeType: "INVALID_TYPE",
					Address:    "123 Test St",
					IsActive:   true,
				}
			})

			It("should return InvalidOfficeType error", func() {
				office, err := officeService.Create(ctx, cmd)

				Expect(office).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when repository create fails", func() {
			BeforeEach(func() {
				cmd = &service.CreateOfficeCommand{
					OfficeName: "Test Office",
					OfficeType: entity.OfficeTypeEVM,
					Address:    "123 Test St",
					IsActive:   true,
				}
			})

			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockRepo.EXPECT().Create(ctx, MatchOffice(cmd)).Return(dbErr).Once()

				office, err := officeService.Create(ctx, cmd)

				Expect(err).To(HaveOccurred())
				Expect(office).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("GetByID", func() {
		var officeID uuid.UUID

		BeforeEach(func() {
			officeID = uuid.New()
		})

		Context("when office is found", func() {
			It("should return the office", func() {
				expectedOffice := &entity.Office{
					ID:         officeID,
					OfficeName: "Test Office",
					OfficeType: entity.OfficeTypeEVM,
					Address:    "123 Test St",
					IsActive:   true,
				}

				mockRepo.EXPECT().FindByID(ctx, officeID).Return(expectedOffice, nil).Once()

				office, err := officeService.GetByID(ctx, officeID)

				Expect(err).NotTo(HaveOccurred())
				Expect(office).NotTo(BeNil())
				Expect(office.ID).To(Equal(expectedOffice.ID))
				Expect(office.OfficeName).To(Equal(expectedOffice.OfficeName))
				Expect(office.OfficeType).To(Equal(expectedOffice.OfficeType))
				Expect(office.Address).To(Equal(expectedOffice.Address))
				Expect(office.IsActive).To(Equal(expectedOffice.IsActive))
			})
		})

		Context("when office is not found", func() {
			It("should return OfficeNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockRepo.EXPECT().FindByID(ctx, officeID).Return(nil, notFoundErr).Once()

				office, err := officeService.GetByID(ctx, officeID)

				Expect(office).To(BeNil())
				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when there is a repository error", func() {
			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockRepo.EXPECT().FindByID(ctx, officeID).Return(nil, dbErr).Once()

				office, err := officeService.GetByID(ctx, officeID)

				Expect(err).To(HaveOccurred())
				Expect(office).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("GetAll", func() {
		Context("when offices are found", func() {
			It("should return all offices", func() {
				expectedOffices := []*entity.Office{
					{
						ID:         uuid.New(),
						OfficeName: "Office 1",
						OfficeType: entity.OfficeTypeEVM,
						Address:    "123 Test St",
						IsActive:   true,
					},
					{
						ID:         uuid.New(),
						OfficeName: "Office 2",
						OfficeType: entity.OfficeTypeSC,
						Address:    "456 Main St",
						IsActive:   false,
					},
				}

				mockRepo.EXPECT().FindAll(ctx).Return(expectedOffices, nil).Once()

				offices, err := officeService.GetAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(offices).NotTo(BeNil())
				Expect(offices).To(HaveLen(2))
				Expect(offices[0].ID).To(Equal(expectedOffices[0].ID))
				Expect(offices[0].OfficeName).To(Equal(expectedOffices[0].OfficeName))
				Expect(offices[1].ID).To(Equal(expectedOffices[1].ID))
				Expect(offices[1].OfficeName).To(Equal(expectedOffices[1].OfficeName))
			})
		})

		Context("when no offices are found", func() {
			It("should return an empty slice", func() {
				mockRepo.EXPECT().FindAll(ctx).Return([]*entity.Office{}, nil).Once()

				offices, err := officeService.GetAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(offices).NotTo(BeNil())
				Expect(offices).To(BeEmpty())
			})
		})

		Context("when there is a repository error", func() {
			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockRepo.EXPECT().FindAll(ctx).Return(nil, dbErr).Once()

				offices, err := officeService.GetAll(ctx)

				Expect(err).To(HaveOccurred())
				Expect(offices).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			officeID       uuid.UUID
			existingOffice *entity.Office
			cmd            *service.UpdateOfficeCommand
		)

		BeforeEach(func() {
			officeID = uuid.New()
			existingOffice = &entity.Office{
				ID:         officeID,
				OfficeName: "Old Office Name",
				OfficeType: entity.OfficeTypeEVM,
				Address:    "Old Address",
				IsActive:   true,
			}
		})

		Context("when office is updated successfully", func() {
			BeforeEach(func() {
				cmd = &service.UpdateOfficeCommand{
					OfficeName: "Updated Office Name",
					OfficeType: entity.OfficeTypeSC,
					Address:    "Updated Address",
					IsActive:   false,
				}
			})

			It("should return nil error", func() {
				mockRepo.EXPECT().FindByID(ctx, officeID).Return(existingOffice, nil).Once()
				mockRepo.EXPECT().Update(ctx, MatchUpdatedOffice(officeID, cmd)).Return(nil).Once()

				err := officeService.Update(ctx, officeID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when office is not found", func() {
			BeforeEach(func() {
				cmd = &service.UpdateOfficeCommand{
					OfficeName: "Updated Office Name",
					OfficeType: entity.OfficeTypeEVM,
					Address:    "Updated Address",
					IsActive:   true,
				}
			})

			It("should return OfficeNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockRepo.EXPECT().FindByID(ctx, officeID).Return(nil, notFoundErr).Once()

				err := officeService.Update(ctx, officeID, cmd)

				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when office type is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UpdateOfficeCommand{
					OfficeName: "Updated Office Name",
					OfficeType: "INVALID_TYPE",
					Address:    "Updated Address",
					IsActive:   false,
				}
			})

			It("should return InvalidOfficeType error", func() {
				mockRepo.EXPECT().FindByID(ctx, officeID).Return(existingOffice, nil).Once()

				err := officeService.Update(ctx, officeID, cmd)

				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when repository update fails", func() {
			BeforeEach(func() {
				cmd = &service.UpdateOfficeCommand{
					OfficeName: "Updated Office Name",
					OfficeType: entity.OfficeTypeEVM,
					Address:    "Updated Address",
					IsActive:   false,
				}
			})

			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockRepo.EXPECT().FindByID(ctx, officeID).Return(existingOffice, nil).Once()
				mockRepo.EXPECT().Update(ctx, MatchUpdatedOffice(officeID, cmd)).Return(dbErr).Once()

				err := officeService.Update(ctx, officeID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("DeleteByID", func() {
		var officeID uuid.UUID

		BeforeEach(func() {
			officeID = uuid.New()
		})

		Context("when office is deleted successfully", func() {
			It("should return nil error", func() {
				mockRepo.EXPECT().SoftDelete(ctx, officeID).Return(nil).Once()

				err := officeService.DeleteByID(ctx, officeID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when repository delete fails", func() {
			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockRepo.EXPECT().SoftDelete(ctx, officeID).Return(dbErr).Once()

				err := officeService.DeleteByID(ctx, officeID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when office is not found", func() {
			It("should return OfficeNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockRepo.EXPECT().SoftDelete(ctx, officeID).Return(notFoundErr).Once()

				err := officeService.DeleteByID(ctx, officeID)

				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})
	})
})

func MatchOffice(cmd *service.CreateOfficeCommand) interface{} {
	return mock.MatchedBy(func(o *entity.Office) bool {
		return o.OfficeName == cmd.OfficeName &&
			o.OfficeType == cmd.OfficeType &&
			o.Address == cmd.Address &&
			o.IsActive == cmd.IsActive
	})
}

func MatchUpdatedOffice(id uuid.UUID, cmd *service.UpdateOfficeCommand) interface{} {
	return mock.MatchedBy(func(o *entity.Office) bool {
		return o.ID == id &&
			o.OfficeName == cmd.OfficeName &&
			o.OfficeType == cmd.OfficeType &&
			o.Address == cmd.Address &&
			o.IsActive == cmd.IsActive
	})
}
