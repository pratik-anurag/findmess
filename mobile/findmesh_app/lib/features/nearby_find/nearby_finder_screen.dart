import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';
import '../../core/models.dart';

class NearbyFinderScreen extends ConsumerWidget {
  const NearbyFinderScreen({required this.tag, super.key});

  final FindMeshTag tag;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final ble = ref.watch(bleServiceProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('Nearby finder')),
      body: StreamBuilder(
        stream: ble.scanFindMeshTags(),
        builder: (context, snapshot) {
          final adv = snapshot.data;
          final rssi = adv?.rssi;
          return Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Text(tag.label, style: Theme.of(context).textTheme.headlineSmall),
                const SizedBox(height: 20),
                LinearProgressIndicator(value: rssi == null ? 0 : ((100 + rssi).clamp(0, 60) / 60)),
                const SizedBox(height: 12),
                Text(rssi == null ? 'Scanning nearby...' : 'Signal: $rssi dBm'),
                const SizedBox(height: 20),
                FilledButton.icon(
                  icon: const Icon(Icons.notifications_active_outlined),
                  label: const Text('Ring tag'),
                  onPressed: () => ble.ringNearbyTag(tag.id),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
