#include "esp_timer.h"
#include "anti_stalking.h"

static int64_t last_owner_seen_us;

void fm_anti_stalking_init(void)
{
    last_owner_seen_us = esp_timer_get_time();
}

void fm_anti_stalking_owner_seen(void)
{
    last_owner_seen_us = esp_timer_get_time();
}

bool fm_anti_stalking_should_sound(void)
{
    const int64_t separation_us = esp_timer_get_time() - last_owner_seen_us;
    return separation_us > (8LL * 60LL * 60LL * 1000000LL);
}
