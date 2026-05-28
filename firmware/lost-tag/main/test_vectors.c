#include <string.h>
#include "ephemeral_id.h"

int fm_ephemeral_id_test_vector(void)
{
    const uint8_t secret[32] = "0123456789abcdef0123456789abcdef";
    const uint8_t expected[16] = {
        0x9e, 0x40, 0xef, 0x0c, 0x67, 0x7a, 0xda, 0xe9,
        0x87, 0x08, 0x09, 0xe1, 0xcd, 0x95, 0x2f, 0xc2,
    };
    uint8_t out[16] = {0};
    fm_ephemeral_id_derive(secret, sizeof(secret), 123456, out);
    return memcmp(out, expected, sizeof(expected)) == 0;
}
