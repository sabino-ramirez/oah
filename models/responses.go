package models

type DbRow struct {
	Auth       string
	OrgId      int
	ProjTempId int
}

type Meta struct {
	CurrentPage  int
	PerPage      int
	TotalEntries int
}

type ProjectRequisition struct {
	Identifier              string
	Requisition_template_id int
	Status                  string
	Accession_status        string
	Processing_status       string
	Reporting_status        string
	Billing_status          string
	CreatedAt               string
	UpdatedAt               string
}

type ProjectRequisitions struct {
	Requisitions []ProjectRequisition
	Meta         Meta
}

type ProjectTemplate struct {
	Id           int
	ProjectName  string
	TemplateName string
}

type ProjectTemplates struct {
	Project_Templates []ProjectTemplate
}
