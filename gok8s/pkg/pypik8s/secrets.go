package pypik8s

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Secrets(name string, secretData map[string]string) *apiv1.Secret{
	return &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		StringData: secretData,
		Type: apiv1.SecretType("Opaque"),
	}
}