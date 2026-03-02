package dto

type CreateUserAssetRequest struct {
	AssetID int64 `json:"asset_id" validate:"required,gt=0"`
}
