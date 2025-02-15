package dto

type Link struct {
	CorrelationID string `json:"-"`
	Short         string `json:"short"`
	Long          string `json:"long"`
}
