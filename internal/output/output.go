package output

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	CheckMark   = "✓"
	CrossMark   = "✗"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

const (
	StatusHealthy   = "Healthy"
	StatusUnhealthy = "Unhealthy"
	StatusUnknown   = "Unknown"
)

// Resource represents a single Kubernetes resource's health status.
type Resource struct {
	Kind      string         `json:"kind"`
	Name      string         `json:"name"`
	Namespace string         `json:"namespace,omitempty"`
	Status    string         `json:"status"`
	Health    map[string]any `json:"health,omitempty"`
	Message   string         `json:"-"` // Short summary for text output
	Issues    []string       `json:"-"` // Detailed issue lines for text output
}

// Summary holds aggregated resource health counts.
type Summary struct {
	Total     int `json:"total"`
	Healthy   int `json:"healthy"`
	Unhealthy int `json:"unhealthy"`
	Unknown   int `json:"unknown"`
}

// OutputResult is the top-level result of a health check.
type OutputResult struct {
	Release   string     `json:"release"`
	Namespace string     `json:"namespace"`
	Status    string     `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
	Duration  string     `json:"duration"`
	Summary   Summary    `json:"summary"`
	Resources []Resource `json:"resources"`
}

// OutputFormat controls how results are rendered.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

// --- Constructors ---

func NewResource(kind, name, namespace string) *Resource {
	return &Resource{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
		Status:    StatusUnknown,
		Health:    make(map[string]any),
		Issues:    []string{},
	}
}

func NewOutputResult(release, namespace string) *OutputResult {
	return &OutputResult{
		Release:   release,
		Namespace: namespace,
		Status:    StatusHealthy,
		Timestamp: time.Now().UTC(),
		Resources: []Resource{},
	}
}

// --- Resource builder methods ---

func (r *Resource) SetStatus(status string) *Resource {
	r.Status = status
	return r
}

func (r *Resource) SetMessage(msg string) *Resource {
	r.Message = msg
	return r
}

func (r *Resource) AddIssue(issue string) *Resource {
	r.Issues = append(r.Issues, issue)
	return r
}

func (r *Resource) SetHealth(key string, value any) *Resource {
	r.Health[key] = value
	return r
}

// --- OutputResult methods ---

func (o *OutputResult) AddResource(r Resource) {
	o.Resources = append(o.Resources, r)
}

// Finalize computes duration and summary from the collected resources.
func (o *OutputResult) Finalize(startTime time.Time) {
	duration := time.Since(startTime)
	o.Duration = fmt.Sprintf("%.1fs", duration.Seconds())

	for _, r := range o.Resources {
		o.Summary.Total++
		switch r.Status {
		case StatusHealthy:
			o.Summary.Healthy++
		case StatusUnhealthy:
			o.Summary.Unhealthy++
			o.Status = StatusUnhealthy
		default:
			o.Summary.Unknown++
		}
	}
}

// Print renders the result in the specified format.
func (o *OutputResult) Print(format OutputFormat) {
	switch format {
	case FormatJSON:
		o.printJSON()
	default:
		o.printText()
	}
}

func (o *OutputResult) printText() {
	statusColor := ColorGreen
	if o.Status != StatusHealthy {
		statusColor = ColorRed
	}

	fmt.Printf("\nRELEASE: %s\n", o.Release)
	fmt.Printf("STATUS: %s%s%s\n", statusColor, o.Status, ColorReset)
	fmt.Printf("NAMESPACE: %s\n", o.Namespace)
	fmt.Println()
	fmt.Println("RESOURCES:")

	for _, r := range o.Resources {
		switch r.Status {
		case StatusHealthy:
			fmt.Printf("%s%s %s/%s%s", ColorGreen, CheckMark, r.Kind, r.Name, ColorReset)
			if r.Message != "" {
				fmt.Printf(" (%s)", r.Message)
			}
			fmt.Println()
		case StatusUnhealthy:
			fmt.Printf("%s%s %s/%s%s", ColorRed, CrossMark, r.Kind, r.Name, ColorReset)
			if r.Message != "" {
				fmt.Printf(" (%s)", r.Message)
			}
			fmt.Println()
			for _, issue := range r.Issues {
				fmt.Printf("  %s\n", issue)
			}
		default:
			fmt.Printf("%s%s %s/%s%s", ColorYellow, CheckMark, r.Kind, r.Name, ColorReset)
			if r.Message != "" {
				fmt.Printf(" (%s)", r.Message)
			}
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Printf("Health check completed in %s\n", o.Duration)
}

func (o *OutputResult) printJSON() {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}
