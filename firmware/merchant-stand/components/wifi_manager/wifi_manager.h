#pragma once
#include <stdbool.h>

void fm_wifi_connect(void);
bool fm_wifi_is_connected(void);
bool fm_wifi_reconnect_tick(void);
