package config

import (
	"context"
	"fmt"
)

var namespaceCtxKey = NamespaceCtxKey("namespace")

type NamespaceCtxKey string

// GetNamespaceFromCtx returns the namespace extracted from the context.
func GetNamespaceFromCtx(ctx context.Context) (string, error) {
	anyNamespace := ctx.Value(namespaceCtxKey)
	if anyNamespace == nil {
		return "", fmt.Errorf("no namespace found in context")
	}

	namespace, ok := anyNamespace.(string)
	if !ok {
		return "", fmt.Errorf("could not cast namespace to string")
	}

	return namespace, nil
}

// MustGetNamespaceFromCtx only returns the namespace extracted from the context.
// It panics if the namespace cannot be extracted.
func MustGetNamespaceFromCtx(ctx context.Context) string {
	namespace, err := GetNamespaceFromCtx(ctx)
	if err != nil {
		panic(err)
	}

	return namespace
}

// AddNamespaceToCtx adds the namespace to the context. The namespace is stored
// in the key `namespace` and can be retrieved with `GetNamespaceFromCtx` or
// `MustGetNamespaceFromCtx`.
func AddNamespaceToCtx(ctx context.Context, namespace string) context.Context {
	ctx = context.WithValue(ctx, namespaceCtxKey, namespace)

	return ctx
}
