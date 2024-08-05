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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HotNewsSpec defines the desired state of HotNews
type HotNewsSpec struct {
	Keywords      []string      `json:"keywords,omitempty"`
	StartDate     string        `json:"date_start,omitempty"`
	EndDate       string        `json:"date_end,omitempty"` // todo maybe should use time.Time
	Sources       []Source      `json:"sources,omitempty"`
	SourceGroups  []string      `json:"source_groups,omitempty"`
	SummaryConfig SummaryConfig `json:"summary_config,omitempty"`
}

type SummaryConfig struct {
	TitlesCount int `json:"titles_count,omitempty"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	ArticlesCount  int      `json:"articles_count,omitempty"`
	NewsLink       string   `json:"news_link,omitempty"`
	ArticlesTitles []string `json:"articles_titles,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotNewsSpec   `json:"spec,omitempty"`
	Status HotNewsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HotNewsList contains a list of HotNews
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
