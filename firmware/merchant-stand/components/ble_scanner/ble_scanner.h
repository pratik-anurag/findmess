#pragma once
#include <stdbool.h>
#include <stdint.h>

typedef struct {
    uint8_t ephemeral_id[16];
    int rssi;
    bool lost_hint;
    bool battery_low;
} fm_ble_tag_adv_t;

typedef void (*fm_ble_scan_callback_t)(const fm_ble_tag_adv_t *adv);

void fm_ble_scanner_init(fm_ble_scan_callback_t callback);
void fm_ble_scanner_start(void);
const char *fm_rssi_bucket(int rssi);
