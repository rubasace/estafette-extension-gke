package api

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	trueValue   = true
	falseValue  = false
	validParams = Params{
		Action:    ActionDeploySimple,
		App:       "myapp",
		Namespace: "mynamespace",
		Autoscale: AutoscaleParams{
			Enabled:       &trueValue,
			MinReplicas:   3,
			MaxReplicas:   100,
			CPUPercentage: 80,
		},
		StrategyType: StrategyTypeRollingUpdate,
		RollingUpdate: RollingUpdateParams{
			MaxSurge:       "25%",
			MaxUnavailable: "25%",
		},
		Container: ContainerParams{
			ImageRepository: "estafette",
			ImageName:       "my-app",
			ImageTag:        "1.0.0",
			ImagePullPolicy: "IfNotPresent",
			Port:            5000,

			CPU: CPUParams{
				Request: "100m",
				Limit:   "150m",
			},
			Memory: MemoryParams{
				Request: "768Mi",
				Limit:   "1024Mi",
			},
			LivenessProbe: ProbeParams{
				Path:                "/liveness",
				Port:                5000,
				InitialDelaySeconds: 30,
				TimeoutSeconds:      1,
				PeriodSeconds:       10,
				FailureThreshold:    3,
				SuccessThreshold:    1,
			},
			ReadinessProbe: ProbeParams{
				Path:                "/readiness",
				Port:                5000,
				InitialDelaySeconds: 0,
				TimeoutSeconds:      1,
				PeriodSeconds:       10,
				FailureThreshold:    3,
				SuccessThreshold:    1,
			},
			Metrics: MetricsParams{
				Scrape: &trueValue,
				Path:   "/metrics",
				Port:   5000,
			},
		},
		Visibility: VisibilityPrivate,
		Hosts:      []string{"gke.estafette.io"},
		Basepath:   "/",
		Sidecar: SidecarParams{
			Type:  SidecarTypeOpenresty,
			Image: "estafette/openresty-sidecar:1.13.6.2-alpine",
			CPU: CPUParams{
				Request: "10m",
				Limit:   "50m",
			},
			Memory: MemoryParams{
				Request: "10Mi",
				Limit:   "50Mi",
			},
		},
		Sidecars: []*SidecarParams{
			&SidecarParams{
				Type:  SidecarTypeOpenresty,
				Image: "estafette/openresty-sidecar:1.13.6.2-alpine",
				CPU: CPUParams{
					Request: "10m",
					Limit:   "50m",
				},
				Memory: MemoryParams{
					Request: "10Mi",
					Limit:   "50Mi",
				},
			},
			&SidecarParams{
				Type:  SidecarTypeESP,
				Image: "estafette/estafette-docker-cache-heater:dev",
				CPU: CPUParams{
					Request: "10m",
					Limit:   "50m",
				},
				Memory: MemoryParams{
					Request: "10Mi",
					Limit:   "50Mi",
				},
			},
		},
		StorageClass:        "standard",
		StorageSize:         "1Gi",
		StorageMountPath:    "/data",
		PodManagementPolicy: "Parallel",
		ProbeService:        &trueValue,
	}
	validCredential = GKECredentials{
		Name: "gke-production",
	}
)

func stringInErrorSlice(a string, list []error) string {
	for _, b := range list {
		if b.Error() == a {
			return a
		}
	}
	return ""
}

