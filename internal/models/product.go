package models

type Product struct {
	ID        int
	Name      string
	GTIN      string
	LabelData string // JSON с данными для этикетки
}
