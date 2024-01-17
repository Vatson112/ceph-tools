package internal

import "github.com/ceph/go-ceph/rados"

func NewCephClient(configPath string) (*rados.Conn, error) {
	client, err := rados.NewConn()
	if err != nil {
		return nil, err
	}
	err = client.ReadConfigFile(configPath)
	if err != nil {
		return nil, err
	}
	err = client.Connect()
	if err != nil {
		return nil, err
	}
	return client, nil
}