func TestSetDefaults(t *testing.T) {

	t.Run("DefaultsAppToGitNameIfAppParamIsEmptyAndAppLabelIsEmpty", func(t *testing.T) {

		params := Params{
			App: "",
		}
		gitName := "mygitrepo"
		appLabel := ""

		// act
		params.SetDefaults("", "", gitName, appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "mygitrepo", params.App)
	})

	t.Run("DefaultsAppToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			App: "",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "myapp", params.App)
	})

	t.Run("KeepsAppIfNotEmpty", func(t *testing.T) {

		params := Params{
			App: "yourapp",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "yourapp", params.App)
	})

	t.Run("DefaultsDisableServiceAccountKeyRotationToTrueIfEmpty", func(t *testing.T) {

		params := Params{
			DisableServiceAccountKeyRotation: nil,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.DisableServiceAccountKeyRotation)
	})

	t.Run("KeepsDisableServiceAccountKeyRotationIfNotEmpty", func(t *testing.T) {

		params := Params{
			DisableServiceAccountKeyRotation: &falseValue,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.DisableServiceAccountKeyRotation)
	})

	t.Run("DefaultsGoogleCloudCredentialsAppToAppIfEmpty", func(t *testing.T) {

		params := Params{
			App:                       "yourapp",
			GoogleCloudCredentialsApp: "",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "yourapp", params.GoogleCloudCredentialsApp)
	})

	t.Run("KeepsGoogleCloudCredentialsAppIfNotEmpty", func(t *testing.T) {

		params := Params{
			App:                       "yourapp",
			GoogleCloudCredentialsApp: "myapp",
		}
		appLabel := "someapp"

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "myapp", params.GoogleCloudCredentialsApp)
	})

	t.Run("DefaultsImageNameToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageName: "",
			},
		}
		appLabel := "myapp"

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "myapp", params.Container.ImageName)
	})

	t.Run("KeepsImageTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageName: "my-app",
			},
		}
		appLabel := "myapp"

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", map[string]string{})

		assert.Equal(t, "my-app", params.Container.ImageName)
	})

	t.Run("DefaultsImageTagToBuildVersionIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageTag: "",
			},
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", "", "", "", buildVersion, "", "", "", map[string]string{})

		assert.Equal(t, "1.0.0", params.Container.ImageTag)
	})

	t.Run("KeepsImagePullPolicyIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImagePullPolicy: "Always",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "Always", params.Container.ImagePullPolicy)
	})

	t.Run("DefaultsImagePullPolicyToIfNotPresentIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImagePullPolicy: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "IfNotPresent", params.Container.ImagePullPolicy)
	})

	t.Run("KeepsImageTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageTag: "2.1.3",
			},
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", "", "", "", buildVersion, "", "", "", map[string]string{})

		assert.Equal(t, "2.1.3", params.Container.ImageTag)
	})

	t.Run("DefaultsLabelsToEstafetteLabelsIfEmpty", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", estafetteLabels)

		assert.Equal(t, 3, len(params.Labels))
		assert.Equal(t, "myapp", params.Labels["app"])
		assert.Equal(t, "myteam", params.Labels["team"])
		assert.Equal(t, "golang", params.Labels["language"])
	})

	t.Run("KeepsLabelsIfNotEmpty", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"app":  "yourapp",
				"team": "yourteam",
			},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", estafetteLabels)

		assert.Equal(t, 2, len(params.Labels))
		assert.Equal(t, "yourapp", params.Labels["app"])
		assert.Equal(t, "yourteam", params.Labels["team"])
	})

	t.Run("AddsAppLabelToLabelsIfNotSet", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"team": "yourteam",
			},
		}
		appLabel := "myapp"
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", estafetteLabels)

		assert.Equal(t, 2, len(params.Labels))
		assert.Equal(t, "myapp", params.Labels["app"])
		assert.Equal(t, "yourteam", params.Labels["team"])
	})

	t.Run("AddsEstafettePipelineLabel", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("github.com", "estafette", "estafette-extensions-gke", "", "", "", "", "", estafetteLabels)

		assert.Equal(t, "github.com-estafette-estafette-extensions-gke", params.Labels["estafette.io/pipeline"])
	})

	t.Run("AddsEstafettePipelineBase64Label", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("github.com", "estafette", "estafette-extensions-gke", "", "", "", "", "", estafetteLabels)

		assert.Equal(t, "Z2l0aHViLmNvbS9lc3RhZmV0dGUvZXN0YWZldHRlLWV4dGVuc2lvbnMtZ2tl", params.Labels["estafette.io/pipeline-base64"])

		decodedLabel, err := base64.StdEncoding.DecodeString(params.Labels["estafette.io/pipeline-base64"])
		assert.Nil(t, err)
		assert.Equal(t, "github.com/estafette/estafette-extensions-gke", string(decodedLabel))
	})

	t.Run("OverwritesAppLabelToAppIfSetFromEstafetteLabels", func(t *testing.T) {

		params := Params{}
		appLabel := "yourapp"
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", appLabel, "", "", "", "", estafetteLabels)

		assert.Equal(t, 3, len(params.Labels))
		assert.Equal(t, "yourapp", params.Labels["app"])
		assert.Equal(t, "myteam", params.Labels["team"])
		assert.Equal(t, "golang", params.Labels["language"])
	})

	t.Run("DefaultsVisibilityToPrivateIfEmpty", func(t *testing.T) {

		params := Params{
			Visibility: VisibilityUnknown,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, VisibilityPrivate, params.Visibility)
	})

	t.Run("KeepsVisibilityIfNotEmpty", func(t *testing.T) {

		params := Params{
			Visibility: VisibilityPublic,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, VisibilityPublic, params.Visibility)
	})

	t.Run("DefaultsCpuRequestTo100MIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "100m", params.Container.CPU.Request)
	})

	t.Run("DefaultsCpuRequestToLimitIfRequestIsEmptyButLimitIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "300m",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "300m", params.Container.CPU.Request)
	})

	t.Run("KeepsCpuRequestIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "250m",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "250m", params.Container.CPU.Request)
	})

	t.Run("KeepsCpuLimitEmptyIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "", params.Container.CPU.Limit)
	})

	t.Run("KeepsCpuLimitEmptyIfLimitIsEmptyButRequestIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "300m",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "", params.Container.CPU.Limit)
	})

	t.Run("KeepsCpuLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "250m",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "250m", params.Container.CPU.Limit)
	})

	t.Run("DefaultsMemoryRequestTo128MiIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "128Mi", params.Container.Memory.Request)
	})

	t.Run("DefaultsMemoryRequestToLimitIfRequestIsEmptyButLimitIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "256Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "256Mi", params.Container.Memory.Request)
	})

	t.Run("KeepsMemoryRequestIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "512Mi",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "512Mi", params.Container.Memory.Request)
	})

	t.Run("DefaultsMemoryLimitTo128MiIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "128Mi", params.Container.Memory.Limit)
	})

	t.Run("DefaultsMemoryLimitToRequestIfLimitIsEmptyButRequestIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "768Mi",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "768Mi", params.Container.Memory.Limit)
	})

	t.Run("KeepsMemoryLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "1024Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "1024Mi", params.Container.Memory.Limit)
	})

	t.Run("DefaultsContainerPortTo5000IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5000, params.Container.Port)
	})

	t.Run("KeepsContainerPortIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 3000,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3000, params.Container.Port)
	})

	t.Run("DefaultsAdditionalPortProtocolToTCPIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				AdditionalPorts: []*AdditionalPortParams{
					&AdditionalPortParams{
						Protocol: "",
					},
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "TCP", params.Container.AdditionalPorts[0].Protocol)
	})

	t.Run("KeepsAdditionalPortProtocolIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				AdditionalPorts: []*AdditionalPortParams{
					{
						Protocol: "UDP",
					},
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "UDP", params.Container.AdditionalPorts[0].Protocol)
	})

	t.Run("DefaultsAdditionalPortVisibilityToApplicationVisibilityIfEmpty", func(t *testing.T) {

		params := Params{
			Visibility: VisibilityPublic,
			Container: ContainerParams{
				AdditionalPorts: []*AdditionalPortParams{
					{
						Visibility: VisibilityUnknown,
					},
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, VisibilityPublic, params.Container.AdditionalPorts[0].Visibility)
	})

	t.Run("KeepsAdditionalPortVisibilityIfNotEmpty", func(t *testing.T) {

		params := Params{
			Visibility: VisibilityPublic,
			Container: ContainerParams{
				AdditionalPorts: []*AdditionalPortParams{
					{
						Visibility: VisibilityPrivate,
					},
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, VisibilityPrivate, params.Container.AdditionalPorts[0].Visibility)
	})

	t.Run("DefaultsAutoscaleMinReplicasTo3IfZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, params.Autoscale.MinReplicas)
	})

	t.Run("KeepsAutoscaleMinReplicasIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 2,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 2, params.Autoscale.MinReplicas)
	})

	t.Run("DefaultsAutoscaleMaxReplicasTo100IfZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MaxReplicas: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 100, params.Autoscale.MaxReplicas)
	})

	t.Run("KeepsAutoscaleMaxReplicasIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MaxReplicas: 50,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 50, params.Autoscale.MaxReplicas)
	})

	t.Run("DefaultsAutoscaleCPUPercentageTo80IfZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				CPUPercentage: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 80, params.Autoscale.CPUPercentage)
	})

	t.Run("KeepsAutoscaleCPUPercentageIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				CPUPercentage: 30,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 30, params.Autoscale.CPUPercentage)
	})

	t.Run("DefaultsAutoscaleSafetyEnabledToFalse", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.False(t, params.Autoscale.Safety.Enabled)
	})

	t.Run("KeepsAutoscaleSafetyEnabled", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					Enabled: true,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.True(t, params.Autoscale.Safety.Enabled)
	})

	t.Run("DefaultsAutoscaleSafetyPromQueryToRequestRateForAppLabelOverLast5Minutes", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "sum(rate(nginx_http_requests_total{app='my-app'}[5m])) by (app)", params.Autoscale.Safety.PromQuery)
	})

	t.Run("KeepsAutoscaleSafetyPromQuery", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					PromQuery: "sum(rate(nginx_http_requests_total{app='your-app'}[5m])) by (app)",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "sum(rate(nginx_http_requests_total{app='your-app'}[5m])) by (app)", params.Autoscale.Safety.PromQuery)
	})

	t.Run("DefaultsAutoscaleSafetyRatioToOne", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "1", params.Autoscale.Safety.Ratio)
	})

	t.Run("KeepsAutoscaleSafetyRatio", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					Ratio: "1.5",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "1.5", params.Autoscale.Safety.Ratio)
	})

	t.Run("DefaultsAutoscaleSafetyScaleDownRatioToOne", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "1", params.Autoscale.Safety.ScaleDownRatio)
	})

	t.Run("KeepsAutoscaleSafetyScaleDownRatio", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					ScaleDownRatio: "0.2",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "0.2", params.Autoscale.Safety.ScaleDownRatio)
	})

	t.Run("DefaultsVerticalPodAutoscalerEnabledToFalse", func(t *testing.T) {

		params := Params{
			VerticalPodAutoscaler: VPAParams{
				Enabled: nil,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.VerticalPodAutoscaler.Enabled)
	})

	t.Run("KeepsVerticalPodAutoscalerEnabled", func(t *testing.T) {

		trueValue := true
		params := Params{
			VerticalPodAutoscaler: VPAParams{
				Enabled: &trueValue,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.VerticalPodAutoscaler.Enabled)
	})

	t.Run("DefaultsVerticalPodAutoscalerUpdateModeToOff", func(t *testing.T) {

		params := Params{
			VerticalPodAutoscaler: VPAParams{
				UpdateMode: UpdateModeUnknown,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, UpdateModeOff, params.VerticalPodAutoscaler.UpdateMode)
	})

	t.Run("KeepsVerticalPodAutoscalerUpdateMode", func(t *testing.T) {

		params := Params{
			VerticalPodAutoscaler: VPAParams{
				UpdateMode: UpdateModeAuto,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, UpdateModeAuto, params.VerticalPodAutoscaler.UpdateMode)
	})

	t.Run("DefaultsRequestTimeoutTo60sIfEmpty", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "60s", params.Request.Timeout)
	})

	t.Run("KeepsRequestTimeoutIfNotEmpty", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "10s",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "10s", params.Request.Timeout)
	})

	t.Run("DefaultsRequestMaxBodySizeTo128mIfEmpty", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				MaxBodySize: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "128m", params.Request.MaxBodySize)
	})

	t.Run("KeepsRequestMaxBodySizeIfNotEmpty", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				MaxBodySize: "16m",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "16m", params.Request.MaxBodySize)
	})

	t.Run("DefaultsRequestProxyBufferSizeTo4kIfEmpty", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				ProxyBufferSize: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "4k", params.Request.ProxyBufferSize)
	})

	t.Run("KeepsRequestProxyBufferSizeIfNotEmpty", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				ProxyBufferSize: "8k",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "8k", params.Request.ProxyBufferSize)
	})

	t.Run("DefaultsLivenessEnabledToTrue", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Enabled: nil,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.Container.LivenessProbe.Enabled)
	})

	t.Run("KeepsLivenessWhenSet", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Enabled: &falseValue,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.Container.LivenessProbe.Enabled)
	})

	t.Run("DefaultsLivenessInitialDelaySecondsTo30IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					InitialDelaySeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 30, params.Container.LivenessProbe.InitialDelaySeconds)
	})

	t.Run("KeepsLivenessInitialDelaySecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					InitialDelaySeconds: 120,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 120, params.Container.LivenessProbe.InitialDelaySeconds)
	})

	t.Run("DefaultsLivenessTimeoutSecondsTo1IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					TimeoutSeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, params.Container.LivenessProbe.TimeoutSeconds)
	})

	t.Run("KeepsLivenessTimeoutSecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					TimeoutSeconds: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.LivenessProbe.TimeoutSeconds)
	})

	t.Run("DefaultsLivenessFailureThresholdTo3IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					FailureThreshold: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, params.Container.LivenessProbe.FailureThreshold)
	})

	t.Run("KeepsLivenessFailureThresholdIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					FailureThreshold: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.LivenessProbe.FailureThreshold)
	})

	t.Run("DefaultsLivenessSuccessThresholdTo1IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					SuccessThreshold: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, params.Container.LivenessProbe.SuccessThreshold)
	})

	t.Run("KeepsLivenessSuccessThresholdIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					SuccessThreshold: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.LivenessProbe.SuccessThreshold)
	})

	t.Run("DefaultsLivenessPathToLivenessIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Path: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/liveness", params.Container.LivenessProbe.Path)
	})

	t.Run("KeepsLivenessPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Path: "/healthz",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/healthz", params.Container.LivenessProbe.Path)
	})

	t.Run("DefaultsLivenessProbePortToContainerPortIfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 8080,
				LivenessProbe: ProbeParams{
					Port: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 8080, params.Container.LivenessProbe.Port)
	})

	t.Run("KeepsLivenessProbePortIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 8080,
				LivenessProbe: ProbeParams{
					Port: 8081,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 8081, params.Container.LivenessProbe.Port)
	})

	t.Run("DefaultsReadinessEnabledToTrue", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Enabled: nil,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.Container.ReadinessProbe.Enabled)
	})

	t.Run("DefaultsReadinessEnabledToFalseIfKindIsHeadlessDeployment", func(t *testing.T) {

		params := Params{
			Kind: KindHeadlessDeployment,
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Enabled: nil,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.Container.ReadinessProbe.Enabled)
	})

	t.Run("KeepsReadinessWhenSet", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Enabled: &falseValue,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.Container.ReadinessProbe.Enabled)
	})

	t.Run("DefaultsReadinessInitialDelaySecondsTo0IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					InitialDelaySeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 0, params.Container.ReadinessProbe.InitialDelaySeconds)
	})

	t.Run("KeepsReadinessInitialDelaySecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					InitialDelaySeconds: 120,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 120, params.Container.ReadinessProbe.InitialDelaySeconds)
	})

	t.Run("DefaultsReadinessTimeoutSecondsTo1IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					TimeoutSeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, params.Container.ReadinessProbe.TimeoutSeconds)
	})

	t.Run("KeepsReadinessTimeoutSecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					TimeoutSeconds: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.ReadinessProbe.TimeoutSeconds)
	})

	t.Run("DefaultsReadinessFailureThresholdTo3IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					FailureThreshold: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, params.Container.ReadinessProbe.FailureThreshold)
	})

	t.Run("KeepsReadinessFailureThresholdIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					FailureThreshold: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.ReadinessProbe.FailureThreshold)
	})

	t.Run("DefaultsReadinessSuccessThresholdTo1IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					SuccessThreshold: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, params.Container.ReadinessProbe.SuccessThreshold)
	})

	t.Run("KeepsReadinessSuccessThresholdIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					SuccessThreshold: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.ReadinessProbe.SuccessThreshold)
	})

	t.Run("DefaultsReadinessPathToReadinessIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/readiness", params.Container.ReadinessProbe.Path)
	})

	t.Run("KeepsReadinessPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "/healthz",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/healthz", params.Container.ReadinessProbe.Path)
	})

	t.Run("DefaultsReadinessProbePortToContainerPortIfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 8080,
				ReadinessProbe: ProbeParams{
					Port: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 8080, params.Container.ReadinessProbe.Port)
	})

	t.Run("KeepsReadinessProbePortIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 8080,
				ReadinessProbe: ProbeParams{
					Port: 8082,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 8082, params.Container.ReadinessProbe.Port)
	})

	t.Run("DefaultsProbeServiceToTrueIfEmpty", func(t *testing.T) {

		params := Params{
			ProbeService: nil,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.ProbeService)
	})

	t.Run("KeepsProbeServiceIfNotEmpty", func(t *testing.T) {

		params := Params{
			ProbeService: &falseValue,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.ProbeService)
	})

	t.Run("DefaultsMetricsPathToMetricsIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Path: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/metrics", params.Container.Metrics.Path)
	})

	t.Run("KeepsMetricsPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Path: "/mymetrics",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/mymetrics", params.Container.Metrics.Path)
	})

	t.Run("DefaultsMetricsPortToContainerPortIfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 5000,
				Metrics: MetricsParams{
					Port: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5000, params.Container.Metrics.Port)
	})

	t.Run("KeepsMetricsPortIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 5000,
				Metrics: MetricsParams{
					Port: 5001,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 5001, params.Container.Metrics.Port)
	})

	t.Run("DefaultsMetricsScrapeToTrueIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Scrape: nil,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.Container.Metrics.Scrape)
	})

	t.Run("KeepsMetricsScrapeIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Scrape: &falseValue,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.Container.Metrics.Scrape)
	})

	t.Run("DefaultsLifecyclePrestopSleepToTrueIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Lifecycle: LifecycleParams{
					PrestopSleep: nil,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, true, *params.Container.Lifecycle.PrestopSleep)
	})

	t.Run("KeepsLifecyclePrestopSleepIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Lifecycle: LifecycleParams{
					PrestopSleep: &falseValue,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, false, *params.Container.Lifecycle.PrestopSleep)
	})

	t.Run("DefaultsLifecyclePrestopSleepSecondsTo20IfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Lifecycle: LifecycleParams{
					PrestopSleepSeconds: nil,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 20, *params.Container.Lifecycle.PrestopSleepSeconds)
	})

	t.Run("KeepsLifecyclePrestopSleepIfNotEmpty", func(t *testing.T) {

		nonDefaultValue := 25

		params := Params{
			Container: ContainerParams{
				Lifecycle: LifecycleParams{
					PrestopSleepSeconds: &nonDefaultValue,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 25, *params.Container.Lifecycle.PrestopSleepSeconds)
	})

	t.Run("AddsDefaultsOpenrestySidecarIfEmptyAndGlobalKindIsDeployment", func(t *testing.T) {

		params := Params{
			Kind: KindDeployment,
			Sidecar: SidecarParams{
				Type: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, SidecarTypeOpenresty, params.Sidecars[0].Type)
	})

	t.Run("DoesntAddDefaultSidecarIfEmptyAndGlobalKindIsHeadlessDeployment", func(t *testing.T) {

		params := Params{
			Kind: KindHeadlessDeployment,
			Sidecar: SidecarParams{
				Type: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 0, len(params.Sidecars))
	})

	t.Run("DoesntAddDefaultSidecarIfInjectFlagIsFalseEvenIfNoSidecarSpecified", func(t *testing.T) {

		falseValue := false
		params := Params{
			Kind:                   KindDeployment,
			InjectHTTPProxySidecar: &falseValue,
			Sidecar: SidecarParams{
				Type: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 0, len(params.Sidecars))
	})

	t.Run("AddsNoDefaultSidecarIfEmptyAndGlobalKindIsJob", func(t *testing.T) {

		params := Params{
			Kind: KindJob,
			Sidecar: SidecarParams{
				Type: SidecarTypeUnknown,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 0, len(params.Sidecars))
		assert.Equal(t, SidecarTypeUnknown, params.Sidecar.Type)
	})

	t.Run("KeepsSidecarTypeIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: SidecarTypeIstio,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, SidecarTypeIstio, params.Sidecar.Type)
	})

	t.Run("DefaultsSidecarImageToEstafetteOpenrestyIfEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Image: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		// digest for estafette/openresty-sidecar:1.5.8.2
		assert.Equal(t, "estafette/openresty-sidecar@sha256:2aa9f2c8c3f506e0f6cc70871701b5ac81aa0f12e8574c7b8213e4d0379d2ddd", params.Sidecars[0].Image)
	})

	t.Run("IfNoOpenrestySidecarPresentThenCustomSidecarsKeptAndOpenrestySidecarAdded", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: SidecarTypeESP,
			},
			Sidecars: []*SidecarParams{
				{
					Type: SidecarTypeIstio,
				},
				{
					Type: SidecarTypeCloudSQLProxy,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 4, len(params.Sidecars))
		assert.Equal(t, SidecarTypeESP, params.Sidecars[0].Type)
		assert.Equal(t, SidecarTypeIstio, params.Sidecars[1].Type)
		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[2].Type)
		assert.Equal(t, SidecarTypeOpenresty, params.Sidecars[3].Type)
	})

	t.Run("SidecarIsOpenrestyThenItsPrependedToSidecars", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: SidecarTypeOpenresty,
			},
			Sidecars: []*SidecarParams{
				{
					Type: SidecarTypeIstio,
				},
				{
					Type: SidecarTypeCloudSQLProxy,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, len(params.Sidecars))
		assert.Equal(t, SidecarTypeOpenresty, params.Sidecars[0].Type)
		assert.Equal(t, SidecarTypeIstio, params.Sidecars[1].Type)
		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[2].Type)
	})

	t.Run("OneOfTheSidecarsIsOpenrestyThenOtherSidecarsAreKeptAndNoExtraSidecarAdded", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: SidecarTypeIstio,
			},
			Sidecars: []*SidecarParams{
				{
					Type: SidecarTypeOpenresty,
				},
				{
					Type: SidecarTypeCloudSQLProxy,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, len(params.Sidecars))
		assert.Equal(t, SidecarTypeIstio, params.Sidecars[0].Type)
		assert.Equal(t, SidecarTypeOpenresty, params.Sidecars[1].Type)
		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[2].Type)
	})

	t.Run("KeepsSidecarImageIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Image: "estafette/openresty-sidecar:latest",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "estafette/openresty-sidecar:latest", params.Sidecar.Image)
	})

	t.Run("DefaultsSidecarHealthCheckPathToContainerReadinessPathIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "/myreadiness",
				},
			},
			Sidecar: SidecarParams{
				Type:            SidecarTypeOpenresty,
				HealthCheckPath: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/myreadiness", params.Sidecar.HealthCheckPath)
	})

	t.Run("KeepsSidecarHealthCheckPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "/myreadiness",
				},
			},
			Sidecar: SidecarParams{
				HealthCheckPath: "/nomyreadiness",
			},
		}
		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/nomyreadiness", params.Sidecar.HealthCheckPath)
	})

	t.Run("KeepsSidecarMemoryLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "1024Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "1024Mi", params.Sidecar.Memory.Limit)
	})

	t.Run("SetsHealthCheckPathDefault", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "testReadinessPath",
				},
			},
			Sidecar: SidecarParams{
				Type: SidecarTypeOpenresty,
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "testReadinessPath", params.Sidecar.HealthCheckPath)
	})

	t.Run("SetsHealthCheckPathDefaultEvenIfImageIsCustomized", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "testReadinessPath",
				},
			},
			Sidecar: SidecarParams{
				Type:  SidecarTypeOpenresty,
				Image: "testImage",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "testReadinessPath", params.Sidecar.HealthCheckPath)
	})

	t.Run("DefaultsBasePathToSlashIfEmpty", func(t *testing.T) {

		params := Params{
			Basepath: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/", params.Basepath)
	})

	t.Run("KeepsBasepathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Basepath: "/api",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/api", params.Basepath)
	})

	t.Run("DefaultsRollingUpdateMaxSurgeTo25PercentIfEmpty", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxSurge: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "25%", params.RollingUpdate.MaxSurge)
	})

	t.Run("KeepsRollingUpdateMaxSurgeIfNotEmpty", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxSurge: "10%",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "10%", params.RollingUpdate.MaxSurge)
	})

	t.Run("DefaultsRollingUpdateMaxUnavailableTo0IfEmpty", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxUnavailable: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "0", params.RollingUpdate.MaxUnavailable)
	})

	t.Run("KeepsRollingUpdateMaxUnavailableIfNotEmpty", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxUnavailable: "20%",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "20%", params.RollingUpdate.MaxUnavailable)
	})

	t.Run("DefaultsRollingUpdateTimeoutTo5MinutesIfEmpty", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				Timeout: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "5m", params.RollingUpdate.Timeout)
	})

	t.Run("KeepsRollingUpdateTimeoutIfNotEmpty", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				Timeout: "10m",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "10m", params.RollingUpdate.Timeout)
	})

	t.Run("SetBuildVersionToBuildVersion", func(t *testing.T) {

		params := Params{}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", "", "", "", buildVersion, "", "", "", map[string]string{})

		assert.Equal(t, "1.0.0", params.BuildVersion)
	})

	t.Run("DefaultsConfigMountPathToSlashConfigsIfEmpty", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				MountPath: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/configs", params.Configs.MountPath)
	})

	t.Run("KeepsConfigMountPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				MountPath: "/etc/app-config",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/etc/app-config", params.Configs.MountPath)
	})

	t.Run("DefaultsSecretMountPathToSlashSecretsIfEmpty", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				MountPath: "",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/secrets", params.Secrets.MountPath)
	})

	t.Run("KeepsSecretMountPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				MountPath: "/etc/app-secret",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "/etc/app-secret", params.Secrets.MountPath)
	})

	t.Run("DefaultsTrustedIPRangesToCloudflareIPsIfEmpty", func(t *testing.T) {

		params := Params{
			TrustedIPRanges: []string{},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 14, len(params.TrustedIPRanges))
		assert.Equal(t, "103.21.244.0/22", params.TrustedIPRanges[0])
		assert.Equal(t, "198.41.128.0/17", params.TrustedIPRanges[13])
	})

	t.Run("KeepsTrustedIPRangesIfNotEmpty", func(t *testing.T) {

		params := Params{
			TrustedIPRanges: []string{
				"0.0.0.0/0",
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, len(params.TrustedIPRanges))
		assert.Equal(t, "0.0.0.0/0", params.TrustedIPRanges[0])
	})

	t.Run("DefaultsActionToDeploySimpleIfEmpty", func(t *testing.T) {

		params := Params{
			Action: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, ActionDeploySimple, params.Action)
	})

	t.Run("KeepsActionIfNotEmpty", func(t *testing.T) {

		params := Params{
			Action: ActionDeployCanary,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, ActionDeployCanary, params.Action)
	})

	t.Run("OverridesActionIfActionIsNotEmptyButReleaseActionIsSet", func(t *testing.T) {

		params := Params{
			Action: ActionDeployCanary,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", ActionRollbackCanary, "", map[string]string{})

		assert.Equal(t, ActionRollbackCanary, params.Action)
	})

	t.Run("SetsActionToReleaseActionIfActionIsEmptyAndReleaseActionIsSet", func(t *testing.T) {

		params := Params{
			Action: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", ActionRollbackCanary, "", map[string]string{})

		assert.Equal(t, ActionRollbackCanary, params.Action)
	})

	t.Run("DefaultsKindToDeploymentIfEmpty", func(t *testing.T) {

		params := Params{
			Kind: KindUnknown,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, KindDeployment, params.Kind)
	})

	t.Run("DefaultsToAllowConcurrencyPolicyForCronJobs", func(t *testing.T) {

		params := Params{
			Kind:              KindCronJob,
			ConcurrencyPolicy: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "Allow", params.ConcurrencyPolicy)
	})

	t.Run("DefaultsRestartPolicyToOnFailure", func(t *testing.T) {

		params := Params{
			RestartPolicy: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "OnFailure", params.RestartPolicy)
	})

	t.Run("KeepsRestartPolicyIfSet", func(t *testing.T) {

		params := Params{
			RestartPolicy: "Never",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "Never", params.RestartPolicy)
	})

	t.Run("DefaultsCompletionsTo1", func(t *testing.T) {

		params := Params{
			Completions: 0,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, params.Completions)
	})

	t.Run("KeepsCompletionsIfSet", func(t *testing.T) {

		params := Params{
			Completions: 3,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, params.Completions)
	})

	t.Run("DefaultsParallelismTo1", func(t *testing.T) {

		params := Params{
			Parallelism: 0,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 1, params.Parallelism)
	})

	t.Run("KeepsParallelismIfSet", func(t *testing.T) {

		params := Params{
			Parallelism: 3,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, params.Parallelism)
	})

	t.Run("DefaultsBackoffLimitTo6", func(t *testing.T) {

		params := Params{
			BackoffLimit: nil,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 6, *params.BackoffLimit)
	})

	t.Run("KeepsBackoffLimitIfSet", func(t *testing.T) {

		backoffLimit := 3
		params := Params{
			BackoffLimit: &backoffLimit,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, *params.BackoffLimit)
	})

	t.Run("KeepsBackoffLimitZeroIfSet", func(t *testing.T) {

		backoffLimit := 0
		params := Params{
			BackoffLimit: &backoffLimit,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 0, *params.BackoffLimit)
	})

	t.Run("KeepsKindIfNotEmpty", func(t *testing.T) {

		params := Params{
			Kind: KindJob,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, KindJob, params.Kind)
	})

	t.Run("DefaultsToParallelPodManagementPolicyForStatefulsets", func(t *testing.T) {

		params := Params{
			Kind:                KindStatefulset,
			PodManagementPolicy: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "Parallel", params.PodManagementPolicy)
	})

	t.Run("DefaultsToStandardStorageClassForStatefulsets", func(t *testing.T) {

		params := Params{
			Kind:         KindStatefulset,
			StorageClass: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "standard", params.StorageClass)
	})

	t.Run("DefaultsTo1GiStorageSizeForStatefulsets", func(t *testing.T) {

		params := Params{
			Kind:        KindStatefulset,
			StorageSize: "",
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, "1Gi", params.StorageSize)
	})

	t.Run("DefaultsTo3VerifyDepthForApigee", func(t *testing.T) {

		params := Params{
			Kind:       KindDeployment,
			Visibility: VisibilityApigee,
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, 3, params.Request.VerifyDepth)
	})
}

func TestValidateRequiredProperties(t *testing.T) {

	t.Run("ReturnsFalseIfAppIsNotSet", func(t *testing.T) {

		params := validParams
		params.App = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAppIsSet", func(t *testing.T) {

		params := validParams
		params.App = "myapp"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfNamespaceIsNotSet", func(t *testing.T) {

		params := validParams
		params.Namespace = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfNamespaceIsSet", func(t *testing.T) {

		params := validParams
		params.Namespace = "mynamespace"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfImageRepositoryIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageRepository = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageRepositoryIsSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageRepository = "myrepository"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfImageNameIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageName = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageNameIsSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageName = "myimage"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfImageTagIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageTag = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageTagIsSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageTag = "1.0.0"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfVisibilityIsNotSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfVisibilityIsSetToUnsupportedValue", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = "everywhere"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfVisibilityIsSetToPublic", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = VisibilityPublic

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfVisibilityIsSetToPrivate", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = VisibilityPrivate

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfVisibilityIsSetToIAP", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = "iap"
		params.IapOauthCredentialsClientID = "123123"
		params.IapOauthCredentialsClientSecret = "somesecret"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfVisibilityIsSetToPublicWhitelist", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = VisibilityPublicWhitelist

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfCpuRequestIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Request = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfCpuRequestIsSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Request = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfCpuLimitIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Limit = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfCpuLimitIsSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Limit = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMemoryRequestIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Request = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMemoryRequestIsSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Request = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMemoryLimitIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Limit = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMemoryLimitIsSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Limit = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfContainerPortIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.Port = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfContainerPortIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.Port = 5000

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfHostsAreNotSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Hosts = []string{}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfOneOrMoreHostsAreSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Hosts = []string{"gke.estafette.io"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfOneOrMoreUppercaseHostsAreSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Hosts = []string{"GKE.ESTAFETTE.IO"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfHostHasLabelsLongerThan63Characters", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Hosts = []string{"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl.estafette.io"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfHostIsLongerThan253Characters", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Hosts = []string{"ab.abcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.estafette.io"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfHostHasOtherCharacterThanAlphaNumericOrHyphen", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Hosts = []string{"gke_site.estafette.io"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfInternalHostsAreNotSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.InternalHosts = []string{}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfOneOrMoreInternalHostsAreSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.InternalHosts = []string{"ci.estafette.internal"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfOneOrMoreUppercaseInternalHostsAreSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.InternalHosts = []string{"CI.ESTAFETTE.INTERNAL"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfInternalHostHasLabelsLongerThan63Characters", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.InternalHosts = []string{"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl.estafette.internal"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfInternalHostIsLongerThan253Characters", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.InternalHosts = []string{"abcdefghijklmnopqrstuvw.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.estafette.internal"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfInternalHostHasOtherCharacterThanAlphaNumericOrHyphen", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.InternalHosts = []string{"gke_site.estafette.internal"}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfAutoscaleMinReplicasIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Autoscale.MinReplicas = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAutoscaleMinReplicasIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Autoscale.MinReplicas = 5

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfAutoscaleMaxReplicasIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Autoscale.MaxReplicas = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAutoscaleMaxReplicasIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Autoscale.MaxReplicas = 15

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfAutoscaleCPUPercentageIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Autoscale.CPUPercentage = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAutoscaleCPUPercentageIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Autoscale.CPUPercentage = 35

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessPathIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.Path = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessPathIsNotEmpty", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.Path = "/liveness"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessProbePortIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.Port = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessProbePortIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.Port = 5000

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessInitialDelaySecondsIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.InitialDelaySeconds = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessInitialDelaySecondsIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.InitialDelaySeconds = 30

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessTimeoutSecondsIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.TimeoutSeconds = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessTimeoutSecondsIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.TimeoutSeconds = 2

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfReadinessProbePathIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.Path = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfReadinessProbePathIsNotEmpty", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.Path = "/readiness"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfReadinessProbePortIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.Port = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfReadinessProbePortIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.Port = 5000

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfReadinessProbeTimeoutSecondsIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.TimeoutSeconds = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfReadinessProbeTimeoutSecondsIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.TimeoutSeconds = 2

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsPathIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Path = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMetricsPathIsNotEmpty", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Path = "/metrics"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfMetricsPathIsEmptyButScrapeIsFalse", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = &falseValue
		params.Container.Metrics.Path = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsPortIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Port = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMetricsPortIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Port = 5000

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfMetricsPortIsZeroOrLessButScrapeIsFalse", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = &falseValue
		params.Container.Metrics.Port = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsScrapeIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = nil

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMetricsScrapeIsTrue", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = &trueValue

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfMetricsScrapeIsFalse", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = &falseValue

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseSidecarTypeIsEmpty", func(t *testing.T) {

		params := validParams

		params.Sidecars = []*SidecarParams{
			&SidecarParams{
				Type:  "",
				Image: "docker",
				CPU: CPUParams{
					Request: "10m",
					Limit:   "50m",
				},
				Memory: MemoryParams{
					Request: "10Mi",
					Limit:   "50Mi",
				},
			},
		}

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSidecarTypeIsSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Type = SidecarTypeOpenresty

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfSidecarImageIsNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Image = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSidecarImageIsSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Image = "estafette/openresty-sidecar:1.13.6.2-alpine"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfSidecarCpuRequestIsNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.CPU.Request = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSidecarCpuRequestIsSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.CPU.Request = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfSidecarCpuLimitIsNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.CPU.Limit = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfSidecarCpuLimitIsSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.CPU.Limit = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfSidecarMemoryRequestIsNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Memory.Request = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSidecarMemoryRequestIsSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Memory.Request = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfSidecarMemoryLimitIsNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Memory.Limit = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSidecarMemoryLimitIsSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Memory.Limit = "100m"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfSqlProxyInstanceNameNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Type = SidecarTypeCloudSQLProxy
		params.Sidecar.DbInstanceConnectionName = ""
		params.Sidecar.SQLProxyPort = 8080

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) == 1)
	})

	t.Run("ReturnsFalseIfSqlProxyPortNotSet", func(t *testing.T) {

		params := validParams
		params.Sidecar.Type = SidecarTypeCloudSQLProxy
		params.Sidecar.DbInstanceConnectionName = "instance"
		params.Sidecar.SQLProxyPort = 0

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) == 1)
	})

	t.Run("ReturnsTrueIfSqlProxyProperlyConfigured", func(t *testing.T) {

		params := validParams
		params.Sidecar.Type = SidecarTypeCloudSQLProxy
		params.Sidecar.DbInstanceConnectionName = "instance"
		params.Sidecar.SQLProxyPort = 8080

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsWarningIfSidecarFieldUsed", func(t *testing.T) {

		params := validParams
		params.Sidecar = SidecarParams{
			Type:  SidecarTypeOpenresty,
			Image: "estafette/openresty-sidecar:1.13.6.2-alpine",
			CPU: CPUParams{
				Request: "10m",
				Limit:   "50m",
			},
			Memory: MemoryParams{
				Request: "10Mi",
				Limit:   "50Mi",
			},
		}

		// act
		valid, errors, warnings := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
		assert.Equal(t, 1, len(warnings))
	})

	t.Run("NoWarningRetuendIfSidecarsCollectionUsed", func(t *testing.T) {

		params := validParams
		params.Sidecar = SidecarParams{
			Type: "",
		}

		params.Sidecars = []*SidecarParams{
			&SidecarParams{
				Type:  SidecarTypeOpenresty,
				Image: "estafette/openresty-sidecar:1.13.6.2-alpine",
				CPU: CPUParams{
					Request: "10m",
					Limit:   "50m",
				},
				Memory: MemoryParams{
					Request: "10Mi",
					Limit:   "50Mi",
				},
			},
		}

		// act
		valid, errors, warnings := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
		assert.Equal(t, 0, len(warnings))
	})

	t.Run("ReturnsFalseIfBasepathIsNotSet", func(t *testing.T) {

		params := validParams
		params.Basepath = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfBasepathIsSet", func(t *testing.T) {

		params := validParams
		params.Basepath = "/"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfRollingUpdateMaxSurgeIsNotSet", func(t *testing.T) {

		params := validParams
		params.RollingUpdate.MaxSurge = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfRollingUpdateMaxSurgeIsSet", func(t *testing.T) {

		params := validParams
		params.RollingUpdate.MaxSurge = "25%"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfRollingUpdateMaxUnavailableIsNotSet", func(t *testing.T) {

		params := validParams
		params.RollingUpdate.MaxUnavailable = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfRollingUpdateMaxUnavailableIsSet", func(t *testing.T) {

		params := validParams
		params.RollingUpdate.MaxUnavailable = "25%"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfScheduleIsNotSetAndKindIsCronjob", func(t *testing.T) {

		params := validParams
		params.Kind = KindCronJob
		params.Schedule = ""
		params.ConcurrencyPolicy = "Allow"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfConcurrencyPolicyIsInvalidAndKindIsCronjob", func(t *testing.T) {

		params := validParams
		params.Kind = KindCronJob
		params.Schedule = "*/5 * * * *"
		params.ConcurrencyPolicy = "InvalidPolicy"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfScheduleIsSetAndConcurrencyPolicyIsValidAndKindIsCronjob", func(t *testing.T) {

		params := validParams
		params.Kind = KindCronJob
		params.Schedule = "*/5 * * * *"
		params.ConcurrencyPolicy = "Allow"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfPodManagementPolicyIsInvalidAndKindIsStatefulset", func(t *testing.T) {

		params := validParams
		params.Kind = KindStatefulset
		params.PodManagementPolicy = "InvalidPolicy"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfPodManagementPolicyIsValidAndKindIsStatefulset", func(t *testing.T) {

		params := validParams
		params.Kind = KindStatefulset
		params.PodManagementPolicy = "Parallel"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfStorageClassIsEmptyAndKindIsStatefulset", func(t *testing.T) {

		params := validParams
		params.Kind = KindStatefulset
		params.StorageClass = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfStorageClassIsSetAndKindIsStatefulset", func(t *testing.T) {

		params := validParams
		params.Kind = KindStatefulset
		params.StorageClass = "standard"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfStorageSizeIsEmptyAndKindIsStatefulset", func(t *testing.T) {

		params := validParams
		params.Kind = KindStatefulset
		params.StorageSize = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfStorageSizeIsSetAndKindIsStatefulset", func(t *testing.T) {

		params := validParams
		params.Kind = KindStatefulset
		params.StorageSize = "1Gi"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLoadBalanceAlgorithmIsNonValidValue", func(t *testing.T) {

		params := validParams
		params.Request.LoadBalanceAlgorithm = "leastconn"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLoadBalanceAlgorithmIsEmpty", func(t *testing.T) {

		params := validParams
		params.Request.LoadBalanceAlgorithm = ""

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfLoadBalanceAlgorithmIsEwma", func(t *testing.T) {

		params := validParams
		params.Request.LoadBalanceAlgorithm = "ewma"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfLoadBalanceAlgorithmIsRoundRobin", func(t *testing.T) {

		params := validParams
		params.Request.LoadBalanceAlgorithm = "round_robin"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseAuthSecretIsNotSetForVisibilityApigee", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = VisibilityApigee

		// act
		valid, errors, _ := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.Equal(t, 1, len(errors))
	})

	t.Run("WarnOnUsingVisibilityPublic", func(t *testing.T) {

		params := validParams
		params.Visibility = VisibilityPublic

		// act
		valid, _, warnings := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.Equal(t, 1, len(warnings))
	})

	t.Run("ReturnsFalseIfEspEndpointsProjectIDNotSet", func(t *testing.T) {

		params := validParams
		params.Kind = KindDeployment
		params.Visibility = VisibilityESP
		params.EspEndpointsProjectID = ""
		error_string := "With visibility 'esp' property espEndpointsProjectID is required; provide id of the 'endpoints' project"

		// act
		valid, errors, _ := params.ValidateRequiredProperties()
		assert.False(t, valid)
		assert.Equal(t, error_string, stringInErrorSlice(error_string, errors))
	})
}

func TestReplaceSidecarTagsWithDigest(t *testing.T) {

	t.Run("ReplacesFirstSidecarImageTagWithDigest", func(t *testing.T) {

		params := validParams

		// act
		params.ReplaceSidecarTagsWithDigest()

		assert.Equal(t, SidecarTypeOpenresty, params.Sidecars[0].Type)
		assert.True(t, strings.HasPrefix(params.Sidecars[0].Image, "estafette/openresty-sidecar@sha256:"))
	})

	t.Run("KeepsFirstSidecarImageTagWithDigest", func(t *testing.T) {

		params := validParams
		params.Sidecars[0].Image = "estafette/openresty-sidecar@sha256:4300dc7d45600c428f4196009ee842c1c3bdd51aaa4f55361479f6fa60e78faf"

		// act
		params.ReplaceSidecarTagsWithDigest()

		assert.Equal(t, SidecarTypeOpenresty, params.Sidecars[0].Type)
		assert.True(t, strings.HasPrefix(params.Sidecars[0].Image, "estafette/openresty-sidecar@sha256:"))
	})

	t.Run("ReplacesLastSidecarImageTagWithDigest", func(t *testing.T) {

		params := validParams

		// act
		params.ReplaceSidecarTagsWithDigest()

		assert.Equal(t, SidecarTypeESP, params.Sidecars[1].Type)
		assert.True(t, strings.HasPrefix(params.Sidecars[1].Image, "estafette/estafette-docker-cache-heater@sha256:"))
	})

	t.Run("KeepsLastSidecarImageTagWithDigest", func(t *testing.T) {

		params := validParams
		params.Sidecars[1].Image = "estafette/estafette-docker-cache-heater@sha256:4300dc7d45600c428f4196009ee842c1c3bdd51aaa4f55361479f6fa60e78faf"

		// act
		params.ReplaceSidecarTagsWithDigest()

		assert.Equal(t, SidecarTypeESP, params.Sidecars[1].Type)
		assert.True(t, strings.HasPrefix(params.Sidecars[1].Image, "estafette/estafette-docker-cache-heater@sha256:"))
	})

	t.Run("DefaultsSQLProxyPortTo5432IfNotSet", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				{
					Type:         SidecarTypeCloudSQLProxy,
					SQLProxyPort: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[0].Type)
		assert.Equal(t, 5432, params.Sidecars[0].SQLProxyPort)
	})

	t.Run("KeepsSQLProxyPortIfSet", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				{
					Type:         SidecarTypeCloudSQLProxy,
					SQLProxyPort: 1601,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[0].Type)
		assert.Equal(t, 1601, params.Sidecars[0].SQLProxyPort)
	})

	t.Run("DefaultsSQLProxyTerminationTimeoutSecondsTo60IfNotSet", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				{
					Type:                              SidecarTypeCloudSQLProxy,
					SQLProxyTerminationTimeoutSeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[0].Type)
		assert.Equal(t, 60, params.Sidecars[0].SQLProxyTerminationTimeoutSeconds)
	})

	t.Run("KeepsSQLProxyTerminationTimeoutSecondsIfSet", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				{
					Type:                              SidecarTypeCloudSQLProxy,
					SQLProxyTerminationTimeoutSeconds: 45,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", "", "", "", "", "", map[string]string{})

		assert.Equal(t, SidecarTypeCloudSQLProxy, params.Sidecars[0].Type)
		assert.Equal(t, 45, params.Sidecars[0].SQLProxyTerminationTimeoutSeconds)
	})
}
