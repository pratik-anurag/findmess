#pragma once
#include <stdbool.h>

void fm_battery_init(void);
int fm_battery_percent(void);
bool fm_battery_is_low(void);
