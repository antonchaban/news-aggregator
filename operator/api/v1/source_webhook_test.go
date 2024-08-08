package v1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

func TestSource_Default(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       SourceSpec
		Status     SourceStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test Default",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
				Status: SourceStatus{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			r.Default()
		})
	}
}

func TestSource_ValidateCreate(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       SourceSpec
		Status     SourceStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Create",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
				Status: SourceStatus{},
			},
			want: nil,
		},
		{
			name: "Test Validate Create Error",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "verybigandlongshortnametocreateerror",
					Link:      "https://test.com",
				},
				Status: SourceStatus{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
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
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       SourceSpec
		Status     SourceStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Update",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
				Status: SourceStatus{
					ID: 1,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
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
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       SourceSpec
		Status     SourceStatus
	}
	type args struct {
		old runtime.Object
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Update",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
				Status: SourceStatus{
					ID: 1,
				},
			},
			want: nil,
		},
		{
			name: "Test Validate Update Error",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "verybigandlongshortnametocreateerror",
					Link:      "https://test.com",
				},
				Status: SourceStatus{
					ID: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got, err := r.ValidateUpdate(tt.args.old)
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
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       SourceSpec
		Status     SourceStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Check Unique Fields",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
				Status: SourceStatus{
					ID: 1,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
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
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       SourceSpec
		Status     SourceStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "Test Validate Source",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "test",
					Link:      "https://test.com",
				},
				Status: SourceStatus{
					ID: 1,
				},
			},
			want: nil,
		},
		{
			name: "Test Validate Source Error",
			fields: fields{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
				Spec: SourceSpec{
					Name:      "test",
					ShortName: "verybigandlongshortnametocreateerror",
					Link:      "https://test.com",
				},
				Status: SourceStatus{
					ID: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Source{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
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
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test Valid URL",
			args: args{
				link: "https://test.com",
			},
			want: true,
		},
		{
			name: "Test Invalid URL",
			args: args{
				link: "test.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.args.link); got != tt.want {
				t.Errorf("isValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
