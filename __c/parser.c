#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mem.h"

#include "scanner.h"
#include "parser.h"

static void advance(struct parser *p)
{
    scan(&p->s, &p->tok);
}

static void consume(struct parser *p, enum token_kind kind)
{
    if (p->tok.kind != kind) {
        fprintf(stderr, "%s:%d:expected %s but got %s\n", p->s.filepath, p->tok.line, token_kind_str(kind), token_kind_str(p->tok.kind));
        exit(1);
    }

    advance(p);
}

void make_parser(struct parser *p, char *filepath)
{
    make_scanner(&p->s, filepath);
    advance(p);
}

static void primary(struct parser *p, struct expr *e)
{
    if (p->tok.kind == TOK_NUM) {
        struct expr_lit *el = malloc(sizeof(struct expr_lit));

        if (el == NULL) {
            perror("primary");
            exit(1);
        }

        memcpy(&el->value, &p->tok, sizeof(struct token));
        e->body = el;

        advance(p);
        return;
    }

    fprintf(stderr, "%s:%d:unknown primary expression %s\n", p->s.filepath, p->tok.line, token_kind_str(p->tok.kind));
    exit(1);
}

static void expression(struct parser *p, struct expr *e)
{
    primary(p, e);
}

static void print_stmt(struct parser *p, struct stmt *s)
{
    struct stmt_print *sp = malloc(sizeof(struct stmt_print));

    if (sp == NULL) {
        perror("print_stmt");
        exit(1);
    }

    advance(p);
    expression(p, &sp->value);
    consume(p, TOK_SEMICOLON);

    s->body = sp;
}

static void statement(struct parser *p, struct stmt *s)
{
    if (p->tok.kind == TOK_PRINT) {
        print_stmt(p, s);
        return;
    }

    fprintf(stderr, "%s:%d:unknown statement %s\n", p->s.filepath, p->tok.line, token_kind_str(p->tok.kind));
    exit(1);
}

void parse(struct parser *p, void *stmts)
{
    while (p->tok.kind != TOK_EOF) {
        struct stmt *s = memnext(stmts);
        statement(p, s);
    }
}
