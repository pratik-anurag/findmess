import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';

class FoundItemScreen extends ConsumerWidget {
  const FoundItemScreen({required this.publicLostToken, super.key});

  final String publicLostToken;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(title: const Text('Found item')),
      body: FutureBuilder<dynamic>(
        future: ref.watch(apiClientProvider).get('/v1/found/$publicLostToken'),
        builder: (context, snapshot) {
          if (!snapshot.hasData) return const Center(child: CircularProgressIndicator());
          final data = snapshot.data as Map<String, dynamic>;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              Text(data['tag_label'] as String? ?? 'Lost item', style: Theme.of(context).textTheme.headlineSmall),
              const SizedBox(height: 12),
              Text(data['safe_message'] as String? ?? 'Contact via FindMesh.'),
              const SizedBox(height: 20),
              FilledButton.icon(
                icon: const Icon(Icons.volunteer_activism_outlined),
                label: const Text('Report found item'),
                onPressed: () => ref.read(apiClientProvider).post('/v1/found/$publicLostToken/report', {}),
              ),
            ],
          );
        },
      ),
    );
  }
}
