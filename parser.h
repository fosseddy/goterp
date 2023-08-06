// #include "scanner.h"

struct parser {
    struct token tok;
    struct scanner s;
};

enum expr_kind {
    EXPR_LIT,
    EXPR_BINARY
};

struct expr {
    enum expr_kind kind;
    void *body;
};

struct expr_lit {
    struct token token;
};

struct expr_binary {
    struct expr x;
    enum token_kind op;
    struct expr y;
};

struct expr_array {
    int len;
    int cap;
    int data_size;
    struct expr *buf;
};

enum stmt_kind {
    STMT_PRINT
};

struct stmt {
    enum stmt_kind kind;
    void *body;
};

struct stmt_print {
    struct expr value;
};

struct stmt_array {
    int len;
    int cap;
    int data_size;
    struct stmt *buf;
};

void parse(struct parser *p, struct stmt_array *arr);
void init_parser(struct parser *p, char *filepath);
