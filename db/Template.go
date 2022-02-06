package db

import (
	"encoding/json"
)

type TemplateType string

const (
	TemplateTask   TemplateType = ""
	TemplateBuild  TemplateType = "build"
	TemplateDeploy TemplateType = "deploy"
)

// Template is a user defined model that is used to run a task
type Template struct {
	ID int `db:"id" json:"id"`

	ProjectID     int  `db:"project_id" json:"project_id"`
	InventoryID   int  `db:"inventory_id" json:"inventory_id"`
	RepositoryID  int  `db:"repository_id" json:"repository_id"`
	EnvironmentID *int `db:"environment_id" json:"environment_id"`

	// Alias as described in https://github.com/ansible-semaphore/semaphore/issues/188
	Alias string `db:"alias" json:"alias"`
	// playbook name in the form of "some_play.yml"
	Playbook string `db:"playbook" json:"playbook"`
	// to fit into []string
	Arguments *string `db:"arguments" json:"arguments"`
	// if true, semaphore will not prepend any arguments to `arguments` like inventory, etc
	OverrideArguments bool `db:"override_args" json:"override_args"`

	Removed bool `db:"removed" json:"-"`

	Description *string `db:"description" json:"description"`

	VaultKeyID *int      `db:"vault_key_id" json:"vault_key_id"`
	VaultKey   AccessKey `db:"-" json:"-"`

	Type            TemplateType `db:"type" json:"type"`
	StartVersion    *string      `db:"start_version" json:"start_version"`
	BuildTemplateID *int         `db:"build_template_id" json:"build_template_id"`

	ViewID *int `db:"view_id" json:"view_id"`

	LastTask *TaskWithTpl `db:"-" json:"last_task"`
	DynamicVars	string	  `db:"dynamic_vars" json:"dynamic_vars"`
}

func (tpl *Template) Validate() error {
	if tpl.Alias == "" {
		return &ValidationError{"template alias can not be empty"}
	}

	if tpl.Playbook == "" {
		return &ValidationError{"template playbook can not be empty"}
	}

	if tpl.Arguments != nil {
		if !json.Valid([]byte(*tpl.Arguments)) {
			return &ValidationError{"template arguments must be valid JSON"}
		}
	}

	return nil
}

func FillTemplates(d Store, templates []Template) (err error) {
	for i := range templates {
		tpl := &templates[i]
		var tasks []TaskWithTpl
		tasks, err = d.GetTemplateTasks(*tpl, RetrieveQueryParams{Count: 1})
		if err != nil {
			return
		}
		if len(tasks) > 0 {
			tpl.LastTask = &tasks[0]
		}
	}

	return
}

func FillTemplate(d Store, template *Template) (err error) {
	if template.VaultKeyID != nil {
		template.VaultKey, err = d.GetAccessKey(template.ProjectID, *template.VaultKeyID)
	}

	if err != nil {
		return
	}

	err = FillTemplates(d, []Template{*template})

	return
}
