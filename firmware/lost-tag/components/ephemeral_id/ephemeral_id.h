#pragma once
#include <stddef.h>
#include <stdint.h>

void fm_ephemeral_id_derive(const uint8_t *secret, size_t secret_len, int64_t epoch, uint8_t out[16]);
int64_t fm_epoch_now(void);
