package handler

import (
	// "log"
	"errors"
	"net/http"
	"strconv"
	"twittir-go/internal/domain"
	"twittir-go/internal/dto"
	"twittir-go/internal/helpers"
	"twittir-go/internal/services"
	"twittir-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	appJSON = "application/json"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService}
}

// Register godoc
// @Summary Register user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param fullname body string true "Fullname"
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param password_verify body string true "Password verification"
// @Success 200 {object} dto.SuccessResponse{data=dto.RegisterSuccess}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if registerRequest.Password != registerRequest.PasswordVerify {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid password verification")
		return
	}

	user, err := h.userService.Register(registerRequest.Username, registerRequest.Email, registerRequest.FullName, registerRequest.Password, registerRequest.PasswordVerify)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, dto.RegisterSuccess{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	})
}

// UserLogin godoc
// @Summary Login user
// @Description Logs in a user and returns a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {object} dto.SuccessResponse{data=dto.SignInSuccess}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest

	// Bind the JSON request body into loginRequest struct
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call the service to perform the login
	token, err := h.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		// If login fails, return unauthorized response
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// If login is successful, return a token
	utils.RespondWithSuccess(c, http.StatusOK, dto.SignInSuccess{
		Token: token,
	})
}

// UpdateProfile godoc
// @Summary Search user by username
// @Description Search user by username with query parameter
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param fullname body string true "Fullname"
// @Param bio body string true "Bio"
// @Param username body string true "Username"
// @Success 200 {object} dto.SuccessResponse{data=dto.FormatUsers} "Users successfully retrieved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or missing parameters"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Router /api/user/update [get]
func (h *UserHandler) UpdateProfile(c *gin.Context) {

	userID, err := h.extractUserID(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Bind input data
	updateUserProfile, err := h.bindUserData(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	// Update user profile
	updatedUser, err := h.userService.UpdateProfile(int(userID), updateUserProfile)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Kirim response sukses
	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"username": updatedUser.Username,
		"fullname": updatedUser.FullName,
		"bio":      updatedUser.Bio,
	})
}

// ShowProfile godoc
// @Summary Show Profile user
// @Description Show user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SuccessResponse{data=dto.ProfileResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/profile [get]
func (h *UserHandler) ShowProfile(c *gin.Context) {

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	user, err := h.userService.GetUserByID(int(userID))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, dto.ProfileResponse{
		ID:             user.ID,
		FullName:       user.FullName,
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
	})
}

// SearchUser godoc
// @Summary Search user by username
// @Description Search user by username with query parameter
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "Username to search for"
// @Success 200 {object} dto.SuccessResponse{data=dto.FormatUsers} "Users successfully retrieved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or missing parameters"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Router /api/search [get]
func (h *UserHandler) SearchUser(c *gin.Context) {

	usernameParam := c.DefaultQuery("username", "")

	users, err := h.userService.SearchUserByUsername(usernameParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var formattedUsers []dto.FormatUsers
	for _, user := range users {
		formattedUsers = append(formattedUsers, dto.FormatUsers{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
		})
	}

	utils.RespondWithSuccess(c, http.StatusOK, formattedUsers)
}

// GetUserByID godoc
// @Summary Show user by ID
// @Description Retrieve the user profile details by the specified user ID
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SuccessResponse{data=dto.ProfileResponse} "User profile successfully retrieved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or missing parameters"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDInt := int(userID)

	user, err := h.userService.GetUserByID(userIDInt)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK,
		dto.ProfileResponse{
			ID:             user.ID,
			FullName:       user.FullName,
			Username:       user.Username,
			ProfilePicture: user.ProfilePicture,
		})
}

// Helper untuk extract user ID dari token
func (h *UserHandler) extractUserID(c *gin.Context) (uint, error) {
	userData, exists := c.Get("userData")
	if !exists {
		return 0, errors.New("userData not found")
	}

	claims, ok := userData.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID format")
	}

	return uint(userID), nil
}

// Helper untuk bind data input dari request
func (h *UserHandler) bindUserData(c *gin.Context) (*domain.User, error) {
	var user domain.User
	contentType := helpers.GetContentType(c)

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&user); err != nil {
			return nil, err
		}
	} else {
		if err := c.ShouldBind(&user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}
