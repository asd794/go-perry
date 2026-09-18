#ifndef URING_H
#define URING_H

#include <stdint.h>

#define BUFFER_SIZE 4096

typedef enum {
    OP_ACCEPT = 0,
    OP_RECV,
    OP_SEND
} op_type;

typedef struct connection connection_t;

typedef struct {
    op_type op;
    connection_t *conn;
    int result;
} completion_t;


int server_socket(int port);

int uring_init(unsigned entries);
void uring_cleanup(void);


int submit_accept(int listen_fd);
int submit_recv(connection_t *conn);
int submit_send(connection_t *conn);

int wait_cqe(completion_t *completion);

connection_t *connection_new(int fd);
void connection_free(connection_t *conn);

int connection_fd(connection_t *conn);
char *connection_buffer(connection_t *conn);

int connection_len(connection_t *conn);
void connection_set_len(connection_t *conn, int len);

void connection_consume(connection_t *conn, int n);
int connection_remaining(connection_t *conn);
#endif 