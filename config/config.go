package config

import (
	"bytes"
	"crypto/tls"

	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"

	"github.com/projecteru2/pistage/errors"
)

const defaultTemplate = `
eru_addr = "127.0.0.1:5001"
eru_deploy_name = "pistage"
eru_deploy_user = "root"

meta_timeout = "1m"
meta_type = "etcd"

etcd_prefix = "/pistage/v1"
etcd_endpoints = ["http://127.0.0.1:2379"]

lambda_cpu = 2
lambda_cpu_bind = false
lambda_memory = 4294967296
lambda_podname = "lambda_pod"
lambda_network = "lambda_network"
lambda_appname = "lambda"

log_level = "info"
graceful_timeout = "32s"
parse_dependencies_timeout = "1m"
bind_grpc_addr = "0.0.0.0:8697"
prof_http_port = 8999
`

// Conf .
var Conf Config

func init() {
	if _, err := toml.Decode(defaultTemplate, &Conf); err != nil {
		panic(err)
	}
}

// Config .
type Config struct {
	EruAddr       string `toml:"eru_addr"`
	EruDeployName string `toml:"eru_deploy_name"`
	EruDeployUser string `toml:"eru_deploy_user"`
	EruUsername   string `toml:"eru_username"`
	EruPassword   string `toml:"eru_password"`

	MetaTimeout Duration `toml:"meta_timeout"`
	MetaType    string   `toml:"meta_type"`

	EtcdPrefix    string   `toml:"etcd_prefix"`
	EtcdEndpoints []string `toml:"etcd_endpoints"`
	EtcdUsername  string   `toml:"etcd_username"`
	EtcdPassword  string   `toml:"etcd_password"`
	EtcdCA        string   `toml:"etcd_ca"`
	EtcdKey       string   `toml:"etcd_key"`
	EtcdCert      string   `toml:"etcd_cert"`

	LambdaCPU     int    `toml:"lambda_cpu"`
	LambdaCPUBind bool   `toml:"lambda_cpu_bind"`
	LambdaMemory  int64  `toml:"lambda_memory"`
	LambdaPodname string `toml:"lambda_podname"`
	LambdaNetwork string `toml:"lambda_network"`
	LambdaAppname string `toml:"lambda_appname"`

	LogLevel string `toml:"log_level"`
	LogFile  string `toml:"log_file"`

	BindGRPCAddr             string   `toml:"bind_grpc_addr"`
	GracefulTimeout          Duration `toml:"graceful_timeout"`
	ParseDependenciesTimeout Duration `toml:"parse_dependencies_timeout"`
	ProfHTTPPort             int      `toml:"prof_http_port"`
}

// Dump .
func (c *Config) Dump() (string, error) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)

	if err := enc.Encode(c); err != nil {
		return "", errors.Trace(err)
	}

	return buf.String(), nil
}

// ParseFile parses a config file.
func (c *Config) ParseFile(filepath string) (err error) {
	_, err = toml.DecodeFile(filepath, c)
	return
}

// ParseFiles parses a bunch of config files.
func (c *Config) ParseFiles(filepaths ...string) (err error) {
	for _, f := range filepaths {
		if err = c.ParseFile(f); err != nil {
			return
		}
	}
	return
}

// NewEtcdConfig .
func (c *Config) NewEtcdConfig() (etcdcnf clientv3.Config, err error) {
	etcdcnf.Endpoints = c.EtcdEndpoints
	etcdcnf.Username = c.EtcdUsername
	etcdcnf.Password = c.EtcdPassword
	etcdcnf.TLS, err = c.newEtcdTLSConfig()
	return
}

func (c *Config) newEtcdTLSConfig() (*tls.Config, error) {
	if len(c.EtcdCA) < 1 || len(c.EtcdKey) < 1 || len(c.EtcdCert) < 1 {
		return nil, nil
	}

	return transport.TLSInfo{
		TrustedCAFile: c.EtcdCA,
		KeyFile:       c.EtcdKey,
		CertFile:      c.EtcdCert,
	}.ClientConfig()
}
