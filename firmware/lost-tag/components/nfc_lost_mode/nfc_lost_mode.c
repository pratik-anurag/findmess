#include "nfc_lost_mode.h"

static const char *payload;

void fm_nfc_lost_mode_init(void) {}
void fm_nfc_lost_mode_set_payload(const char *value) { payload = value; (void)payload; }
