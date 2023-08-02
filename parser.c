#include <stdlib.h>
#include <assert.h>

#include "lib/mem.h"

#include "scanner.h"
#include "parser.h"

static void advance(struct parser *p) {
    scan(&p->s, &p->tok);

    // TODO(art): error handling
    assert(p->tok.kind != TOK_ERR);
}

static void consume(struct parser *p, enum token_kind kind) {
    // TODO(art): error handling
    assert(p->tok.kind == kind && "expected something else");

    advance(p);
}

static void primary(struct parser *p, struct expr *e)
{
    if (p->tok.kind == TOK_NUM) {
        struct expr_lit *lit = malloc(sizeof(struct expr_lit));
        // TODO(art): error handling
        assert(lit != NULL);

        lit->kind = p->tok.kind;

        e->kind = EXPR_LIT;
        e->body = lit;

        advance(p);
        return;
    }

    // TODO(art): should not be there after all
    assert(0 && "should not be there");
}

static void expression(struct parser *p, struct expr *e)
{
    primary(p, e);
}

static void print_stmt(struct parser *p, struct stmt *s)
{
    struct stmt_print *print = malloc(sizeof(struct stmt_print));
    // TODO(art): error handling
    assert(print != NULL);

    consume(p, TOK_PRINT);
    expression(p, &print->value);
    consume(p, TOK_SEMICOLON);

    s->kind = STMT_PRINT;
    s->body = print;
}

static void statement(struct parser *p, struct stmt *s)
{
    if (p->tok.kind == TOK_PRINT) {
        print_stmt(p, s);
        return;
    }

    // TODO(art): should not be there after all
    assert(0 && "should not be there");
}

void parse(struct parser *p, struct stmt_array *arr)
{
    while (p->tok.kind != TOK_EOF) {
        struct stmt *s = memnext(arr);

        s->kind = STMT_ERR;
        statement(p, s);

        // TODO(art): error handling
        assert(s->kind != STMT_ERR);
    }
}

void init_parser(struct parser *p, char *filepath)
{
    init_scanner(&p->s, filepath);
    advance(p);
}
