package dto

type CreateAssetRequest struct {
	Symbol             string `json:"symbol" validate:"required,min=1,max=20"`
	Name               string `json:"name" validate:"required,min=1,max=200"`
	InstrumentType     string `json:"instrument_type" validate:"required,instrument_type"`
	ISIN               string `json:"isin" validate:"required,min=10"`
	Exchange           string `json:"exchange" validate:"required,max=100"`
	Currency           string `json:"currency" validate:"omitempty,max=10"`
	ExternalPlatformID string `json:"external_platform_id" validate:"required,max=100"`
}

type UpdateAssetRequest struct {
	Symbol             *string `json:"symbol" validate:"omitempty,min=1,max=20"`
	Name               *string `json:"name" validate:"omitempty,min=1,max=200"`
	InstrumentType     *string `json:"instrument_type" validate:"omitempty,instrument_type"`
	ISIN               *string `json:"isin" validate:"omitempty,min=10"`
	Exchange           *string `json:"exchange" validate:"omitempty,max=100"`
	Currency           *string `json:"currency" validate:"omitempty,max=10"`
	ExternalPlatformID *string `json:"external_platform_id" validate:"omitempty,max=100"`
}
