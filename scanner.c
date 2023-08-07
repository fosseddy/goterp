#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "lib/os.h"

#include "scanner.h"

char *token_kind_str(enum token_kind kind)
{
    switch (kind) {
    case TOK_IDENT:
		return "identifier";
    case TOK_NUM:
		return "number";

    case TOK_SEMICOLON:
		return ";";

    case TOK_PRINT:
        return "print";

    case TOK_EOF:
        return "end of file";

    default:
        return "invalid token";
    }
}

struct kwd_entry {
    char *str;
    enum token_kind kind;
};

static struct kwd_entry keywords[] = {
    {.str = "print", .kind = TOK_PRINT},
    {.str = "", .kind = TOK_ERR} // art: end of array
};

static enum token_kind lookupKeyword(struct scanner *s)
{
    char *ident = s->src + s->pos;
    int len = s->cur - s->pos;

    for (struct kwd_entry *kwd = keywords; kwd->kind != TOK_ERR; kwd++) {
        if ((int) strlen(kwd->str) == len &&
                memcmp(kwd->str, ident, len) == 0) {
            return kwd->kind;
        }
    }

    return TOK_IDENT;
}

void make_scanner(struct scanner *s, char *filepath)
{
    s->src_len = read_file(filepath, &s->src);
    if (s->src_len < 0) {
        perror("failed to make scanner");
        exit(1);
    }

    s->filepath = filepath;
    s->pos = 0;
    s->cur = 0;
    s->line = 1;
    s->ch = s->src[0];
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

static int has_src(struct scanner *s)
{
    return s->cur < s->src_len;
}

static void advance(struct scanner *s)
{
    s->cur++;

    if (has_src(s) == 0) {
        s->ch--;
        return;
    }

    s->ch = s->src[s->cur];
}

static char *get_literal(struct scanner *s)
{
    int len = s->cur - s->pos;
    char *str = s->src + s->pos;
    char *buf = malloc(len + 1);

    if (buf == NULL) {
        perror("get_literal");
        exit(1);
    }

    memcpy(buf, str, len);
    buf[len] = '\0';

    return buf;
}

static void make_token(struct token *tok, enum token_kind kind, int line,
                       char *lit)
{
    tok->kind = kind;
    tok->line = line;
    tok->lit = lit;
}

static void make_ident_token(struct scanner *s, struct token *tok)
{
    tok->line = s->line;
    tok->kind = lookupKeyword(s);

    if (tok->kind == TOK_IDENT) {
        tok->lit = get_literal(s);
    }
}

void scan(struct scanner *s, struct token *tok)
{
scan_again:
    if (has_src(s) == 0) {
        make_token(tok, TOK_EOF, s->line, NULL);
        return;
    }

    s->pos = s->cur;

    switch (s->ch) {
    case ' ':
    case '\t':
    case '\r':
    case '\n':
        if (s->ch == '\n') {
            s->line++;
        }
        advance(s);
        goto scan_again;

    case ';':
        make_token(tok, TOK_SEMICOLON, s->line, NULL);
        advance(s);
        break;

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
            make_token(tok, TOK_NUM, s->line, get_literal(s));
        } else if (is_char(s->ch) == 1) {
            while (is_alnum(s->ch) == 1) {
                advance(s);
            }
            make_ident_token(s, tok);
        } else {
            fprintf(stderr, "%s:%d:unexpected character %c\n", s->filepath, s->line, s->ch);
            exit(1);
        }
    }
}
