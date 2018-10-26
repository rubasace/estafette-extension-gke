package main

import (
	"fmt"
)

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// control params
	Credentials string `json:"credentials,omitempty"`
	DryRun      bool   `json:"dryrun,string,omitempty"`

	// app params
	App        string            `json:"app,omitempty"`
	Namespace  string            `json:"namespace,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	Visibility string            `json:"visibility,omitempty"`
	Hosts      []string          `json:"hosts,omitempty"`
	Autoscale  AutoscaleParams   `json:"autoscale,omitempty"`

	// container params
	Container ContainerParams `json:"container,omitempty"`
}

// ContainerParams defines the container image to deploy
type ContainerParams struct {
	ImageRepository      string            `json:"repository,omitempty"`
	ImageName            string            `json:"name,omitempty"`
	ImageTag             string            `json:"tag,omitempty"`
	Port                 int               `json:"port,string,omitempty"`
	EnvironmentVariables map[string]string `json:"env,omitempty"`

	CPU            CPUParams     `json:"cpu,omitempty"`
	Memory         MemoryParams  `json:"memory,omitempty"`
	LivenessProbe  ProbeParams   `json:"liveness,omitempty"`
	ReadinessProbe ProbeParams   `json:"readiness,omitempty"`
	Metrics        MetricsParams `json:"metrics,omitempty"`
}

// CPUParams sets cpu request and limit values
type CPUParams struct {
	Request string `json:"request,omitempty"`
	Limit   string `json:"limit,omitempty"`
}

// MemoryParams sets memory request and limit values
type MemoryParams struct {
	Request string `json:"request,omitempty"`
	Limit   string `json:"limit,omitempty"`
}

// AutoscaleParams controls autoscaling
type AutoscaleParams struct {
	MinReplicas   int `json:"min,string,omitempty"`
	MaxReplicas   int `json:"max,string,omitempty"`
	CPUPercentage int `json:"cpu,string,omitempty"`
}

// ProbeParams sets params for liveness or readiness probe
type ProbeParams struct {
	Path                string `json:"path,omitempty"`
	InitialDelaySeconds int    `json:"delay,string,omitempty"`
	TimeoutSeconds      int    `json:"timeout,string,omitempty"`
}

// MetricsParams sets params for scraping prometheus metrics
type MetricsParams struct {
	Path string `json:"path,omitempty"`
	Port int    `json:"port,string,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults(appLabel, buildVersion, releaseName string, estafetteLabels map[string]string) {

	// default app to estafette app label if no override in stage params
	if p.App == "" && appLabel != "" {
		p.App = appLabel
	}

	// default image name to estafette app label if no override in stage params
	if p.Container.ImageName == "" && p.App != "" {
		p.Container.ImageName = p.App
	}

	// default image tag to estafette build version if no override in stage params
	if p.Container.ImageTag == "" && buildVersion != "" {
		p.Container.ImageTag = buildVersion
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

	// default visibility to private if no override in stage params
	if p.Visibility == "" {
		p.Visibility = "private"
	}

	// set cpu defaults
	cpuRequestIsEmpty := p.Container.CPU.Request == ""
	if cpuRequestIsEmpty {
		if p.Container.CPU.Limit != "" {
			p.Container.CPU.Request = p.Container.CPU.Limit
		} else {
			p.Container.CPU.Request = "100m"
		}
	}
	if p.Container.CPU.Limit == "" {
		if !cpuRequestIsEmpty {
			p.Container.CPU.Limit = p.Container.CPU.Request
		} else {
			p.Container.CPU.Limit = "125m"
		}
	}

	// set memory defaults
	memoryRequestIsEmpty := p.Container.Memory.Request == ""
	if memoryRequestIsEmpty {
		if p.Container.Memory.Limit != "" {
			p.Container.Memory.Request = p.Container.Memory.Limit
		} else {
			p.Container.Memory.Request = "128Mi"
		}
	}
	if p.Container.Memory.Limit == "" {
		if !memoryRequestIsEmpty {
			p.Container.Memory.Limit = p.Container.Memory.Request
		} else {
			p.Container.Memory.Limit = "128Mi"
		}
	}

	// set container port defaults
	if p.Container.Port <= 0 {
		p.Container.Port = 5000
	}

	// set autoscale defaults
	if p.Autoscale.MinReplicas <= 0 {
		p.Autoscale.MinReplicas = 3
	}
	if p.Autoscale.MaxReplicas <= 0 {
		p.Autoscale.MaxReplicas = 100
	}
	if p.Autoscale.CPUPercentage <= 0 {
		p.Autoscale.CPUPercentage = 80
	}

	// set liveness probe defaults
	if p.Container.LivenessProbe.Path == "" {
		p.Container.LivenessProbe.Path = "/liveness"
	}
	if p.Container.LivenessProbe.InitialDelaySeconds <= 0 {
		p.Container.LivenessProbe.InitialDelaySeconds = 30
	}
	if p.Container.LivenessProbe.TimeoutSeconds <= 0 {
		p.Container.LivenessProbe.TimeoutSeconds = 1
	}

	// set readiness probe defaults
	if p.Container.ReadinessProbe.Path == "" {
		p.Container.ReadinessProbe.Path = "/readiness"
	}
	if p.Container.ReadinessProbe.TimeoutSeconds <= 0 {
		p.Container.ReadinessProbe.TimeoutSeconds = 1
	}

	// set metrics defaults
	if p.Container.Metrics.Path == "" {
		p.Container.Metrics.Path = "/metrics"
	}
	if p.Container.Metrics.Port <= 0 {
		p.Container.Metrics.Port = p.Container.Port
	}
}

