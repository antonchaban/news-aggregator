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

// HotNewsSpec defines the desired state of HotNews
// Fields:
type HotNewsSpec struct {
	// - Keywords: A list of keywords to filter news, must be always required.
	Keywords []string `json:"keywords"`
	// - DateStart: The start date for the news filter, can be empty.
	DateStart string `json:"date_start,omitempty"`
	// - DateEnd: The end date for the news filter, can be empty.
	DateEnd string `json:"date_end,omitempty"`
	// - Sources: All source names in the current namespace, if empty, will watch ALL available feeds. This should be names of Source resources.
	Sources []string `json:"sources,omitempty"`
	// - FeedGroups: Available sections of feeds from the 'feed-group-source' ConfigMap.
	FeedGroups []string `json:"feed_groups,omitempty"`
	// - SummaryConfig: Configuration for how the status will show the summary of observed hot news.
	SummaryConfig SummaryConfig `json:"summary_config"`
}

// SummaryConfig defines the summary configuration
type SummaryConfig struct {
	// - TitlesCount: The number of article titles to include in the summary.
	TitlesCount int `json:"titles_count"`
}

// HNewsConditionType is a type to define the condition type for HotNews.
type HNewsConditionType string

const (
	HNewsAdded   HNewsConditionType = "Added"
	HNewsUpdated HNewsConditionType = "Updated"
)

type HotNewsCondition struct {
	Type           HNewsConditionType     `json:"type"`
	Status         metav1.ConditionStatus `json:"status"`
	LastUpdateTime metav1.Time            `json:"lastUpdateTime,omitempty"`
	Reason         string                 `json:"reason,omitempty"`
	Message        string                 `json:"message,omitempty"`
}

// HotNewsStatus shows info about the HotNews resource, such
// as the count of articles, a link to the news-aggregator HTTPS server,
// the first specified amount of article titles, and conditions.
type HotNewsStatus struct {
	// - ArticlesCount: The count of articles by the criteria.
	ArticlesCount int `json:"articles_count,omitempty"`
	// - NewsLink: A link to the news-aggregator HTTPS server to get all news by the criteria in JSON format.
	NewsLink string `json:"news_link,omitempty"`
	// - ArticlesTitles: The first 'spec.summaryConfig.titlesCount' article titles, sorted by feed name.
	ArticlesTitles []string `json:"articles_titles,omitempty"`
	// - Conditions: Conditions of the HotNews resource.
	Conditions []HotNewsCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
// It represents a hot news resource in the Kubernetes cluster.
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

func (s *HotNewsStatus) SetCondition(condition HotNewsCondition) {
	for i, c := range s.Conditions {
		if c.Type == condition.Type {
			s.Conditions[i] = condition
			return
		}
	}
	s.Conditions = append(s.Conditions, condition)
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
