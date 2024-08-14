package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

// setupFakeClient initializes a fake k8sClient with preloaded Source objects
func setupFakeClient() client.Client {
	s := runtime.NewScheme()
	_ = AddToScheme(s)

	existingSourceList := &SourceList{
		Items: []Source{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "test",
					Link:      "https://test.com",
					ShortName: "test",
				},
			},
		},
	}

	return fake.NewClientBuilder().WithScheme(s).WithLists(existingSourceList).Build()
}

func TestSource_Default(t *testing.T) {
	tests := []struct {
		name   string
		fields Source
	}{
		{
			name: "Test Default",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			r.Default()
		})
	}
}

func TestSource_ValidateCreate(t *testing.T) {
	k8sClient = setupFakeClient()

	tests := []struct {
		name    string
		fields  Source
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Create",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
			},
			want: nil,
		},
		{
			name: "Test Validate Create Error",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "verybigandlongshortnametocreateerror",
					Link:      "https://test.com",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			got, err := r.ValidateCreate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateCreate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_ValidateDelete(t *testing.T) {
	k8sClient = setupFakeClient()

	tests := []struct {
		name    string
		fields  Source
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Delete",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			got, err := r.ValidateDelete()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateDelete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_ValidateUpdate(t *testing.T) {
	k8sClient = setupFakeClient()

	oldSource := &Source{
		Spec: SourceSpec{
			Name:      "test",
			ShortName: "test",
			Link:      "https://test.com",
		},
	}

	tests := []struct {
		name    string
		newSpec SourceSpec
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Update",
			newSpec: SourceSpec{
				Name:      "test",
				ShortName: "test",
				Link:      "https://test.com",
			},
			want: nil,
		},
		{
			name: "Test Validate Update Error",
			newSpec: SourceSpec{
				Name:      "test",
				ShortName: "verybigandlongshortnametocreateerror",
				Link:      "https://test.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: tt.newSpec,
			}
			got, err := r.ValidateUpdate(oldSource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_checkUniqueFields(t *testing.T) {
	k8sClient = setupFakeClient()

	tests := []struct {
		name    string
		fields  Source
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Check Unique Fields",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			got, err := r.checkUniqueFields()
			if (err != nil) != tt.wantErr {
				t.Errorf("checkUniqueFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkUniqueFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_validateSource(t *testing.T) {
	k8sClient = setupFakeClient()

	tests := []struct {
		name    string
		fields  Source
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Source",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
			},
			want: nil,
		},
		{
			name: "Test Validate Source Error",
			fields: Source{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "verybigandlongshortnametocreateerror",
					Link:      "https://test.com",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			got, err := r.validateSource()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidURL(t *testing.T) {
	tests := []struct {
		name string
		link string
		want bool
	}{
		{
			name: "Test Valid URL",
			link: "https://test.com",
			want: true,
		},
		{
			name: "Test Invalid URL",
			link: "test.com",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.link); got != tt.want {
				t.Errorf("isValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
