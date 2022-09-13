package models

import "net/http"

type Client struct {
	Http              *http.Client
	OrganizationId    any
	ProjectTemplateId any
	Bearer            string
}

func NewClient(httpClient *http.Client, orgId, projTempId any, bearer string) *Client {
	return &Client{httpClient, orgId, projTempId, bearer}
}
