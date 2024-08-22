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
// - Keywords: A list of keywords to filter news, must be always required.
// - DateStart: The start date for the news filter, can be empty.
// - DateEnd: The end date for the news filter, can be empty.
// - Sources: All source names in the current namespace, if empty, will watch ALL available feeds. This should be names of Source resources.
// - FeedGroups: Available sections of feeds from the 'feed-group-source' ConfigMap.
// - SummaryConfig: Configuration for how the status will show the summary of observed hot news.
type HotNewsSpec struct {
	Keywords      []string      `json:"keywords"`
	DateStart     string        `json:"date_start,omitempty"`
	DateEnd       string        `json:"date_end,omitempty"`
	Sources       []string      `json:"sources,omitempty"`
	FeedGroups    []string      `json:"feed_groups,omitempty"`
	SummaryConfig SummaryConfig `json:"summary_config"`
}

// SummaryConfig defines the summary configuration
// Fields:
// - TitlesCount: The number of article titles to include in the summary.
type SummaryConfig struct {
	TitlesCount int `json:"titles_count"`
}

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

// HotNewsStatus defines the observed state of HotNews
// Fields:
// - ArticlesCount: The count of articles by the criteria.
// - NewsLink: A link to the news-aggregator HTTPS server to get all news by the criteria in JSON format.
// - ArticlesTitles: The first 'spec.summaryConfig.titlesCount' article titles, sorted by feed name.
type HotNewsStatus struct {
	ArticlesCount  int                `json:"articles_count,omitempty"`
	NewsLink       string             `json:"news_link,omitempty"`
	ArticlesTitles []string           `json:"articles_titles,omitempty"`
	Conditions     []HotNewsCondition `json:"conditions,omitempty"`
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
// It is a list of HotNews resources.
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
