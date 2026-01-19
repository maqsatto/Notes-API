package service

import (
	"context"

	"github.com/maqsatto/Notes-API/internal/auth"
	"github.com/maqsatto/Notes-API/internal/domain"
	"github.com/maqsatto/Notes-API/internal/repository"
	"github.com/maqsatto/Notes-API/internal/validator"
)

type UserService struct {
	UserRepository repository.UserRepo
	NoteRepository repository.NoteRepo
	JWT *auth.JWTManager
}

func NewUserService(UserRepository repository.UserRepo,NoteRepository repository.NoteRepo ,JWT *auth.JWTManager) *UserService {
	return &UserService{
		UserRepository: UserRepository,
		NoteRepository: NoteRepository,
		JWT: JWT,
	}
}

func (u *UserService) 	Register(ctx context.Context, username, email, password string) (*domain.User, string,  error) {
	if err := validator.ValidateUserRegister(email, username, password); err != nil {
		return nil, "", err
	}
	if exists, _ := u.UserRepository.ExistsByEmail(ctx, email); exists {
		return nil, "", domain.ErrUserAlreadyExists
	}
	if exists, _ := u.UserRepository.ExistsByUsername(ctx, username); exists {
		return nil, "", domain.ErrUserAlreadyExists
	}
	hashedPass, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	newUser := &domain.User{
		Email: email,
		Username: username,
		Password: hashedPass,
	}

	if err := u.UserRepository.Create(ctx, newUser); err != nil {
		return nil, "", err
	}

	token, err := u.JWT.GenerateToken(newUser.ID)
	if err != nil {
		return nil, "", err
	}

	return newUser, token, nil
}

func (u *UserService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	if err := validator.ValidateUserLogin(email, password); err != nil {
		return nil, "", err
	}
	user, err := u.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", domain.ErrUserNotFound
	}

	if !auth.CheckPassword(user.Password, password){
		return nil, "", domain.ErrPasswordMismatch
	}

	token, err := u.JWT.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, userID uint64, username, email string) (*domain.User, error) {
	if err := validator.ValidateUserUpdate(email, username); err != nil {
		return nil, err
	}
	existsEmail, err := u.UserRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existsEmail {
		return nil, domain.ErrUserAlreadyExists
	}

	existsUsername, err := u.UserRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existsUsername {
		return nil, domain.ErrUserAlreadyExists
	}

	updatedUser := &domain.User{
		Email: email,
		Username: username,
	}

	if err := u.UserRepository.Update(ctx, userID, updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserService) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	user, err := u.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !auth.CheckPassword(user.Password, oldPassword) {
		return domain.ErrPasswordMismatch
	}

	if _, err := validator.IsValidPassword(newPassword); err != nil {
		return err
	}
	if oldPassword == newPassword {
		return domain.ErrPasswordSameAsOld
	}
	hashed, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return u.UserRepository.UpdatePassword(ctx, userID, hashed)
}

func (u *UserService) DeleteAccount(ctx context.Context, userID uint64) error {
	user, err := u.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.DeletedAt != nil {
		return nil
	}
	if err := u.NoteRepository.SoftDelete(ctx, userID); err != nil {
		return err
	}
	return u.UserRepository.SoftDelete(ctx, userID)
}

func (u *UserService) PermanentDeleteAccount(ctx context.Context, userID uint64) error {
	if _, err := u.UserRepository.GetByIDAny(ctx, userID); err != nil {
		return err
	}
	if err := u.NoteRepository.HardDelete(ctx, userID); err != nil {
		return err
	}
	return u.UserRepository.HardDelete(ctx, userID)
}

func (u *UserService) GetByID(ctx context.Context, userID uint64) (*domain.User, error) {
	user, err := u.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (u *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := u.UserRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (u *UserService) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	exists, err := u.UserRepository.ExistsByEmail(ctx, email);
	if err != nil {
		return false, domain.ErrEmailAlreadyExists
	}
	return exists, nil
}

func (u *UserService) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	exists, err := u.UserRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	users, total, err := u.UserRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u *UserService) GetTotalUserCount(ctx context.Context) (uint64, error) {
	count, err := u.UserRepository.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
