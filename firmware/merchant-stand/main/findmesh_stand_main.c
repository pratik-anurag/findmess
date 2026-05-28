#include <stdbool.h>
#include <stdint.h>
#include "esp_log.h"
#include "esp_timer.h"
#include "ble_scanner.h"
#include "ble_advertiser.h"
#include "wifi_manager.h"
#include "nfc_manager.h"
#include "sighting_buffer.h"
#include "crypto_signer.h"
#include "ota_manager.h"
#include "led_status.h"
#include "power_manager.h"
#include "provisioning.h"

static const char *TAG = "findmesh_stand";

typedef enum {
    FM_STAND_UNCLAIMED,
    FM_STAND_PROVISIONING,
    FM_STAND_ONLINE,
    FM_STAND_OFFLINE_BUFFERING,
    FM_STAND_RECOVERY_MODE,
    FM_STAND_ERROR,
    FM_STAND_OTA_UPDATING,
} fm_stand_state_t;

static fm_stand_state_t state = FM_STAND_UNCLAIMED;
static fm_sighting_buffer_t buffer;

static void handle_scan_result(const fm_ble_tag_adv_t *adv)
{
    fm_signed_sighting_t sighting = {0};
    sighting.protocol_version = 1;
    sighting.time_bucket = fm_time_bucket_now();
    sighting.rssi_bucket = fm_rssi_bucket(adv->rssi);
    sighting.nonce = fm_nonce_next();
    for (int i = 0; i < 16; i++) {
        sighting.tag_ephemeral_id[i] = adv->ephemeral_id[i];
    }
    fm_crypto_sign_sighting(&sighting);
    fm_sighting_buffer_push(&buffer, &sighting);
    fm_led_double_blink();
}

void app_main(void)
{
    ESP_LOGI(TAG, "booting merchant stand firmware");
    fm_led_init();
    fm_power_init();
    fm_sighting_buffer_init(&buffer);
    fm_crypto_signer_init();
    fm_nfc_init();
    fm_ble_advertiser_init();
    fm_ble_scanner_init(handle_scan_result);

    if (!fm_provisioning_is_claimed()) {
        state = FM_STAND_PROVISIONING;
        fm_led_fast_blink();
        fm_provisioning_start();
        fm_nfc_set_payload("findmesh://stand/setup?s=local-claim-token");
    } else {
        state = FM_STAND_ONLINE;
        fm_wifi_connect();
        fm_ble_advertise_zone();
        fm_ble_scanner_start();
        fm_led_slow_blink();
    }

    while (true) {
        if (state == FM_STAND_ONLINE && !fm_wifi_is_connected()) {
            state = FM_STAND_OFFLINE_BUFFERING;
            fm_led_error_pattern();
        }
        if (state == FM_STAND_OFFLINE_BUFFERING && fm_wifi_reconnect_tick()) {
            state = FM_STAND_ONLINE;
            fm_led_slow_blink();
        }
        if (state == FM_STAND_ONLINE) {
            fm_sighting_buffer_upload_batch(&buffer);
            fm_ota_check_manifest();
        }
        fm_power_apply_scan_duty_cycle();
        int64_t next_tick_us = 1000 * 1000;
        esp_timer_delay_us(next_tick_us);
    }
}
