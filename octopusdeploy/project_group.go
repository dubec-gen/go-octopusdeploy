package octopusdeploy

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"gopkg.in/go-playground/validator.v9"
)

type ProjectGroupService struct {
	sling *sling.Sling
}

func NewProjectGroupService(sling *sling.Sling) *ProjectGroupService {
	return &ProjectGroupService{
		sling: sling,
	}
}

type ProjectGroups struct {
	Items []ProjectGroup `json:"Items"`
	PagedResults
}

type ProjectGroup struct {
	Description       string   `json:"Description,omitempty"`
	EnvironmentIds    []string `json:"EnvironmentIds"`
	ID                string   `json:"Id,omitempty"`
	LastModifiedBy    string   `json:"LastModifiedBy,omitempty"`
	LastModifiedOn    string   `json:"LastModifiedOn,omitempty"`
	Links             Links    `json:"Links,omitempty"`
	Name              string   `json:"Name,omitempty" validate:"required"`
	RetentionPolicyID string   `json:"RetentionPolicyId,omitempty"`
}

func (p *ProjectGroup) Validate() error {
	validate := validator.New()

	err := validate.Struct(p)

	if err != nil {
		return err
	}

	return nil
}

func NewProjectGroup(name string) *ProjectGroup {
	return &ProjectGroup{
		Name: name,
	}
}

func (s *ProjectGroupService) Get(projectGroupID string) (*ProjectGroup, error) {
	path := fmt.Sprintf("projectgroups/%s", projectGroupID)
	resp, err := apiGet(s.sling, new(ProjectGroup), path)

	if err != nil {
		return nil, err
	}

	return resp.(*ProjectGroup), nil
}

func (s *ProjectGroupService) GetAll() (*[]ProjectGroup, error) {
	var listOfProjectGroups []ProjectGroup
	path := fmt.Sprintf("projectgroups")

	for {
		var projectGroups ProjectGroups
		octopusDeployError := new(APIError)

		resp, err := s.sling.New().Get(path).Receive(&projectGroups, &octopusDeployError)

		apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)

		if apiErrorCheck != nil {
			return nil, apiErrorCheck
		}

		for _, projectGroup := range projectGroups.Items {
			listOfProjectGroups = append(listOfProjectGroups, projectGroup)
		}

		if projectGroups.PagedResults.Links.PageNext != "" {
			path = projectGroups.PagedResults.Links.PageNext
		} else {
			break
		}
	}

	return &listOfProjectGroups, nil // no more pages to go through
}

func (s *ProjectGroupService) Add(projectGroup *ProjectGroup) (*ProjectGroup, error) {
	resp, err := apiAdd(s.sling, projectGroup, new(ProjectGroup), "projectgroups")

	if err != nil {
		return nil, err
	}

	return resp.(*ProjectGroup), nil
}

func (s *ProjectGroupService) Delete(projectGroupID string) error {
	path := fmt.Sprintf("projectgroups/%s", projectGroupID)
	err := apiDelete(s.sling, path)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectGroupService) Update(projectGroup *ProjectGroup) (*ProjectGroup, error) {
	err := projectGroup.Validate()

	if err != nil {
		return nil, err
	}

	var updated ProjectGroup
	octopusDeployError := new(APIError)
	path := fmt.Sprintf("projectgroups/%s", projectGroup.ID)

	resp, err := s.sling.New().Put(path).BodyJSON(projectGroup).Receive(&updated, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)

	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return &updated, nil
}
