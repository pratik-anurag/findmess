#include "nfc_manager.h"

static const char *current_payload;

void fm_nfc_init(void) {}
void fm_nfc_set_payload(const char *payload) { current_payload = payload; (void)current_payload; }