// SetDefaultsFromCredentials sets defaults based on the credentials fetched with first-run defaults
func (p *Params) SetDefaultsFromCredentials(credentials GKECredentials) {

	// default namespace to credential default namespace if no override in stage params
	if p.Namespace == "" && credentials.AdditionalProperties.DefaultNamespace != "" {
		p.Namespace = credentials.AdditionalProperties.DefaultNamespace
	}

	// default image repository to credential project if no override in stage params
	if p.Container.ImageRepository == "" && credentials.AdditionalProperties.Project != "" {
		p.Container.ImageRepository = credentials.AdditionalProperties.Project
	}
}

// ValidateRequiredProperties checks whether all needed properties are set
func (p *Params) ValidateRequiredProperties() (bool, []error) {

	errors := []error{}

	// validate control params
	if p.Credentials == "" {
		errors = append(errors, fmt.Errorf("Credentials property is required; set it via credentials property on this stage"))
	}

	// validate app params
	if p.App == "" {
		errors = append(errors, fmt.Errorf("Application name is required; either define an app label or use app property on this stage"))
	}
	if p.Namespace == "" {
		errors = append(errors, fmt.Errorf("Namespace is required; either use credentials with a defaultNamespace or set it via namespace property on this stage"))
	}
	if p.Visibility == "" || (p.Visibility != "private" && p.Visibility != "public") {
		errors = append(errors, fmt.Errorf("Visibility property is required; set it via visibility property on this stage; allowed values are private or public"))
	}
	if len(p.Hosts) == 0 {
		errors = append(errors, fmt.Errorf("At least one host is required; set it via hosts array property on this stage"))
	}

	// validate autoscale params
	if p.Autoscale.MinReplicas <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling min replicas must be larger than zero; set it via autoscale.min property on this stage"))
	}
	if p.Autoscale.MaxReplicas <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling max replicas must be larger than zero; set it via autoscale.max property on this stage"))
	}
	if p.Autoscale.CPUPercentage <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling cpu percentage must be larger than zero; set it via autoscale.cpu property on this stage"))
	}

	// validate container params
	if p.Container.ImageRepository == "" {
		errors = append(errors, fmt.Errorf("Image repository is required; set it via container.repository property on this stage"))
	}
	if p.Container.ImageName == "" {
		errors = append(errors, fmt.Errorf("Image name is required; set it via container.name property on this stage"))
	}
	if p.Container.ImageTag == "" {
		errors = append(errors, fmt.Errorf("Image tag is required; set it via container.tag property on this stage"))
	}
	if p.Container.Port <= 0 {
		errors = append(errors, fmt.Errorf("Container port must be larger than zero; set it via container.port property on this stage"))
	}

	// validate cpu params
	if p.Container.CPU.Request == "" {
		errors = append(errors, fmt.Errorf("Cpu request is required; set it via container.cpu.request property on this stage"))
	}
	if p.Container.CPU.Limit == "" {
		errors = append(errors, fmt.Errorf("Cpu limit is required; set it via container.cpu.limit property on this stage"))
	}

	// validate memory params
	if p.Container.Memory.Request == "" {
		errors = append(errors, fmt.Errorf("Memory request is required; set it via container.memory.request property on this stage"))
	}
	if p.Container.Memory.Limit == "" {
		errors = append(errors, fmt.Errorf("Memory limit is required; set it via container.memory.limit property on this stage"))
	}

	// validate liveness params
	if p.Container.LivenessProbe.Path == "" {
		errors = append(errors, fmt.Errorf("Liveness path is required; set it via container.liveness.path property on this stage"))
	}
	if p.Container.LivenessProbe.InitialDelaySeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness initial delay must be larger than zero; set it via container.liveness.delay property on this stage"))
	}
	if p.Container.LivenessProbe.TimeoutSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness timeout must be larger than zero; set it via container.liveness.timeout property on this stage"))
	}

	// validate readiness params
	if p.Container.ReadinessProbe.Path == "" {
		errors = append(errors, fmt.Errorf("Readiness path is required; set it via container.readiness.path property on this stage"))
	}
	if p.Container.ReadinessProbe.TimeoutSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Readiness timeout must be larger than zero; set it via container.readiness.timeout property on this stage"))
	}

	// validate metrics params
	if p.Container.Metrics.Path == "" {
		errors = append(errors, fmt.Errorf("Metrics path is required; set it via container.metrics.path property on this stage"))
	}
	if p.Container.Metrics.Port <= 0 {
		errors = append(errors, fmt.Errorf("Metrics port must be larger than zero; set it via container.metrics.port property on this stage"))
	}

	return len(errors) == 0, errors
}
