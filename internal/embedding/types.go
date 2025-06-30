package embedding

type Request struct {
	Text string `json:"inputText"`
}

type Response struct {
	Embedding []float64 `json:"embedding"`
}

type TitanRequest struct {
	InputText  string `json:"inputText"`
	Dimensions int    `json:"dimensions,omitempty"`
	Normalize  bool   `json:"normalize,omitempty"`
}

type TitanResponse struct {
	Embedding []float64 `json:"embedding"`
}
