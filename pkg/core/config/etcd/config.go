package etcd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	logs "github.com/bhojpur/logger/pkg/engine"
	"github.com/bhojpur/web/pkg/core/config"
)

type EtcdConfigure struct {
	prefix string
	client *clientv3.Client
	config.BaseConfigure
}

func newEtcdConfigure(client *clientv3.Client, prefix string) *EtcdConfigure {
	res := &EtcdConfigure{
		client: client,
		prefix: prefix,
	}

	res.BaseConfigure = config.NewBaseConfigure(res.reader)
	return res
}

// reader is an general implementation that read config from etcd.
func (e *EtcdConfigure) reader(ctx context.Context, key string) (string, error) {
	resp, err := get(e.client, e.prefix+key)
	if err != nil {
		return "", err
	}

	if resp.Count > 0 {
		return string(resp.Kvs[0].Value), nil
	}

	return "", nil
}

// Set do nothing and return an error
// I think write data to remote config center is not a good practice
func (e *EtcdConfigure) Set(key, val string) error {
	return errors.New("Unsupported operation")
}

// DIY return the original response from etcd
// be careful when you decide to use this
func (e *EtcdConfigure) DIY(key string) (interface{}, error) {
	return get(e.client, key)
}

// GetSection in this implementation, we use section as prefix
func (e *EtcdConfigure) GetSection(section string) (map[string]string, error) {
	var (
		resp *clientv3.GetResponse
		err  error
	)

	resp, err = e.client.Get(context.TODO(), e.prefix+section, clientv3.WithPrefix())

	if err != nil {
		return nil, errors.WithMessage(err, "GetSection failed")
	}
	res := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		res[string(kv.Key)] = string(kv.Value)
	}
	return res, nil
}

func (e *EtcdConfigure) SaveConfigFile(filename string) error {
	return errors.New("Unsupported operation")
}

// Unmarshaler is not very powerful because we lost the type information when we get configuration from etcd
// for example, when we got "5", we are not sure whether it's int 5, or it's string "5"
// TODO(support more complicated decoder)
func (e *EtcdConfigure) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	res, err := e.GetSection(prefix)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("could not read config with prefix: %s", prefix))
	}

	prefixLen := len(e.prefix + prefix)
	m := make(map[string]string, len(res))
	for k, v := range res {
		m[k[prefixLen:]] = v
	}
	return mapstructure.Decode(m, obj)
}

// Sub return an sub configure.
func (e *EtcdConfigure) Sub(key string) (config.Configure, error) {
	return newEtcdConfigure(e.client, e.prefix+key), nil
}

// TODO remove this before release v2.0.0
func (e *EtcdConfigure) OnChange(key string, fn func(value string)) {

	buildOptsFunc := func() []clientv3.OpOption {
		return []clientv3.OpOption{}
	}

	rch := e.client.Watch(context.Background(), e.prefix+key, buildOptsFunc()...)
	go func() {
		for {
			for resp := range rch {
				if err := resp.Err(); err != nil {
					logs.Error("listen to key but got error callback", err)
					break
				}

				for _, e := range resp.Events {
					if e.Kv == nil {
						continue
					}
					fn(string(e.Kv.Value))
				}
			}
			time.Sleep(time.Second)
			rch = e.client.Watch(context.Background(), e.prefix+key, buildOptsFunc()...)
		}
	}()

}

type EtcdConfigureProvider struct {
}

// Parse = ParseData([]byte(key))
// key must be json
func (provider *EtcdConfigureProvider) Parse(key string) (config.Configure, error) {
	return provider.ParseData([]byte(key))
}

// ParseData try to parse key as clientv3.Config, using this to build etcdClient
func (provider *EtcdConfigureProvider) ParseData(data []byte) (config.Configure, error) {
	cfg := &clientv3.Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "parse data to etcd config failed, please check your input")
	}

	cfg.DialOptions = []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	}
	client, err := clientv3.New(*cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "create etcd client failed")
	}

	return newEtcdConfigure(client, ""), nil
}

func get(client *clientv3.Client, key string) (*clientv3.GetResponse, error) {
	var (
		resp *clientv3.GetResponse
		err  error
	)
	resp, err = client.Get(context.Background(), key)

	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("read config from etcd with key %s failed", key))
	}
	return resp, err
}

func init() {
	config.Register("json", &EtcdConfigureProvider{})
}
