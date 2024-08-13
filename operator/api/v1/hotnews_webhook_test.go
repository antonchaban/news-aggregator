package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// setupFakeClient initializes a fake k8sClient with preloaded Source objects
func setupFakeClient() client.Client {
	s := runtime.NewScheme()
	_ = AddToScheme(s)

	existingSourceList := &SourceList{
		Items: []Source{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bbc-source",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "BBC News",
					Link:      "https://feeds.bbci.co.uk/news/rss.xml",
					ShortName: "bbc",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "abc-source",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "ABC News",
					Link:      "https://abcnews.go.com/abcnews/internationalheadlines",
					ShortName: "abcnews",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fox-source",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "FOX NEWS",
					Link:      "https://moxie.foxnews.com/google-publisher/world.xml",
					ShortName: "foxnews",
				},
			},
		},
	}

	return fake.NewClientBuilder().WithScheme(s).WithLists(existingSourceList).Build()
}

func TestHotNews_Default(t *testing.T) {
	tests := []struct {
		name   string
		fields HotNews
		want   int
	}{
		{
			name: "Test Default TitlesCount",
			fields: HotNews{
				Spec: HotNewsSpec{
					SummaryConfig: SummaryConfig{},
				},
			},
			want: 10,
		},
		{
			name: "Test TitlesCount is set",
			fields: HotNews{
				Spec: HotNewsSpec{
					SummaryConfig: SummaryConfig{TitlesCount: 5},
				},
			},
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNews{
				Spec: tt.fields.Spec,
			}
			r.Default()
			if got := r.Spec.SummaryConfig.TitlesCount; got != tt.want {
				t.Errorf("Default() TitlesCount = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHotNews_ValidateCreate(t *testing.T) {
	k8sClient = setupFakeClient()

	tests := []struct {
		name    string
		fields  HotNews
		wantErr bool
	}{
		{
			name: "Test Validate Create - Valid Input",
			fields: HotNews{
				Spec: HotNewsSpec{
					Keywords: []string{"news", "hot"},
					Sources:  []string{"bbc"},
				},
			},
			wantErr: false,
		},
		{
			name: "Test Validate Create - Missing Keywords",
			fields: HotNews{
				Spec: HotNewsSpec{
					Sources: []string{"source1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Test Validate Create - Invalid Dates",
			fields: HotNews{
				Spec: HotNewsSpec{
					DateStart: "2024-12-31",
					DateEnd:   "2024-01-01",
					Keywords:  []string{"news"},
				},
			},
			wantErr: true,
		},
		{
			name: "Test Validate Create - Sources and FeedGroups Conflict",
			fields: HotNews{
				Spec: HotNewsSpec{
					Sources:    []string{"source1"},
					FeedGroups: []string{"group1"},
					Keywords:   []string{"news"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNews{
				Spec: tt.fields.Spec,
			}
			_, err := r.ValidateCreate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHotNews_ValidateUpdate(t *testing.T) {
	k8sClient = setupFakeClient()

	oldHotNews := &HotNews{
		Spec: HotNewsSpec{
			DateStart: "2024-01-01",
			DateEnd:   "2024-12-31",
			Keywords:  []string{"news"},
			Sources:   []string{"source1"},
		},
	}

	tests := []struct {
		name    string
		newSpec HotNewsSpec
		wantErr bool
	}{
		{
			name: "Test Validate Update - Valid Input",
			newSpec: HotNewsSpec{
				DateStart: "2024-01-01",
				DateEnd:   "2024-12-31",
				Keywords:  []string{"news"},
				Sources:   []string{"bbc"},
			},
			wantErr: false,
		},
		{
			name: "Test Validate Update - Invalid Date Format",
			newSpec: HotNewsSpec{
				DateStart: "invalid-date",
				DateEnd:   "2024-12-31",
				Keywords:  []string{"news"},
			},
			wantErr: true,
		},
		{
			name: "Test Validate Update - DateStart After DateEnd",
			newSpec: HotNewsSpec{
				DateStart: "2024-12-31",
				DateEnd:   "2024-01-01",
				Keywords:  []string{"news"},
			},
			wantErr: true,
		},
		{
			name: "Test Validate Update - Sources and FeedGroups Conflict",
			newSpec: HotNewsSpec{
				DateStart:  "2024-01-01",
				DateEnd:    "2024-12-31",
				Keywords:   []string{"news"},
				Sources:    []string{"source1"},
				FeedGroups: []string{"group1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNews{
				Spec: tt.newSpec,
			}
			_, err := r.ValidateUpdate(oldHotNews)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHotNews_ValidateDelete(t *testing.T) {
	k8sClient = setupFakeClient()

	tests := []struct {
		name    string
		fields  HotNews
		wantErr bool
	}{
		{
			name: "Test Validate Delete - No Error",
			fields: HotNews{
				Spec: HotNewsSpec{
					Keywords: []string{"news"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HotNews{
				Spec: tt.fields.Spec,
			}
			_, err := r.ValidateDelete()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
