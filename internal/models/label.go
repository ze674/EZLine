package models

// LabelData содержит структурированные данные этикетки
type LabelData struct {
	Article     string `json:"article"`
	BoxQuantity string `json:"box_quantity"`
	BoxWeight   string `json:"box_weight"`
	GTIN        string `json:"gtin"`
	Header      string `json:"header"`
	LabelName   string `json:"label_name"`
	Standard    string `json:"standard"`
	UnitWeight  string `json:"unit_weight"`
}
