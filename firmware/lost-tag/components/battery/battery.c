#include "battery.h"

void fm_battery_init(void) {}
int fm_battery_percent(void) { return 100; }
bool fm_battery_is_low(void) { return fm_battery_percent() < 15; }
