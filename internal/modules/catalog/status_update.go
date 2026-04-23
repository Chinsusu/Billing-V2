package catalog

type UpdateProductStatusInput struct {
	ID     ProductID
	Status ProductStatus
}

func (input UpdateProductStatusInput) Normalize() UpdateProductStatusInput {
	output := input
	output.ID = ProductID(trim(string(output.ID)))
	output.Status = ProductStatus(trim(string(output.Status)))
	return output
}

func (input UpdateProductStatusInput) Validate() error {
	if input.ID.Empty() {
		return ErrProductIDMissing
	}
	if !input.Status.Valid() {
		return ErrProductStatusInvalid
	}
	return nil
}

type UpdatePlanStatusInput struct {
	ID     PlanID
	Status PlanStatus
}

func (input UpdatePlanStatusInput) Normalize() UpdatePlanStatusInput {
	output := input
	output.ID = PlanID(trim(string(output.ID)))
	output.Status = PlanStatus(trim(string(output.Status)))
	return output
}

func (input UpdatePlanStatusInput) Validate() error {
	if input.ID.Empty() {
		return ErrPlanIDMissing
	}
	if !input.Status.Valid() {
		return ErrPlanStatusInvalid
	}
	return nil
}

type UpdateProviderSourceStatusInput struct {
	ID     ProviderSourceID
	Status ProviderSourceStatus
}

func (input UpdateProviderSourceStatusInput) Normalize() UpdateProviderSourceStatusInput {
	output := input
	output.ID = ProviderSourceID(trim(string(output.ID)))
	output.Status = ProviderSourceStatus(trim(string(output.Status)))
	return output
}

func (input UpdateProviderSourceStatusInput) Validate() error {
	if input.ID.Empty() {
		return ErrSourceIDMissing
	}
	if !input.Status.Valid() {
		return ErrSourceStatusInvalid
	}
	return nil
}
