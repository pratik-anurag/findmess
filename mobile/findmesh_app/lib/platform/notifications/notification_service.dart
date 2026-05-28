abstract class NotificationService {
  Future<void> showLocal(String title, String body);
}

class LogNotificationService implements NotificationService {
  @override
  Future<void> showLocal(String title, String body) async {}
}
