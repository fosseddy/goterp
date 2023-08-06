#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "lib/mem.h"

#include "scanner.h"
#include "parser.h"

struct value {
    enum { VAL_NUM } kind;
    union { double num; } as;
};

void eval(struct expr *e, struct value *res)
{
    struct expr_lit *lit;
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
    default: assert(0 && "unreachable");
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
        struct stmt *s = stmts.buf + i;
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
            default: assert(0 && "unreachable");
            }
            break;
        default: assert(0 && "unreachable");
        }
    }

    return 0;
}
