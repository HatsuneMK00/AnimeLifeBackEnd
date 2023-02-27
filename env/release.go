//go:build release
// +build release

package env

var (
	Mode = "release"
	Etcd = etcd{
		Key:      "/config/release",
		Username: "xxx",
		Password: "xxxxxx",
	}
)

type etcd struct {
	Key      string
	Username string
	Password string
}
