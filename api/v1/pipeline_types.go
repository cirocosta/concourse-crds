package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PipelineSpec defines the desired state of the pipeline.
//
type PipelineSpec struct {
	// +kubebuilder:validation:MinLength=1
	Team string `json:"team"`

	// Config is
	Config runtime.RawExtension `json:"config"`

	// +optional
	Paused *bool `json:"paused,omitempty"`

	// +optional
	Exposed *bool `json:"exposed,omitempty"`

	// +optional
	CheckCreds *bool `json:"checkCreds,omitempty"`

	// +optional
	Vars *runtime.RawExtension `json:"vars,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// is public
	// is paused

	LastSetTime     *metav1.Time `json:"lastSetTime,omitempty"`
	LastUnpauseTime *metav1.Time `json:"lastUnpauseTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Team",type=string,JSONPath=`.spec.team`
// +kubebuilder:printcolumn:name="Public",type=string,JSONPath=`.spec.exposed`
// +kubebuilder:printcolumn:name="Paused",type=string,JSONPath=`.spec.paused`

// Pipeline is the Schema for the pipelines API
//
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
//
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
