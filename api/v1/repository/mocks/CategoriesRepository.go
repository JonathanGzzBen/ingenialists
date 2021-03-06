// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/JonathanGzzBen/ingenialists/api/v1/models"
	mock "github.com/stretchr/testify/mock"
)

// CategoriesRepository is an autogenerated mock type for the CategoriesRepository type
type CategoriesRepository struct {
	mock.Mock
}

// CreateCategory provides a mock function with given fields: _a0
func (_m *CategoriesRepository) CreateCategory(_a0 *models.Category) (*models.Category, error) {
	ret := _m.Called(_a0)

	var r0 *models.Category
	if rf, ok := ret.Get(0).(func(*models.Category) *models.Category); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Category) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCategory provides a mock function with given fields: _a0
func (_m *CategoriesRepository) DeleteCategory(_a0 uint) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllCategories provides a mock function with given fields:
func (_m *CategoriesRepository) GetAllCategories() ([]models.Category, error) {
	ret := _m.Called()

	var r0 []models.Category
	if rf, ok := ret.Get(0).(func() []models.Category); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCategory provides a mock function with given fields: _a0
func (_m *CategoriesRepository) GetCategory(_a0 uint) (*models.Category, error) {
	ret := _m.Called(_a0)

	var r0 *models.Category
	if rf, ok := ret.Get(0).(func(uint) *models.Category); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCategory provides a mock function with given fields: _a0
func (_m *CategoriesRepository) UpdateCategory(_a0 *models.Category) (*models.Category, error) {
	ret := _m.Called(_a0)

	var r0 *models.Category
	if rf, ok := ret.Get(0).(func(*models.Category) *models.Category); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Category) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
