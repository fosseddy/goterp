#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "lib/mem.h"

#include "scanner.h"
#include "parser.h"

struct value {
    enum {
        VAL_NUM
    } kind;

    union {
        double num;
    } as;
};

void eval(struct expr *e, struct value *res)
{
    struct expr_lit *lit;
    struct expr_binary *binary;
    struct value x, y;
    char *s;

    switch (e->kind) {
    case EXPR_LIT:
        lit = e->body;
        switch (lit->token.kind) {
        case TOK_NUM:
            s = malloc(lit->token.lit_len + 1);
            // TODO(art): error
            assert(s != NULL);
            memcpy(s, lit->token.lit, lit->token.lit_len);
            s[lit->token.lit_len] = '\0';
            // TODO(art): handle conversion error
            res->as.num = strtod(s, NULL);
            res->kind = VAL_NUM;
            free(s);
            break;
        default: assert(0 && "unreachable");
        }
        break;
    case EXPR_BINARY:
        binary = e->body;
        eval(&binary->x, &x);
        eval(&binary->y, &y);

        // TODO(art): validate values type
        assert(x.kind == VAL_NUM);
        assert(y.kind == VAL_NUM);
        res->kind = VAL_NUM;

        switch(binary->op) {
        case TOK_PLUS:
            res->as.num = x.as.num + y.as.num;
            break;
        case TOK_MINUS:
            res->as.num = x.as.num - y.as.num;
            break;
        case TOK_STAR:
            res->as.num = x.as.num * y.as.num;
            break;
        case TOK_SLASH:
            res->as.num = x.as.num / y.as.num;
            break;
        default: assert(0 && "unreachable");
        }
        break;
    default: assert(0 && "unreachable");
    }
}

void execute(struct stmt *s)
{
    struct stmt_print *print;
    struct value res;

    switch (s->kind) {
    case STMT_PRINT:
        print = s->body;
        eval(&print->value, &res);

        switch (res.kind) {
        case VAL_NUM:
            printf("%g\n", res.as.num);
            break;
        default:
            assert(0 && "unreachable");
        }
        break;
    default:
        assert(0 && "unreachable");
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

    init_parser(&p, *argv);
    meminit(&stmts, sizeof(struct stmt), 32);

    parse(&p, &stmts);

    for (int i = 0; i < stmts.len; ++i) {
        execute(stmts.buf + i);
    }

    return 0;
}
