package client

import (
	"github.com/rancher/norman/types"
)

const (
	ApplicationConfigurationTemplateType                      = "applicationConfigurationTemplate"
	ApplicationConfigurationTemplateFieldAnnotations          = "annotations"
	ApplicationConfigurationTemplateFieldAppTraits            = "appTraits"
	ApplicationConfigurationTemplateFieldComponents           = "components"
	ApplicationConfigurationTemplateFieldCreated              = "created"
	ApplicationConfigurationTemplateFieldCreatorID            = "creatorId"
	ApplicationConfigurationTemplateFieldLabels               = "labels"
	ApplicationConfigurationTemplateFieldName                 = "name"
	ApplicationConfigurationTemplateFieldOwnerReferences      = "ownerReferences"
	ApplicationConfigurationTemplateFieldParameters           = "parameters"
	ApplicationConfigurationTemplateFieldRemoved              = "removed"
	ApplicationConfigurationTemplateFieldState                = "state"
	ApplicationConfigurationTemplateFieldStatus               = "status"
	ApplicationConfigurationTemplateFieldTransitioning        = "transitioning"
	ApplicationConfigurationTemplateFieldTransitioningMessage = "transitioningMessage"
	ApplicationConfigurationTemplateFieldUUID                 = "uuid"
)

type ApplicationConfigurationTemplate struct {
	types.Resource
	Annotations          map[string]string               `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	AppTraits            *AppTraits                      `json:"appTraits,omitempty" yaml:"appTraits,omitempty"`
	Components           []Component                     `json:"components,omitempty" yaml:"components,omitempty"`
	Created              string                          `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string                          `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels               map[string]string               `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string                          `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences      []OwnerReference                `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Parameters           []Parameter                     `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Removed              string                          `json:"removed,omitempty" yaml:"removed,omitempty"`
	State                string                          `json:"state,omitempty" yaml:"state,omitempty"`
	Status               *ApplicationConfigurationStatus `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string                          `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string                          `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string                          `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ApplicationConfigurationTemplateCollection struct {
	types.Collection
	Data   []ApplicationConfigurationTemplate `json:"data,omitempty"`
	client *ApplicationConfigurationTemplateClient
}

type ApplicationConfigurationTemplateClient struct {
	apiClient *Client
}

type ApplicationConfigurationTemplateOperations interface {
	List(opts *types.ListOpts) (*ApplicationConfigurationTemplateCollection, error)
	Create(opts *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error)
	Update(existing *ApplicationConfigurationTemplate, updates interface{}) (*ApplicationConfigurationTemplate, error)
	Replace(existing *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error)
	ByID(id string) (*ApplicationConfigurationTemplate, error)
	Delete(container *ApplicationConfigurationTemplate) error
}

func newApplicationConfigurationTemplateClient(apiClient *Client) *ApplicationConfigurationTemplateClient {
	return &ApplicationConfigurationTemplateClient{
		apiClient: apiClient,
	}
}

func (c *ApplicationConfigurationTemplateClient) Create(container *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error) {
	resp := &ApplicationConfigurationTemplate{}
	err := c.apiClient.Ops.DoCreate(ApplicationConfigurationTemplateType, container, resp)
	return resp, err
}

func (c *ApplicationConfigurationTemplateClient) Update(existing *ApplicationConfigurationTemplate, updates interface{}) (*ApplicationConfigurationTemplate, error) {
	resp := &ApplicationConfigurationTemplate{}
	err := c.apiClient.Ops.DoUpdate(ApplicationConfigurationTemplateType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ApplicationConfigurationTemplateClient) Replace(obj *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error) {
	resp := &ApplicationConfigurationTemplate{}
	err := c.apiClient.Ops.DoReplace(ApplicationConfigurationTemplateType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ApplicationConfigurationTemplateClient) List(opts *types.ListOpts) (*ApplicationConfigurationTemplateCollection, error) {
	resp := &ApplicationConfigurationTemplateCollection{}
	err := c.apiClient.Ops.DoList(ApplicationConfigurationTemplateType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ApplicationConfigurationTemplateCollection) Next() (*ApplicationConfigurationTemplateCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ApplicationConfigurationTemplateCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ApplicationConfigurationTemplateClient) ByID(id string) (*ApplicationConfigurationTemplate, error) {
	resp := &ApplicationConfigurationTemplate{}
	err := c.apiClient.Ops.DoByID(ApplicationConfigurationTemplateType, id, resp)
	return resp, err
}

func (c *ApplicationConfigurationTemplateClient) Delete(container *ApplicationConfigurationTemplate) error {
	return c.apiClient.Ops.DoResourceDelete(ApplicationConfigurationTemplateType, &container.Resource)
}
