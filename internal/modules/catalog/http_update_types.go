package catalog

type updateProductStatusRequest struct {
	Status ProductStatus `json:"status"`
}

func (request updateProductStatusRequest) toInput(id ProductID) UpdateProductStatusInput {
	return UpdateProductStatusInput{
		ID:     id,
		Status: request.Status,
	}
}

type updatePlanStatusRequest struct {
	Status PlanStatus `json:"status"`
}

func (request updatePlanStatusRequest) toInput(id PlanID) UpdatePlanStatusInput {
	return UpdatePlanStatusInput{
		ID:     id,
		Status: request.Status,
	}
}

type updateProviderSourceStatusRequest struct {
	Status ProviderSourceStatus `json:"status"`
}

func (request updateProviderSourceStatusRequest) toInput(id ProviderSourceID) UpdateProviderSourceStatusInput {
	return UpdateProviderSourceStatusInput{
		ID:     id,
		Status: request.Status,
	}
}
