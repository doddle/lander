package endpoints

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_isAnnotatedForLander(t *testing.T) {
	type args struct {
		ingress    v1beta1.Ingress
		annotation string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true for annotation 1",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander": "false",
							"foo.acmecorp.org/lander":    "true",
						},
					},
				},
				annotation: "foo.acmecorp.org/lander",
			},
			want: true,
		},
		{
			name: "should return false for annotation 1",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander": "true",
							"foo.acmecorp.org/lander":    "false",
						},
					},
				},
				annotation: "foo.acmecorp.org/lander",
			},
			want: false,
		},
		{
			name: "should return true for annotation 1",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander": "false",
							"foo.acmecorp.org/lander":    "true",
						},
					},
				},
				annotation: "foo.acmecorp.org/notLander",
			},
			want: false,
		},
		{
			name: "should return false for annotation 1",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander": "true",
							"foo.acmecorp.org/lander":    "false",
						},
					},
				},
				annotation: "foo.acmecorp.org/notLander",
			},
			want: true,
		},
		{
			name: "false for non-existant annotation",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander": "true",
							"foo.acmecorp.org/lander":    "false",
						},
					},
				},
				annotation: "foo.acmecorp.org/doesNotExist",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isAnnotatedForLander(tt.args.ingress, tt.args.annotation))
		})
	}
}

func Test_annotationKeyExists(t *testing.T) {
	type args struct {
		ingress v1beta1.Ingress
		key     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "annotation exists",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"annotation1": "I exist",
						},
					},
				},
				key: "annotation1",
			},
			want: true,
		},
		{
			name: "annotation exists",
			args: args{
				ingress: v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"annotation1": "I exist",
						},
					},
				},
				key: "annotation2",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, annotationKeyExists(tt.args.ingress, tt.args.key))
		})
	}
}
