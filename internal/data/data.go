package data

// Resource kind constants

var (
	DEPLOYMENT         = "Deployment"
	STATEFULSET        = "StatefulSet"
	DAEMONSEET         = "DaemonSet"
	JOB                = "Job"
	POD                = "Pod"
	SERVICE            = "Service"
	ReplicaSet         = "ReplicaSet"
	PVC                = "PersistentVolumeClaim"
	CRONJOB            = "CronJob"
	INGRESS            = "Ingress"
	NETWORKPOLICY      = "NetworkPolicy"
	CONFIGMAP          = "ConfigMap"
	SECRET             = "Secret"
	SERVICEACCOUNT     = "ServiceAccount"
	ROLE               = "Role"
	ROLEBINDING        = "RoleBinding"
	CLUSTERROLE        = "ClusterRole"
	CLUSTERROLEBINDING = "ClusterRoleBinding"
)
var (
	HighPriorityResources   = []string{DEPLOYMENT, STATEFULSET, DAEMONSEET, SERVICE, ReplicaSet, JOB, POD, PVC, CRONJOB}
	MediumPriorityResources = []string{INGRESS, NETWORKPOLICY}
	LowPriorityResources    = []string{CONFIGMAP, SECRET, SERVICEACCOUNT, ROLE, ROLEBINDING, CLUSTERROLE, CLUSTERROLEBINDING}
)
