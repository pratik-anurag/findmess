import '../../core/models.dart';

class BleAdvertisement {
  const BleAdvertisement({
    required this.type,
    required this.ephemeralId,
    required this.rssi,
    this.flags = const {},
    this.deviceId,
    this.name,
  });

  final String type;
  final String ephemeralId;
  final int rssi;
  final Map<String, bool> flags;
  final String? deviceId;
  final String? name;
}

abstract class BleService {
  Future<bool> isAvailable();
  Future<void> requestPermissions();
  Stream<BleAdvertisement> scanFindMeshTags();
  Stream<BleAdvertisement> scanMerchantZones();
  Future<void> connectToTag(String ephemeralId);
  Future<void> connectToStandProvisioning(String standId);
  Future<void> ringNearbyTag(String tagId);
  Future<void> startHackathonTagAdvertiser(String ephemeralId);
  Future<void> startHackathonZoneAdvertiser(String zoneEphemeralId);
  Future<void> stopAdvertising();
}

class MockBleService implements BleService {
  @override
  Future<bool> isAvailable() async => true;

  @override
  Future<void> requestPermissions() async {}

  @override
  Stream<BleAdvertisement> scanFindMeshTags() async* {
    yield const BleAdvertisement(type: 'FM_TAG', ephemeralId: '00112233445566778899aabbccddeeff', rssi: -58);
  }

  @override
  Stream<BleAdvertisement> scanMerchantZones() async* {
    yield const BleAdvertisement(type: 'FM_ZONE', ephemeralId: 'ffeeddccbbaa99887766554433221100', rssi: -65);
  }

  @override
  Future<void> connectToTag(String ephemeralId) async {}

  @override
  Future<void> connectToStandProvisioning(String standId) async {}

  @override
  Future<void> ringNearbyTag(String tagId) async {}

  @override
  Future<void> startHackathonTagAdvertiser(String ephemeralId) async {}

  @override
  Future<void> startHackathonZoneAdvertiser(String zoneEphemeralId) async {}

  @override
  Future<void> stopAdvertising() async {}
}

class NearbySignal {
  const NearbySignal(this.tag, this.rssi);

  final FindMeshTag tag;
  final int rssi;

  String get hint {
    if (rssi >= -55) return 'Very close';
    if (rssi >= -70) return 'Nearby';
    return 'Move slowly around the area';
  }
}
