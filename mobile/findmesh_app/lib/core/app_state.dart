import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../platform/ble/ble_service.dart';
import '../platform/ble/real_ble_service.dart';
import '../platform/location/location_service.dart';
import '../platform/nfc/nfc_service.dart';
import '../platform/nfc/real_nfc_service.dart';
import '../platform/notifications/notification_service.dart';
import '../platform/secure_storage/secure_storage_service.dart';
import 'api_client.dart';
import 'models.dart';
import '../features/anti_stalking/anti_stalking_logic.dart';

final apiClientProvider = Provider<ApiClient>((ref) => ApiClient());
final secureStorageProvider = Provider<SecureStorageService>((ref) => InMemorySecureStorageService());
const useRealRadio = bool.fromEnvironment('FINDMESH_USE_REAL_RADIO');

final bleServiceProvider = Provider<BleService>((ref) => useRealRadio ? RealBleService() : MockBleService());
final nfcServiceProvider = Provider<NfcService>((ref) => useRealRadio ? RealNfcService() : MockNfcService());
final locationServiceProvider = Provider<LocationService>((ref) => MockLocationService());
final notificationServiceProvider = Provider<NotificationService>((ref) => LogNotificationService());
final antiStalkingProvider = Provider<AntiStalkingDetector>((ref) => AntiStalkingDetector());

final sessionProvider = StateNotifierProvider<SessionController, SessionState>((ref) {
  return SessionController(ref.watch(apiClientProvider), ref.watch(secureStorageProvider));
});

class SessionController extends StateNotifier<SessionState> {
  SessionController(this.api, this.storage) : super(const SessionState());

  final ApiClient api;
  final SecureStorageService storage;

  Future<void> startOtp(String phone) => api.post('/v1/auth/otp/start', {'phone': phone});

  Future<void> verifyOtp(String phone, String otp) async {
    final response = await api.post('/v1/auth/otp/verify', {'phone': phone, 'otp': otp});
    final token = response['token'] as String;
    api.token = token;
    await storage.write('session_token', token);
    state = SessionState(token: token, userId: (response['user'] as Map<String, dynamic>)['id'] as String?);
  }

  Future<void> logout() async {
    if (state.token != null) {
      await api.post('/v1/auth/logout', {});
    }
    api.token = null;
    await storage.delete('session_token');
    state = const SessionState();
  }
}
