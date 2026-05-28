class SessionState {
  const SessionState({this.token, this.userId});

  final String? token;
  final String? userId;

  SessionState copyWith({String? token, String? userId}) => SessionState(
        token: token ?? this.token,
        userId: userId ?? this.userId,
      );
}

class FindMeshTag {
  const FindMeshTag({
    required this.id,
    required this.status,
    required this.label,
    this.batteryLevel,
    this.firmwareVersion,
  });

  final String id;
  final String status;
  final String label;
  final int? batteryLevel;
  final String? firmwareVersion;

  factory FindMeshTag.fromJson(Map<String, dynamic> json) => FindMeshTag(
        id: json['id'] as String,
        status: json['status'] as String? ?? 'active',
        label: json['public_label'] as String? ?? 'Item',
        batteryLevel: json['battery_level'] as int?,
        firmwareVersion: json['firmware_version'] as String?,
      );
}

class LastSeenSummary {
  const LastSeenSummary({
    required this.displayArea,
    required this.confidenceLevel,
    required this.lastSeenAt,
  });

  final String displayArea;
  final String confidenceLevel;
  final DateTime lastSeenAt;

  factory LastSeenSummary.fromJson(Map<String, dynamic> json) => LastSeenSummary(
        displayArea: json['display_area'] as String? ?? 'coarse participating area',
        confidenceLevel: json['confidence_level'] as String? ?? 'low',
        lastSeenAt: DateTime.parse(json['last_seen_at'] as String),
      );
}

class MerchantStand {
  const MerchantStand({
    required this.id,
    required this.status,
    this.firmwareVersion,
    this.lastHeartbeatAt,
  });

  final String id;
  final String status;
  final String? firmwareVersion;
  final DateTime? lastHeartbeatAt;

  factory MerchantStand.fromJson(Map<String, dynamic> json) => MerchantStand(
        id: json['id'] as String,
        status: json['status'] as String? ?? 'unknown',
        firmwareVersion: json['firmware_version'] as String?,
        lastHeartbeatAt: json['last_heartbeat_at'] == null ? null : DateTime.parse(json['last_heartbeat_at'] as String),
      );
}

class LocalObservation {
  const LocalObservation({
    required this.ephemeralId,
    required this.firstSeen,
    required this.lastSeen,
    required this.count,
  });

  final String ephemeralId;
  final DateTime firstSeen;
  final DateTime lastSeen;
  final int count;

  LocalObservation seenAgain(DateTime now) => LocalObservation(
        ephemeralId: ephemeralId,
        firstSeen: firstSeen,
        lastSeen: now,
        count: count + 1,
      );
}
