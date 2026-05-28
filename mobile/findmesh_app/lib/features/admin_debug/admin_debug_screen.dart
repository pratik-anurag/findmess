import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';
import '../../platform/ble/ble_service.dart';
import '../../platform/ble/findmesh_ble_protocol.dart';

class AdminDebugScreen extends ConsumerStatefulWidget {
  const AdminDebugScreen({super.key});

  @override
  ConsumerState<AdminDebugScreen> createState() => _AdminDebugScreenState();
}

class _AdminDebugScreenState extends ConsumerState<AdminDebugScreen> {
  final tagId = TextEditingController(text: FindMeshBleProtocol.demoEphemeralId('hackathon-tag'));
  final zoneId = TextEditingController(text: FindMeshBleProtocol.demoEphemeralId('hackathon-zone'));
  final nfcPayload = TextEditingController(text: 'findmesh://tag-found?t=demo-lost-token');
  final List<BleAdvertisement> advertisements = [];
  String? status;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final ble = ref.watch(bleServiceProvider);
    final nfc = ref.watch(nfcServiceProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('Debug')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          ListTile(
            leading: const Icon(Icons.api_outlined),
            title: const Text('API base'),
            subtitle: Text(ref.watch(apiClientProvider).baseUrl),
          ),
          ListTile(
            leading: Icon(useRealRadio ? Icons.bluetooth_connected : Icons.bluetooth_disabled),
            title: const Text('Radio mode'),
            subtitle: Text(useRealRadio ? 'Real BLE/NFC plugins enabled' : 'Mock BLE/NFC services enabled'),
          ),
          TextField(
            controller: tagId,
            decoration: const InputDecoration(labelText: 'Demo tag ephemeral ID'),
          ),
          TextField(
            controller: zoneId,
            decoration: const InputDecoration(labelText: 'Demo zone ephemeral ID'),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              FilledButton.icon(
                icon: const Icon(Icons.radar_outlined),
                label: const Text('Scan tags'),
                onPressed: () => _scan(ble.scanFindMeshTags()),
              ),
              FilledButton.icon(
                icon: const Icon(Icons.storefront_outlined),
                label: const Text('Scan zones'),
                onPressed: () => _scan(ble.scanMerchantZones()),
              ),
              OutlinedButton.icon(
                icon: const Icon(Icons.settings_bluetooth_outlined),
                label: const Text('Advertise tag'),
                onPressed: () => _run(() async {
                  await ble.startHackathonTagAdvertiser(tagId.text);
                  return 'Advertising demo FindMesh tag';
                }),
              ),
              OutlinedButton.icon(
                icon: const Icon(Icons.wifi_tethering),
                label: const Text('Advertise zone'),
                onPressed: () => _run(() async {
                  await ble.startHackathonZoneAdvertiser(zoneId.text);
                  return 'Advertising demo merchant zone';
                }),
              ),
              OutlinedButton.icon(
                icon: const Icon(Icons.stop_circle_outlined),
                label: const Text('Stop advertise'),
                onPressed: () => _run(() async {
                  await ble.stopAdvertising();
                  return 'Stopped BLE advertising';
                }),
              ),
            ],
          ),
          const Divider(height: 32),
          TextField(
            controller: nfcPayload,
            decoration: const InputDecoration(labelText: 'NFC payload'),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              FilledButton.icon(
                icon: const Icon(Icons.nfc_outlined),
                label: const Text('Read NFC'),
                onPressed: () => _run(() async {
                  final payload = await nfc.readPayload();
                  if (payload != null) nfcPayload.text = payload;
                  return payload == null ? 'No NFC payload read' : 'Read NFC payload';
                }),
              ),
              OutlinedButton.icon(
                icon: const Icon(Icons.edit_note_outlined),
                label: const Text('Write NFC'),
                onPressed: () => _run(() async {
                  await nfc.writePayload(nfcPayload.text);
                  return 'Wrote NFC payload';
                }),
              ),
            ],
          ),
          if (status != null) Padding(padding: const EdgeInsets.only(top: 12), child: Text(status!)),
          const SizedBox(height: 12),
          for (final adv in advertisements.take(10))
            ListTile(
              leading: Icon(adv.type == 'FM_TAG' ? Icons.sell_outlined : Icons.storefront_outlined),
              title: Text('${adv.type} ${adv.rssi} dBm'),
              subtitle: Text('${adv.ephemeralId}${adv.name == null ? '' : '\n${adv.name}'}'),
            ),
          const Divider(height: 32),
          FilledButton.icon(
            icon: const Icon(Icons.health_and_safety_outlined),
            label: const Text('Send demo abuse report'),
            onPressed: () => ref.read(apiClientProvider).post('/v1/abuse/reports', {
              'category': 'debug',
              'description': 'Debug report',
            }),
          ),
        ],
      ),
    );
  }

  void _scan(Stream<BleAdvertisement> stream) {
    setState(() {
      status = 'Scanning BLE for 20 seconds';
      advertisements.clear();
    });
    final subscription = stream.listen((adv) {
      if (!mounted) return;
      setState(() {
        advertisements.insert(0, adv);
        status = 'Found ${adv.type}';
      });
    }, onError: (Object error) {
      if (!mounted) return;
      setState(() => status = error.toString());
    });
    Future<void>.delayed(const Duration(seconds: 22), subscription.cancel);
  }

  Future<void> _run(Future<String> Function() action) async {
    try {
      final message = await action();
      if (!mounted) return;
      setState(() => status = message);
    } catch (error) {
      if (!mounted) return;
      setState(() => status = error.toString());
    }
  }
}
