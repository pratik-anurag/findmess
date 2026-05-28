#include <stdbool.h>
#include "ble_scanner.h"

static fm_ble_scan_callback_t scan_callback;

void fm_ble_scanner_init(fm_ble_scan_callback_t callback)
{
    scan_callback = callback;
}

void fm_ble_scanner_start(void)
{
    (void)scan_callback;
}

const char *fm_rssi_bucket(int rssi)
{
    if (rssi >= -60) return "near";
    if (rssi >= -78) return "medium";
    return "far";
}
