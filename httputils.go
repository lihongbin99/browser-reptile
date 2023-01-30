package main

import (
	"fmt"
	"net"
	"strings"
)

func ReadRequest(conn net.Conn) (*Req, error) {
	method, err := ReadByStop(conn, []byte{' '})
	if err != nil {
		return nil, err
	}
	path, err := ReadByStop(conn, []byte{' '})
	if err != nil {
		return nil, err
	}
	proto, err := ReadByStop(conn, []byte{'\r', '\n'})
	if err != nil {
		return nil, err
	}
	header, err := ReadHeader(conn)
	if err != nil {
		return nil, err
	}
	var body []byte = nil
	// TODO body
	return &Req{
		Method: string(method),
		Path:   string(path),
		Proto:  string(proto),
		Header: header,
		Body:   body,
	}, nil
}

func SendRequest(req *Req, conn net.Conn) error {
	if _, err := conn.Write([]byte(req.Method)); err != nil {
		return err
	}
	if _, err := conn.Write([]byte{' '}); err != nil {
		return err
	}
	if _, err := conn.Write([]byte(req.Path)); err != nil {
		return err
	}
	if _, err := conn.Write([]byte{' '}); err != nil {
		return err
	}
	if _, err := conn.Write([]byte(req.Proto)); err != nil {
		return err
	}
	if _, err := conn.Write([]byte{'\r', '\n'}); err != nil {
		return err
	}
	for k, vl := range req.Header {
		for _, v := range vl {
			if _, err := conn.Write([]byte(k)); err != nil {
				return err
			}
			if _, err := conn.Write([]byte{':', ' '}); err != nil {
				return err
			}
			if _, err := conn.Write([]byte(v)); err != nil {
				return err
			}
			if _, err := conn.Write([]byte{'\r', '\n'}); err != nil {
				return err
			}
		}
	}
	if _, err := conn.Write([]byte{'\r', '\n'}); err != nil {
		return err
	}
	if req.Body != nil {
		// TODO body 协议
		if _, err := conn.Write(req.Body); err != nil {
			return err
		}
	}
	return nil
}

func ReadResponse(conn net.Conn) (*Rep, error) {
	proto, err := ReadByStop(conn, []byte{' '})
	if err != nil {
		return nil, err
	}
	status, err := ReadByStop(conn, []byte{'\r', '\n'})
	if err != nil {
		return nil, err
	}
	header, err := ReadHeader(conn)
	if err != nil {
		return nil, err
	}
	var body []byte = nil
	//
	return &Rep{
		Proto:  string(proto),
		Status: string(status),
		Header: header,
		Body:   body,
	}, nil
}

func SendResponse(rep *Rep, conn net.Conn) error {
	return fmt.Errorf("not supported")
}

func ReadHeader(conn net.Conn) (map[string][]string, error) {
	header := make(map[string][]string)
	for {
		hb, err := ReadByStop(conn, []byte{'\r', '\n'})
		if err != nil {
			return header, err
		}
		if len(hb) == 0 {
			break
		}
		index := strings.Index(string(hb), ": ")
		name := string(hb[:index])
		value := string(hb[index+2:])
		if _, exist := header[name]; exist {
			header[name] = append(header[name], value)
		} else {
			header[name] = []string{value}
		}
	}
	return header, nil
}

func ReadByStop(conn net.Conn, stop []byte) ([]byte, error) {
	data := make([]byte, 0)
	var buf [1]byte
	for {
		_, err := conn.Read(buf[:])
		if err != nil {
			return data, err
		}

		data = append(data, buf[0])

		if len(data) >= len(stop) {
			ret := true
			for i := 1; i <= len(stop); i++ {
				if data[len(data)-i] != stop[len(stop)-i] {
					ret = false
					break
				}
			}

			if ret {
				return data[:len(data)-len(stop)], nil
			}
		}
	}
}
