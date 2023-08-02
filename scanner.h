enum token_kind {
    TOK_ERR = 0,

	TOK_IDENT,
	TOK_NUM,
	//TOK_STR,

	//TOK_PLUS,
	//TOK_MINUS,
	//TOK_STAR,
	TOK_SLASH,
	//TOK_COMMA,
	//TOK_BANG,
	//TOK_EQ,
	//TOK_EQ_EQ,
	//TOK_BANG_EQ,
	//TOK_LESS,
	//TOK_LESS_EQ,
	//TOK_GREATER,
	//TOK_GREATER_EQ,
	//TOK_AND,
	//TOK_OR,

	//TOK_LPAREN,
	//TOK_RPAREN,
	//TOK_LBRACE,
	//TOK_RBRACE,
	TOK_SEMICOLON,

	TOK_PRINT,
	//TOK_LET,
	//TOK_IF,
	//TOK_ELSE,
	//TOK_WHILE,
	//TOK_FN,

	//TOK_TRUE,
	//TOK_FALSE,
	//TOK_NIL,

	TOK_EOF
};

struct scanner {
    char *src;
    int cur;
    int pos;
    char ch;
};

struct token {
    enum token_kind kind;
};

void scan(struct scanner *s, struct token *tok);
void init_scanner(struct scanner *s, char *filepath);
