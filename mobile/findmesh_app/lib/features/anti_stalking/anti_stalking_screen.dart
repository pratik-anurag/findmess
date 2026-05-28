import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';

class AntiStalkingAlertsScreen extends ConsumerWidget {
  const AntiStalkingAlertsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detector = ref.watch(antiStalkingProvider);
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        const ListTile(
          leading: Icon(Icons.shield_outlined),
          title: Text('Unknown tag alerts'),
          subtitle: Text('FindMesh looks for repeated nearby unknown tag advertisements over time.'),
        ),
        for (final obs in detector.observations)
          ListTile(
            title: Text(obs.ephemeralId),
            subtitle: Text('Seen ${obs.count} times'),
          ),
        OutlinedButton.icon(
          icon: const Icon(Icons.report_outlined),
          label: const Text('Report safety concern'),
          onPressed: () => ref.read(apiClientProvider).post('/v1/abuse/reports', {
            'category': 'unknown_tracker_alert',
            'description': 'User reported an unwanted tracker concern from mobile alert flow.',
          }),
        ),
      ],
    );
  }
}
