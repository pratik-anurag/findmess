import 'dart:convert';
import 'dart:typed_data';

import 'ble_service.dart';

class FindMeshBleProtocol {
  static const serviceUuid = 'f17d0001-2d9a-4c8a-a4f1-f0d641b90f10';
  static const manufacturerId = 0xF17D;
  static const tagAdvertisementType = 0x01;
  static const zoneAdvertisementType = 0x02;

  static Uint8List encodeTag(String ephemeralId, {bool lostHint = false, bool batteryLow = false}) {
    return _encode(tagAdvertisementType, ephemeralId, lostHint: lostHint, batteryLow: batteryLow);
  }

  static Uint8List encodeZone(String zoneEphemeralId) {
    return _encode(zoneAdvertisementType, zoneEphemeralId);
  }

  static BleAdvertisement? parse({
    required Iterable<List<int>> serviceData,
    required Iterable<List<int>> manufacturerData,
    required int rssi,
    String? deviceId,
    String? name,
  }) {
    for (final data in [...manufacturerData, ...serviceData]) {
      final parsed = parseBytes(data, rssi: rssi, deviceId: deviceId, name: name);
      if (parsed != null) return parsed;
    }
    final parsedName = parseLocalName(name, rssi: rssi, deviceId: deviceId);
    if (parsedName != null) return parsedName;
    return null;
  }

  static BleAdvertisement? parseBytes(List<int> data, {required int rssi, String? deviceId, String? name}) {
    if (data.length < 18) return null;
    final type = data[0];
    final flags = data[1];
    final ephemeralId = data.sublist(2, 18).map((byte) => byte.toRadixString(16).padLeft(2, '0')).join();
    if (type == tagAdvertisementType) {
      return BleAdvertisement(
        type: 'FM_TAG',
        ephemeralId: ephemeralId,
        rssi: rssi,
        deviceId: deviceId,
        name: name,
        flags: {
          'lost_hint': flags & 0x01 == 0x01,
          'battery_low': flags & 0x02 == 0x02,
        },
      );
    }
    if (type == zoneAdvertisementType) {
      return BleAdvertisement(type: 'FM_ZONE', ephemeralId: ephemeralId, rssi: rssi, deviceId: deviceId, name: name);
    }
    return null;
  }

  static BleAdvertisement? parseLocalName(String? name, {required int rssi, String? deviceId}) {
    if (name == null || name.isEmpty) return null;
    if (name.startsWith('FM_TAG:')) {
      final id = _normalizeHex(name.substring('FM_TAG:'.length));
      if (id != null) return BleAdvertisement(type: 'FM_TAG', ephemeralId: id, rssi: rssi, deviceId: deviceId, name: name);
    }
    if (name.startsWith('FM_ZONE:')) {
      final id = _normalizeHex(name.substring('FM_ZONE:'.length));
      if (id != null) return BleAdvertisement(type: 'FM_ZONE', ephemeralId: id, rssi: rssi, deviceId: deviceId, name: name);
    }
    return null;
  }

  static String demoEphemeralId(String seed) {
    final bytes = utf8.encode(seed);
    final out = List<int>.filled(16, 0);
    for (var i = 0; i < bytes.length; i++) {
      out[i % out.length] = (out[i % out.length] + bytes[i] + i) & 0xff;
    }
    return out.map((byte) => byte.toRadixString(16).padLeft(2, '0')).join();
  }

  static String normalizeEphemeralId(String value) {
    final normalized = _normalizeHex(value);
    if (normalized == null) {
      throw ArgumentError.value(value, 'value', 'must contain at least one hex character and at most 16 bytes');
    }
    return normalized;
  }

  static Uint8List _encode(int type, String id, {bool lostHint = false, bool batteryLow = false}) {
    final hex = _normalizeHex(id);
    if (hex == null) {
      throw ArgumentError.value(id, 'id', 'must contain 16 bytes of hex');
    }
    final bytes = <int>[type, (lostHint ? 0x01 : 0x00) | (batteryLow ? 0x02 : 0x00)];
    for (var i = 0; i < hex.length; i += 2) {
      bytes.add(int.parse(hex.substring(i, i + 2), radix: 16));
    }
    return Uint8List.fromList(bytes);
  }

  static String? _normalizeHex(String value) {
    final normalized = value.toLowerCase().replaceAll(RegExp('[^a-f0-9]'), '');
    if (normalized.length == 32) return normalized;
    if (normalized.length < 32 && normalized.isNotEmpty) return normalized.padRight(32, '0');
    return null;
  }
}
