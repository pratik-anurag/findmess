import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';

class SettingsScreen extends ConsumerStatefulWidget {
  const SettingsScreen({super.key});

  @override
  ConsumerState<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends ConsumerState<SettingsScreen> {
  bool finderParticipation = true;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        SwitchListTile(
          value: finderParticipation,
          onChanged: (value) => setState(() => finderParticipation = value),
          title: const Text('Anonymous finder participation'),
          subtitle: const Text('Upload private lost-item sightings without exposing your identity to owners.'),
        ),
        ListTile(
          leading: const Icon(Icons.file_download_outlined),
          title: const Text('Export account data'),
          onTap: () {},
        ),
        ListTile(
          leading: const Icon(Icons.logout),
          title: const Text('Log out'),
          onTap: () => ref.read(sessionProvider.notifier).logout(),
        ),
        ListTile(
          leading: const Icon(Icons.delete_forever_outlined),
          title: const Text('Delete account'),
          onTap: () => ref.read(apiClientProvider).delete('/v1/me'),
        ),
      ],
    );
  }
}
