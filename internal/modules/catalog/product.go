package catalog

import "time"

type ProductType string

const (
	ProductTypeVPS          ProductType = "vps"
	ProductTypeProxy        ProductType = "proxy"
	ProductTypeServiceAddon ProductType = "service_addon"
)

func (productType ProductType) Valid() bool {
	switch productType {
	case ProductTypeVPS, ProductTypeProxy, ProductTypeServiceAddon:
		return true
	default:
		return false
	}
}

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusDisabled ProductStatus = "disabled"
	ProductStatusArchived ProductStatus = "archived"
)

func (status ProductStatus) Valid() bool {
	switch status {
	case ProductStatusDraft, ProductStatusActive, ProductStatusDisabled, ProductStatusArchived:
		return true
	default:
		return false
	}
}

type Product struct {
	ID           ProductID
	DisplayID    int64
	Type         ProductType
	Name         string
	Description  string
	Status       ProductStatus
	DisplayOrder int
	CreatedBy    UserID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateProductInput struct {
	Type         ProductType
	Name         string
	Description  string
	Status       ProductStatus
	DisplayOrder int
	CreatedBy    UserID
}

func (input CreateProductInput) Normalize() CreateProductInput {
	output := input
	output.Name = trim(output.Name)
	output.Description = trim(output.Description)
	output.CreatedBy = UserID(trim(string(output.CreatedBy)))
	if output.Status == "" {
		output.Status = ProductStatusDraft
	}
	return output
}

func (input CreateProductInput) Validate() error {
	if !input.Type.Valid() {
		return ErrProductTypeInvalid
	}
	if input.Name == "" {
		return ErrProductNameMissing
	}
	if !input.Status.Valid() {
		return ErrProductStatusInvalid
	}
	if input.CreatedBy == "" {
		return ErrCreatedByMissing
	}
	return nil
}
