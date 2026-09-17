#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

SEC("xdp")
int packet_filter(struct xdp_md *ctx)
{
    void *data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;

    struct ethhdr *eth = data;

    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;

    bpf_printk("XDP: ethernet packet");

    if (bpf_ntohs(eth->h_proto) == ETH_P_IP)
        bpf_printk("XDP: IPv4 packet");

    return XDP_PASS;
}

char LICENSE[] SEC("license") = "GPL";