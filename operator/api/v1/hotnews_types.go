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

// HotNewsSpec defines the desired state of HotNews.
// It contains specifications for filtering news articles based on keywords,
// dates, sources, and feed groups, as well as summary configuration.
type HotNewsSpec struct {
	// Keywords is a list of keywords to filter news articles.
	// This field is required and must contain at least one keyword.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Keywords []string `json:"keywords"`

	// DateStart is the start date for filtering news articles.
	// Format must be YYYY-MM-DD.
	// This field is optional.
	// +kubebuilder:validation:Pattern=`^\d{4}-\d{2}-\d{2}$`
	// +optional
	DateStart string `json:"date_start,omitempty"`

	// DateEnd is the end date for filtering news articles.
	// Format must be YYYY-MM-DD.
	// This field is optional.
	// +kubebuilder:validation:Pattern=`^\d{4}-\d{2}-\d{2}$`
	// +optional
	DateEnd string `json:"date_end,omitempty"`

	// Sources specifies the short names of Source resources to filter news from.
	// These should match the 'spec.shortName' field of Source resources in the same namespace.
	// If empty, all available feeds will be used.
	// This field is optional.
	// +optional
	Sources []string `json:"sources,omitempty"`

	// FeedGroups specifies the feed groups defined in the 'feed-group-source' ConfigMap
	// in the 'news-alligator' namespace. The feeds in these groups will be used to filter news.
	// The specified feed groups must exist in the ConfigMap; otherwise, validation will fail.
	// This field is optional.
	// +optional
	FeedGroups []string `json:"feed_groups,omitempty"`

	// SummaryConfig defines the configuration for the summary in the status.
	SummaryConfig SummaryConfig `json:"summary_config"`
}

// SummaryConfig defines the summary configuration for the HotNews resource.
type SummaryConfig struct {
	// TitlesCount is the number of article titles to include in the summary.
	// Must be greater than or equal to 1.
	// Defaults to 10 if not specified.
	TitlesCount int `json:"titles_count"`
}

// HNewsConditionType is a string representing the condition type for HotNews.
type HNewsConditionType string

const (
	// HNewsAdded indicates that the HotNews resource has been added.
	HNewsAdded HNewsConditionType = "Added"
	// HNewsUpdated indicates that the HotNews resource has been updated.
	HNewsUpdated HNewsConditionType = "Updated"
)

// HotNewsCondition represents the condition of a HotNews resource.
type HotNewsCondition struct {
	Type           HNewsConditionType     `json:"type"`                     // Type of HotNews condition.
	Status         metav1.ConditionStatus `json:"status"`                   // Status of the condition, one of True, False, or Unknown.
	LastUpdateTime metav1.Time            `json:"lastUpdateTime,omitempty"` // LastUpdateTime is the last time the condition was updated.
	Reason         string                 `json:"reason,omitempty"`         // Reason is a brief explanation for the condition's last transition.
	Message        string                 `json:"message,omitempty"`        // Message is a human-readable message indicating details about the transition.
}

// HotNewsStatus represents the observed state of HotNews.
// It contains information about the count of articles, a link to the news aggregator,
// a summary of article titles, and conditions.
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

// SetCondition updates or adds a condition to the HotNewsStatus.
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
