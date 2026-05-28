#include "wifi_manager.h"

static bool connected;

void fm_wifi_connect(void) { connected = true; }
bool fm_wifi_is_connected(void) { return connected; }
bool fm_wifi_reconnect_tick(void) { connected = true; return true; }
