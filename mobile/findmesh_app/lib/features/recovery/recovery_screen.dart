import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';

class RecoveryRequestsScreen extends ConsumerWidget {
  const RecoveryRequestsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return FutureBuilder<dynamic>(
      future: ref.watch(apiClientProvider).get('/v1/recovery/requests'),
      builder: (context, snapshot) {
        if (!snapshot.hasData) return const Center(child: CircularProgressIndicator());
        final requests = snapshot.data as List<dynamic>;
        return ListView(
          padding: const EdgeInsets.all(16),
          children: [
            for (final req in requests.cast<Map<String, dynamic>>())
              Card(
                child: ListTile(
                  leading: const Icon(Icons.forum_outlined),
                  title: Text(req['status'] as String? ?? 'requested'),
                  subtitle: Text(req['masked_thread_id'] as String? ?? 'masked thread'),
                ),
              ),
            if (requests.isEmpty) const ListTile(title: Text('No active recovery requests')),
          ],
        );
      },
    );
  }
}
