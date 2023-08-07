// #include "lib/mem.h"
// #include "scanner.h"

enum expr_kind {
    EXPR_LIT = 0
};

struct expr_lit {
    struct token value;
};

struct expr {
    enum expr_kind kind;
    void *body;
};

enum stmt_kind {
    STMT_PRINT = 0
};

struct stmt_print {
    struct expr value;
};

struct stmt {
    enum stmt_kind kind;
    void *body;
};

struct parser {
    struct token tok;
    struct scanner s;
};

void make_parser(struct parser *p, char *filepath);
void parse(struct parser *p, void *stmts);
