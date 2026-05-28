import '../../core/models.dart';

class AntiStalkingDetector {
  AntiStalkingDetector({
    this.minObservations = 4,
    this.minDuration = const Duration(minutes: 45),
  });

  final int minObservations;
  final Duration minDuration;
  final Map<String, LocalObservation> _observations = {};

  List<LocalObservation> get observations => _observations.values.toList(growable: false);

  bool observe(String ephemeralId, DateTime now) {
    final current = _observations[ephemeralId];
    if (current == null) {
      _observations[ephemeralId] = LocalObservation(ephemeralId: ephemeralId, firstSeen: now, lastSeen: now, count: 1);
      return false;
    }
    final updated = current.seenAgain(now);
    _observations[ephemeralId] = updated;
    return updated.count >= minObservations && updated.lastSeen.difference(updated.firstSeen) >= minDuration;
  }

  void prune(DateTime now) {
    _observations.removeWhere((_, value) => now.difference(value.lastSeen) > const Duration(hours: 24));
  }
}
