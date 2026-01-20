package services

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type TCPState uint8

const (
	TCP_ESTABLISHED TCPState = 1
	TCP_SYN_SENT    TCPState = 2
	TCP_SYN_RECV    TCPState = 3
	TCP_FIN_WAIT1   TCPState = 4
	TCP_FIN_WAIT2   TCPState = 5
	TCP_TIME_WAIT   TCPState = 6
	TCP_CLOSE       TCPState = 7
	TCP_CLOSE_WAIT  TCPState = 8
	TCP_LAST_ACK    TCPState = 9
	TCP_LISTEN      TCPState = 10
	TCP_CLOSING     TCPState = 11
)

var stateNames = map[TCPState]string{
	TCP_ESTABLISHED: "ESTABLISHED",
	TCP_SYN_SENT:    "SYN_SENT",
	TCP_SYN_RECV:    "SYN_RECV",
	TCP_FIN_WAIT1:   "FIN_WAIT1",
	TCP_FIN_WAIT2:   "FIN_WAIT2",
	TCP_TIME_WAIT:   "TIME_WAIT",
	TCP_CLOSE:       "CLOSE",
	TCP_CLOSE_WAIT:  "CLOSE_WAIT",
	TCP_LAST_ACK:    "LAST_ACK",
	TCP_LISTEN:      "LISTEN",
	TCP_CLOSING:     "CLOSING",
}

func (s TCPState) String() string {
	if name, ok := stateNames[s]; ok {
		return name
	}
	return "UNKNOWN"
}

type TCPSocketEntry struct {
	LocalAddress  string `json:"local_address"`
	LocalPort     uint16 `json:"local_port"`
	RemoteAddress string `json:"remote_address"`
	RemotePort    uint16 `json:"remote_port"`
	State         string `json:"state"`
	Protocol      string `json:"protocol"`
}

type NetworkInspectorService interface {
	GetPodNetworkConnections(ctx context.Context, namespace, podName, container string) ([]TCPSocketEntry, error)
}

type networkInspectorService struct {
	clientset kubernetes.Interface
	config    *rest.Config
}

func NewNetworkInspectorService(clientset kubernetes.Interface, config *rest.Config) NetworkInspectorService {
	return &networkInspectorService{
		clientset: clientset,
		config:    config,
	}
}

func (s *networkInspectorService) GetPodNetworkConnections(ctx context.Context, namespace, podName, container string) ([]TCPSocketEntry, error) {
	cmd := []string{
		"/bin/sh",
		"-c",
		"cat /proc/net/tcp /proc/net/tcp6 2>/dev/null",
	}

	req := s.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container).
		Param("stdin", "false").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	for _, c := range cmd {
		req.Param("command", c)
	}

	executor, err := remotecommand.NewSPDYExecutor(s.config, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	pr, pw := io.Pipe()
	var errBuf strings.Builder

	go func() {
		err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: pw,
			Stderr: &errBuf,
		})
		pw.CloseWithError(err)
	}()

	return s.parseNetworkStream(pr)
}

func (s *networkInspectorService) parseNetworkStream(r io.Reader) ([]TCPSocketEntry, error) {
	var entries []TCPSocketEntry
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "sl") {
			continue
		}

		entry, err := s.parseLine(line)
		if err == nil {
			entries = append(entries, *entry)
		}
	}

	return entries, nil
}

func (s *networkInspectorService) parseLine(line string) (*TCPSocketEntry, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil, fmt.Errorf("invalid line")
	}

	localIP, localPort, err := s.parseIPPort(fields[1])
	if err != nil {
		return nil, err
	}

	remoteIP, remotePort, err := s.parseIPPort(fields[2])
	if err != nil {
		return nil, err
	}

	stateHex, err := strconv.ParseUint(fields[3], 16, 8)
	if err != nil {
		stateHex = 0
	}

	protocol := "tcp"
	if localIP.To4() == nil {
		protocol = "tcp6"
	}

	return &TCPSocketEntry{
		LocalAddress:  localIP.String(),
		LocalPort:     localPort,
		RemoteAddress: remoteIP.String(),
		RemotePort:    remotePort,
		State:         TCPState(stateHex).String(),
		Protocol:      protocol,
	}, nil
}

func (s *networkInspectorService) parseIPPort(field string) (net.IP, uint16, error) {
	parts := strings.Split(field, ":")
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid format")
	}

	ipHex := parts[0]
	portHex := parts[1]

	ipBytes, err := hex.DecodeString(ipHex)
	if err != nil {
		return nil, 0, err
	}

	for i := 0; i < len(ipBytes); i += 4 {
		if i+3 < len(ipBytes) {
			ipBytes[i], ipBytes[i+1], ipBytes[i+2], ipBytes[i+3] =
				ipBytes[i+3], ipBytes[i+2], ipBytes[i+1], ipBytes[i]
		}
	}

	port, err := strconv.ParseUint(portHex, 16, 16)
	if err != nil {
		return nil, 0, err
	}

	return net.IP(ipBytes), uint16(port), nil
}
