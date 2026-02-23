package data

type DeploymentStatus struct {
	Kind              string `json:"kind"`
	Name              string `json:"name"`
	DesiredReplicas   int32  `json:"desired_replicas"`
	ReadyReplicas     int32  `json:"ready_replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
	UpdatedReplicas   int32  `json:"updated_replicas"`
	Healthy           bool   `json:"healthy"`
	ShortMessage      string `json:"short_message"`
	LongMessage       string `json:"long_message"`
}

// POD STATUS
type PodStatus struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Phase        string `json:"phase"`
	Ready        bool   `json:"ready"`
	Containers   int    `json:"containers"`
	Healthy      bool   `json:"healthy"`
	ShortMessage string `json:"short_message"`
	LongMessage  string `json:"long_message"`
}

// StatefulSet STATUS
type StatefulSetStatus struct {
	Kind            string `json:"kind"`
	Name            string `json:"name"`
	DesiredReplicas int32  `json:"desired_replicas"`
	ReadyReplicas   int32  `json:"ready_replicas"`
	CurrentReplicas int32  `json:"current_replicas"`
	UpdatedReplicas int32  `json:"updated_replicas"`
	Healthy         bool   `json:"healthy"`
	ShortMessage    string `json:"short_message"`
	LongMessage     string `json:"long_message"`
}

// DEMONSETATUS
type DaemonSetStatus struct {
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	DesiredNumber      int32  `json:"desired_number"`
	ReadyNumber        int32  `json:"ready_number"`
	NumberMisscheduled int32  `json:"number_misscheduled"`
	AvailableNumber    int32  `json:"available_number"`
	Healthy            bool   `json:"healthy"`
	ShortMessage       string `json:"short_message"`
	LongMessage        string `json:"long_message"`
}

// JOB STATUS
type JobStatus struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Active       int32  `json:"active"`
	Succeeded    int32  `json:"succeeded"`
	Failed       int32  `json:"failed"`
	Completions  int32  `json:"completions"`
	Healthy      bool   `json:"healthy"`
	ShortMessage string `json:"short_message"`
	LongMessage  string `json:"long_message"`
}

// PVC STATUS
type PVCStatus struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Phase        string `json:"phase"`
	Healthy      bool   `json:"healthy"`
	ShortMessage string `json:"short_message"`
	LongMessage  string `json:"long_message"`
}

type ServiceStatus struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Healthy      bool   `json:"healthy"`
	ShortMessage string `json:"short_message"`
	LongMessage  string `json:"long_message"`
}

// high priority resources

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
