package internal

import "github.com/ceph/go-ceph/rados"

// Implement rados connection interface.
type RadosConn struct {
	CephConfig string
	Client     *rados.Conn
}

// NewRadosConnection return RadosConn interface.
func NewRadosConnection(cephConfig string) *RadosConn {
	return &RadosConn{
		CephConfig: cephConfig,
	}
}

func (c *RadosConn) Connect() error {
	var err error
	c.Client, err = NewCephClient(c.CephConfig)
	return err
}
