package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	jwttoken "github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parser: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}
	err = user.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to encrypt password: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	user.Password = string(hashPassword)

	repository.InsertNewUser(ctx.UserContext(), user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert user: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := user
	resp.Password = "" // Clear password before sending response

	return response.SendSuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	loginReq := new(models.LoginRequest)
	resp := new(models.LoginResponse)
	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResponse := fmt.Errorf("failed to parser: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}
	err = loginReq.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}
	user, err := repository.GetUserByUsername(ctx.UserContext(), loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get user: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username / password salah", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResponse := fmt.Errorf("failed to check password: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "username / password salah", nil)
	}
	token, err := jwttoken.GenerateToken(ctx.UserContext(), user.Username, user.FullName, "token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	refreshToken, err := jwttoken.GenerateToken(ctx.UserContext(), user.Username, user.FullName, "refresh_token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %w", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	resp.Username = user.Username
	resp.FullName = user.FullName
	resp.Token = token
	resp.RefreshToken = refreshToken
	return response.SendSuccessResponse(ctx, resp)
}
