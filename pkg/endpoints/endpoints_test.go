package endpoints

import (
	"testing"

	"github.com/doddle/lander/pkg/util"
	"github.com/stretchr/testify/assert"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_isAnnotatedForLander(t *testing.T) {
	type args struct {
		ingress    networkingv1.Ingress
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
				ingress: networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander":   "false",
							"foo.acmecorp.org/lander/show": "true",
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
				ingress: networkingv1.Ingress{
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
				ingress: networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander":   "false",
							"foo.acmecorp.org/lander/show": "true",
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
				ingress: networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"foo.acmecorp.org/notLander/show": "true",
							"foo.acmecorp.org/lander/show":    "false",
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
				ingress: networkingv1.Ingress{
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
		ingress networkingv1.Ingress
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
				ingress: networkingv1.Ingress{
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
				ingress: networkingv1.Ingress{
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

type WarnImp struct {
	util.LoggerIFace
	warnCount int
}

func (m *WarnImp) Warnf(_ string, _ ...interface{}) {
	m.warnCount++
}

func Test_getIngressClass(t *testing.T) {
	type args struct {
		logger  *WarnImp
		ingress networkingv1.Ingress
	}
	ingressClassName := "testClass"
	tests := []struct {
		name      string
		args      args
		want      string
		warnCount int
	}{
		{
			name: "should return value",
			args: args{
				logger: new(WarnImp),
				ingress: networkingv1.Ingress{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec: networkingv1.IngressSpec{
						IngressClassName: &ingressClassName,
					},
					Status: networkingv1.IngressStatus{},
				},
			},
			want: "testClass",
		},
		{
			name: "should warn",
			args: args{
				logger: new(WarnImp),
				ingress: networkingv1.Ingress{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{},
					},
					Spec:   networkingv1.IngressSpec{},
					Status: networkingv1.IngressStatus{},
				},
			},
			want:      "",
			warnCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getIngressClass(tt.args.logger, tt.args.ingress))
			assert.Equal(t, tt.warnCount, tt.args.logger.warnCount)
		})
	}
}
