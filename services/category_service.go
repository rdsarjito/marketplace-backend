package services

import (
	"errors"

	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/request"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/domain/model"
	"github.com/rdsarjito/marketplace-backend/repositories"
)

type CategoryService interface {
	GetListCategory() ([]response.CategoryResponse, error)
	GetDetailCategory(id int) (*response.CategoryResponse, error)
	CreateCategory(req *request.CreateCategoryRequest) (*response.CategoryResponse, error)
	UpdateCategory(id int, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error)
	DeleteCategory(id int) error
    EnsureDefaultCategories(names []string) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetListCategory() ([]response.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var categoryResponses []response.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.CategoryResponse{
			ID:        category.ID,
			Nama:      category.Nama,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return categoryResponses, nil
}

func (s *categoryService) GetDetailCategory(id int) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrCategoryNotFound)
	}

	categoryResponse := &response.CategoryResponse{
		ID:        category.ID,
		Nama:      category.Nama,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return categoryResponse, nil
}

func (s *categoryService) CreateCategory(req *request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	category := &model.Category{
		Nama: req.Nama,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	categoryResponse := &response.CategoryResponse{
		ID:        category.ID,
		Nama:      category.Nama,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return categoryResponse, nil
}

func (s *categoryService) UpdateCategory(id int, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(constants.ErrCategoryNotFound)
	}

	category.Nama = req.Nama

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	categoryResponse := &response.CategoryResponse{
		ID:        category.ID,
		Nama:      category.Nama,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return categoryResponse, nil
}

func (s *categoryService) DeleteCategory(id int) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New(constants.ErrCategoryNotFound)
	}

	return s.categoryRepo.Delete(id)
}

// EnsureDefaultCategories creates categories if they don't exist
func (s *categoryService) EnsureDefaultCategories(names []string) error {
    for _, name := range names {
        if name == "" { continue }
        if _, err := s.categoryRepo.GetByName(name); err == nil {
            // already exists
            continue
        }
        // create
        if err := s.categoryRepo.Create(&model.Category{Nama: name}); err != nil {
            return err
        }
    }
    return nil
}
