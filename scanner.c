#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "lib/os.h"

#include "scanner.h"

struct kwd {
    char *lit;
    enum token_kind tok;
};

static struct kwd keywords[] = {
    {.lit = "print", .tok = TOK_PRINT},

    {.lit = "", .tok = TOK_ERR} // end value
};

static enum token_kind lookup_kwd(struct scanner *s)
{
    char *ident = s->src + s->pos;
    int ident_len = s->cur - s->pos;
    int i = 0;

    for (;;) {
        struct kwd *kwd = keywords + i;
        i++;

        if (kwd->tok == TOK_ERR) {
            break;
        }

        if ((int) strlen(kwd->lit) == ident_len &&
                memcmp(ident, kwd->lit, ident_len) == 0) {
            return kwd->tok;
        }
    }

    return TOK_ERR;
}

static int has_src(struct scanner *s)
{
    return s->cur < s->src_len;
}

static int next(struct scanner *s, char ch)
{
    int next = s->cur + 1;

    if (next < s->src_len) {
        return s->src[next] == ch;
    }

    return 0;
}

static void advance(struct scanner *s)
{
    if (has_src(s) == 0) {
        s->ch = '\0';
        return;
    }

    s->cur++;
    s->ch = s->src[s->cur];
}

static int is_digit(char ch)
{
    return ch >= '0' && ch <= '9';
}

static int is_char(char ch)
{
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_';
}

static int is_alnum(char ch)
{
    return is_char(ch) == 1 || is_digit(ch) == 1;
}

void scan(struct scanner *s, struct token *tok)
{
    tok->kind = TOK_ERR;
    tok->lit = NULL;
    tok->lit_len = 0;

scan_again:

    s->pos = s->cur;

    if (has_src(s) == 0) {
        tok->kind = TOK_EOF;
        return;
    }

    switch (s->ch) {
    case ' ':
    case '\t':
    case '\r':
    case '\n':
        advance(s);
        goto scan_again;

    case '/':
        if (next(s, '/') == 1) {
            while (has_src(s) == 1 && s->ch != '\n') {
                advance(s);
            }
            goto scan_again;
        }
        tok->kind = TOK_SLASH;
        advance(s);
        return;

    case '+':
        tok->kind = TOK_PLUS;
        advance(s);
        return;

    case '-':
        tok->kind = TOK_MINUS;
        advance(s);
        return;

    case '*':
        tok->kind = TOK_STAR;
        advance(s);
        return;

    case '(':
        tok->kind = TOK_LPAREN;
        advance(s);
        return;

    case ')':
        tok->kind = TOK_RPAREN;
        advance(s);
        return;

    case ';':
        tok->kind = TOK_SEMICOLON;
        advance(s);
        return;

    default:
        if (is_digit(s->ch) == 1) {
            while (is_digit(s->ch) == 1) {
                advance(s);
            }

            if (s->ch == '.') {
                advance(s);
                while (is_digit(s->ch) == 1) {
                    advance(s);
                }
            }

            tok->kind = TOK_NUM;
            tok->lit = s->src + s->pos;
            tok->lit_len = s->cur - s->pos;
        } else if (is_char(s->ch) == 1) {
            while (is_alnum(s->ch) == 1) {
                advance(s);
            }
            tok->kind = lookup_kwd(s);

            // TODO(art): remove when TOK_IDENT
            if (tok->kind == TOK_ERR) {
                fprintf(stderr, "found identifier");
            }
        } else {
            fprintf(stderr, "unknown char %c\n", s->ch);
            advance(s);
        }
    }
}

void init_scanner(struct scanner *s, char *filepath)
{
    s->src = read_file(filepath);
    if (s->src == NULL) {
        perror("failed to read file");
        exit(1);
    }

    s->src_len = (int) strlen(s->src);
    s->cur = 0;
    s->pos = 0;
    s->ch = s->src[0];
}
