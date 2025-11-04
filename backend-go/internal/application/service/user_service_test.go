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

var _ = Describe("UserService", func() {
	var (
		mockUserRepo   *mocks.UserRepository
		mockOfficeRepo *mocks.OfficeRepository
		userService    service.UserService
		ctx            context.Context
	)

	BeforeEach(func() {
		mockUserRepo = mocks.NewUserRepository(GinkgoT())
		mockOfficeRepo = mocks.NewOfficeRepository(GinkgoT())
		userService = service.NewUserService(mockUserRepo, mockOfficeRepo)
		ctx = context.Background()
	})

	Describe("Create", func() {
		var cmd *service.UserCreateCommand

		Context("when user is created successfully", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entity.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return created user", func() {
				office := entity.NewOffice("Test Office", entity.OfficeTypeEVM, "123 Test St", true)
				office.ID = cmd.OfficeID
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Create(ctx, mock.MatchedBy(func(u *entity.User) bool {
					return u.Name == cmd.Name &&
						u.Email == cmd.Email &&
						u.Role == cmd.Role &&
						u.IsActive == cmd.IsActive &&
						u.OfficeID == cmd.OfficeID
				})).Return(nil).Once()

				user, err := userService.Create(ctx, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Name).To(Equal(cmd.Name))
				Expect(user.Email).To(Equal(cmd.Email))
				Expect(user.Role).To(Equal(cmd.Role))
			})
		})

		Context("when name is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "A",
					Email:    "john@example.com",
					Role:     entity.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := userService.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when email is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "invalid-email",
					Role:     entity.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := userService.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when password is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entity.UserRoleAdmin,
					Password: "weak",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := userService.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when role is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     "INVALID_ROLE",
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := userService.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when office is not found", func() {
			It("should return OfficeNotFound error", func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entity.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}

				officeErr := apperror.ErrNotFoundError
				mockOfficeRepo.EXPECT().FindByID(mock.Anything, cmd.OfficeID).Return(nil, officeErr).Once()

				user, err := userService.Create(ctx, cmd)

				Expect(err).To(HaveOccurred())
				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when role does not match office type", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entity.UserRoleScStaff,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidOfficeType error", func() {
				office := &entity.Office{
					ID:         cmd.OfficeID,
					OfficeType: entity.OfficeTypeEVM,
				}
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()

				user, err := userService.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when repository create fails", func() {
			BeforeEach(func() {
				cmd = &service.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entity.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return DBOperationError", func() {
				office := &entity.Office{
					ID:         cmd.OfficeID,
					OfficeType: entity.OfficeTypeEVM,
				}
				dbErr := apperror.ErrDBOperation
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Create(ctx, mock.AnythingOfType("*entity.User")).Return(dbErr).Once()

				user, err := userService.Create(ctx, cmd)

				Expect(err).To(HaveOccurred())
				Expect(user).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("FindByID", func() {
		var userID uuid.UUID

		BeforeEach(func() {
			userID = uuid.New()
		})

		Context("when user is found", func() {
			It("should return the user", func() {
				expectedUser := &entity.User{
					ID:       userID,
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entity.UserRoleAdmin,
					IsActive: true,
				}

				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(expectedUser, nil).Once()

				user, err := userService.GetByID(ctx, userID)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.ID).To(Equal(expectedUser.ID))
				Expect(user.Name).To(Equal(expectedUser.Name))
				Expect(user.Email).To(Equal(expectedUser.Email))
			})
		})

		Context("when user is not found", func() {
			It("should return UserNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, notFoundErr).Once()

				user, err := userService.GetByID(ctx, userID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})
	})

	Describe("GetAll", func() {
		Context("when users are found", func() {
			It("should return all users", func() {
				expectedUsers := []*entity.User{
					{
						ID:    uuid.New(),
						Name:  "User 1",
						Email: "user1@example.com",
						Role:  entity.UserRoleAdmin,
					},
					{
						ID:    uuid.New(),
						Name:  "User 2",
						Email: "user2@example.com",
						Role:  entity.UserRoleEvmStaff,
					},
				}

				mockUserRepo.EXPECT().FindAll(ctx).Return(expectedUsers, nil).Once()

				users, err := userService.GetAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).NotTo(BeNil())
				Expect(users).To(HaveLen(2))
				Expect(users[0].ID).To(Equal(expectedUsers[0].ID))
				Expect(users[1].ID).To(Equal(expectedUsers[1].ID))
			})
		})

		Context("when no users are found", func() {
			It("should return an empty slice", func() {
				mockUserRepo.EXPECT().FindAll(ctx).Return([]*entity.User{}, nil).Once()

				users, err := userService.GetAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).NotTo(BeNil())
				Expect(users).To(BeEmpty())
			})
		})

		Context("when there is a repository error", func() {
			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockUserRepo.EXPECT().FindAll(ctx).Return(nil, dbErr).Once()

				users, err := userService.GetAll(ctx)

				Expect(err).To(HaveOccurred())
				Expect(users).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			userID       uuid.UUID
			existingUser *entity.User
			cmd          *service.UserUpdateCommand
		)

		BeforeEach(func() {
			userID = uuid.New()
			existingUser = &entity.User{
				ID:       userID,
				Name:     "Old Name",
				Role:     entity.UserRoleAdmin,
				IsActive: true,
				OfficeID: uuid.New(),
			}
		})

		Context("when user is updated successfully", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entity.UserRoleEvmStaff,
					IsActive: false,
					OfficeID: uuid.New(),
				}
			})

			It("should return nil error", func() {
				office := &entity.Office{
					ID:         cmd.OfficeID,
					OfficeType: entity.OfficeTypeEVM,
				}
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, mock.MatchedBy(func(u *entity.User) bool {
					return u.ID == userID &&
						u.Name == cmd.Name &&
						u.Role == cmd.Role &&
						u.IsActive == cmd.IsActive &&
						u.OfficeID == cmd.OfficeID
				})).Return(nil).Once()

				err := userService.Update(ctx, userID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when user is not found", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entity.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return UserNotFound error", func() {
				notFoundErr := apperror.ErrNotFoundError
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, notFoundErr).Once()

				err := userService.Update(ctx, userID, cmd)

				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when name is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "A",
					Role:     entity.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()

				err := userService.Update(ctx, userID, cmd)

				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when role is invalid", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     "INVALID_ROLE",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()

				err := userService.Update(ctx, userID, cmd)

				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when office is not found", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entity.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return OfficeNotFound error", func() {
				officeErr := apperror.ErrNotFoundError
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(nil, officeErr).Once()

				err := userService.Update(ctx, userID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(officeErr))
			})
		})

		Context("when role does not match office type", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entity.UserRoleScStaff,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidOfficeType error", func() {
				office := &entity.Office{
					ID:         cmd.OfficeID,
					OfficeType: entity.OfficeTypeEVM,
				}
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()

				err := userService.Update(ctx, userID, cmd)

				ExpectAppError(err, apperror.ErrInvalidInput.ErrorCode)
			})
		})

		Context("when repository update fails", func() {
			BeforeEach(func() {
				cmd = &service.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entity.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return DBOperationError", func() {
				office := &entity.Office{
					ID:         cmd.OfficeID,
					OfficeType: entity.OfficeTypeEVM,
				}
				dbErr := apperror.ErrDBOperation
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, mock.AnythingOfType("*entity.User")).Return(dbErr).Once()

				err := userService.Update(ctx, userID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Delete", func() {
		var userID uuid.UUID

		BeforeEach(func() {
			userID = uuid.New()
		})

		Context("when user is deleted successfully", func() {
			It("should return nil error", func() {
				mockUserRepo.EXPECT().SoftDelete(ctx, userID).Return(nil).Once()

				err := userService.Delete(ctx, userID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when repository delete fails", func() {
			It("should return DBOperationError", func() {
				dbErr := apperror.ErrDBOperation
				mockUserRepo.EXPECT().SoftDelete(ctx, userID).Return(dbErr).Once()

				err := userService.Delete(ctx, userID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})
})
