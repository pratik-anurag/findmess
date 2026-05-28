#pragma once
#include <stdbool.h>
#include <stdint.h>

void fm_ble_advertiser_init(void);
void fm_ble_advertise_tag(const uint8_t ephemeral_id[16], bool lost_hint, bool battery_low);
