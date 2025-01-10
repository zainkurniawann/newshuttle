package dto

type SchoolRequestDTO struct {
	Name        string `json:"name" validate:"required,max=255"`
	Address     string `json:"address" validate:"required,max=255"`
	Contact     string `json:"contact" validate:"required,phone"`
	Email       string `json:"email" validate:"required,email"`
	Description string `json:"description" validate:"omitempty,max=255"`
	Point    map[string]interface{} `json:"point"`
}

type SchoolResponseDTO struct {
	UUID           string `json:"school_uuid"`
	Name           string `json:"school_name"`
	AdminUUID      string `json:"admin_uuid,omitempty"`
	AdminName      string `json:"school_admin_name,omitempty"`
	Address        string `json:"school_address"`
	Contact        string `json:"school_contact"`
	Email          string `json:"school_email"`
	Description    string `json:"school_description,omitempty"`
	Point          string `json:"school_point,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	UpdatedBy      string `json:"updated_by,omitempty"`
}
