package main

import (
	"context"
	"fmt"
	"time"

	"github.com/maqsatto/Notes-API/internal/auth"
	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/database"
	"github.com/maqsatto/Notes-API/internal/domain"
	"github.com/maqsatto/Notes-API/internal/logger"
	"github.com/maqsatto/Notes-API/internal/repository"
	"github.com/maqsatto/Notes-API/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logg, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer logg.Close()

	ctx := context.Background()

	db, err := database.NewPostgresDB(ctx, cfg.Database)
	if err != nil {
		logg.Error("failed to connect to database", err)
		return
	}
	defer db.Close()
	fmt.Println("DB connected")

	userRepo := repository.NewUserRepo(db)
	noteRepo := repository.NewNoteRepo(db)

	ttl := time.Duration(cfg.JWT.ExpiryHour) * time.Hour
	issuer := "notes-api" // <- just a constant

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, issuer, ttl)

	userService := service.NewUserService(*userRepo, *noteRepo, jwtManager)

	runUserServiceSmoke(ctx, *userService)

}


func runUserServiceSmoke(ctx context.Context, userService service.UserService) {
	// 0) Test data (unique email әр запускта conflict болмас үшін)
	ts := time.Now().UnixNano()
	username := fmt.Sprintf("test_user_%d", ts)
	email := fmt.Sprintf("test_%d@example.com", ts)

	oldPass := "OldPass1!"
	newPass := "NewPass2!"

	var (
		createdUser *domain.User
		token1      string
	)

	type testCase struct {
		name string
		fn   func() error
	}

	tests := []testCase{
		{
			name: "Register",
			fn: func() error {
				u, t, err := userService.Register(ctx, username, email, oldPass)
				if err != nil {
					return err
				}
				createdUser = u
				token1 = t
				fmt.Println("Registered:", u.ID, u.Email, "token_len:", len(t))
				return nil
			},
		},
		{
			name: "Login",
			fn: func() error {
				u, t, err := userService.Login(ctx, email, oldPass)
				if err != nil {
					return err
				}
				fmt.Println("Logged in:", u.ID, u.Email, "token_len:", len(t), "token_changed:", t != token1)
				return nil
			},
		},
		{
			name: "IsEmailTaken (should be true)",
			fn: func() error {
				taken, err := userService.IsEmailTaken(ctx, email)
				if err != nil {
					return err
				}
				fmt.Println("IsEmailTaken:", taken)
				if !taken {
					return fmt.Errorf("expected email to be taken")
				}
				return nil
			},
		},
		{
			name: "IsUsernameTaken (should be true)",
			fn: func() error {
				taken, err := userService.IsUsernameTaken(ctx, username)
				if err != nil {
					return err
				}
				fmt.Println("IsUsernameTaken:", taken)
				if !taken {
					return fmt.Errorf("expected username to be taken")
				}
				return nil
			},
		},
		{
			name: "GetByID",
			fn: func() error {
				u, err := userService.GetByID(ctx, createdUser.ID)
				if err != nil {
					return err
				}
				fmt.Println("GetByID:", u.ID, u.Email, u.Username)
				return nil
			},
		},
		{
			name: "GetByEmail",
			fn: func() error {
				u, err := userService.GetByEmail(ctx, email)
				if err != nil {
					return err
				}
				fmt.Println("GetByEmail:", u.ID, u.Email)
				return nil
			},
		},
		{
			name: "GetByUsername",
			fn: func() error {
				u, err := userService.GetByUsername(ctx, username)
				if err != nil {
					return err
				}
				fmt.Println("GetByUsername:", u.ID, u.Username)
				return nil
			},
		},
		{
			name: "UpdateProfile",
			fn: func() error {
				newUsername := username + "_upd"
				newEmail := fmt.Sprintf("upd_%d@example.com", ts)

				u, err := userService.UpdateProfile(ctx, createdUser.ID, newUsername, newEmail)
				if err != nil {
					return err
				}

				// local state update (келесі тесттер жаңа мәнді қолдансын)
				username = newUsername
				email = newEmail

				fmt.Println("Updated:", u.Username, u.Email)
				return nil
			},
		},
		{
			name: "ChangePassword",
			fn: func() error {
				if err := userService.ChangePassword(ctx, createdUser.ID, oldPass, newPass); err != nil {
					return err
				}
				// келесі login үшін
				oldPass = newPass
				fmt.Println("Password changed")
				return nil
			},
		},
		{
			name: "Login (after password change)",
			fn: func() error {
				u, _, err := userService.Login(ctx, email, oldPass)
				if err != nil {
					return err
				}
				fmt.Println("Logged in after pw change:", u.ID)
				return nil
			},
		},
		{
			name: "ListUsers",
			fn: func() error {
				users, total, err := userService.ListUsers(ctx, 10, 0)
				if err != nil {
					return err
				}
				fmt.Println("ListUsers total:", total, "returned:", len(users))
				return nil
			},
		},
		{
			name: "GetTotalUserCount",
			fn: func() error {
				total, err := userService.GetTotalUserCount(ctx)
				if err != nil {
					return err
				}
				fmt.Println("GetTotalUserCount:", total)
				return nil
			},
		},
		{
			name: "DeleteAccount (soft)",
			fn: func() error {
				if err := userService.DeleteAccount(ctx, createdUser.ID); err != nil {
					return err
				}
				fmt.Println("Soft deleted")
				return nil
			},
		},
		{
			name: "PermanentDeleteAccount (hard)",
			fn: func() error {
				if err := userService.PermanentDeleteAccount(ctx, createdUser.ID); err != nil {
					return err
				}
				fmt.Println("Hard deleted")
				return nil
			},
		},
	}

	for _, tc := range tests {
		fmt.Println("===", tc.name, "===")
		if err := tc.fn(); err != nil {
			fmt.Println("FAILED:", tc.name, "err:", err)
			return
		}
		fmt.Println("OK:", tc.name)
	}
}
