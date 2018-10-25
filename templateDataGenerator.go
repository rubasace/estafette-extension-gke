package main

func generateTemplateData(params Params) TemplateData {

	data := TemplateData{
		Name:             params.App,
		Namespace:        params.Namespace,
		Labels:           params.Labels,
		AppLabelSelector: params.App,

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

		Container: ContainerData{
			Repository: params.Image.ImageRepository,
			Name:       params.Image.ImageName,
			Tag:        params.Image.ImageTag,

			CPURequest:    params.CPU.Request,
			CPULimit:      params.CPU.Limit,
			MemoryRequest: params.Memory.Request,
			MemoryLimit:   params.Memory.Limit,

			Liveness: ProbeData{
				// Path                string
				// InitialDelaySeconds int
				// TimeoutSeconds      int
			},
			Readiness: ProbeData{
				// Path                string
				// InitialDelaySeconds int
				// TimeoutSeconds      int
			},
		},
	}

	if params.Visibility == "private" {
		data.ServiceType = "ClusterIP"
	}

	return data
}
