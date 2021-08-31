package getrequest

import (
	"net/url"
)

// Manage the HTTP GET request parameters
type GetRequest struct {
	urls url.Values
}

// Initializer
func (p *GetRequest) Init() *GetRequest {
	p.urls = url.Values{}
	return p
}

// Initialized from another instance
func (p *GetRequest) InitFrom(reqParams *GetRequest) *GetRequest {
	if reqParams != nil {
		p.urls = reqParams.urls
	} else {
		p.urls = url.Values{}
	}
	return p
}

// Add URL escape property and value pair
func (p *GetRequest) AddParam(property string, value string) *GetRequest {
	if property != "" && value != "" {
		p.urls.Add(property, value)
	}
	return p
}

// Concat the property and value pair
func (p *GetRequest) BuildParams() string {
	return p.urls.Encode()
}
