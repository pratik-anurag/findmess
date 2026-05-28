#pragma once
#include <stdbool.h>

void fm_pairing_init(void);
bool fm_pairing_button_held(void);
void fm_pairing_start(void);
bool fm_pairing_complete(void);
bool fm_pairing_ring_requested(void);
