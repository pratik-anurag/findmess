import 'package:flutter_test/flutter_test.dart';
import 'package:findmesh_app/features/anti_stalking/anti_stalking_logic.dart';

void main() {
  test('alerts after repeated co-presence across time', () {
    final detector = AntiStalkingDetector(minObservations: 3, minDuration: const Duration(minutes: 30));
    final start = DateTime.utc(2026, 5, 28, 10);
    expect(detector.observe('abc', start), isFalse);
    expect(detector.observe('abc', start.add(const Duration(minutes: 10))), isFalse);
    expect(detector.observe('abc', start.add(const Duration(minutes: 31))), isTrue);
  });
}
