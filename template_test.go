package fs

import (
	"reflect"
	"testing"
	"text/template"
)

func Test_fieldString(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				field: "{{ .Foo }}",
			},
			want: "Foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fieldString(tt.args.field); got != tt.want {
				t.Errorf("fieldString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listTemplateFields(t *testing.T) {
	type args struct {
		t *template.Template
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				t: template.Must(template.New("test").Parse("{{ .Foo }} {{ .Bar }} {{ .Baz }}")),
			},
			want: []string{"Foo", "Bar", "Baz"},
		},
		{
			name: "",
			args: args{
				t: template.Must(template.New("test").Parse("")),
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listTemplateFields(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listTemplateFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnippet_ListTemplateFields(t *testing.T) {
	type fields struct {
		Name    string
		Content string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "",
			fields: fields{
				Name:    "{{ .Foo }}.txt",
				Content: "{{ .Foo }} {{ .Bar }} {{ .Baz }}",
			},
			want: []string{"Foo", "Bar", "Baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snippet{
				Name:    tt.fields.Name,
				Content: tt.fields.Content,
			}
			if got := s.ListTemplateFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Snippet.ListTemplateFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
