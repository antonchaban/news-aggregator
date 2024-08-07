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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SourceSpec defines Source fields
type SourceSpec struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Link      string `json:"link,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

// SourceConditionType is a valid value for SourceCondition.Type
type SourceConditionType string

const (
	SourceAdded   SourceConditionType = "Added"
	SourceUpdated SourceConditionType = "Updated"
)

// SourceCondition defines the observed state of Source at a certain point.
type SourceCondition struct {
	// Type of condition.
	Type SourceConditionType `json:"type"`
	// Status of the condition, one of True or False.
	Status metav1.ConditionStatus `json:"status"`
	// Last time the condition was checked.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

// SourceStatus defines the observed state of Source
type SourceStatus struct {
	ID         int               `json:"id,omitempty"`
	Conditions []SourceCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Source is the Schema for the sources API
type Source struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SourceSpec   `json:"spec,omitempty"`
	Status SourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SourceList contains a list of Source
type SourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Source `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Source{}, &SourceList{})
}
