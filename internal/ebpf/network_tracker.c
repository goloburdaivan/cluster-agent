// +build ignore

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>

char LICENSE[] SEC("license") = "GPL";

struct event_t {
    u32 pid;
    u32 saddr;
    u32 daddr;
    u16 sport;
    u16 dport;
    u32 old_state;
    u32 new_state;
    char comm[16];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 18); // 256 KB
} events SEC(".maps");

const struct event_t *unused __attribute__((unused));

SEC("tracepoint/sock/inet_sock_set_state")
int handle_set_state(struct trace_event_raw_inet_sock_set_state *ctx) {
    if (ctx->protocol != 6) {
        return 0;
    }

    if (ctx->newstate != 1) {
            return 0;
    }

    struct event_t *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 0;
    }

    e->pid = bpf_get_current_pid_tgid() >> 32;

    e->saddr = *((u32 *)ctx->saddr);
    e->daddr = *((u32 *)ctx->daddr);

    e->sport = bpf_ntohs(ctx->sport);
    e->dport = bpf_ntohs(ctx->dport);

    e->old_state = ctx->oldstate;
    e->new_state = ctx->newstate;

    bpf_get_current_comm(&e->comm, sizeof(e->comm));

    bpf_ringbuf_submit(e, 0);

    return 0;
}