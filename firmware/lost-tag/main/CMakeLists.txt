idf_component_register(
  SRCS "findmesh_tag_main.c"
  INCLUDE_DIRS "."
  REQUIRES ble_advertiser ephemeral_id nfc_lost_mode buzzer battery pairing anti_stalking esp_timer
)
