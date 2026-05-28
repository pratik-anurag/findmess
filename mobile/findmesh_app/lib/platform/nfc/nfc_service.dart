abstract class NfcService {
  Future<bool> isAvailable();
  Future<String?> readPayload();
  Future<void> writePayload(String payload);
}

class MockNfcService implements NfcService {
  String? payload;

  @override
  Future<bool> isAvailable() async => true;

  @override
  Future<String?> readPayload() async => payload ?? 'findmesh://tag-found?t=demo-lost-token';

  @override
  Future<void> writePayload(String payload) async {
    this.payload = payload;
  }
}
