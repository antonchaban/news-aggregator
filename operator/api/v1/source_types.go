/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SourceSpec defines the desired state of Source,
// includes fields for source creation and management in the news aggregator.
type SourceSpec struct {
	// - Id: Unique identifier for the source.
	Id int `json:"id,omitempty"`
	// - Name: full name of the source, as it appears in the news aggregator.
	Name string `json:"name,omitempty"`
	// - Link: URL link to the source.
	Link string `json:"link,omitempty"`
	// - ShortName: Shortened name of the source for search purposes.
	ShortName string `json:"short_name,omitempty"`
}

// SourceConditionType represents the type of condition for Source
// Valid values are:
// - "Added": Indicates that the source has been added.
// - "Updated": Indicates that the source has been updated.
type SourceConditionType string

const (
	SourceAdded   SourceConditionType = "Added"
	SourceUpdated SourceConditionType = "Updated"
)

// SourceCondition defines the current condition of Source, at the moment of condition creation.
type SourceCondition struct {
	// - Type: stands for type of condition (e.g., Added, Updated).
	Type SourceConditionType `json:"type"`
	// - Status: shows status of the condition, one of True or False depending on its success.
	Status metav1.ConditionStatus `json:"status"`
	// - LastUpdateTime: represents the last time the condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// - Reason: indicates the reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Message: provides a human-readable message indicating details about the last transition.
	Message string `json:"message,omitempty"`
}

// SourceStatus defines the observed state of Source
// It contains the current status of the Source resource, including conditions.
// Fields:
// - ID: Unique identifier for the source.
// - Conditions: List of conditions associated with the source.
type SourceStatus struct {
	ID         int               `json:"id,omitempty"`
	Conditions []SourceCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Source is the Schema for the sources API
// It represents a source resource from news aggregator in the k8s cluster.
type Source struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SourceSpec   `json:"spec,omitempty"`
	Status SourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SourceList contains a list of Source resources.
type SourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Source `json:"items"`
}

func (s *SourceStatus) SetCondition(condition SourceCondition) {
	for i, c := range s.Conditions {
		if c.Type == condition.Type {
			s.Conditions[i] = condition
			return
		}
	}
	s.Conditions = append(s.Conditions, condition)
}

func init() {
	SchemeBuilder.Register(&Source{}, &SourceList{})
}
