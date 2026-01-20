package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cluster-agent/internal/ebpf"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

var tcpStates = map[uint32]string{
	1:  "ESTABLISHED",
	2:  "SYN_SENT",
	3:  "SYN_RECV",
	4:  "FIN_WAIT1",
	5:  "FIN_WAIT2",
	6:  "TIME_WAIT",
	7:  "CLOSE",
	8:  "CLOSE_WAIT",
	9:  "LAST_ACK",
	10: "LISTEN",
	11: "CLOSING",
	12: "NEW_SYN_RECV",
}

type bpfEvent struct {
	Pid      uint32
	Saddr    uint32
	Daddr    uint32
	Sport    uint16
	Dport    uint16
	OldState uint32
	NewState uint32
	Comm     [16]byte
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("failed to remove memlock: %v", err)
	}

	objs := ebpf.NetworkTrackerObjects{}
	if err := ebpf.LoadNetworkTrackerObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	tp, err := link.Tracepoint("sock", "inet_sock_set_state", objs.HandleSetState, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %v", err)
	}
	defer tp.Close()

	rd, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("creating ringbuf reader: %v", err)
	}
	defer rd.Close()

	log.Println("Cluster Agent started. Waiting for TCP events...")
	log.Println("Try running 'curl google.com' in another terminal.")

	go func() {
		<-ctx.Done()
		rd.Close()
	}()

	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("Ring buffer closed, exiting...")
				return
			}
			log.Printf("error reading from ringbuf: %v", err)
			continue
		}

		var event bpfEvent
		if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %v", err)
			continue
		}

		comm := string(bytes.TrimRight(event.Comm[:], "\x00"))
		srcIP := intToIP(event.Saddr)
		dstIP := intToIP(event.Daddr)
		oldSt := getStateName(event.OldState)
		newSt := getStateName(event.NewState)

		fmt.Printf("[%s] PID: %d | %s:%d -> %s:%d | %s -> %s\n",
			comm, event.Pid,
			srcIP, event.Sport,
			dstIP, event.Dport,
			oldSt, newSt)
	}
}

func intToIP(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.LittleEndian.PutUint32(ip, nn)
	return ip
}

func getStateName(state uint32) string {
	if name, ok := tcpStates[state]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN(%d)", state)
}
