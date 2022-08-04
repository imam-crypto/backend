package user

import (
	"backend/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	IsEmailAvailable(input string) (bool, error)
	FindById(ID int) (User, error)
	Login(LoginInput) (User, error)
	DeleteUser(id int) (User, error)
	UpdateUser(ID int, input EditInput) (User, error)
	//User(Limit int) (User, error)
	PaginationUser(c *gin.Context, pagination *helper.Pagination) (helper.Response, error)
	ChangePassword(ID int, input InputChangesPassword) (User, error)
	GetUsers() (error, [][]string)
}
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, nil
	}
	user.Password = string(passwordHash)
	// file,err := user.Form.
	newUser, er := s.repository.Save(user)

	if er != nil {
		return newUser, er
	}
	return newUser, nil
}
func (s *service) IsEmailAvailable(email string) (bool, error) {
	//email = email

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}
	return false, nil

}
func (s *service) FindById(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)

	if err != nil {
		return user, err
	}

	return user, nil

}
func (s *service) Login(input LoginInput) (User, error) {

	email := input.Email
	password := input.Password

	user, err := s.repository.FindEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User Not Found")
	}

	er := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if er != nil {
		return user, er
	}
	return user, nil
}
func (s *service) DeleteUser(id int) (User, error) {

	user, err := s.repository.DeleteUser(id)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (s *service) UpdateUser(ID int, input EditInput) (User, error) {

	us, err := s.repository.FindByID(ID)

	if err != nil {
		return us, err
	}
	us.Name = input.Name
	us.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return us, nil
	}
	password := string(passwordHash)

	us.Password = password
	UpdatedUser, er := s.repository.Update(us)
	if er != nil {
		return UpdatedUser, err
	}
	return UpdatedUser, nil
}

func (s *service) PaginationUser(c *gin.Context, pagination *helper.Pagination) (helper.Response, error) {

	err, datapaginations, totalPages := s.repository.PaginationUser(pagination)

	if err != nil {

		result := helper.APIResponse(http.StatusNotFound, false, "Data not found", nil)
		c.JSON(http.StatusNotFound, result)

	}
	urlPath := c.Request.URL.Path
	datapaginations.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, 1, pagination.Sort)
	datapaginations.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, totalPages, pagination.Sort)

	if datapaginations.Page > 0 {

		datapaginations.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, datapaginations.Page-1, pagination.Sort)

	}
	if datapaginations.Page < totalPages {
		datapaginations.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, datapaginations.Page+1, pagination.Sort)
	}
	if datapaginations.Page > totalPages {
		datapaginations.PreviousPage = " "
	}

	//return helper.Response{http.StatusOK, true, "data", datapaginations}, err
	return helper.Response{http.StatusOK, true, "OK", datapaginations.Rows, helper.Paginate{Sort: datapaginations.Sort, Limit: datapaginations.Limit, FirstPage: datapaginations.FirstPage, NextPage: datapaginations.NextPage, LastPage: datapaginations.LastPage, PreviousPage: datapaginations.PreviousPage, Page: datapaginations.Page, TotalPages: datapaginations.TotalPages, TotalRows: datapaginations.TotalRows}}, err

}
func (s *service) ChangePassword(ID int, input InputChangesPassword) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, nil
	}
	user.Password = string(passwordHash)
	change, err := s.repository.ChangePassword(user)
	if err != nil {
		return user, err
	}
	return change, err
}
func (s *service) GetUsers() (error, [][]string) {
	//var user []User
	rows := [][]string{}
	totalUser := s.repository.CountUsers()
	convert := int(totalUser)
	fmt.Println(convert)
	//for i := 1; i <= convert; i++ {
	users, err := s.repository.GetUsers()

	if err != nil {
		return err, nil
	}

	for _, users := range users {

		id := strconv.Itoa(users.ID)
		Name := users.Name
		Email := users.Email
		//Created_At := users.CreatedAt.String()
		rows = append(rows, []string{id, Name, Email})
	}

	header := []string{"ID", "Name", "Email"}
	pdf := helper.SetToPDF()
	pdf = helper.Header(pdf, header)
	pdf = helper.Table(pdf, rows)
	helper.SaveFile(pdf)
	//if pdf.Err() {
	//	return nil
	//}
	//err := helper.SaveFile(pdf)
	//if err != nil {
	//	return helper.Response{Status: http.StatusOK, Is_Success: true}
	//}
	//}
	//fmt.Println(rows)

	//result := helper.APIResponse(http.StatusOK, true, "http://"+c.Request.Host+"/assets/"+"DataUser.pdf", user)
	//c.JSON(http.StatusOK, result)

	return err, rows

}
