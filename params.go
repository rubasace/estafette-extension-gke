package main

import (
	"fmt"
)

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	Credentials     string            `json:"credentials,omitempty"`
	App             string            `json:"app,omitempty"`
	AppContainerTag string            `json:"tag,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`

	// used for seeing the rendered template without executing it but testing it with a dryrun
	DryRun bool `json:"dryrun,omitempty"`

	// AppLabelSelector    string
	// Hosts               []string
	// HostsJoined         string
	// IngressPath         string
	// UseNginxIngress     bool
	// UseGCEIngress       bool
	// ServiceType         string
	// MinReplicas         int
	// MaxReplicas         int
	// TargetCPUPercentage int
	// PreferPreemptibles  bool
	// Container           ContainerData
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults(appLabel, buildVersion, releaseName string, estafetteLabels map[string]string) {

	// default app to estafette app label if no override in stage params
	if p.App == "" && appLabel != "" {
		p.App = appLabel
	}

	// default container tag to estafette build version if no override in stage params
	if p.AppContainerTag == "" && buildVersion != "" {
		p.AppContainerTag = buildVersion
	}

	// default credentials to release name if no override in stage params
	if p.Credentials == "" && releaseName != "" {
		p.Credentials = fmt.Sprintf("gke-%v", releaseName)
	}

	// default labels to estafette labels if no override in stage params
	if p.Labels == nil {
		p.Labels = map[string]string{}
	}
	if len(p.Labels) == 0 && estafetteLabels != nil && len(estafetteLabels) != 0 {
		p.Labels = estafetteLabels
	}
	// ensure the app label is set and equals the app label or app override in stage params if present
	if p.App != "" {
		p.Labels["app"] = p.App
	}
}

// SetDefaultNamespace sets default namespace separately from other defaults, because the credential fetched with params has the default value
func (p *Params) SetDefaultNamespace(defaultNamespace string) {

	if p.Namespace == "" && defaultNamespace != "" {
		p.Namespace = defaultNamespace
	}

}

// ValidateRequiredProperties checks whether all needed properties are set
func (p *Params) ValidateRequiredProperties() (bool, []error) {

	errors := []error{}

	if p.App == "" {
		errors = append(errors, fmt.Errorf("Application name is required; either define an app label or use app property on this stage"))
	}
	if p.Namespace == "" {
		errors = append(errors, fmt.Errorf("Namespace is required; either use credentials with a defaultNamespace or set it via namespace property on this stage"))
	}
	if p.AppContainerTag == "" {
		errors = append(errors, fmt.Errorf("App container tag is required; set it via tag property on this stage"))
	}
	if p.Credentials == "" {
		errors = append(errors, fmt.Errorf("Credentials property is required; set it via credentials property on this stage"))
	}

	return len(errors) == 0, errors
}