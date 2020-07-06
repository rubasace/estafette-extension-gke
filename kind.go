package main

type Kind string

const (
	KindDeployment         Kind = "deployment"
	KindHeadlessDeployment Kind = "headless-deployment"
	KindStatefulset        Kind = "statefulset"
	KindJob                Kind = "job"
	KindCronJob            Kind = "cronjob"
	KindConfig             Kind = "config"

	KindUnknown Kind = ""
)
