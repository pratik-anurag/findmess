#include <stdbool.h>
#include <stdint.h>
#include "esp_log.h"
#include "esp_timer.h"
#include "ble_advertiser.h"
#include "ephemeral_id.h"
#include "nfc_lost_mode.h"
#include "buzzer.h"
#include "battery.h"
#include "pairing.h"
#include "anti_stalking.h"

static const char *TAG = "findmesh_tag";

typedef enum {
    FM_TAG_UNPAIRED,
    FM_TAG_PAIRING,
    FM_TAG_PAIRED_NORMAL,
    FM_TAG_LOST_MODE,
    FM_TAG_SEPARATED,
    FM_TAG_RINGING,
    FM_TAG_DISABLED,
} fm_tag_state_t;

static fm_tag_state_t state = FM_TAG_UNPAIRED;
static uint8_t tag_secret[32] = {
    0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
    0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
    0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
    0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
};

static void advertise_current_id(bool lost_hint)
{
    uint8_t ephemeral_id[16];
    int64_t epoch = fm_epoch_now();
    fm_ephemeral_id_derive(tag_secret, sizeof(tag_secret), epoch, ephemeral_id);
    fm_ble_advertise_tag(ephemeral_id, lost_hint, fm_battery_is_low());
}

void app_main(void)
{
    ESP_LOGI(TAG, "booting lost tag firmware");
    fm_ble_advertiser_init();
    fm_nfc_lost_mode_init();
    fm_buzzer_init();
    fm_battery_init();
    fm_pairing_init();
    fm_anti_stalking_init();

    if (fm_pairing_button_held()) {
        state = FM_TAG_PAIRING;
        fm_pairing_start();
    } else {
        state = FM_TAG_PAIRED_NORMAL;
    }

    while (true) {
        if (state == FM_TAG_PAIRING && fm_pairing_complete()) {
            state = FM_TAG_PAIRED_NORMAL;
        }
        if (state == FM_TAG_PAIRED_NORMAL || state == FM_TAG_LOST_MODE) {
            advertise_current_id(state == FM_TAG_LOST_MODE);
        }
        if (fm_anti_stalking_should_sound()) {
            state = FM_TAG_SEPARATED;
            fm_buzzer_pattern_unwanted_tracker();
        }
        if (fm_pairing_ring_requested()) {
            state = FM_TAG_RINGING;
            fm_buzzer_ring();
            state = FM_TAG_PAIRED_NORMAL;
        }
        fm_nfc_lost_mode_set_payload("findmesh://tag-found?t=local-lost-token");
        esp_timer_delay_us(60 * 1000 * 1000);
    }
}
