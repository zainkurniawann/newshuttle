package services

import (
	"database/sql"
	"encoding/json"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
)

type ChildernServiceInterface interface {
	GetAllChilderns(id string) ([]dto.StudentResponseDTO, int, error)
	GetSpecChildern(id string) (dto.StudentResponseDTO, error)
	UpdateChildern(id string, req dto.StudentRequestByParentDTO, username string) error
	UpdateChildernStatus(id string, req dto.StudentStatusRequestByParentDTO, username string) error
}

type ChildernService struct {
	ChildernRepository repositories.ChildernRepositoryInterface
}

func NewChildernService(childernRepository repositories.ChildernRepositoryInterface) ChildernServiceInterface {
	return &ChildernService{
		ChildernRepository: childernRepository,
	}
}

func (service *ChildernService) GetAllChilderns(id string) ([]dto.StudentResponseDTO, int, error) {
	childerns, err := service.ChildernRepository.FetchAllChilderns(id)
	if err != nil {
		return nil, 0, err
	}

	var childernsDTO []dto.StudentResponseDTO

	for _, childern := range childerns {
		childernsDTO = append(childernsDTO, dto.StudentResponseDTO{
			UUID:           childern.UUID.String(),
			FirstName:      childern.FirstName,
			LastName:       childern.LastName,
			Grade:          childern.Grade,
			StudentAddress: childern.StudentAddress.String,
			PickupPoint:    childern.StudentPickupPoint.String,
			Status:         childern.Status,
			Gender:         childern.Gender,
			SchoolUUID:     childern.SchoolUUID.String(),
			SchoolName:     childern.SchoolName,
		})
	}
	total := len(childernsDTO)
	return childernsDTO, total, nil
}

func (service *ChildernService) GetSpecChildern(id string) (dto.StudentResponseDTO, error) {
	childern, err := service.ChildernRepository.FetchSpecChildern(id)
	if err != nil {
		return dto.StudentResponseDTO{}, err
	}

	Address := ""
	if childern.StudentAddress.Valid {
		Address = childern.StudentAddress.String
	}

	pickupPoint := ""
	if childern.StudentPickupPoint.Valid {
		pickupPoint = childern.StudentPickupPoint.String
	}

	studentDTO := dto.StudentResponseDTO{
		UUID:           childern.UUID.String(),
		FirstName:      childern.FirstName,
		LastName:       childern.LastName,
		Gender:         childern.Gender,
		StudentAddress: Address,
		PickupPoint:    pickupPoint,
		Grade:          childern.Grade,
		Status:         childern.Status,
		SchoolUUID:     childern.SchoolUUID.String(),
	}

	return studentDTO, nil
}

func (service *ChildernService) UpdateChildern(id string, req dto.StudentRequestByParentDTO, username string) error {
	var pickupPointJSON []byte
	var err error
	if req.StudentPickupPoint != nil {
		pickupPointJSON, err = json.Marshal(req.StudentPickupPoint)
		if err != nil {
			return err
		}
	}

	student := entity.Student{
		FirstName:         req.StudentFirstName,
		LastName:          req.StudentLastName,
		Gender:            string(req.StudentGender),
		StudentAddress:    sql.NullString{String: req.StudentAddress, Valid: req.StudentAddress != ""},
		StudentPickupPoint: sql.NullString{String: string(pickupPointJSON), Valid: req.StudentPickupPoint != nil},
		Status:            req.StudentStatus,
		UpdatedBy:         sql.NullString{String: username, Valid: username != ""},
	}

	return service.ChildernRepository.UpdateChildern(student, id)
}

func (service *ChildernService) UpdateChildernStatus(id string, req dto.StudentStatusRequestByParentDTO, username string) error {
	student := entity.Student{
		Status:    req.StudentStatus,
		UpdatedBy: sql.NullString{String: username, Valid: username != ""},
	}

	return service.ChildernRepository.UpdateChildernStatus(student, id)
}