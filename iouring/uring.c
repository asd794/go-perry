#include "uring.h"

#include <liburing.h>

#include<arpa/inet.h>
#include<errno.h>
#include<netinet/in.h>
#include<stdlib.h>
#include<string.h>
#include<sys/socket.h>
#include<unistd.h>
#include <stdio.h>

struct connection {
    int fd;

    char buffer[BUFFER_SIZE];

    int len;
    int offset;
};

struct request {
    op_type op;
    connection_t *conn;
};

static struct io_uring ring;

static struct request *request_new(op_type op, connection_t *conn) {
    struct request *req = calloc(1, sizeof(*req));
    if (!req) {
        return NULL;
    }
    req->op = op;
    req->conn = conn;
    return req;
}

int server_socket(int port) {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) {
        perror("socket");
        return -1;
    }

    int opt = 1;
    if (setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) < 0) {
        perror("setsockopt");
        close(fd);
        return -1;
    }

    struct sockaddr_in addr = {
        .sin_family = AF_INET,
        .sin_addr.s_addr = htonl(INADDR_LOOPBACK),
        .sin_port = htons(port),
    };

    if (bind(fd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        perror("bind");
        close(fd);
        return -1;
    }

    if (listen(fd, SOMAXCONN) < 0) {
        perror("listen");
        close(fd);
        return -1;
    }

    return fd;
}
int uring_init(unsigned entries) {
    return io_uring_queue_init(entries, &ring, 0);
}

void uring_cleanup(void) {
    io_uring_queue_exit(&ring);
}

int submit_accept(int server_fd) {
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring); // get submission queue entry
    if (!sqe) {
        return -1;
    }

    struct request *req = request_new(OP_ACCEPT, NULL);
    if (!req) {
        return -1;
    }

    io_uring_prep_accept(sqe, server_fd, NULL, NULL, 0);
    io_uring_sqe_set_data(sqe, req);

    return io_uring_submit(&ring);
}

int submit_recv(connection_t *conn) {
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
    if (!sqe) {
        return -1;
    }

    struct request *req = request_new(OP_RECV, conn);
    if (!req) {
        return -1;
    }

    conn->len = 0;
    conn->offset = 0;

    io_uring_prep_recv(sqe, conn->fd, conn->buffer, BUFFER_SIZE, 0);
    io_uring_sqe_set_data(sqe, req);

    return io_uring_submit(&ring);
}

int submit_send(connection_t *conn) {
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
    if (!sqe) {
        fprintf(stderr, "submit_send: no SQE\n");
        return -1;
    }

    struct request *req = request_new(OP_SEND, conn);
    if (!req) {
        fprintf(stderr, "submit_send: request_new failed\n");
        return -1;
    }

    int remaining = conn->len - conn->offset;

    fprintf(
        stderr,
        "submit_send: fd=%d len=%d offset=%d remaining=%d\n",
        conn->fd,
        conn->len,
        conn->offset,
        remaining
    );

    if (remaining <= 0) {
        fprintf(stderr, "submit_send: no data to send\n");
        free(req);
        return -1;
    }

    io_uring_prep_send(
        sqe,
        conn->fd,
        conn->buffer + conn->offset,
        remaining,
        0
    );

    io_uring_sqe_set_data(sqe, req);

    int ret = io_uring_submit(&ring);

    fprintf(stderr, "submit_send: io_uring_submit=%d\n", ret);

    if (ret < 0) {
        fprintf(stderr, "submit_send: submit error=%s\n", strerror(-ret));
        free(req);
        return ret;
    }

    return ret;
}

int wait_cqe(completion_t *completion) {
    struct io_uring_cqe *cqe;
    int ret = io_uring_wait_cqe(&ring, &cqe);
    if (ret < 0) {
        return ret;
    }

    struct request *req = io_uring_cqe_get_data(cqe);
    completion->op = req->op;
    completion->conn = req->conn;
    completion->result = cqe->res;

    io_uring_cqe_seen(&ring, cqe);
    free(req);

    return 0;
}

connection_t *connection_new(int fd) {
    connection_t *conn = calloc(1, sizeof(*conn));
    if (!conn) {
        return NULL;
    }
    conn->fd = fd;
    return conn;
}

void connection_free(connection_t *conn) {
    if (conn) {
        close(conn->fd);
        free(conn);
    }
}

int connection_fd(connection_t *conn) {
    return conn->fd;
}

char *connection_buffer(connection_t *conn) {
    return conn->buffer;
}

int connection_len(connection_t *conn) {
    return conn->len;
}

void connection_set_len(connection_t *conn, int len) {
    conn->len = len;
    conn->offset = 0;
}

void connection_consume(connection_t *conn, int n) {
    if (n <= 0 || n > conn->len - conn->offset) {
        return;
    }
    conn->offset += n;
}

int connection_remaining(connection_t *conn) {
    return conn->len - conn->offset;
}