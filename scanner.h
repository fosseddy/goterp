enum token_kind {
    TOK_ERR = 0,

	TOK_NUM,

    TOK_PLUS,
    TOK_MINUS,

	TOK_SLASH,
	TOK_SEMICOLON,

	TOK_PRINT,

	TOK_EOF
};

struct scanner {
    char *src;
    int src_len;
    int cur;
    int pos;
    char ch;
};

struct token {
    enum token_kind kind;
    char *lit;
    int lit_len;
};

void scan(struct scanner *s, struct token *tok);
void init_scanner(struct scanner *s, char *filepath);
