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

var _ = Describe("UserService", func() {
	var (
		mockUserRepo   *mocks.UserRepository
		mockOfficeRepo *mocks.OfficeRepository
		service        services.UserService
		ctx            context.Context
	)

	BeforeEach(func() {
		mockUserRepo = mocks.NewUserRepository(GinkgoT())
		mockOfficeRepo = mocks.NewOfficeRepository(GinkgoT())
		service = services.NewUserService(mockUserRepo, mockOfficeRepo)
		ctx = context.Background()
	})

	Describe("Create", func() {
		var cmd *services.UserCreateCommand

		Context("when user is created successfully", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entities.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return created user", func() {
				office := &entities.Office{
					ID:         cmd.OfficeID,
					OfficeType: entities.OfficeTypeEVM,
				}
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Create(ctx, mock.MatchedBy(func(u *entities.User) bool {
					return u.Name == cmd.Name &&
						u.Email == cmd.Email &&
						u.Role == cmd.Role &&
						u.IsActive == cmd.IsActive &&
						u.OfficeID == cmd.OfficeID
				})).Return(nil).Once()

				user, err := service.Create(ctx, cmd)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Name).To(Equal(cmd.Name))
				Expect(user.Email).To(Equal(cmd.Email))
				Expect(user.Role).To(Equal(cmd.Role))
			})
		})

		Context("when name is invalid", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "A",
					Email:    "john@example.com",
					Role:     entities.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := service.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when email is invalid", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "invalid-email",
					Role:     entities.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := service.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when password is invalid", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entities.UserRoleAdmin,
					Password: "weak",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := service.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when role is invalid", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     "INVALID_ROLE",
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				user, err := service.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when office is not found", func() {
			It("should return OfficeNotFound error", func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entities.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}

				officeErr := apperrors.NewOfficeNotFound()
				mockOfficeRepo.EXPECT().FindByID(mock.Anything, cmd.OfficeID).Return(nil, officeErr).Once()

				user, err := service.Create(ctx, cmd)

				Expect(err).To(HaveOccurred())
				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeOfficeNotFound)
			})
		})

		Context("when role does not match office type", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entities.UserRoleScStaff,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidOfficeType error", func() {
				office := &entities.Office{
					ID:         cmd.OfficeID,
					OfficeType: entities.OfficeTypeEVM,
				}
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()

				user, err := service.Create(ctx, cmd)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeInvalidOfficeType)
			})
		})

		Context("when repository create fails", func() {
			BeforeEach(func() {
				cmd = &services.UserCreateCommand{
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entities.UserRoleAdmin,
					Password: "Password123!",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return DBOperationError", func() {
				office := &entities.Office{
					ID:         cmd.OfficeID,
					OfficeType: entities.OfficeTypeEVM,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Create(ctx, mock.AnythingOfType("*entities.User")).Return(dbErr).Once()

				user, err := service.Create(ctx, cmd)

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
				expectedUser := &entities.User{
					ID:       userID,
					Name:     "John Doe",
					Email:    "john@example.com",
					Role:     entities.UserRoleAdmin,
					IsActive: true,
				}

				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(expectedUser, nil).Once()

				user, err := service.GetByID(ctx, userID)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.ID).To(Equal(expectedUser.ID))
				Expect(user.Name).To(Equal(expectedUser.Name))
				Expect(user.Email).To(Equal(expectedUser.Email))
			})
		})

		Context("when user is not found", func() {
			It("should return UserNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeUserNotFound, errors.New("user not found"))
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, notFoundErr).Once()

				user, err := service.GetByID(ctx, userID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeUserNotFound)
			})
		})
	})

	Describe("GetAll", func() {
		Context("when users are found", func() {
			It("should return all users", func() {
				expectedUsers := []*entities.User{
					{
						ID:    uuid.New(),
						Name:  "User 1",
						Email: "user1@example.com",
						Role:  entities.UserRoleAdmin,
					},
					{
						ID:    uuid.New(),
						Name:  "User 2",
						Email: "user2@example.com",
						Role:  entities.UserRoleEvmStaff,
					},
				}

				mockUserRepo.EXPECT().FindAll(ctx).Return(expectedUsers, nil).Once()

				users, err := service.GetAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).NotTo(BeNil())
				Expect(users).To(HaveLen(2))
				Expect(users[0].ID).To(Equal(expectedUsers[0].ID))
				Expect(users[1].ID).To(Equal(expectedUsers[1].ID))
			})
		})

		Context("when no users are found", func() {
			It("should return an empty slice", func() {
				mockUserRepo.EXPECT().FindAll(ctx).Return([]*entities.User{}, nil).Once()

				users, err := service.GetAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).NotTo(BeNil())
				Expect(users).To(BeEmpty())
			})
		})

		Context("when there is a repository error", func() {
			It("should return DBOperationError", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockUserRepo.EXPECT().FindAll(ctx).Return(nil, dbErr).Once()

				users, err := service.GetAll(ctx)

				Expect(err).To(HaveOccurred())
				Expect(users).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Update", func() {
		var (
			userID       uuid.UUID
			existingUser *entities.User
			cmd          *services.UserUpdateCommand
		)

		BeforeEach(func() {
			userID = uuid.New()
			existingUser = &entities.User{
				ID:       userID,
				Name:     "Old Name",
				Role:     entities.UserRoleAdmin,
				IsActive: true,
				OfficeID: uuid.New(),
			}
		})

		Context("when user is updated successfully", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entities.UserRoleEvmStaff,
					IsActive: false,
					OfficeID: uuid.New(),
				}
			})

			It("should return nil error", func() {
				office := &entities.Office{
					ID:         cmd.OfficeID,
					OfficeType: entities.OfficeTypeEVM,
				}
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, mock.MatchedBy(func(u *entities.User) bool {
					return u.ID == userID &&
						u.Name == cmd.Name &&
						u.Role == cmd.Role &&
						u.IsActive == cmd.IsActive &&
						u.OfficeID == cmd.OfficeID
				})).Return(nil).Once()

				err := service.Update(ctx, userID, cmd)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when user is not found", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entities.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return UserNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeUserNotFound, errors.New("user not found"))
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, notFoundErr).Once()

				err := service.Update(ctx, userID, cmd)

				ExpectAppError(err, apperrors.ErrorCodeUserNotFound)
			})
		})

		Context("when name is invalid", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "A",
					Role:     entities.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()

				err := service.Update(ctx, userID, cmd)

				ExpectAppError(err, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when role is invalid", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     "INVALID_ROLE",
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidUserInput error", func() {
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()

				err := service.Update(ctx, userID, cmd)

				ExpectAppError(err, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when office is not found", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entities.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return OfficeNotFound error", func() {
				officeErr := apperrors.New(404, apperrors.ErrorCodeOfficeNotFound, errors.New("office not found"))
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(nil, officeErr).Once()

				err := service.Update(ctx, userID, cmd)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(officeErr))
			})
		})

		Context("when role does not match office type", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entities.UserRoleScStaff,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return InvalidOfficeType error", func() {
				office := &entities.Office{
					ID:         cmd.OfficeID,
					OfficeType: entities.OfficeTypeEVM,
				}
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()

				err := service.Update(ctx, userID, cmd)

				ExpectAppError(err, apperrors.ErrorCodeInvalidOfficeType)
			})
		})

		Context("when repository update fails", func() {
			BeforeEach(func() {
				cmd = &services.UserUpdateCommand{
					Name:     "Updated Name",
					Role:     entities.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
			})

			It("should return DBOperationError", func() {
				office := &entities.Office{
					ID:         cmd.OfficeID,
					OfficeType: entities.OfficeTypeEVM,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existingUser, nil).Once()
				mockOfficeRepo.EXPECT().FindByID(ctx, cmd.OfficeID).Return(office, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, mock.AnythingOfType("*entities.User")).Return(dbErr).Once()

				err := service.Update(ctx, userID, cmd)

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

				err := service.Delete(ctx, userID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when repository delete fails", func() {
			It("should return DBOperationError", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockUserRepo.EXPECT().SoftDelete(ctx, userID).Return(dbErr).Once()

				err := service.Delete(ctx, userID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})
	})
})
