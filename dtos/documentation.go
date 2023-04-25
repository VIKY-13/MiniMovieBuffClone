package dtos

//these structs are for APIDocumentations
type EndpointsHead struct{
	Movie []EndpointDescriptions `json:"movie"`
	User []EndpointDescriptions `json:"user"`
	Favourite []EndpointDescriptions `json:"favourite"`
	Watchlist []EndpointDescriptions `json:"watchlist"`
}

type EndpointDescriptions struct{
	Method string `json:"method"`
	Endpoints string	`json:"endpoints"`
	Description []string	`json:"description"`
	Parameters []string	`json:"parameters"`
	Samplereqres SampleReqRes `json:"samplereqres"`
}

type SampleReqRes struct{
	ReqURL string `json:"requrl"`
	ReqBody string `json:"reqbody"`
	Response string `json:"response"`
}

type documentationparsedata struct{
	Title string
	Endpointsdata EndpointsHead
}