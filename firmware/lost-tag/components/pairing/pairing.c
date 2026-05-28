#include "pairing.h"

static bool pairing;

void fm_pairing_init(void) {}
bool fm_pairing_button_held(void) { return false; }
void fm_pairing_start(void) { pairing = true; }
bool fm_pairing_complete(void) { pairing = false; return true; }
bool fm_pairing_ring_requested(void) { return false; }
