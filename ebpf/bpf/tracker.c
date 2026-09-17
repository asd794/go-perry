// +build ignore
#include<linux/bpf.h>
#include<bpf/bpf_helpers.h>


SEC("tracepoint/syscalls/sys_enter_execve")
int trace_execve(void *ctx) {
    unsigned long long pid_tgid = bpf_get_current_pid_tgid();
    unsigned int pid = pid_tgid >> 32;

    char comm[16] = {0};
    bpf_get_current_comm(&comm, sizeof(comm));

    char log_msg[] = "Process executed: PID=%d, Command=%s\n";

    bpf_trace_printk(log_msg, sizeof(log_msg), pid, comm);
    return 0;
}


char _license[] SEC("license") = "GPL";