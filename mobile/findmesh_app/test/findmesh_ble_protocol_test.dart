import 'package:flutter_test/flutter_test.dart';
import 'package:findmesh_app/platform/ble/findmesh_ble_protocol.dart';

void main() {
  test('encodes and parses hackathon tag advertisements', () {
    const id = '00112233445566778899aabbccddeeff';
    final encoded = FindMeshBleProtocol.encodeTag(id, lostHint: true);
    final parsed = FindMeshBleProtocol.parseBytes(encoded, rssi: -55);

    expect(parsed, isNotNull);
    expect(parsed!.type, 'FM_TAG');
    expect(parsed.ephemeralId, id);
    expect(parsed.flags['lost_hint'], isTrue);
    expect(parsed.rssi, -55);
  });

  test('parses local-name fallback advertisements', () {
    final parsed = FindMeshBleProtocol.parseLocalName('FM_ZONE:abc123', rssi: -70);

    expect(parsed, isNotNull);
    expect(parsed!.type, 'FM_ZONE');
    expect(parsed.ephemeralId, 'abc12300000000000000000000000000');
  });
}
