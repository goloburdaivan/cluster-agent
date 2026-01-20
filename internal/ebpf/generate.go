package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -cc clang -cflags "-I./headers" NetworkTracker network_tracker.c
