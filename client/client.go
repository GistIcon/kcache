package client

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
)

type ListFn func(context.Context, api.ListOptions) (runtime.Object, error)
type WatchFn func(context.Context, api.ListOptions) (watch.Interface, error)

type Client interface {
	List(context.Context, api.ListOptions) (runtime.Object, error)
	Watch(context.Context, api.ListOptions) (watch.Interface, error)
}

type client struct {
	list  ListFn
	watch WatchFn
}

func NewClient(list ListFn, watch WatchFn) Client {
	return &client{list, watch}
}

func (c *client) List(ctx context.Context, opts api.ListOptions) (runtime.Object, error) {
	return c.list(ctx, opts)
}

func (c *client) Watch(ctx context.Context, opts api.ListOptions) (watch.Interface, error) {
	return c.watch(ctx, opts)
}

type restRequester interface {
	Get() *rest.Request
}

func ClientForResource(
	c restRequester, res string, ns string, fsel fields.Selector) Client {
	return NewClient(
		ListFnForResource(c, res, ns, fsel),
		WatchFnForResource(c, res, ns, fsel),
	)
}

func ListFnForResource(
	c restRequester, res string, ns string, fsel fields.Selector) ListFn {
	return func(ctx context.Context, opts api.ListOptions) (runtime.Object, error) {
		return c.Get().
			Context(ctx).
			Namespace(ns).
			Resource(res).
			VersionedParams(&opts, api.ParameterCodec).
			FieldsSelectorParam(fsel).
			Do().
			Get()
	}
}

func WatchFnForResource(
	c restRequester, res string, ns string, fsel fields.Selector) WatchFn {

	return func(ctx context.Context, opts api.ListOptions) (watch.Interface, error) {
		return c.Get().
			Context(ctx).
			Prefix("watch").
			Namespace(ns).
			Resource(res).
			VersionedParams(&opts, api.ParameterCodec).
			FieldsSelectorParam(fsel).
			Watch()
	}
}
