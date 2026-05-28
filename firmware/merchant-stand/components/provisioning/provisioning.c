#include "provisioning.h"

static bool claimed;

bool fm_provisioning_is_claimed(void) { return claimed; }
void fm_provisioning_start(void) { claimed = false; }
