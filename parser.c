#include <stdlib.h>
#include <stdio.h>
#include <string.h>
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
    if (p->tok.kind != kind) {
        // TODO(art): translate token enum to text
        fprintf(stderr, "expected %i, but got %i\n", kind, p->tok.kind);
        exit(1);
    }

    advance(p);
}

static void primary(struct parser *p, struct expr *e)
{
    if (p->tok.kind == TOK_NUM) {
        struct expr_lit *lit = malloc(sizeof(struct expr_lit));
        // TODO(art): error handling
        assert(lit != NULL);

        memcpy(&lit->token, &p->tok, sizeof(struct token));

        e->kind = EXPR_LIT;
        e->body = lit;

        advance(p);
        return;
    }

    // TODO(art): should not be there after all
    assert(0 && "should not be there");
}

static void factor(struct parser *p, struct expr *e)
{
    primary(p, e);

    if (p->tok.kind == TOK_STAR || p->tok.kind == TOK_SLASH) {
        struct expr_binary *b = malloc(sizeof(struct expr_binary));
        // TODO(art): error
        assert(b != NULL);

        b->op = p->tok.kind;
        memcpy(&b->x, e, sizeof(struct expr));

        advance(p);
        factor(p, &b->y);

        e->kind = EXPR_BINARY;
        e->body = b;
    }
}
static void term(struct parser *p, struct expr *e)
{
    factor(p, e);

    if (p->tok.kind == TOK_PLUS || p->tok.kind == TOK_MINUS) {
        struct expr_binary *b = malloc(sizeof(struct expr_binary));
        // TODO(art): error
        assert(b != NULL);

        b->op = p->tok.kind;
        memcpy(&b->x, e, sizeof(struct expr));

        advance(p);
        term(p, &b->y);

        e->kind = EXPR_BINARY;
        e->body = b;
    }
}

static void expression(struct parser *p, struct expr *e)
{
    term(p, e);
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
        statement(p, s);
    }
}

void init_parser(struct parser *p, char *filepath)
{
    init_scanner(&p->s, filepath);
    advance(p);
}
