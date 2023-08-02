// #include "scanner.h"

struct parser {
    struct token tok;
    struct scanner s;
};

enum expr_kind {
    EXPR_ERR,
    EXPR_LIT
};

struct expr_lit {
    enum token_kind kind;
    char *value;
};

struct expr {
    enum expr_kind kind;
    void *body;
};

struct expr_array {
    int len;
    int cap;
    int data_size;
    struct expr *buf;
};

enum stmt_kind {
    STMT_ERR,
    STMT_PRINT
};

struct stmt_print {
    struct expr value;
};

struct stmt {
    enum stmt_kind kind;
    void *body;
};

struct stmt_array {
    int len;
    int cap;
    int data_size;
    struct stmt *buf;
};

void parse(struct parser *p, struct stmt_array *arr);
void init_parser(struct parser *p, char *filepath);
