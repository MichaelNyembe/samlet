package controllers

import (
	"context"

	samletv1 "github.com/bison-cloud-platform/samlet/api/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	userKey = "username"
	passKey = "password"
)

func (r *Saml2AwsReconciler) targetSecret(saml *samletv1.Saml2Aws, data []byte) (*v1.Secret, error) {
	secretData := map[string][]byte{
		"credentials": data,
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saml.Spec.TargetSecretName,
			Namespace: saml.Namespace,
		},
		Data: secretData,
	}

	err := controllerutil.SetControllerReference(saml, secret, r.Scheme)
	if err != nil {
		return nil, err
	}
	return secret, err
}

func getLoginData(secret *v1.Secret) (string, string) {
	user := string(secret.Data[userKey])
	pass := string(secret.Data[passKey])
	return user, pass
}

func (r *Saml2AwsReconciler) readSecret(name, namespace string) (*v1.Secret, error) {
	loginSecret := &v1.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, loginSecret)
	if err != nil {
		log.Error(err, "Failed to read secret")
		return nil, err
	}
	return loginSecret, nil
}
