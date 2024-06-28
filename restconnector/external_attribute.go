package restconnector

type ExternalAttributeValue struct {
	AccountIDPath          string                  `json:"accountIdPath"`
	ResponseColsToPropsMap map[string]any          `json:"responseColsToPropsMap"`
	Call                   []ExternalAttributeCall `json:"call"`
}

type ExternalAttributeCall struct {
	Name             string                         `json:"name"`
	Connection       string                         `json:"connection"`
	URL              string                         `json:"url"`
	HTTPMethod       string                         `json:"httpMethod"`
	HTTPParams       string                         `json:"httpParams"`
	HTTPHeaders      map[string]string              `json:"httpHeaders"`
	HTTPContentType  string                         `json:"httpContentType"`
	SuccessResponses ExternalAttributeCallResponses `json:"successResponses"`
}

type ExternalAttributeCallResponses struct {
	StatusCodes []int `json:"statusCode"`
}
