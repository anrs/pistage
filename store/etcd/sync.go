package etcd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"

	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/sync"
)

// Mutex .
type Mutex struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
}

// NewMutex .
func NewMutex(cli *clientv3.Client, key string) (sync.Locker, error) {
	var sess, err = concurrency.NewSession(cli)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Mutex{
		mutex:   concurrency.NewMutex(sess, key),
		session: sess,
	}, nil
}

// Lock .
func (m *Mutex) Lock(ctx context.Context) (sync.Unlocker, error) {
	if err := m.mutex.Lock(ctx); err != nil {
		return nil, errors.Trace(err)
	}
	return m.Unlock, nil
}

// Unlock .
func (m *Mutex) Unlock(ctx context.Context) (err error) {
	defer func() {
		if e := m.session.Close(); e != nil {
			err = errors.Wrap(err, e)
		}
	}()
	return m.mutex.Unlock(ctx)
}
