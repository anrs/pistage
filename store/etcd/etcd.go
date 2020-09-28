package etcd

import (
	"context"
	"strconv"
	"sync"

	"github.com/coreos/etcd/clientv3"

	"github.com/projecteru2/aa/codec/json"
	"github.com/projecteru2/aa/config"
	"github.com/projecteru2/aa/errors"
	"github.com/projecteru2/aa/log"
	sync2 "github.com/projecteru2/aa/sync"
)

// Etcd .
type Etcd struct {
	sync.Mutex
	cli *clientv3.Client
}

// New .
func New() (*Etcd, error) {
	etcdcnf, err := config.Conf.NewEtcdConfig()
	if err != nil {
		return nil, errors.Trace(err)
	}

	cli, err := clientv3.New(etcdcnf)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Etcd{cli: cli}, nil
}

// IncrUint32 .
func (e *Etcd) IncrUint32(ctx context.Context, key string) (n uint32, err error) {
	var mutex sync2.Locker
	if mutex, err = e.NewMutex(key); err != nil {
		return
	}

	var unlock sync2.Unlocker
	if unlock, err = mutex.Lock(ctx); err != nil {
		return
	}
	defer func() {
		if ue := unlock(ctx); ue != nil {
			err = errors.Wrap(err, ue)
		}
	}()

	var data = map[string]string{}
	var ver int64

	switch ver, err = e.Get(ctx, key, &n); {
	case errors.Contain(err, errors.ErrKeyNotExists):
		data[key] = "1"
		if err = e.Create(ctx, data); err != nil {
			return
		}
		return 1, nil

	case err != nil:
		return
	}

	n++
	data[key] = strconv.FormatInt(int64(n), 10)
	err = e.Update(ctx, data, map[string]int64{key: ver})

	return // nolint
}

// Create .
func (e *Etcd) Create(ctx context.Context, data map[string]string, opts ...clientv3.OpOption) error {
	var ev = newTxnEvent()
	ev.data = data
	ev.opts = opts
	ev.txnErr = errors.ErrKeyExists
	ev.vers = map[string]int64{}

	for k := range ev.data {
		ev.vers[k] = 0
	}

	return e.batchPut(ctx, ev)
}

// Update .
func (e *Etcd) Update(ctx context.Context, data map[string]string, vers map[string]int64, opts ...clientv3.OpOption) error {
	var ev = newTxnEvent()
	ev.data = data
	ev.opts = opts
	ev.txnErr = errors.ErrKeyBadVersion
	ev.vers = vers

	return e.batchPut(ctx, ev)
}

func (e *Etcd) batchPut(ctx context.Context, ev *txnEvent) error {
	var ops, cmps = ev.generate()

	switch succ, err := e.BatchOperate(ctx, ops, cmps...); {
	case err != nil:
		return errors.Trace(err)

	case !succ:
		return ev.txnErr
	}

	return nil
}

// Delete .
func (e *Etcd) Delete(ctx context.Context, keys []string, vers map[string]int64, opts ...clientv3.OpOption) error {
	var ev = newDelTxnEvent(keys, vers, opts...)
	var ops, cmps = ev.generate()

	switch succ, err := e.BatchOperate(ctx, ops, cmps...); {
	case err != nil:
		return errors.Trace(err)

	case !succ:
		return errors.Trace(errors.ErrKeyBadVersion)
	}

	return nil
}

// BatchOperate .
func (e *Etcd) BatchOperate(ctx context.Context, ops []clientv3.Op, cmps ...clientv3.Cmp) (bool, error) {
	e.Lock()
	defer e.Unlock()

	var txn = e.cli.Txn(ctx)
	var resp, err = txn.If(cmps...).Then(ops...).Commit()
	if err != nil {
		return false, errors.Trace(err)
	}

	return resp.Succeeded, nil
}

// GetPrefix .
func (e *Etcd) GetPrefix(ctx context.Context, prefix string, limit int64) (map[string][]byte, map[string]int64, error) {
	e.Lock()
	defer e.Unlock()

	var resp, err = e.cli.Get(ctx, prefix, clientv3.WithLimit(limit), clientv3.WithPrefix())
	switch {
	case err != nil:
		return nil, nil, errors.Trace(err)
	case resp.Count < 1:
		return nil, nil, errors.Annotatef(errors.ErrKeyNotExists, prefix)
	}

	var data = map[string][]byte{}
	var vers = map[string]int64{}

	for _, kv := range resp.Kvs {
		var key = string(kv.Key)
		data[key] = kv.Value
		vers[key] = kv.Version
	}

	return data, vers, nil
}

// Exists .
func (e *Etcd) Exists(ctx context.Context, keys []string) (map[string]bool, error) {
	var exists = map[string]bool{}

	for _, k := range keys {
		var resp, err = e.cli.Get(ctx, k, clientv3.WithKeysOnly())
		if err != nil {
			return nil, errors.Trace(err)
		}
		exists[k] = resp.Count > 0
	}

	return exists, nil
}

// Get .
func (e *Etcd) Get(ctx context.Context, key string, obj interface{}, opts ...clientv3.OpOption) (int64, error) {
	e.Lock()
	defer e.Unlock()

	switch resp, err := e.cli.Get(ctx, key, opts...); {
	case err != nil:
		return 0, errors.Trace(err)

	case resp.Count != 1:
		return 0, errors.Annotatef(errors.ErrKeyNotExists, key)

	default:
		return resp.Kvs[0].Version, json.Decode(resp.Kvs[0].Value, obj)
	}
}

// NewMutex .
func (e *Etcd) NewMutex(key string) (sync2.Locker, error) {
	return NewMutex(e.cli, key)
}

// Close .
func (e *Etcd) Close() error {
	e.Lock()
	defer e.Unlock()
	return e.cli.Close()
}

// RetryTimedOut .
func RetryTimedOut(fn func() error, retryTimes int) error {
	for retried := 0; ; retried++ {
		if err := fn(); err != nil {
			if retried < retryTimes && IsETCDServerTimedOutErr(err) {
				log.Warnf("etcdserver: request timed out, retry it")
				continue
			}

			return err
		}

		return nil
	}
}

// IsETCDServerTimedOutErr .
func IsETCDServerTimedOutErr(err error) bool {
	return err.Error() == "etcdserver: request timed out"
}
