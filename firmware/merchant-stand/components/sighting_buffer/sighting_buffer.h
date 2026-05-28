#pragma once
#include <stdbool.h>
#include <stdint.h>

#define FM_SIGHTING_BUFFER_CAPACITY 256

typedef struct {
    uint8_t tag_ephemeral_id[16];
    int protocol_version;
    int64_t time_bucket;
    const char *rssi_bucket;
    uint64_t nonce;
    uint8_t signature[64];
} fm_signed_sighting_t;

typedef struct {
    fm_signed_sighting_t entries[FM_SIGHTING_BUFFER_CAPACITY];
    int head;
    int count;
} fm_sighting_buffer_t;

void fm_sighting_buffer_init(fm_sighting_buffer_t *buffer);
bool fm_sighting_buffer_push(fm_sighting_buffer_t *buffer, const fm_signed_sighting_t *sighting);
int fm_sighting_buffer_upload_batch(fm_sighting_buffer_t *buffer);
int64_t fm_time_bucket_now(void);
uint64_t fm_nonce_next(void);
