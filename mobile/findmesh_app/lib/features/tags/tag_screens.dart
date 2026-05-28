import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api_client.dart';
import '../../core/app_state.dart';
import '../../core/models.dart';
import '../lost_mode/lost_mode_screen.dart';
import '../nearby_find/nearby_finder_screen.dart';

final tagsProvider = FutureProvider<List<FindMeshTag>>((ref) async {
  final api = ref.watch(apiClientProvider);
  final raw = await api.get('/v1/tags') as List<dynamic>;
  return raw.cast<Map<String, dynamic>>().map(FindMeshTag.fromJson).toList();
});

class MyTagsScreen extends ConsumerWidget {
  const MyTagsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tags = ref.watch(tagsProvider);
    return Scaffold(
      body: tags.when(
        data: (items) => ListView(
          padding: const EdgeInsets.all(16),
          children: [
            FilledButton.icon(
              icon: const Icon(Icons.add_link),
              label: const Text('Pair tag'),
              onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const PairTagScreen())),
            ),
            const SizedBox(height: 12),
            for (final tag in items)
              Card(
                child: ListTile(
                  leading: const Icon(Icons.sell_outlined),
                  title: Text(tag.label),
                  subtitle: Text('Status: ${tag.status}'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => TagDetailScreen(tag: tag))),
                ),
              ),
          ],
        ),
        error: (error, _) => Center(child: Text(error.toString())),
        loading: () => const Center(child: CircularProgressIndicator()),
      ),
    );
  }
}

class PairTagScreen extends ConsumerStatefulWidget {
  const PairTagScreen({super.key});

  @override
  ConsumerState<PairTagScreen> createState() => _PairTagScreenState();
}

class _PairTagScreenState extends ConsumerState<PairTagScreen> {
  final serial = TextEditingController(text: 'FM-TAG-DEV-1');
  final label = TextEditingController(text: 'Keys');
  String? status;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Pair tag')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          const Text('Put the tag in pairing mode, then confirm the serial shown by BLE or NFC.'),
          TextField(controller: serial, decoration: const InputDecoration(labelText: 'Tag serial')),
          TextField(controller: label, decoration: const InputDecoration(labelText: 'Item label')),
          const SizedBox(height: 16),
          FilledButton.icon(
            icon: const Icon(Icons.link),
            label: const Text('Complete pairing'),
            onPressed: () async {
              final api = ref.read(apiClientProvider);
              await api.post('/v1/tags/pair/complete', {
                'serial': serial.text,
                'public_label': label.text,
                'firmware_version': 'tag-mobile-dev',
              });
              ref.invalidate(tagsProvider);
              setState(() => status = 'Paired');
            },
          ),
          if (status != null) Padding(padding: const EdgeInsets.only(top: 12), child: Text(status!)),
        ],
      ),
    );
  }
}

class TagDetailScreen extends StatelessWidget {
  const TagDetailScreen({required this.tag, super.key});

  final FindMeshTag tag;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(tag.label)),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          ListTile(title: const Text('Status'), subtitle: Text(tag.status)),
          ListTile(title: const Text('Firmware'), subtitle: Text(tag.firmwareVersion ?? 'unknown')),
          FilledButton.icon(
            icon: const Icon(Icons.location_searching),
            label: const Text('Last seen'),
            onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => LastSeenScreen(tag: tag))),
          ),
          const SizedBox(height: 8),
          FilledButton.icon(
            icon: const Icon(Icons.campaign_outlined),
            label: const Text('Mark lost'),
            onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => LostModeScreen(tag: tag))),
          ),
          const SizedBox(height: 8),
          OutlinedButton.icon(
            icon: const Icon(Icons.radar_outlined),
            label: const Text('Nearby finder'),
            onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => NearbyFinderScreen(tag: tag))),
          ),
        ],
      ),
    );
  }
}

class LastSeenScreen extends ConsumerWidget {
  const LastSeenScreen({required this.tag, super.key});

  final FindMeshTag tag;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final future = ref.watch(apiClientProvider).get('/v1/tags/${tag.id}/last-seen');
    return Scaffold(
      appBar: AppBar(title: const Text('Last seen')),
      body: FutureBuilder<dynamic>(
        future: future,
        builder: (context, snapshot) {
          if (!snapshot.hasData) return const Center(child: CircularProgressIndicator());
          final summary = LastSeenSummary.fromJson(snapshot.data as Map<String, dynamic>);
          return Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(summary.displayArea, style: Theme.of(context).textTheme.headlineSmall),
                const SizedBox(height: 8),
                Text('Confidence: ${summary.confidenceLevel}'),
                const SizedBox(height: 16),
                AspectRatio(
                  aspectRatio: 16 / 10,
                  child: DecoratedBox(
                    decoration: BoxDecoration(border: Border.all(color: Theme.of(context).dividerColor), borderRadius: BorderRadius.circular(8)),
                    child: const Center(child: Text('Coarse map placeholder')),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
