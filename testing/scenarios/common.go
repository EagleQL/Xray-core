package scenarios

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/eagleql/xray-core/app/dispatcher"
	"github.com/eagleql/xray-core/app/proxyman"
	"github.com/eagleql/xray-core/common"
	"github.com/eagleql/xray-core/common/errors"
	"github.com/eagleql/xray-core/common/log"
	"github.com/eagleql/xray-core/common/net"
	"github.com/eagleql/xray-core/common/retry"
	"github.com/eagleql/xray-core/common/serial"
	core "github.com/eagleql/xray-core/core"
	"github.com/golang/protobuf/proto"
)

func xor(b []byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		r[i] = v ^ 'c'
	}
	return r
}

func readFrom(conn net.Conn, timeout time.Duration, length int) []byte {
	b := make([]byte, length)
	deadline := time.Now().Add(timeout)
	conn.SetReadDeadline(deadline)
	n, err := io.ReadFull(conn, b[:length])
	if err != nil {
		fmt.Println("Unexpected error from readFrom:", err)
	}
	return b[:n]
}

func readFrom2(conn net.Conn, timeout time.Duration, length int) ([]byte, error) {
	b := make([]byte, length)
	deadline := time.Now().Add(timeout)
	conn.SetReadDeadline(deadline)
	n, err := io.ReadFull(conn, b[:length])
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}

func InitializeServerConfigs(configs ...*core.Config) ([]*exec.Cmd, error) {
	servers := make([]*exec.Cmd, 0, 10)

	for _, config := range configs {
		server, err := InitializeServerConfig(config)
		if err != nil {
			CloseAllServers(servers)
			return nil, err
		}
		servers = append(servers, server)
	}

	time.Sleep(time.Second * 2)

	return servers, nil
}

func InitializeServerConfig(config *core.Config) (*exec.Cmd, error) {
	err := BuildXray()
	if err != nil {
		return nil, err
	}

	config = withDefaultApps(config)
	configBytes, err := proto.Marshal(config)
	if err != nil {
		return nil, err
	}
	proc := RunXrayProtobuf(configBytes)

	if err := proc.Start(); err != nil {
		return nil, err
	}

	return proc, nil
}

var (
	testBinaryPath    string
	testBinaryPathGen sync.Once
)

func genTestBinaryPath() {
	testBinaryPathGen.Do(func() {
		var tempDir string
		common.Must(retry.Timed(5, 100).On(func() error {
			dir, err := ioutil.TempDir("", "xray")
			if err != nil {
				return err
			}
			tempDir = dir
			return nil
		}))
		file := filepath.Join(tempDir, "xray.test")
		if runtime.GOOS == "windows" {
			file += ".exe"
		}
		testBinaryPath = file
		fmt.Printf("Generated binary path: %s\n", file)
	})
}

func GetSourcePath() string {
	return filepath.Join("github.com", "xtls", "xray-core", "main")
}

func CloseAllServers(servers []*exec.Cmd) {
	log.Record(&log.GeneralMessage{
		Severity: log.Severity_Info,
		Content:  "Closing all servers.",
	})
	for _, server := range servers {
		if runtime.GOOS == "windows" {
			server.Process.Kill()
		} else {
			server.Process.Signal(syscall.SIGTERM)
		}
	}
	for _, server := range servers {
		server.Process.Wait()
	}
	log.Record(&log.GeneralMessage{
		Severity: log.Severity_Info,
		Content:  "All server closed.",
	})
}

func withDefaultApps(config *core.Config) *core.Config {
	config.App = append(config.App, serial.ToTypedMessage(&dispatcher.Config{}))
	config.App = append(config.App, serial.ToTypedMessage(&proxyman.InboundConfig{}))
	config.App = append(config.App, serial.ToTypedMessage(&proxyman.OutboundConfig{}))
	return config
}

func testTCPConn(port net.Port, payloadSize int, timeout time.Duration) func() error {
	return func() error {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(port),
		})
		if err != nil {
			return err
		}
		defer conn.Close()

		return testTCPConn2(conn, payloadSize, timeout)()
	}
}

func testUDPConn(port net.Port, payloadSize int, timeout time.Duration) func() error {
	return func() error {
		conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(port),
		})
		if err != nil {
			return err
		}
		defer conn.Close()

		return testTCPConn2(conn, payloadSize, timeout)()
	}
}

func testTCPConn2(conn net.Conn, payloadSize int, timeout time.Duration) func() error {
	return func() error {
		payload := make([]byte, payloadSize)
		common.Must2(rand.Read(payload))

		nBytes, err := conn.Write(payload)
		if err != nil {
			return err
		}
		if nBytes != len(payload) {
			return errors.New("expect ", len(payload), " written, but actually ", nBytes)
		}

		response, err := readFrom2(conn, timeout, payloadSize)
		if err != nil {
			return err
		}
		_ = response

		if r := bytes.Compare(response, xor(payload)); r != 0 {
			return errors.New(r)
		}

		return nil
	}
}
