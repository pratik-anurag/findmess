class CoarseLocation {
  const CoarseLocation({required this.displayArea, this.latitude, this.longitude, this.precisionMeters = 500});

  final String displayArea;
  final double? latitude;
  final double? longitude;
  final int precisionMeters;
}

abstract class LocationService {
  Future<CoarseLocation> coarseLocation();
}

class MockLocationService implements LocationService {
  @override
  Future<CoarseLocation> coarseLocation() async => const CoarseLocation(displayArea: 'current coarse area');
}
