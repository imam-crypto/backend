package handler

import (
	"backend/auth"
	"backend/helper"
	"backend/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type userhandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userhandler {
	return &userhandler{userService, authService}
}
func (h *userhandler) RegisterUser(c *gin.Context) {
	// menangkap input dari user
	// mapping input dari user ke struct RegisterUser
	// struct di atas di parsing menjadi parameter ke service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse(http.StatusHTTPVersionNotSupported, false, "failed", errorMessage)
		c.JSON(http.StatusHTTPVersionNotSupported, response)
		return
	}

	//cek email
	newCek, r := h.userService.IsEmailAvailable(input.Email)
	//bataas cek
	if r != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, false, "user  not registered", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if newCek != true {
		response := helper.APIResponse(http.StatusUnprocessableEntity, newCek, "email has been registered", newCek)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	NewUser, er := h.userService.RegisterUser(input)
	if er != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, false, er.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//data := gin.H{
	//	"is_available": isEmailAvailable,
	//}

	token, er := h.authService.GenerateToken(NewUser.ID)
	if er != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, false, er.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(NewUser, token)
	result := helper.APIResponse(http.StatusOK, true, "success", formatter)
	c.JSON(http.StatusOK, result)
}
func (h *userhandler) Login(c *gin.Context) {

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		//errors := helper.FormatError(err)
		//errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse(http.StatusNotFound, false, "Login Failed", "login failed")
		c.JSON(http.StatusNotFound, response)
		return
	}
	loogedinUser, er := h.userService.Login(input)

	if er != nil {
		//errorMessage := gin.H{"errors": er.Error()}
		response := helper.APIResponse(http.StatusNotFound, false, "Login Failed", "login failed")
		c.JSON(http.StatusNotFound, response)
		return
	}
	token, er := h.authService.GenerateToken(loogedinUser.ID)
	if er != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, false, er.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatUser(loogedinUser, token)
	response := helper.APIResponse(http.StatusOK, true, "Login Succesfully", formatter)
	c.JSON(http.StatusOK, response)
}
func (h *userhandler) DeleteUser(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	_, er := h.userService.DeleteUser(id)

	if er != nil {
		//errorMessage := gin.H{"errors": er.Error()}
		response := helper.APIResponse(http.StatusUnprocessableEntity, false, "delete failed", "delete failed")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	//formatter := user.FormatUser(deleted)
	response := helper.APIResponse(http.StatusOK, true, "user deleted", nil)
	c.JSON(http.StatusOK, response)
}

func (h *userhandler) Update(c *gin.Context) {
	// menangkap input dari user
	// mapping input dari user ke struct RegisterUser
	// struct di atas di parsing menjadi parameter ke service

	var input user.EditInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse(http.StatusHTTPVersionNotSupported, false, "failed", errorMessage)
		c.JSON(http.StatusHTTPVersionNotSupported, response)
		return
	}

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	NewUser, er := h.userService.UpdateUser(id, input)
	if er != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, false, er.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(NewUser, idString)
	result := helper.APIResponse(http.StatusOK, true, "edited success", formatter)
	c.JSON(http.StatusOK, result)
}

//func (h *userhandler) Paginate(c *gin.Context) {
//	page, _ := strconv.Atoi(c.Query("page"))
//	//perPage := 2
//
//	fmt.Println(page)
//
//	NewUser, er := h.userService.User(page)
//	if er != nil {
//		response := helper.APIResponse(http.StatusUnprocessableEntity, false, er.Error(), nil)
//		c.JSON(http.StatusBadRequest, response)
//		return
//	}
//	formatter := user.FormatUser(NewUser)
//	result := helper.APIResponse(http.StatusOK, true, "edited success", formatter)
//	c.JSON(http.StatusOK, result)
//}
func (h *userhandler) PaginationUser(c *gin.Context) {

	pagination := helper.GeneratePaginationRequest(c)
	//if pagination.Page == 0 {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"errors": nil,
	//	})
	//	return
	//}

	responseData, _ := h.userService.PaginationUser(c, pagination)

	if responseData.Is_Success == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": responseData.Message,
		})
	}

	response := responseData
	c.JSON(http.StatusOK, response)
	//c.JSON(http.StatusOK, gin.H{
	//	"data": response,
	//})
}
func (h *userhandler) ChangesPassword(c *gin.Context) {

	var input user.InputChangesPassword

	err := c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(http.StatusBadRequest, false, "failed", errorMessage)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("current_user").(user.User)
	UserID := currentUser.ID

	change, er := h.userService.ChangePassword(UserID, input)

	if er != nil {
		return
	}
	result := helper.APIResponse(http.StatusOK, true, "Changes Password Succesfully", change)
	c.JSON(http.StatusOK, result)
}
func (h *userhandler) GetUsers(c *gin.Context) {

	err, user := h.userService.GetUsers()
	if err != nil {
		result := helper.APIResponse(http.StatusBadRequest, true, "failed", err)
		c.JSON(http.StatusBadRequest, result)
		return
	}
	//fmt.Println("user :", user)

	result := helper.APIResponse(http.StatusOK, true, "link data users :"+" "+"http://"+c.Request.Host+"/assets/"+"DataUser.pdf", user)
	c.JSON(http.StatusOK, result)
	fmt.Println("tes branch dev")
	//return nil
}
