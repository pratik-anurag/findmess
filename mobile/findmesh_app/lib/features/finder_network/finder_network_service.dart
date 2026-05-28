import '../../core/api_client.dart';
import '../../platform/ble/ble_service.dart';

class FinderNetworkService {
  FinderNetworkService(this.api, this.ble);

  final ApiClient api;
  final BleService ble;

  Future<void> uploadOne(BleAdvertisement adv) async {
    await api.post('/v1/sightings', {
      'protocol_version': 1,
      'source_type': 'user_app',
      'tag_ephemeral_id': adv.ephemeralId,
      'time_bucket': DateTime.now().toUtc().toIso8601String(),
      'rssi_bucket': adv.rssi >= -60 ? 'near' : (adv.rssi >= -78 ? 'medium' : 'far'),
      'nonce': DateTime.now().microsecondsSinceEpoch.toString(),
    });
  }
}
