package output

import "fmt"

type OutputText struct {
	ReleaseName string         `json:"release_name"`
	Namespace   string         `json:"namespace"`
	Status      string         `json:"status"`
	Health      string         `json:"health"`
	Message     string         `json:"message"`
	Resources   map[string]any `json:"resources"`
}

type Output interface {
	Print()
	UpdateStatus(status string)
	UpdateHealth(health string)
}

func NewOutputText(releaseName, namespace string) *OutputText {
	return &OutputText{
		ReleaseName: releaseName,
		Namespace:   namespace,
		Status:      "",
		Health:      "",
		Message:     "",
		Resources:   make(map[string]any),
	}
}

func (o *OutputText) Print() {
	fmt.Printf("RELEASE:%s\n", o.ReleaseName)
	fmt.Printf("NAMESPACE:%s\n", o.Namespace)
	fmt.Printf("STATUS:%s\n", o.Status)
	fmt.Printf("HEALTH:%s\n", o.Health)
}

func (o *OutputText) UpdateStatus(status string) {
	o.Status = status
}

func (o *OutputText) UpdateHealth(health string) {
	o.Health = health
}
