#include <stdlib.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

#include "os.h"

char *read_file(char *filepath)
{
    struct stat statbuf;
    char *buf;
    int fd;

    fd = open(filepath, O_RDONLY);
    if (fd < 0) {
        return NULL;
    }

    if (fstat(fd, &statbuf) < 0) {
        return NULL;
    }

    buf = malloc(statbuf.st_size + 1);
    if (buf == NULL) {
        return NULL;
    }

    if (read(fd, buf, statbuf.st_size) < 0) {
        free(buf);
        return NULL;
    }

    close(fd);

    buf[statbuf.st_size] = '\0';
    return buf;
}
