#include <string.h>
#include "sighting_buffer.h"

void fm_sighting_buffer_init(fm_sighting_buffer_t *buffer)
{
    memset(buffer, 0, sizeof(*buffer));
}

bool fm_sighting_buffer_push(fm_sighting_buffer_t *buffer, const fm_signed_sighting_t *sighting)
{
    int idx = (buffer->head + buffer->count) % FM_SIGHTING_BUFFER_CAPACITY;
    buffer->entries[idx] = *sighting;
    if (buffer->count == FM_SIGHTING_BUFFER_CAPACITY) {
        buffer->head = (buffer->head + 1) % FM_SIGHTING_BUFFER_CAPACITY;
        return false;
    }
    buffer->count++;
    return true;
}

int fm_sighting_buffer_upload_batch(fm_sighting_buffer_t *buffer)
{
    int uploaded = buffer->count;
    buffer->head = 0;
    buffer->count = 0;
    return uploaded;
}

int64_t fm_time_bucket_now(void)
{
    return 0;
}

uint64_t fm_nonce_next(void)
{
    static uint64_t nonce = 1;
    return nonce++;
}
