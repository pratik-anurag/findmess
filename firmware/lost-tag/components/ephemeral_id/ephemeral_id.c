#include <string.h>
#include "esp_timer.h"
#include "mbedtls/md.h"
#include "ephemeral_id.h"

#define FM_EPOCH_SECONDS 900

static void hmac_sha256(const uint8_t *key, size_t key_len, const uint8_t *input, size_t input_len, uint8_t out[32])
{
    const mbedtls_md_info_t *info = mbedtls_md_info_from_type(MBEDTLS_MD_SHA256);
    mbedtls_md_hmac(info, key, key_len, input, input_len, out);
}

void fm_ephemeral_id_derive(const uint8_t *secret, size_t secret_len, int64_t epoch, uint8_t out[16])
{
    static const uint8_t label[] = "findmesh-ephemeral-id";
    const size_t label_len = sizeof(label) - 1;
    uint8_t salt[32] = {0};
    uint8_t prk[32] = {0};
    uint8_t info[29] = {0};
    uint8_t hmac_input[30] = {0};
    uint8_t okm[32] = {0};

    hmac_sha256(salt, sizeof(salt), secret, secret_len, prk);
    memcpy(info, label, label_len);
    for (int i = 0; i < 8; i++) {
        info[label_len + i] = (uint8_t)((uint64_t)epoch >> (56 - (i * 8)));
    }
    memcpy(hmac_input, info, sizeof(info));
    hmac_input[sizeof(info)] = 1;
    hmac_sha256(prk, sizeof(prk), hmac_input, sizeof(hmac_input), okm);
    memcpy(out, okm, 16);
}

int64_t fm_epoch_now(void)
{
    int64_t seconds = esp_timer_get_time() / 1000000;
    return seconds / FM_EPOCH_SECONDS;
}
