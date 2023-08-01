#include <stdio.h>

#include "scanner.h"

int main(int argc, char **argv)
{
    struct scanner s;

    argc--;
    argv++;

    if (argc < 1) {
        fprintf(stderr, "Provide file to execute\n");
        return 1;
    }

    init_scanner(&s, *argv);

    for (;;) {
        struct token t;

        scan(&s, &t);

        printf("Token Kind: %d\n", t.kind);

        if (t.kind == TOK_EOF) {
            break;
        }
    }

    return 0;
}
