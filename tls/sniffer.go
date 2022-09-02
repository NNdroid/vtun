package tls

import (
	"bytes"
	"errors"
	"io"
	"net"
)

var defaultHttpResponse = []byte(`HTTP/1.1 200 OK
Server: nginx/1.18.0 (Ubuntu)
Content-Type: text/html
Content-Length: 612
Connection: keep-alive
Accept-Ranges: bytes

<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>`)

type SniffConn struct {
	net.Conn
	rout         io.Reader
	peeked, read bool
	Type         int
	preData      []byte
	path         string
}

const (
	TypeHttp = iota
	TypeHttp2
	TypeUnknown
)

var (
	httpMethods = [...][]byte{
		[]byte("GET"),
		[]byte("POST"),
		[]byte("HEAD"),
		[]byte("PUT"),
		[]byte("DELETE"),
		[]byte("OPTIONS"),
		[]byte("CONNECT"),
	}
	http2Header = []byte("PRI * HTTP/2.0")
	sep         = []byte(" ")
)

func NewPeekPreDataConn(c net.Conn) *SniffConn {
	s := &SniffConn{Conn: c, rout: c}
	s.Type = s.sniff()
	return s
}

func (c *SniffConn) peekPreData(n int) ([]byte, error) {
	if c.read {
		return nil, errors.New("pre-data must be peek before read")
	}
	if c.peeked {
		return nil, errors.New("can only peek once")
	}
	c.peeked = true
	preDate := make([]byte, n)
	n, err := c.Conn.Read(preDate)
	return preDate[:n], err
}

func (c *SniffConn) Read(p []byte) (int, error) {
	if !c.read {
		c.read = true
		c.rout = io.MultiReader(bytes.NewReader(c.preData), c.Conn)
	}
	return c.rout.Read(p)
}

func (c *SniffConn) sniff() int {
	var err error
	c.preData, err = c.peekPreData(64)
	if err != nil && err != io.EOF {
		return TypeUnknown
	}

	if c.sniffHttp() {
		return TypeHttp
	}

	if c.sniffHttp2() {
		return TypeHttp2
	}

	return TypeUnknown
}

func (c *SniffConn) sniffHttp() bool {
	preDataParts := bytes.Split(c.preData, sep)

	if len(preDataParts) < 2 {
		return false
	}

	for _, m := range httpMethods {
		if bytes.Compare(preDataParts[0], m) == 0 {
			c.path = string(preDataParts[1])
			return true
		}
	}
	return false
}

func (c *SniffConn) sniffHttp2() bool {

	return len(c.preData) >= len(http2Header) &&
		bytes.Compare(c.preData[:len(http2Header)], http2Header) == 0
}

func (c *SniffConn) SetPath(path string) {
	preDataParts := bytes.Split(c.preData, sep)
	preDataParts[1] = []byte(path)
	c.preData = bytes.Join(preDataParts, sep)
	c.path = path
}

func (c *SniffConn) GetPath() string {
	return c.path
}

func (c *SniffConn) Handle() bool {
	write, err := c.Write(defaultHttpResponse)
	if err != nil {
		return false
	}
	return write > 0
}
