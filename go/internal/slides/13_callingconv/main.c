#include <stdio.h>

__attribute__ ((noinline)) int add(int x, int y) {
    return x + y;
}

int main() {
    int x = 1;
    int y = 2;

    int r = add(x, y);

    printf("%i\n", r);
}

