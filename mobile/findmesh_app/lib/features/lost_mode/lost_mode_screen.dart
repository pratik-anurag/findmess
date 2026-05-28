import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';
import '../../core/models.dart';

class LostModeScreen extends ConsumerStatefulWidget {
  const LostModeScreen({required this.tag, super.key});

  final FindMeshTag tag;

  @override
  ConsumerState<LostModeScreen> createState() => _LostModeScreenState();
}

class _LostModeScreenState extends ConsumerState<LostModeScreen> {
  final message = TextEditingController(text: 'If found, contact me via FindMesh.');
  String? token;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Lost mode')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          const Text('This creates anonymous recovery options. Your phone number is not shared by default.'),
          TextField(controller: message, maxLines: 3, decoration: const InputDecoration(labelText: 'Safe message')),
          const SizedBox(height: 16),
          FilledButton.icon(
            icon: const Icon(Icons.campaign_outlined),
            label: const Text('Enable lost mode'),
            onPressed: () async {
              final response = await ref.read(apiClientProvider).post('/v1/tags/${widget.tag.id}/lost-mode', {'safe_message': message.text});
              setState(() => token = response['public_lost_token'] as String?);
            },
          ),
          if (token != null) ListTile(title: const Text('Found item token'), subtitle: Text(token!)),
        ],
      ),
    );
  }
}
