// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package user

import (
	"strings"
	"testing"

	"github.com/elastic/cloud-on-k8s/pkg/utils/k8s"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetSecret gets the first secret in a list that matches the namespace and the name.
func GetSecret(list corev1.SecretList, namespacedName types.NamespacedName) *corev1.Secret {
	for _, secret := range list.Items {
		if secret.Namespace == namespacedName.Namespace && secret.Name == namespacedName.Name {
			return &secret
		}
	}
	return nil
}

// ChecksUser checks that a secret contains the required fields expected by the user reconciler.
func ChecksUser(t *testing.T, secret *corev1.Secret, expectedUsername string, expectedRoles []string) {
	assert.NotNil(t, secret)
	currentUsername, ok := secret.Data["name"]
	assert.True(t, ok)
	assert.Equal(t, expectedUsername, string(currentUsername))
	passwordHash, ok := secret.Data["passwordHash"]
	assert.True(t, ok)
	assert.NotEmpty(t, passwordHash)
	currentRoles, ok := secret.Data["userRoles"]
	assert.True(t, ok)
	assert.ElementsMatch(t, expectedRoles, strings.Split(string(currentRoles), ","))
}

// DeleteUser deletes the user Secrets using the provided label selector.
func DeleteUser(c k8s.Client, opts ...client.ListOption) error {
	var secrets corev1.SecretList
	if err := c.List(&secrets, opts...); err != nil {
		return err
	}
	for _, s := range secrets.Items {
		if err := c.Delete(&s); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}
