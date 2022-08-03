package user

import (
	"backend/helper"
	"gorm.io/gorm"
	"math"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(id int) (User, error)
	Update(user User) (User, error)
	FindEmail(email string) (User, error)
	DeleteUser(id int) (User, error)
	//Pagination(limit int) (User, error)
	//Pagination(pagination *Pagination) (User, error)
	//Pagination(pagination Paginations) ([]User, error)
	PaginationUser(pagination *helper.Pagination) (error, *helper.Pagination, int)
	ChangePassword(user User) (User, error)
	GetUsers() ([]User, error)
	CountUsers() (total int64)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err

	}

	return user, nil
}
func (r *repository) FindByEmail(email string) (User, error) {

	var user User

	err := r.db.Where("email=?", email).Find(&user).Error
	if err != nil {
		return user, nil
	}
	return user, nil

}
func (r *repository) FindByID(id int) (User, error) {

	var user User

	err := r.db.Where("id=?", id).Find(&user).Error
	if err != nil {
		return user, nil
	}
	return user, nil

}

func (r *repository) FindEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email=?", email).Find(&user).Error
	if err != nil {
		return user, nil
	}
	return user, nil
}
func (r *repository) DeleteUser(id int) (User, error) {
	var user User

	//r.db.First(&user, id)
	r.db.Delete(&user, id)

	//err := r.db.Where("id", id).Find(&user).Error
	//if err != nil {
	//	return user, nil
	//}
	////er := r.db.Delete(id)
	//if deleted != nil {
	//	return user, nil
	//}
	return user, nil
}
func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

//
//func (r *repository) Pagination(limit int) (User, error) {
//	var user User
//
//	err := r.db.Limit(limit).Offset(0).Find(&user)
//	if err != nil {
//		return user, nil
//	}
//
//	return user, nil
//}
//func (r *repository) Pagination(pagination *Paginations) ([]User, *Paginations, error, int) {
//	var users []User
//
//	var totalPages, totalRows int64
//	fromRow, toRow := 1, 1
//
//	offset := pagination.Page * pagination.Limit
//
//	// get data with limit, offset & order
//	errfind := r.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Find(&users).Error
//
//	var userResponse []DisplayUsers
//	for _, u := range users {
//		userRespond := DisplayUsers{
//			ID:    u.ID,
//			Name:  u.Name,
//			Email: u.Email,
//		}
//		userResponse = append(userResponse, userRespond)
//	}
//	pagination.Rows = userResponse
//
//	errCount := r.db.Model(User{}).Count(&totalRows).Error
//
//	//if errCount != nil {
//	//	return
//	//}
//
//	pagination.TotalRows = int(totalRows)
//
//	totalPages = int64(int(math.Ceil(float64(totalRows)/float64(pagination.Limit))) - 1)
//
//	if pagination.Page == 0 {
//		fromRow = 1
//		toRow = pagination.Limit
//	} else {
//		if pagination.Page <= totalPages {
//			fromRow = pagination.Page*pagination.Limit + 1
//			toRow = (pagination.Page + 1) * pagination.Limit
//		}
//	}
//	if toRow > int(totalRows) {
//		toRow = int(totalRows)
//	}
//
//	pagination.FromRow = fromRow
//	pagination.ToRow = toRow
//
//	return errCount, pagination, totalPages
//
//}
func (r *repository) PaginationUser(pagination *helper.Pagination) (error, *helper.Pagination, int) {
	var users []User

	var totalRows int64

	var totalPages int
	//var toRow int
	//var fromRow int

	offset := (pagination.Page - 1) * pagination.Limit

	errFind := r.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Find(&users).Error
	var usersResponse []DisplayUsers

	for _, u := range users {
		userRespond := DisplayUsers{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}

		usersResponse = append(usersResponse, userRespond)
	}

	if errFind != nil {
		return errFind, pagination, totalPages
	}
	pagination.Rows = usersResponse

	errCount := r.db.Model(&User{}).Count(&totalRows).Error

	if errCount != nil {
		return errCount, pagination, totalPages
	}

	pagination.TotalRows = int(totalRows)

	totalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	if pagination.Page == 0 {
		pagination.Page += 1
		//fromRow = 1
		//toRow = pagination.Limit
	}
	if pagination.Page <= totalPages {
		//fromRow = pagination.Page*pagination.Limit + 1
		//toRow = (pagination.Page + 1) * pagination.Limit
	}
	//if toRow > int(totalRows) {
	//	toRow = int(totalRows)
	//}

	//pagination.FromRow = fromRow
	//pagination.ToRow = toRow
	pagination.TotalPages = totalPages

	return errCount, pagination, totalPages
}
func (r *repository) ChangePassword(user User) (User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
func (r *repository) GetUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return users, nil
	}
	//fmt.Println("data user", users)
	return users, nil
}
func (r *repository) CountUsers() (total int64) {

	r.db.Model(&User{}).Distinct("id").Count(&total)
	return total

}
