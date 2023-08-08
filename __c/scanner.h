enum token_kind {
    TOK_ERR = 0,

    TOK_IDENT,
    TOK_NUM,

    TOK_SEMICOLON,

    TOK_PRINT,

    TOK_EOF
};

struct token {
    enum token_kind kind;
    int line;
    char *lit;
};

struct scanner {
    char *src;
    char *filepath;
    int src_len;
    int pos;
    int cur;
    int line;
    char ch;
};

void make_scanner(struct scanner *s, char *filepath);
void scan(struct scanner *s, struct token *tok);
char *token_kind_str(enum token_kind kind);
