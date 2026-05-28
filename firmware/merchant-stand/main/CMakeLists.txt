idf_component_register(
  SRCS "findmesh_stand_main.c"
  INCLUDE_DIRS "."
  REQUIRES ble_scanner ble_advertiser wifi_manager nfc_manager sighting_buffer crypto_signer ota_manager led_status power_manager provisioning esp_timer
)
