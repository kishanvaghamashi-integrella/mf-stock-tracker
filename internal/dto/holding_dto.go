package dto

type HoldingResponseDto struct {
	ID                  int64   `json:"id"`
	AssetName           string  `json:"asset_name"`
	AssetInstrumentType string  `json:"asset_instrument_type"`
	Quantity            float64 `json:"quantity"`
	AveragePrice        float64 `json:"average_price"`
	InvestedPrice       float64 `json:"invested_price"`
}
