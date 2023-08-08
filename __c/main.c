#include <stdio.h>
#include <stdlib.h>

#include "lib/mem.h"

#include "scanner.h"
#include "parser.h"

struct stmt_array {
    int len;
    int cap;
    int data_size;
    struct stmt *buf;
};

struct value {
    enum {
        VAL_NUM = 0
    } kind;

    union {
        double num;
    } as;
};

void eval(struct expr *e, struct value *val)
{
    switch (e->kind) {
    case EXPR_LIT: {
        struct expr_lit *el = e->body;

        switch (el->value.kind) {
        case TOK_NUM:
            // TODO(art): errors
            val->as.num = strtod(el->value.lit, NULL);
            val->kind = VAL_NUM;
            free(el->value.lit);
            break;
        default:
            fprintf(stderr, "unknown literal value kind\n");
            exit(1);
        }
    } break;
    default:
        fprintf(stderr, "unknown statement kind\n");
        exit(1);
    }
}

void execute(struct stmt *s)
{
    switch (s->kind) {
    case STMT_PRINT: {
        struct stmt_print *sp = s->body;
        struct value res;

        eval(&sp->value, &res);

        switch (res.kind) {
        case VAL_NUM:
            printf("%g\n", res.as.num);
            break;
        default:
            fprintf(stderr, "unknown value kind\n");
            exit(1);
        }
    } break;
    default:
        fprintf(stderr, "unknown statement kind\n");
        exit(1);
    }
}

int main(int argc, char **argv)
{
    struct parser p;
    struct stmt_array stmts;

    argc--;
    argv++;

    if (argc < 1) {
        fprintf(stderr, "Provide file to execute\n");
        return 1;
    }

    meminit(&stmts, sizeof(struct stmt), 64);
    make_parser(&p, *argv);

    parse(&p, &stmts);

    for (int i = 0; i < stmts.len; ++i) {
        execute(stmts.buf + i);
    }

    return 0;
}
